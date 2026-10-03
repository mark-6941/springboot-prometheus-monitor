// Optional Linux eBPF extension point.
// Packet capture remains portable through libpcap in the reference build.
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_sendto")
int trace_udp_sendto(void *ctx) {
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
