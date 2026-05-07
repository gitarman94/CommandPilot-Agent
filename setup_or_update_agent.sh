#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pilot-agent"
SERVICE_NAME="pilot-agent.service"

APP_DIR="/opt/commandpilot-agent"
SRC_DIR="/tmp/commandpilot-agent-src"

REPO_URL="https://github.com/gitarman94/CommandPilot.git"

CLIENT_BINARY="${APP_DIR}/pilot-agent"

VERBOSE=false

cleanup() {
    rm -rf "${SRC_DIR}"
    rm -f /tmp/go.tar.gz
}

trap cleanup EXIT

echo "CommandPilot Agent Installer"

for arg in "$@"; do
    case "$arg" in
        --verbose|-v)
            VERBOSE=true
            ;;
        *)
            echo "Usage: $0 [--verbose]"
            exit 1
            ;;
    esac
done

run() {
    if [[ "$VERBOSE" == true ]]; then
        echo "[RUN] $*"
        "$@"
    else
        "$@" >/dev/null 2>&1
    fi
}

pass() {
    echo "[PASS] $1"
}

fail() {
    echo "[FAIL] $1"
    exit 1
}

if [[ -f /etc/os-release ]]; then
    . /etc/os-release

    case "$ID" in
        debian|ubuntu|linuxmint|pop|raspbian)
            ;;
        *)
            fail "Unsupported distribution"
            ;;
    esac
else
    fail "/etc/os-release missing"
fi

echo
echo "== Dependencies =="

export DEBIAN_FRONTEND=noninteractive

run apt-get update -y || fail "apt update failed"

run apt-get install -y \
    curl \
    wget \
    git \
    unzip \
    sqlite3 \
    build-essential \
    ca-certificates || fail "dependency install failed"

pass "Dependencies installed"

echo
echo "== Go Installation =="

if ! command -v go >/dev/null 2>&1; then

    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        *)
            fail "Unsupported architecture: ${ARCH}"
            ;;
    esac

    wget -q \
        "https://go.dev/dl/go1.25.0.linux-${GO_ARCH}.tar.gz" \
        -O /tmp/go.tar.gz \
        || fail "Go download failed"

    rm -rf /usr/local/go

    tar -C /usr/local -xzf /tmp/go.tar.gz \
        || fail "Go extraction failed"

    echo 'export PATH=$PATH:/usr/local/go/bin' \
        >/etc/profile.d/golang.sh
fi

export PATH=$PATH:/usr/local/go/bin

go version >/dev/null 2>&1 \
    || fail "Go not available"

pass "Go ready"

echo
echo "== User =="

if ! id -u commandpilot >/dev/null 2>&1; then

    useradd \
        --system \
        --no-create-home \
        --shell /usr/sbin/nologin \
        commandpilot \
        || fail "failed to create service user"
fi

pass "User ready"

echo
echo "== Directories =="

mkdir -p "${APP_DIR}"
mkdir -p "${APP_DIR}/logs"
mkdir -p "${APP_DIR}/updates"

chmod 755 "${APP_DIR}"

pass "Directories ready"

echo
echo "== Server Configuration =="

SERVER_FILE="${APP_DIR}/server_url.txt"

echo
echo "Enter CommandPilot server hostname or IP"
echo "Examples:"
echo "  192.168.1.10"
echo "  commandpilot.local"
echo "  https://commandpilot.example.com"
echo

while true; do

    printf "Server: "

    IFS= read -r SERVER_INPUT

    SERVER_INPUT="$(echo "$SERVER_INPUT" | xargs)"

    if [[ -n "$SERVER_INPUT" ]]; then
        break
    fi

    echo
    echo "[FAIL] Server address required"
    echo
done

if [[ "$SERVER_INPUT" =~ ^https?:// ]]; then
    SERVER_URL="$SERVER_INPUT"
else
    SERVER_URL="http://${SERVER_INPUT}"
fi

if [[ ! "$SERVER_URL" =~ :[0-9]+$ ]]; then
    SERVER_URL="${SERVER_URL}:8080"
fi

echo "$SERVER_URL" > "$SERVER_FILE"

pass "Server URL saved"

echo
echo "== Source =="

rm -rf "${SRC_DIR}"

run git clone "$REPO_URL" "$SRC_DIR" \
    || fail "git clone failed"

AGENT_SRC=""

if [[ -f "${SRC_DIR}/go.mod" ]]; then
    AGENT_SRC="${SRC_DIR}"

elif [[ -f "${SRC_DIR}/pilot-agent/go.mod" ]]; then
    AGENT_SRC="${SRC_DIR}/pilot-agent"

else

    FOUND_GOMOD=$(
        find "${SRC_DIR}" \
            -maxdepth 3 \
            -type f \
            -name go.mod \
            | head -n 1
    )

    if [[ -n "$FOUND_GOMOD" ]]; then
        AGENT_SRC="$(dirname "$FOUND_GOMOD")"
    fi
fi

[[ -n "$AGENT_SRC" ]] \
    || fail "pilot-agent source missing"

echo "Using source: ${AGENT_SRC}"

pass "Repository synced"

echo
echo "== Build =="

cd "$AGENT_SRC"

run go mod tidy \
    || fail "go mod tidy failed"

rm -f "${AGENT_SRC}/pilot-agent"

run go build -o pilot-agent . \
    || fail "go build failed"

[[ -f "${AGENT_SRC}/pilot-agent" ]] \
    || fail "compiled binary missing"

chmod +x "${AGENT_SRC}/pilot-agent"

pass "Build succeeded"

echo
echo "== Install =="

systemctl stop "${SERVICE_NAME}" >/dev/null 2>&1 || true

install -m 755 \
    "${AGENT_SRC}/pilot-agent" \
    "${CLIENT_BINARY}" \
    || fail "binary install failed"

chown -R commandpilot:commandpilot "${APP_DIR}"

[[ -x "${CLIENT_BINARY}" ]] \
    || fail "binary not executable"

pass "Binary installed"

echo
echo "== Service =="

cat > "/etc/systemd/system/${SERVICE_NAME}" <<EOF
[Unit]
Description=CommandPilot Agent
After=network.target

[Service]
Type=simple
User=commandpilot
Group=commandpilot
WorkingDirectory=${APP_DIR}
ExecStart=${CLIENT_BINARY}
Restart=always
RestartSec=5
Environment=COMMANDPILOT_SERVER=${SERVER_URL}

[Install]
WantedBy=multi-user.target
EOF

run systemctl daemon-reload \
    || fail "systemd reload failed"

run systemctl enable "${SERVICE_NAME}" \
    || fail "service enable failed"

run systemctl restart "${SERVICE_NAME}" \
    || fail "service restart failed"

sleep 5

systemctl is-active --quiet "${SERVICE_NAME}" \
    || fail "service failed to start"

pass "Service running"

echo
echo "== Validation =="

pgrep -f pilot-agent >/dev/null 2>&1 \
    || fail "agent process not running"

pass "Process active"

curl -fsS "${SERVER_URL}" >/dev/null 2>&1 \
    || fail "server unreachable"

pass "Server reachable"

[[ -f "${SERVER_FILE}" ]] \
    || fail "server configuration missing"

pass "Configuration persisted"

echo
echo "======================================"
echo " CommandPilot Agent Installed"
echo "======================================"
echo "Service: ${SERVICE_NAME}"
echo "Binary: ${CLIENT_BINARY}"
echo "Server: ${SERVER_URL}"
echo "======================================"