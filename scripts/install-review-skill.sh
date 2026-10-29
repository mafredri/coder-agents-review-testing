#!/bin/sh
# Installs deep-review from a coder-agents-review checkout into the test bot's
# review workspace, as the only skill its chats see, and removes the
# coder/coder context sources the workspace template registers.
#
# Usage: CODER_SESSION_TOKEN=<bot user's token> \
#          scripts/install-review-skill.sh <coder-agents-review checkout> [workspace]
#
# Run it again after the workspace restarts: the PATH link lives in
# /tmp/coder-script-data/bin, the only PATH directory the chat's shell has
# that the user can write, and the workspace clears /tmp on start.
set -eu
src=$1
ws=${2:-review}
make -C "$src" skill-build
coder ssh "$ws" -- rm -rf .coder/skills/deep-review
coder ssh "$ws" -- mkdir -p .coder/skills
tar -C "$src/skills" -czf - deep-review | coder ssh "$ws" -- tar -C .coder/skills -xzf -
# coder ssh joins its arguments into one command line, so the remaining
# steps go to the remote shell on stdin.
coder ssh "$ws" -- sh -eu <<'REMOTE'
ln -sfn ~/.coder/skills/deep-review/deep-review /tmp/coder-script-data/bin/deep-review
for p in ~/coder/.agents/skills ~/coder/.mcp.json; do
	coder exp chat context remove "$p" >/dev/null 2>&1 || true
done
coder exp chat context list
coder exp chat context show ~/.coder/skills
grep -m1 "^version:" ~/.coder/skills/deep-review/SKILL.md
command -v deep-review
test -n "${REVIEW_BOT_URL:-}" && test -n "${REVIEW_BOT_TOKEN:-}" && echo "REVIEW_BOT_URL and REVIEW_BOT_TOKEN are set"
REMOTE
