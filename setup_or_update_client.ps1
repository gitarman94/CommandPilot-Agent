Param(
    [string]$ServerUrl,
    [switch]$Install,
    [switch]$Update,
    [switch]$Uninstall,
    [switch]$VerboseMode,
    [switch]$Silent
)

$ErrorActionPreference="Stop"

$GitHubUser="gitarman94"
$GitHubRepo="CommandPilot"
$Branch="main"

$InstallDir="C:\CommandPilot_Agent"
$BinaryName="pilot-agent.exe"
$ServiceName="CommandPilotAgent"

$BinaryUrl="https://github.com/$GitHubUser/$GitHubRepo/raw/$Branch/windows-agent/$BinaryName"

$BinaryPath=Join-Path $InstallDir $BinaryName
$ServerFile=Join-Path $InstallDir "server_url.txt"
$VersionFile=Join-Path $InstallDir "version.txt"
$AgentIdFile=Join-Path $InstallDir "agent_id.txt"

function Write-Stage($msg) {
    if (-not $Silent) {
        Write-Host ""
        Write-Host "== $msg =="
    }
}

function Write-Pass($msg) {
    if (-not $Silent) {
        Write-Host "[PASS] $msg"
    }
}

function Write-Fail($msg) {
    Write-Host "[FAIL] $msg"
    exit 1
}

function Invoke-Download {
    param(
        [string]$Url,
        [string]$Destination
    )

    if ($VerboseMode -and (-not $Silent)) {
        Write-Host "[DOWNLOAD] $Url"
    }

    Invoke-WebRequest -Uri $Url -OutFile $Destination -UseBasicParsing
}

function Get-FileHashString {
    param([string]$Path)

    if (Test-Path $Path) {
        return (Get-FileHash -Path $Path -Algorithm SHA256).Hash
    }

    return ""
}

function Uninstall-Agent {

    Write-Stage "Uninstall"

    if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {
        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        sc.exe delete $ServiceName | Out-Null
        Write-Pass "Windows service removed"
    }

    if (Test-Path $InstallDir) {
        Remove-Item -Path $InstallDir -Recurse -Force
    }

    Write-Pass "Installation directory removed"

    if (-not $Silent) {
        Write-Host ""
        Write-Host "CommandPilot Agent removed"
    }
}

function Install-Dependencies {

    Write-Stage "Dependencies"

    if (-not (Get-Command curl.exe -ErrorAction SilentlyContinue)) {
        Write-Fail "curl.exe missing"
    }

    Write-Pass "Dependencies ready"
}

function Prepare-Directories {

    Write-Stage "Directories"

    if (-not (Test-Path $InstallDir)) {
        New-Item -Path $InstallDir -ItemType Directory | Out-Null
    }

    Write-Pass "Directories ready"
}

function Configure-Server {

    Write-Stage "Server Configuration"

    if ((-not $ServerUrl) -and (-not (Test-Path $ServerFile))) {

        if ($Silent) {
            Write-Fail "Silent mode requires -ServerUrl"
        }

        Write-Host ""
        Write-Host "Enter CommandPilot server hostname or IP"
        Write-Host "Examples:"
        Write-Host "  192.168.1.10"
        Write-Host "  commandpilot.local"
        Write-Host "  https://commandpilot.example.com"
        Write-Host ""

        $ServerUrl=Read-Host "Server"
    }

    if (-not $ServerUrl) {
        $ServerUrl=Get-Content $ServerFile
    }

    $ServerUrl=$ServerUrl.Trim()

    if ($ServerUrl -notmatch '^https?://') {
        $ServerUrl="http://$ServerUrl"
    }

    if ($ServerUrl -notmatch ':\d+$') {
        $ServerUrl="$ServerUrl`:8080"
    }

    $ServerUrl | Out-File -FilePath $ServerFile -Encoding ASCII -Force

    Write-Pass "Server URL saved"

    if (-not $Silent) {
        Write-Host "Server: $ServerUrl"
    }
}

function Initialize-AgentId {

    Write-Stage "Agent Identity"

    if (-not (Test-Path $AgentIdFile)) {

        [guid]::NewGuid().ToString() |
            Out-File -FilePath $AgentIdFile -Encoding ASCII -Force

        Write-Pass "Agent ID created"
    }
    else {
        Write-Pass "Agent ID already exists"
    }
}

function Install-Binary {

    Write-Stage "Binary"

    $TempBinary=Join-Path $env:TEMP "pilot-agent.exe"

    Invoke-Download -Url $BinaryUrl -Destination $TempBinary

    $RemoteHash=Get-FileHashString $TempBinary
    $LocalHash=Get-FileHashString $BinaryPath

    if ($RemoteHash -ne $LocalHash) {

        Copy-Item -Path $TempBinary -Destination $BinaryPath -Force

        Write-Pass "Binary updated"
    }
    else {
        Write-Pass "Binary already current"
    }

    Remove-Item -Path $TempBinary -Force -ErrorAction SilentlyContinue
}

function Install-Service {

    Write-Stage "Windows Service"

    if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {

        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue

        sc.exe delete $ServiceName | Out-Null

        Start-Sleep -Seconds 2
    }

    $ServiceCommand="`"$BinaryPath`""

    sc.exe create `
        $ServiceName `
        binPath= $ServiceCommand `
        start= auto | Out-Null

    sc.exe description `
        $ServiceName `
        "CommandPilot Agent" | Out-Null

    Start-Service $ServiceName

    Start-Sleep -Seconds 3

    $svc=Get-Service $ServiceName

    if ($svc.Status -ne "Running") {
        Write-Fail "Service failed to start"
    }

    Write-Pass "Service running"
}

function Validate-Install {

    Write-Stage "Validation"

    if (-not (Test-Path $BinaryPath)) {
        Write-Fail "Binary missing"
    }

    Write-Pass "Binary exists"

    if (-not (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue)) {
        Write-Fail "Service missing"
    }

    Write-Pass "Service exists"

    try {

        $ServerUrl=Get-Content $ServerFile

        Invoke-WebRequest `
            -Uri $ServerUrl `
            -UseBasicParsing `
            -TimeoutSec 10 | Out-Null

        Write-Pass "Server reachable"
    }
    catch {
        Write-Fail "Server unreachable"
    }
}

if ($Uninstall) {
    Uninstall-Agent
    exit 0
}

if (-not $Silent) {
    Write-Host ""
    Write-Host "======================================"
    Write-Host " CommandPilot Agent Installer"
    Write-Host "======================================"
}

Install-Dependencies
Prepare-Directories
Configure-Server
Initialize-AgentId
Install-Binary
Install-Service
Validate-Install

if (-not $Silent) {
    Write-Host ""
    Write-Host "======================================"
    Write-Host " CommandPilot Agent Installed"
    Write-Host "======================================"
    Write-Host "Install Path : $InstallDir"
    Write-Host "Binary       : $BinaryPath"
    Write-Host "Service      : $ServiceName"
    Write-Host "======================================"
}