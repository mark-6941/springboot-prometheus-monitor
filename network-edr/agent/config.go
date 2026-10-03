package main

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AgentID   string
	Interface string
	ServerURL string
	SnapLen   int
	Promisc   bool
	BPF       string
	FlowTTL   time.Duration
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	v, err := strconv.Atoi(env(k, strconv.Itoa(def)))
	if err != nil {
		return def
	}
	return v
}

func envBool(k string, def bool) bool {
	v := env(k, strconv.FormatBool(def))
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func LoadConfig() Config {
	host, _ := os.Hostname()
	return Config{
		AgentID:   env("EDR_AGENT_ID", host),
		Interface: env("EDR_INTERFACE", "any"),
		ServerURL: env("EDR_SERVER", "http://127.0.0.1:8080"),
		SnapLen:   envInt("EDR_SNAPLEN", 2048),
		Promisc:   envBool("EDR_PROMISC", true),
		BPF:       env("EDR_BPF", "tcp or udp"),
		FlowTTL:   time.Duration(envInt("EDR_FLOW_TTL_SEC", 60)) * time.Second,
	}
}
