package capture

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/example/network-edr/agent/flow"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Config struct {
	Interface string
	SnapLen   int
	Promisc   bool
	BPF       string
	AgentID   string
}

func Run(ctx context.Context, cfg Config, emit func(flow.Event)) error {
	iface := cfg.Interface
	if iface == "" {
		iface = "any"
	}

	handle, err := pcap.OpenLive(iface, int32(cfg.SnapLen), cfg.Promisc, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("pcap open %s: %w", iface, err)
	}
	defer handle.Close()

	if cfg.BPF != "" {
		if err := handle.SetBPFFilter(cfg.BPF); err != nil {
			return fmt.Errorf("bpf filter: %w", err)
		}
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true

	for {
		select {
		case <-ctx.Done():
			return nil
		case packet, ok := <-source.Packets():
			if !ok {
				return nil
			}
			e, ok := decode(packet, iface, cfg.AgentID)
			if ok {
				emit(e)
			}
		}
	}
}

func decode(packet gopacket.Packet, iface, cfgAgentID string) (flow.Event, bool) {
	ip4, ok4 := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	ip6, ok6 := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
	if !ok4 && !ok6 {
		return flow.Event{}, false
	}

	var src, dst string
	if ok4 {
		src, dst = ip4.SrcIP.String(), ip4.DstIP.String()
	} else {
		src, dst = ip6.SrcIP.String(), ip6.DstIP.String()
	}

	var proto string
	var sp, dp uint16
	var payload []byte

	if l := packet.Layer(layers.LayerTypeTCP); l != nil {
		t := l.(*layers.TCP)
		proto, sp, dp = "TCP", uint16(t.SrcPort), uint16(t.DstPort)
		payload = t.Payload
	} else if l := packet.Layer(layers.LayerTypeUDP); l != nil {
		u := l.(*layers.UDP)
		proto, sp, dp = "UDP", uint16(u.SrcPort), uint16(u.DstPort)
		payload = u.Payload
	} else {
		return flow.Event{}, false
	}

	if net.ParseIP(src) == nil || net.ParseIP(dst) == nil {
		return flow.Event{}, false
	}

	now := packet.Metadata().Timestamp
	if now.IsZero() {
		now = timeNow()
	}

	return flow.Event{
		Timestamp: now,
		AgentID: cfgAgentID,
		Interface: iface,
		SrcIP: src, DstIP: dst,
		SrcPort: sp, DstPort: dp,
		Protocol: proto,
		Bytes:    len(packet.Data()),
		Payload:  len(payload),
		L7:       parseL7(payload, dp),
	}, true
}

// isolated for easy testing/mocking
var timeNow = func() time.Time { return time.Now() }
