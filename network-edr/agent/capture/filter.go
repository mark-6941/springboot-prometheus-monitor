package capture

import (
	"encoding/binary"
	"strings"

	"github.com/example/network-edr/agent/flow"
)

func parseL7(payload []byte, dstPort uint16) flow.L7 {
	if len(payload) == 0 {
		return flow.L7{}
	}

	s := string(payload)
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")

	if strings.HasPrefix(s, "GET ") || strings.HasPrefix(s, "POST ") ||
		strings.HasPrefix(s, "PUT ") || strings.HasPrefix(s, "DELETE ") ||
		strings.HasPrefix(s, "PATCH ") || strings.HasPrefix(s, "HEAD ") {
		parts := strings.Fields(lines[0])
		l := flow.L7{Protocol: "HTTP"}
		if len(parts) >= 2 {
			l.Method = parts[0]
			l.Path = parts[1]
		}
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), "host:") {
				l.Host = strings.TrimSpace(line[5:])
			}
		}
		return l
	}

	// Minimal DNS query metadata for UDP/53.
	if dstPort == 53 && len(payload) >= 12 {
		qd := binary.BigEndian.Uint16(payload[4:6])
		if qd > 0 && len(payload) > 12 {
			i := 12
			var labels []string
			for i < len(payload) && payload[i] != 0 {
				n := int(payload[i])
				i++
				if n == 0 || i+n > len(payload) {
					break
				}
				labels = append(labels, string(payload[i:i+n]))
				i += n
			}
			if len(labels) > 0 {
				return flow.L7{Protocol: "DNS", DNSName: strings.Join(labels, ".")}
			}
		}
	}

	// TLS ClientHello SNI detection is intentionally conservative.
	if dstPort == 443 && len(payload) >= 5 && payload[0] == 0x16 && payload[1] == 0x03 {
		return flow.L7{Protocol: "TLS"}
	}

	return flow.L7{}
}
