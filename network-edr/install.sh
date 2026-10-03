#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PREFIX="/opt/network-edr"

echo "[1/7] Checking tools..."
command -v podman >/dev/null || { echo "ERROR: podman is required"; exit 1; }
command -v go >/dev/null || { echo "ERROR: Go is required to build the host agent"; exit 1; }

echo "[2/7] Building agent..."
mkdir -p "$PREFIX/bin" /etc/network-edr
go build -trimpath -ldflags="-s -w" -o "$PREFIX/bin/network-edr-agent" "$ROOT/agent"

echo "[3/7] Detecting interface..."
IFACE="${EDR_INTERFACE:-}"
if [[ -z "$IFACE" ]]; then
  IFACE="$(ip -o route show default 2>/dev/null | awk '{print $5; exit}')"
fi
IFACE="${IFACE:-any}"

cat > /etc/network-edr/agent.env <<EOF
EDR_AGENT_ID=$(hostname)
EDR_INTERFACE=$IFACE
EDR_SERVER=http://127.0.0.1:8080
EDR_SNAPLEN=2048
EDR_PROMISC=true
EDR_BPF=tcp or udp
EDR_FLOW_TTL_SEC=60
EOF

echo "[4/7] Installing systemd service..."
install -Dm0644 "$ROOT/deploy/systemd/network-edr-agent.service" \
  /etc/systemd/system/network-edr-agent.service

echo "[5/7] Starting Podman control plane..."
podman compose -f "$ROOT/deploy/compose.prod.yml" up -d --build

echo "[6/7] Starting agent..."
systemctl daemon-reload
systemctl enable --now network-edr-agent.service

echo "[7/7] Health check..."
curl -fsS http://127.0.0.1:8088/api/v1/health >/dev/null || {
  echo "WARNING: Web health check failed; inspect: podman compose -f $ROOT/deploy/compose.prod.yml ps"
  exit 1
}

echo
echo "=============================================="
echo " Network EDR installed"
echo " Dashboard: http://$(hostname -I | awk '{print $1}'):8088/"
echo " Agent:     systemctl status network-edr-agent"
echo " Logs:      journalctl -u network-edr-agent -f"
echo "=============================================="
