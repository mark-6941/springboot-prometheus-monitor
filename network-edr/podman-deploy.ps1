$ErrorActionPreference = "Stop"

Write-Host "=================================================="
Write-Host " Network-EDR Podman Quick Deployment (PowerShell)"
Write-Host "=================================================="

# Check if Podman is installed
if (-not (Get-Command podman -ErrorAction SilentlyContinue)) {
    Write-Error "Podman is not installed or not in PATH. Please install Podman first."
    exit 1
}

# Check if Podman machine is running (especially on Windows/macOS)
$machineStatus = podman machine inspect --format "{{.State}}" 2>$null
if ($machineStatus -and $machineStatus -ne "running") {
    Write-Host "Starting Podman machine..."
    podman machine start
}

$ComposeFile = "deploy/compose.prod.yml"
if (-not (Test-Path $ComposeFile)) {
    # Fallback path if run from root vs deploy
    if (Test-Path "compose.prod.yml") {
        $ComposeFile = "compose.prod.yml"
    } else {
        throw "Compose file not found at deploy/compose.prod.yml"
    }
}

Write-Host "[1/2] Building and starting containers with Podman Compose..."
if (Get-Command podman-compose -ErrorAction SilentlyContinue) {
    podman-compose -f $ComposeFile up -d --build
} else {
    # Fallback to podman compose (newer podman plugin)
    podman compose -f $ComposeFile up -d --build
}

Write-Host "[2/2] Deployment complete!"
Write-Host "Access the Network-EDR dashboard at http://localhost:8080"
