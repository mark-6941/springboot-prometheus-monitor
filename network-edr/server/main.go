package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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
	L7        any       `json:"l7,omitempty"`
}

type Store struct {
	mu    sync.RWMutex
	flows map[string]Flow
}

var store = &Store{flows: map[string]Flow{}}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", health)
	mux.HandleFunc("/api/v1/flows", flows)
	mux.HandleFunc("/api/v1/ingest", ingest)
	mux.HandleFunc("/api/v1/stream", stream)

	addr := os.Getenv("EDR_SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("network-edr server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC(),
	})
}

func flows(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	out := make([]Flow, 0, len(store.flows))
	for _, f := range store.flows {
		out = append(out, f)
	}
	store.mu.RUnlock()
	writeJSON(w, map[string]any{"time": time.Now().UTC(), "flows": out})
}

func ingest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var f Flow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if f.Key == "" {
		f.Key = f.SrcIP + "|" + f.DstIP + "|" + f.Protocol
	}

	store.mu.Lock()
	store.flows[f.Key] = f
	store.mu.Unlock()

	w.WriteHeader(http.StatusAccepted)
}

func stream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for range t.C {
		store.mu.RLock()
		out := make([]Flow, 0, len(store.flows))
		for _, f := range store.flows {
			out = append(out, f)
		}
		store.mu.RUnlock()

		if err := conn.WriteJSON(map[string]any{
			"time":  time.Now().UTC(),
			"flows": out,
		}); err != nil {
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
