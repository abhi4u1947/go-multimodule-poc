#!/usr/bin/env bash
# Creates every release tag this PoC references, at the commit that
# introduced each version. Safe to re-run (skips tags that already exist).
#
# See docs/NOTE-on-tags.md for why this script exists separately from
# `git push` in this repo's own history: this sandbox can push branches but
# not tags, so tags are created here for you (or your own CI/local clone)
# to push with `git push origin --tags`.
set -euo pipefail
cd "$(dirname "$0")/.."

tag() {
	local name="$1" rev="$2"
	if git rev-parse -q --verify "refs/tags/$name" >/dev/null; then
		echo "skip (exists): $name"
	else
		git tag "$name" "$rev"
		echo "created: $name -> $(git rev-parse --short "$rev")"
	fi
}

tag entities/shared-lib/v1.0.0                   "$(git log --format=%H --grep='^shared-lib: initial')"
tag entities/shared-lib/v1.1.0                    "$(git log --format=%H --grep='^shared-lib: add utils.Contains')"
tag entities/shopping-svc/v1.0.0                  "$(git log --format=%H --grep='^shopping-svc: initial')"
tag entities/shopping-svc/v1.1.0                  "$(git log --format=%H --grep='^shopping-svc: bump shared-lib to v1.1.0')"
tag entities/shopping-svc/estargz/v0.18.1         "$(git log --format=%H --grep='^shopping-svc: bump shared-lib to v1.1.0')"
tag entities/shopping-svc/estargz/v0.18.2         "$(git log --format=%H --grep='^estargz: add TOCDigest')"
tag entities/shopping-svc/ipfs/v0.18.1            "$(git log --format=%H --grep='^estargz: add TOCDigest')"
tag entities/shopping-svc/ipfs/v0.18.2            "$(git log --format=%H --grep='^ipfs: add Pin')"

echo
echo "All tags:"
git tag | sort
echo
echo "Next: git push origin --tags"
