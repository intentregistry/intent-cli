#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
PORT="${SMOKE_PORT:-18080}"
API_URL="http://127.0.0.1:${PORT}"
INTENT_BIN="${ROOT_DIR}/dist/intent-smoke"
MOCK_API_BIN="${ROOT_DIR}/dist/mock-api-smoke"
MOCK_API_LOG="${TMP_DIR}/mock-api.log"

cleanup() {
  if [[ -n "${MOCK_API_PID:-}" ]] && kill -0 "${MOCK_API_PID}" 2>/dev/null; then
    kill "${MOCK_API_PID}" || true
    wait "${MOCK_API_PID}" 2>/dev/null || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

mkdir -p "${ROOT_DIR}/dist"

go build -o "${INTENT_BIN}" "${ROOT_DIR}/cmd/intent"
go build -o "${MOCK_API_BIN}" "${ROOT_DIR}/cmd/mock-api"

PORT="${PORT}" "${MOCK_API_BIN}" >"${MOCK_API_LOG}" 2>&1 &
MOCK_API_PID=$!

for _ in {1..50}; do
  if curl -sf "${API_URL}/health" >/dev/null; then
    break
  fi
  sleep 0.2
done

if ! curl -sf "${API_URL}/health" >/dev/null; then
  echo "error: mock API did not become healthy" >&2
  cat "${MOCK_API_LOG}" >&2 || true
  exit 1
fi

export INTENT_API_URL="${API_URL}"
export INTENT_TOKEN="smoke-local-token"

pushd "${TMP_DIR}" >/dev/null

"${INTENT_BIN}" init smoke-intent --app
cd smoke-intent

"${INTENT_BIN}" package . --unsigned --out dist
PACKAGE_PATH="$(ls dist/*.itpkg | head -n1)"

"${INTENT_BIN}" verify "${PACKAGE_PATH}" --allow-unsigned
"${INTENT_BIN}" publish "${PACKAGE_PATH}" --message "smoke publish"
"${INTENT_BIN}" search smoke-intent
"${INTENT_BIN}" install scope-smoke-intent@0.1.0 --dest installed

test -f "installed/scope-smoke-intent/.installed.json"

popd >/dev/null

echo "smoke test passed: publish/search/install/verify against local mock registry"
