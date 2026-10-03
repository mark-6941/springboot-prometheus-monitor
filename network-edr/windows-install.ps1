$ErrorActionPreference = "Stop"

# Ensure running as Administrator
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "This script must be run as Administrator."
    exit 1
}

$InstallDir = "$env:ProgramFiles\NetworkEDR"
Write-Host "[1/6] Creating installation directory: $InstallDir"
New-Item -ItemType Directory -Force $InstallDir | Out-Null

Write-Host "[2/6] Checking Npcap..."
$npcap = Get-Service npcap -ErrorAction SilentlyContinue
if (-not $npcap) {
    throw "Npcap is required for packet capture on Windows. Please install Npcap first."
}
if ($npcap.Status -ne "Running") {
    Start-Service npcap
}

Write-Host "[3/6] Building and copying agent..."
go build -trimpath -ldflags="-s -w" -o "$InstallDir\network-edr-agent.exe" .\agent

Write-Host "[4/6] Creating environment and config files..."
$EnvConfig = @"
EDR_AGENT_ID=$env:COMPUTERNAME
EDR_INTERFACE=1
EDR_SERVER=http://127.0.0.1:8080
EDR_SNAPLEN=2048
EDR_PROMISC=true
EDR_BPF=tcp or udp
EDR_FLOW_TTL_SEC=60
"@
$EnvConfig | Set-Content "$InstallDir\agent.env"

# ACL Locking: Restrict access to administrators and SYSTEM only
Write-Host "[5/6] Securing configuration files with ACL locking..."
$Acl = Get-Acl $InstallDir
$Acl.SetAccessRuleProtection($true, $false) # Disable inheritance, remove inherited rules
$AdminRule = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrators","FullControl","ContainerInherit,ObjectInherit","None","Allow")
$SystemRule = New-Object System.Security.AccessControl.FileSystemAccessRule("NT AUTHORITY\SYSTEM","FullControl","ContainerInherit,ObjectInherit","None","Allow")
$Acl.SetAccessRule($AdminRule)
$Acl.SetAccessRule($SystemRule)
Set-Acl $InstallDir $Acl

Write-Host "[6/6] Configuring Windows Scheduled Task for SYSTEM boot start & auto-restart..."
$TaskName = "NetworkEDRAgent"
Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue

$Action = New-ScheduledTaskAction -Execute "$InstallDir\network-edr-agent.exe" -WorkingDirectory $InstallDir
$Trigger = New-ScheduledTaskTrigger -AtStartup
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -RestartCount 3 -RestartInterval (New-TimeSpan -Minutes 1) -ExecutionTimeLimit (New-TimeSpan -Days 365)
$Principal = New-ScheduledTaskPrincipal -UserId "NT AUTHORITY\SYSTEM" -LogonType ServiceAccount -RunLevel Highest

Register-ScheduledTask -TaskName $TaskName -Action $Action -Trigger $Trigger -Settings $Settings -Principal $Principal -Force

# Start the scheduled task immediately
Start-ScheduledTask -TaskName $TaskName

Write-Host "========================================================"
Write-Host "Installation completed successfully!"
Write-Host "Service '$TaskName' is running under NT AUTHORITY\SYSTEM."
Write-Host "========================================================"
