#!/usr/bin/env bash
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "error: gh CLI is required" >&2
  exit 1
fi

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <tag>" >&2
  echo "example: $0 v0.3.14" >&2
  exit 1
fi

TAG="$1"
VERSION="${TAG#v}"
REPO="intentregistry/intent-cli"
TAP_REPO="intentregistry/homebrew-tap"

echo "==> Verifying release ${TAG}"

IS_DRAFT="$(gh release view "${TAG}" --repo "${REPO}" --json isDraft --jq '.isDraft')"
IS_PRERELEASE="$(gh release view "${TAG}" --repo "${REPO}" --json isPrerelease --jq '.isPrerelease')"
if [[ "${IS_DRAFT}" != "false" || "${IS_PRERELEASE}" != "false" ]]; then
  echo "error: release ${TAG} is draft/prerelease (draft=${IS_DRAFT}, prerelease=${IS_PRERELEASE})" >&2
  exit 1
fi

EXPECTED_ASSETS=(
  "checksums.txt"
  "intent-darwin-amd64.tar.gz"
  "intent-darwin-arm64.tar.gz"
  "intent-linux-amd64.tar.gz"
  "intent-linux-arm64.tar.gz"
)

for asset in "${EXPECTED_ASSETS[@]}"; do
  if ! gh release view "${TAG}" --repo "${REPO}" --json assets --jq ".assets[].name | select(. == \"${asset}\")" >/dev/null; then
    echo "error: missing release asset ${asset}" >&2
    exit 1
  fi
done

echo "ok: expected release assets are present"

TMP_CHECKSUMS="$(mktemp)"
trap 'rm -f "${TMP_CHECKSUMS}"' EXIT

gh release download "${TAG}" \
  --repo "${REPO}" \
  --pattern checksums.txt \
  --clobber \
  --output "${TMP_CHECKSUMS}"

for asset in "${EXPECTED_ASSETS[@]:1}"; do
  if ! grep -q " ${asset}\$" "${TMP_CHECKSUMS}"; then
    echo "error: checksums.txt missing entry for ${asset}" >&2
    exit 1
  fi
done

echo "ok: checksums.txt includes expected artifacts"

FORMULA_CONTENT="$(gh api "repos/${TAP_REPO}/contents/Formula/intent-cli.rb" --jq .content | tr -d '\n' | base64 --decode)"
if ! grep -q "version \"${VERSION}\"" <<<"${FORMULA_CONTENT}"; then
  echo "error: homebrew formula version does not match ${VERSION}" >&2
  exit 1
fi
if ! grep -q "releases/download/${TAG}/intent-darwin-arm64.tar.gz" <<<"${FORMULA_CONTENT}"; then
  echo "error: homebrew formula URLs do not reference ${TAG}" >&2
  exit 1
fi

echo "ok: homebrew formula updated for ${TAG}"
echo "==> Release verification passed"
