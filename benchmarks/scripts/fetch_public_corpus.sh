#!/usr/bin/env bash
set -euo pipefail

export LC_ALL=C

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT_DIR="$ROOT_DIR/benchmarks/corpus/public"
CACHE_DIR="$OUT_DIR/_cache"

mkdir -p "$OUT_DIR" "$CACHE_DIR"

clone_or_update() {
  local repo_url="$1"
  local repo_dir="$2"

  if [[ -d "$repo_dir/.git" ]]; then
    git -C "$repo_dir" pull --ff-only >/dev/null
  else
    git clone --depth 1 "$repo_url" "$repo_dir" >/dev/null
  fi
}

normalize_file() {
  local in="$1"
  local out="$2"

  tr -d '\000' < "$in" \
    | sed 's/\r$//' \
    | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' \
    | awk 'length($0) > 0' \
    | awk '$0 !~ /^#/ && $0 !~ /^\/\//' \
    > "$out"
}

merge_file() {
  local src="$1"
  local out="$2"
  if [[ -f "$src" ]]; then
    cat "$src" >> "$out"
  fi
}

RENWA_DIR="$CACHE_DIR/RenwaX23-XSS-Payloads"
PGAIJIN_DIR="$CACHE_DIR/pgaijin66-XSS-Payloads"
PATTT_DIR="$CACHE_DIR/PayloadsAllTheThings"

clone_or_update "https://github.com/RenwaX23/XSS-Payloads.git" "$RENWA_DIR"
clone_or_update "https://github.com/pgaijin66/XSS-Payloads.git" "$PGAIJIN_DIR"
clone_or_update "https://github.com/swisskyrepo/PayloadsAllTheThings.git" "$PATTT_DIR"

normalize_file "$RENWA_DIR/Payloads.txt" "$OUT_DIR/renwax23.txt"
normalize_file "$PGAIJIN_DIR/payload/payload.txt" "$OUT_DIR/pgaijin66.txt"

PATTT_TMP="$(mktemp)"
for f in \
  "$PATTT_DIR/XSS Injection/Intruders/IntrudersXSS.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/JHADDIX_XSS.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/MarioXSSVectors.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/RSNAKE_XSS.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/XSSDetection.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/XSS_Polyglots.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/xss_alert.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/xss_alert_identifiable.txt" \
  "$PATTT_DIR/XSS Injection/Intruders/xss_payloads_quick.txt"
do
  if [[ -f "$f" ]]; then
    cat "$f" >> "$PATTT_TMP"
    printf '\n' >> "$PATTT_TMP"
  fi
done
normalize_file "$PATTT_TMP" "$OUT_DIR/payloadsallthethings.txt"
rm -f "$PATTT_TMP"

MERGED_TMP="$(mktemp)"
: > "$MERGED_TMP"
merge_file "$OUT_DIR/renwax23.txt" "$MERGED_TMP"
merge_file "$OUT_DIR/pgaijin66.txt" "$MERGED_TMP"
merge_file "$OUT_DIR/payloadsallthethings.txt" "$MERGED_TMP"

awk 'length($0) > 0' "$MERGED_TMP" | sort -u > "$OUT_DIR/merged.txt"
rm -f "$MERGED_TMP"

{
  echo "renwax23 $(wc -l < "$OUT_DIR/renwax23.txt" | tr -d ' ')"
  echo "pgaijin66 $(wc -l < "$OUT_DIR/pgaijin66.txt" | tr -d ' ')"
  echo "payloadsallthethings $(wc -l < "$OUT_DIR/payloadsallthethings.txt" | tr -d ' ')"
  echo "merged_unique $(wc -l < "$OUT_DIR/merged.txt" | tr -d ' ')"
} > "$OUT_DIR/stats.txt"

echo "Public corpus updated in: $OUT_DIR"
cat "$OUT_DIR/stats.txt"
