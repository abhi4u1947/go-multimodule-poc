#!/usr/bin/env bash
# Re-runs the go-tool commands behind docs/experiments.md against a scratch
# module, so you can reproduce (and diff against) the captured transcripts
# for yourself once this repo's tags exist on GitHub.
set -euo pipefail
MODPATH="github.com/abhi4u1947/go-multimodule-poc"
WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

cd "$WORKDIR"
go mod init example.com/experiments-scratch

echo "== 1. correct dependency resolution =="
go get "$MODPATH/entities/shared-lib@v1.0.0"
go get "$MODPATH/entities/shopping-svc@v1.1.0"
go get "$MODPATH/entities/shopping-svc/estargz@v0.18.2"
go get "$MODPATH/entities/shopping-svc/ipfs@v0.18.2"
go mod tidy -v

echo
echo "== 2. minimal version selection =="
go mod graph
go list -m all

echo
echo "== 5. pseudo-version (replace <branch> with a real untagged commit-ish) =="
echo "  go get $MODPATH/entities/shared-lib@claude/go-module-resolution-poc-1780jc"

echo
echo "Done. See docs/experiments.md for the full annotated walkthrough,"
echo "including the module-path-mismatch, incorrect-tag, and replace demos"
echo "that require extra setup beyond a plain 'go get'."
