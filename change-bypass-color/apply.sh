#!/bin/sh
# Patch the active Claude Code binary so "bypass permissions on" draws in a
# neutral color. Works on a copy, verifies it in a pseudo-terminal, and only
# then swaps it in, keeping the unpatched binary as <version>.orig.
#
# usage: apply.sh [color]   (default #888; a theme key or #rgb that fits)
set -eu

COLOR="${1:-#888}"
SKILL_DIR="$(cd "$(dirname "$0")" && pwd)"
LINK="$HOME/.local/bin/claude"
TARGET="$(readlink -f "${CLAUDE_BIN:-$LINK}")"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/bypass-color.XXXXXX")"
CANDIDATE="$TARGET.patching"
trap 'rm -rf "$WORK" "$CANDIDATE"' EXIT

echo "claude binary: $TARGET ($("$TARGET" --version))"

(cd "$SKILL_DIR/patcher" && go build -o "$WORK/patcher" .)

if "$WORK/patcher" -dry-run -color "$COLOR" "$TARGET" | grep -q "already patched"; then
  echo "already patched with $COLOR; nothing to do"
  exit 0
fi

cp -p "$TARGET" "$CANDIDATE"
"$WORK/patcher" -color "$COLOR" "$CANDIDATE"
codesign --force -s - "$CANDIDATE" 2>/dev/null
"$CANDIDATE" --version >/dev/null

# Start the patched copy in bypass mode in a pseudo-terminal, from the current
# directory (the invoking Claude session's, so already trusted), and read the
# footer's color and any startup error.
OUT="$WORK/screen.txt"
((sleep 14; printf '\003') |
  TERM=xterm-256color script -q -F "$OUT" "$CANDIDATE" --dangerously-skip-permissions >/dev/null 2>&1) &
RUNNER=$!
sleep 16
pkill -f "$CANDIDATE --dangerously-skip-permissions" 2>/dev/null || true
kill "$RUNNER" 2>/dev/null || true
wait "$RUNNER" 2>/dev/null || true

if LC_ALL=C grep -a -q -E 'SyntaxError|ReferenceError|TypeError' "$OUT"; then
  echo "patched copy fails to start; installed binary left unchanged:" >&2
  LC_ALL=C tr -c '[:print:]\n' '.' <"$OUT" | head -20 >&2
  exit 1
fi
FOOTER="$(LC_ALL=C grep -a -o -E '\[38;2;[0-9;]+m.{0,8}bypass permissio' "$OUT" | tail -1 | LC_ALL=C tr -c '[:print:]' '.')"
if [ -z "$FOOTER" ]; then
  echo "could not see the bypass footer in the patched copy; installed binary left unchanged" >&2
  KEPT="${TMPDIR:-/tmp}/bypass-color-screen.txt"
  cp "$OUT" "$KEPT"
  echo "(a trust prompt or onboarding screen may have shown instead; screen capture: $KEPT)" >&2
  exit 1
fi
case "$FOOTER" in
  *"38;2;255;107;128m"*) echo "footer is still red: patch had no effect; installed binary left unchanged" >&2; exit 1 ;;
esac
echo "verified footer: $FOOTER"

[ -e "$TARGET.orig" ] || cp -p "$TARGET" "$TARGET.orig"
mv "$CANDIDATE" "$TARGET"
echo "installed. backup: $TARGET.orig"
echo "run /restart in each Claude session to pick it up"
