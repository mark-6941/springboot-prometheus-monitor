package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/network-edr/agent/capture"
	"github.com/example/network-edr/agent/exporter"
	"github.com/example/network-edr/agent/flow"
)

func main() {
	cfg := LoadConfig()
	log.Printf("network-edr agent=%s interface=%s server=%s", cfg.AgentID, cfg.Interface, cfg.ServerURL)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	agg := flow.NewAggregator(cfg.FlowTTL)
	go agg.Run(ctx)

	capCfg := capture.Config{
		Interface: cfg.Interface,
		SnapLen:   cfg.SnapLen,
		Promisc:   cfg.Promisc,
		BPF:       cfg.BPF,
		AgentID:   cfg.AgentID,
	}

	go func() {
		if err := capture.Run(ctx, capCfg, agg.Add); err != nil {
			log.Printf("capture stopped: %v", err)
			cancel()
		}
	}()

	exp := exporter.NewHTTP(cfg.ServerURL)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, f := range agg.Snapshot() {
				if err := exp.Send(ctx, f); err != nil {
					log.Printf("export %s: %v", f.Key, err)
				}
			}
		}
	}
}
