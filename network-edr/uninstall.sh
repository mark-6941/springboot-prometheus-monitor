#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

systemctl disable --now network-edr-agent.service 2>/dev/null || true
rm -f /etc/systemd/system/network-edr-agent.service
systemctl daemon-reload

if command -v podman >/dev/null; then
  podman compose -f "$ROOT/deploy/compose.prod.yml" down 2>/dev/null || true
fi

rm -rf /opt/network-edr /etc/network-edr
echo "Network EDR removed. Existing OTel/Loki/Prometheus configuration was not modified."
