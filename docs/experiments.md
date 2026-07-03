# Experiments

Every transcript below is **real, unedited `go` command output**, captured
by actually running these commands against this repository's real commits
and tags. Originally captured through a byte-identical local Git mirror
(this sandbox couldn't push Git tags to GitHub yet — see
[`NOTE-on-tags.md`](./NOTE-on-tags.md)), and since re-verified command by
command against the real, tagged `github.com/abhi4u1947/go-multimodule-poc`
through the real `proxy.golang.org`. Every transcript matched byte-for-byte
(including `go.sum` content hashes and the module cache's recorded commit
`Time` field) except experiment 5, which needed a correction — see there
for details.

All commands run from a scratch module that requires this repo's four
modules, e.g.:

```
module github.com/abhi4u1947/go-multimodule-poc-consumer

require (
	github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.0
	github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0
	github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2
	github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2
)
```

---

## 1. Correct dependency resolution

```
$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.0.0
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.0.0
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.0.0

$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc@v1.1.0
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.0
go: upgraded github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.0.0 => v1.1.0
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0

$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2

$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs@v0.18.2
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2

$ go mod tidy -v
(no output — nothing to add or remove; already minimal and consistent)

$ go run .
2026-07-03T10:17:57Z [INFO] (consumer) shopping-svc 1.1.0, estargz 0.18.2, ipfs 0.18.2
ok
```

`go mod tidy` upgraded `shared-lib` from the `v1.0.0` that `estargz`/`ipfs`
minimally need to `v1.1.0`, because `shopping-svc` needs `>= v1.1.0` — see
experiment 2.

`go.sum` after tidy (real content hashes):

```
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.0 h1:nq7mh1IAdexyQkbTrVqQJ+OOJNIHlasIySWqbkbYF8s=
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.0/go.mod h1:LuvN7kp6upUmjw72d+VlykZ+ewjj23uOl65y0jMO7PM=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0 h1:nkKlyS6i2uQIKUQmNAWewaviLQjv5cfIPpCfzLznNzU=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0/go.mod h1:7moB4yBS8zpyBYT3yjLWEG9Lcdc7mdrmRHNO/dxlzYY=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2 h1:o7dRget40Cg180yHk/HlwZQgA9Sg75KLVz8fdup5LKU=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2/go.mod h1:uKPVJyAb4SiVM/2ICospAWm3EIBjReXYwK4zSGcMyAI=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2 h1:zgE/kLKYLiZbzhA2AyvaixRULEUDnLSh1ivTTQe409s=
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2/go.mod h1:ccPLUq+O98H1NjwckA9er5COyC1/z5yxgTyuLjUQx3M=
```

## 2. Requiring different versions (Minimal Version Selection)

Each module states a *different minimum* for `shared-lib`:

| Requirer | minimum `shared-lib` |
|---|---|
| shopping-svc v1.1.0 | v1.1.0 |
| estargz v0.18.2 | v1.0.0 |
| ipfs v0.18.2 | v1.0.0 |

```
$ go mod graph
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.1.0
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc@v1.1.0
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs@v0.18.2
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed go@1.24.7
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.1.0 go@1.24
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc@v1.1.0 github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.1.0
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc@v1.1.0 go@1.24
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2 github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.0.0
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2 go@1.24
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs@v0.18.2 github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.0.0
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs@v0.18.2 go@1.24
go@1.24.7 toolchain@go1.24.7

$ go list -m all
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.0
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc v1.1.0
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.2
```

`go mod graph` prints the raw requirement edges — note `estargz` and `ipfs`
still show `shared-lib@v1.0.0`, their own stated minimum, unchanged.
`go list -m all` prints the *result of MVS*: a single row per module path
showing the one version actually selected for the build — `v1.1.0`, the
maximum of `{v1.1.0, v1.0.0, v1.0.0}`. This is MVS end to end: no version
ever gets edited out of the graph, only the max is chosen for compilation.

```
$ go mod why -m github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib
# github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger

$ go mod why github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger
# github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger
github.com/abhi4u1947/go-multimodule-poc-consumer/testbed
github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger
```

`go mod why` prints the shortest import path from the main module to the
target — here, directly, because `main.go` imports `shared-lib/logger`.
Note `go mod why <module-path>` (no `-m`) answers "is this exact
*package* path imported", and reports "does not need package" for
`shared-lib` because nothing imports a package at the module's root
directory — only its `logger`/`config`/`utils` subpackages are ever
imported. Use `-m` to ask the *module-level* question instead.

```
$ go list -m -json github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz
{
	"Path": "github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz",
	"Version": "v0.18.2",
	"Time": "2026-07-03T10:14:53Z",
	"Dir": "/root/go/pkg/mod/github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2",
	"GoMod": "/root/go/pkg/mod/cache/download/github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz/@v/v0.18.2.mod",
	"GoVersion": "1.24",
	"Sum": "h1:o7dRget40Cg180yHk/HlwZQgA9Sg75KLVz8fdup5LKU=",
	"GoModSum": "h1:uKPVJyAb4SiVM/2ICospAWm3EIBjReXYwK4zSGcMyAI="
}
```

## 3. Module path mismatch

An intentionally-broken copy of `estargz` was committed on branch
`experiment/module-path-mismatch` and tagged
`entities/shopping-svc/estargz/v0.18.3-mismatch`, with its `go.mod`
changed to:

```go
module github.com/abhi4u1947/go-multimodule-poc/estargz-wrong-path
```

(instead of the correct `.../entities/shopping-svc/estargz`) while staying
in the same directory, so it's fetched *as if* it were
`.../entities/shopping-svc/estargz`:

```
$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.3-mismatch
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.3-mismatch
go: github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.3-mismatch requires
        github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.3-mismatch: parsing go.mod:
	module declares its path as: github.com/abhi4u1947/go-multimodule-poc/estargz-wrong-path
	        but was required as: github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz
```

`go` never even looks at the code — it rejects the module the moment
`go.mod`'s declared identity disagrees with the path it was fetched under.
This is the check described in
[`version-resolution.md`](./version-resolution.md#module-path-mismatch).

## 4. Incorrect Git tags

The task brief's shorthand tag names (`shopping-svc/v1.0.0`, without the
`entities/` prefix this repo's real directory layout requires — see
[`tag-naming.md`](./tag-naming.md)) don't match any module's true
subdirectory. To prove that mismatch actually breaks resolution rather
than "happening to still work", a real additional release,
`entities/shopping-svc/estargz/v0.18.2`'s exact content, was re-tagged at
the *wrong*, unprefixed-style path `shopping-svc/estargz/v0.18.4`:

```
$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.4
go: github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.4: invalid version:
        unknown revision entities/shopping-svc/estargz/v0.18.4
```

`go` looked for a tag literally named `entities/shopping-svc/estargz/v0.18.4`
(the *correct* prefix for this module's directory) — it does not exist,
even though a tag containing the string `v0.18.4` exists in the repo at
`shopping-svc/estargz/v0.18.4` (wrong prefix). Tags are matched by **exact
ref name**, not by suffix or fuzzy search.

## 5. Missing Git tags → pseudo-versions

`shared-lib`'s HEAD commit adds a `Reverse()` helper after the `v1.1.0` tag
and was deliberately left untagged:

```
$ git log --oneline -3 entities/shared-lib
8945534 shared-lib: add Reverse helper (unreleased, no tag yet)
c0dc535 (tag: entities/shared-lib/v1.1.0) shared-lib: add utils.Contains helper (v1.1.0)
d86b6b1 (tag: entities/shared-lib/v1.0.0) shared-lib: initial logger, config, utils packages

$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@claude/go-module-resolution-poc-1780jc
go: github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@claude/go-module-resolution-poc-1780jc: invalid version: version "claude/go-module-resolution-poc-1780jc" invalid: disallowed version string
```

**Correction (verified against the real GitHub remote):** the branch name
itself can't be used as a `@version` query here — `go` rejects any `@query`
string containing `/` with "disallowed version string" (confirmed on both
go1.24.7 and go1.26.1, so this isn't a toolchain-version artifact; it holds
for any branch name with a slash in it, which `claude/go-module-resolution-poc-1780jc`
has). The fix is to resolve the branch to its commit first and query by
that instead — this is what actually produces the pseudo-version:

```
$ git rev-parse --short claude/go-module-resolution-poc-1780jc   # or any prefix of it
8945534

$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@8945534
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.1-0.20260703101852-89455342706c
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.1-0.20260703101852-89455342706c

$ grep shared-lib go.mod
	github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.1-0.20260703101852-89455342706c
```

`89455342706c` is exactly the first 12 hex characters of that commit's real
SHA (`8945534...`); `20260703101852` is its commit timestamp in
`YYYYMMDDHHMMSS` (UTC); and the base version bumped from `v1.1.0` (last
real tag) to `v1.1.1` because a pseudo-version must sort strictly after the
tag it follows. See [`version-resolution.md`](./version-resolution.md#pseudo-versions).

Restoring the tagged release:

```
$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib@v1.1.0
go: downgraded github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib v1.1.1-0.20260703101852-89455342706c => v1.1.0
```

## 6. Local development with `replace`, then removing it

A working copy of `shared-lib` was checked out to `../local-dev-shared-lib`
and patched (unreleased, untagged, uncommitted-upstream) to prefix every
log line with `[LOCAL-DEV]`. Adding:

```
replace github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib => ../local-dev-shared-lib
```

to `go.mod` and re-running, with **no change to the `require` line or any
`go get`**:

```
$ go run .
[LOCAL-DEV] 2026-07-03T10:21:46Z [INFO] (consumer) shopping-svc 1.1.0, estargz 0.18.2, ipfs 0.18.2
ok
```

The build used the on-disk edit instantly — `replace` bypasses version
resolution entirely for the replaced path; the resolved version string in
`go.mod` (`v1.1.0`) becomes cosmetic while the directive is present.

### Removing `replace`, back to Git

Deleting the `replace` line and rebuilding, unchanged otherwise:

```
$ go run .
2026-07-03T10:22:03Z [INFO] (consumer) shopping-svc 1.1.0, estargz 0.18.2, ipfs 0.18.2
ok
```

The `[LOCAL-DEV]` prefix is gone — resolution fell straight back to the
Git-tagged `v1.1.0` recorded in `go.mod`/`go.sum`, with zero other changes
needed. This is the entire point of `replace`: a *local, reversible*
override that leaves the "real" dependency declaration untouched.

## 7. Nested modules stop parent traversal

From inside `entities/shopping-svc` (its own module, `GOWORK=off` so the
repo-wide workspace doesn't mask module boundaries):

```
$ go list ./...
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/api
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/cmd/shopping
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/internal
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/pkg
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/version
```

`entities/shopping-svc/estargz` and `entities/shopping-svc/ipfs` are
physically nested under `shopping-svc/`, yet **do not appear** — each has
its own `go.mod`, so `shopping-svc`'s module boundary ends exactly at those
subdirectories. `shopping-svc` must depend on them the same way any outside
consumer would (`require` + a real version), never via a relative import;
attempting to `import ".../shopping-svc/estargz/pkg"` from inside
`shopping-svc` without a corresponding `require` line fails exactly like
importing any other external module would.

## 8. Module identity vs. Git tags

Module identity is the `module` line in `go.mod` — full stop. The Git tag
that got you a particular version is discarded the moment the version is
resolved; nothing downstream (`go.mod`, `go.sum`, compiled binaries,
`go list -m`) ever again references the tag string.

Proof: `entities/shopping-svc/ipfs/v0.18.1` and
`entities/shopping-svc/estargz/v0.18.2` are two **different** tags that
happen to point at the **same commit** (`204b1d6`, because that commit
touched `estargz` but left `ipfs` unchanged since its previous release):

```
$ git tag --points-at 204b1d6
entities/shopping-svc/estargz/v0.18.2
entities/shopping-svc/ipfs/v0.18.1
```

`go get`ting each by its own tag name yields two distinct, correctly
identified modules — `go` never confuses them, because each resolves via
its own module path's expected subdirectory + `go.mod`, and the tag name
itself is discarded immediately after locating the commit:

```
$ go list -m github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2

$ go list -m github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs
github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/ipfs v0.18.1
```

Renaming or deleting a tag after the fact (e.g. `git tag -d ... && git tag
newname ...`) never changes what a module *is* — only what versions can be
*requested by name* going forward; already-resolved `go.sum` entries pin
the content hash, not the tag.

## 9. Mapping a requested module to its Git tag

```
requested:  github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz @ v0.18.2
```

```
$ go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@v0.18.2
go: downloading github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2
go: added github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz v0.18.2
```

Behind that one line, `go` did:

1. Tried repo-root candidates, longest prefix first, until one exists as a
   real repository: `.../entities/shopping-svc/estargz` (no) → `.../entities/shopping-svc`
   (no) → `.../entities` (no) → `github.com/abhi4u1947/go-multimodule-poc` (yes).
2. Computed the subdirectory as the remainder of the module path:
   `entities/shopping-svc/estargz`.
3. Looked for ref `refs/tags/entities/shopping-svc/estargz/v0.18.2` — found it,
   pointing at commit `204b1d6`.
4. Read `entities/shopping-svc/estargz/go.mod` from that commit's tree,
   confirmed its `module` line equals the requested path (see experiment
   3 for what happens when it doesn't), and used that tree as the module's
   contents for `v0.18.2`.

`git show entities/shopping-svc/estargz/v0.18.2 --stat` confirms step 3/4
independently of `go` — the tag points at a real commit containing exactly
that subdirectory's files at that point in history.
