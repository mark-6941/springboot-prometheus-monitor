package ebpf

// This package is an extension boundary for Linux eBPF process/socket
// correlation. The portable capture engine deliberately does not load eBPF
// yet, so Windows and Linux share the same first-release binary behavior.
func Available() bool { return false }
