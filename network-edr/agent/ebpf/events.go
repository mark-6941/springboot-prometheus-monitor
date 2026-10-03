package ebpf

type NetworkEvent struct {
	PID       uint32
	UID       uint32
	SrcIP     string
	DstIP     string
	SrcPort   uint16
	DstPort   uint16
	Protocol  uint8
	Timestamp uint64
}
