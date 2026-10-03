package flow

import "time"

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	AgentID   string    `json:"agent_id"`
	Interface string    `json:"interface"`
	SrcIP     string    `json:"src_ip"`
	DstIP     string    `json:"dst_ip"`
	SrcPort   uint16    `json:"src_port"`
	DstPort   uint16    `json:"dst_port"`
	Protocol  string    `json:"protocol"`
	Bytes     int       `json:"bytes"`
	Payload   int       `json:"payload_bytes"`
	L7        L7        `json:"l7"`
}

type L7 struct {
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Path     string `json:"path,omitempty"`
	Method   string `json:"method,omitempty"`
	SNI      string `json:"sni,omitempty"`
	DNSName  string `json:"dns_name,omitempty"`
}

type Flow struct {
	Key       string    `json:"key"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	AgentID   string    `json:"agent_id"`
	Interface string    `json:"interface"`
	SrcIP     string    `json:"src_ip"`
	DstIP     string    `json:"dst_ip"`
	SrcPort   uint16    `json:"src_port"`
	DstPort   uint16    `json:"dst_port"`
	Protocol  string    `json:"protocol"`
	Packets   uint64    `json:"packets"`
	Bytes     uint64    `json:"bytes"`
	L7        L7        `json:"l7,omitempty"`
}
