# Network EDR

Cross-platform reference implementation for real host network visibility.

## Scope

- Go host agent
- Layer 3/4 packet metadata
- Conservative Layer 7 metadata detection for HTTP, DNS and TLS
- Go HTTP API + WebSocket
- Agent → server HTTP flow export
- React dashboard
- Podman control plane
- Linux systemd installation
- Windows/Npcap build path

## Important architecture decision

The portable reference agent uses GoPacket/libpcap because this provides a common capture path on Linux and Windows. On Linux, the next production stage can replace/augment capture with eBPF; on Windows, eBPF-for-Windows is not feature-equivalent to Linux and Linux eBPF binaries are not portable to Windows.

The agent intentionally does NOT modify an existing OpenTelemetry Collector, Loki, Prometheus, Tempo, rsyslog or firewall configuration.

## Linux

Requirements:

- Linux
- Go 1.24+
- libpcap development/runtime package
- Podman + podman-compose support
- root for packet capture

Example:

```bash
git clone <repo>
cd network-edr
./install.sh
```

Then open:

`http://<host>:8088/`

Check:

```bash
systemctl status network-edr-agent
journalctl -u network-edr-agent -f
podman ps
curl http://127.0.0.1:8088/api/v1/health
```

## Windows

Install Npcap first, then run PowerShell as Administrator:

```powershell
.\windows-install.ps1
```

The Windows agent uses Npcap through GoPacket/libpcap. The interface value may need to be changed to the adapter index/name visible to the installed Npcap/libpcap environment.

## Real-data rule

The dashboard never fabricates flows. If there is no packet data it displays zero/empty state.

## Production hardening still required

Before internet exposure:

- TLS
- OIDC/JWT authentication
- RBAC
- persistent database
- signed releases
- SBOM
- image scanning
- secrets management
- rate limiting
- audit logging
- packet sampling/backpressure
- strict origin policy for WebSocket
- process/PID correlation
- full TLS SNI/ALPN parser
- protocol parsers with privacy controls

Do not capture or persist credentials, cookies, authorization headers, or full application payloads by default.


## Windows Installation & Deployment

1. **Prerequisites**: Ensure **Npcap** is installed and running on the Windows host.
2. **Run Installation Script**: Open PowerShell as Administrator and run:
   ```powershell
   .\windows-install.ps1
   ```
   This will:
   - Build and install the agent to `C:\ProgramFiles\NetworkEDR`.
   - Set up strict ACL locking on configuration files (restricted to Administrators and SYSTEM).
   - Register a Windows Scheduled Task (`NetworkEDRAgent`) configured to run under `NT AUTHORITY\SYSTEM` at startup with automatic failure recovery.

## OpenTelemetry Integration

- An OTel Collector snippet is provided at `deploy/otel-collector-snippet.yaml` to ingest metrics, alerts, and TLS logs via OTLP (gRPC/HTTP).


## Podman Quick Deployment

You can quickly deploy and run Network-EDR using **Podman** across Windows, macOS, or Linux using either PowerShell or CMD scripts provided in the root directory:

- **PowerShell**:
  ```powershell
  ./podman-deploy.ps1
  ```
- **CMD (Command Prompt)**:
  ```cmd
  podman-deploy.cmd
  ```
