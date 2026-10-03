package flow

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Aggregator struct {
	mu   sync.RWMutex
	data map[string]*Flow
	ttl  time.Duration
}

func NewAggregator(ttl time.Duration) *Aggregator {
	return &Aggregator{data: make(map[string]*Flow), ttl: ttl}
}

func key(e Event) string {
	return fmt.Sprintf("%s|%s|%d|%d|%s", e.SrcIP, e.DstIP, e.SrcPort, e.DstPort, e.Protocol)
}

func (a *Aggregator) Add(e Event) {
	k := key(e)
	a.mu.Lock()
	defer a.mu.Unlock()

	f := a.data[k]
	if f == nil {
		f = &Flow{
			Key: k, FirstSeen: e.Timestamp, AgentID: e.AgentID,
			Interface: e.Interface, SrcIP: e.SrcIP, DstIP: e.DstIP,
			SrcPort: e.SrcPort, DstPort: e.DstPort, Protocol: e.Protocol,
		}
		a.data[k] = f
	}
	f.LastSeen = e.Timestamp
	f.Packets++
	f.Bytes += uint64(e.Bytes)
	if e.L7.Protocol != "" {
		f.L7 = e.L7
	}
}

func (a *Aggregator) Run(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			a.mu.Lock()
			for k, f := range a.data {
				if now.Sub(f.LastSeen) > a.ttl {
					delete(a.data, k)
				}
			}
			a.mu.Unlock()
		}
	}
}

func (a *Aggregator) Snapshot() []Flow {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Flow, 0, len(a.data))
	for _, f := range a.data {
		out = append(out, *f)
	}
	return out
}
