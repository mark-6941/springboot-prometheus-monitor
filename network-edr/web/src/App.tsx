import React, { useEffect, useMemo, useState } from "react";

type Flow = {
  key: string; first_seen: string; last_seen: string;
  agent_id: string; interface: string;
  src_ip: string; dst_ip: string; src_port: number; dst_port: number;
  protocol: string; packets: number; bytes: number;
  l7?: { protocol?: string; host?: string; path?: string; method?: string; sni?: string; dns_name?: string };
};

const API = import.meta.env.VITE_API_BASE || "";

function fmtBytes(n:number) {
  if (n < 1024) return `${n} B`;
  if (n < 1024*1024) return `${(n/1024).toFixed(1)} KB`;
  if (n < 1024*1024*1024) return `${(n/1024/1024).toFixed(1)} MB`;
  return `${(n/1024/1024/1024).toFixed(1)} GB`;
}

export default function App() {
  const [flows,setFlows] = useState<Flow[]>([]);
  const [health,setHealth] = useState("checking");

  useEffect(() => {
    fetch(`${API}/api/v1/health`).then(r => r.ok ? setHealth("online") : setHealth("error"))
      .catch(() => setHealth("offline"));

    const wsURL = (location.protocol === "https:" ? "wss://" : "ws://") +
      location.host + `${API}/api/v1/stream`;

    const ws = new WebSocket(wsURL);
    ws.onmessage = e => {
      try { setFlows(JSON.parse(e.data).flows || []); } catch {}
    };
    ws.onerror = () => setHealth("error");
    return () => ws.close();
  }, []);

  const totalBytes = useMemo(() => flows.reduce((n,f)=>n+f.bytes,0), [flows]);

  return <div className="app">
    <header>
      <div>
        <h1>Network EDR</h1>
        <div className="sub">Layer 3 / Layer 4 / Layer 7 Network Visibility</div>
      </div>
      <div className={`status ${health}`}>● {health}</div>
    </header>

    <section className="cards">
      <div className="card"><span>Active Flows</span><b>{flows.length}</b></div>
      <div className="card"><span>Total Bytes</span><b>{fmtBytes(totalBytes)}</b></div>
      <div className="card"><span>Agents</span><b>{new Set(flows.map(f=>f.agent_id)).size}</b></div>
      <div className="card"><span>Protocols</span><b>{new Set(flows.map(f=>f.protocol)).size}</b></div>
    </section>

    <section className="panel">
      <h2>Live Network Flow</h2>
      <div className="table">
        <div className="tr th"><span>Source</span><span>Destination</span><span>L4</span><span>L7</span><span>Packets</span><span>Bytes</span></div>
        {flows.map(f =>
          <div className="tr" key={f.key}>
            <span>{f.src_ip}:{f.src_port}</span>
            <span>{f.dst_ip}:{f.dst_port}</span>
            <span>{f.protocol}</span>
            <span>{f.l7?.protocol || "-"} {f.l7?.host || f.l7?.dns_name || ""}</span>
            <span>{f.packets}</span>
            <span>{fmtBytes(f.bytes)}</span>
          </div>
        )}
        {flows.length === 0 && <div className="empty">目前沒有收到真實封包資料</div>}
      </div>
    </section>

    <footer>
      Data source: host packet capture → Go Agent → Go API. No mock traffic.
    </footer>
  </div>
}
