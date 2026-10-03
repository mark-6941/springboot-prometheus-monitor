package process

// Linux process/socket correlation is intentionally kept separate from packet
// capture so the portable core can run on Windows with Npcap.
// A production Linux implementation can add /proc + eBPF socket correlation here.
func PIDForConnection(srcIP, dstIP string, srcPort, dstPort uint16) (int, string) {
	return 0, ""
}
