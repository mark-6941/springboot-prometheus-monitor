// Optional Linux eBPF extension point.
// The portable reference agent currently uses libpcap.
// Do not load this program until kernel/version/attach strategy is validated.
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("tracepoint/sock/inet_sock_set_state")
int trace_tcp_state(void *ctx) {
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
