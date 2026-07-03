# GOPROXY-ready build artifacts

Every push to `main` runs [`tools/goproxy-gen`](../tools/goproxy-gen) and
uploads its output as a GitHub Actions build artifact named
`goproxy-<commit-sha>` (see
[`.github/workflows/goproxy-artifacts.yml`](../.github/workflows/goproxy-artifacts.yml)).
That output is a **static file tree implementing the real Go module proxy
protocol** — the same wire format `proxy.golang.org` and any GOPROXY server
speak — built directly from this repo's own Git tags for all four modules.

## Why

This repo is a teaching example about Go module resolution, and "here's a
directory you can point `GOPROXY` at and watch `go get` resolve our modules
with zero network access, zero GitHub, zero tag-push permissions" is the
most direct way to *see* the tag → module → version mapping described in
[`tag-naming.md`](./tag-naming.md) and [`version-resolution.md`](./version-resolution.md),
rather than take it on faith. It's also just useful: air-gapped CI, offline
development, or vendoring this repo's modules into an environment that
can't reach GitHub all work by serving this artifact instead.

## What `tools/goproxy-gen` does

It's a small, generic Go program (not tied to any of this repo's four
modules — it's its own module under `tools/`, kept out of their dependency
graphs) built entirely on the official `golang.org/x/mod` packages:

1. **Discovers modules** by walking the repo for `go.mod` files
   (`golang.org/x/mod/modfile`), skipping `tools/` and `.git/`.
2. **Finds each module's real versions** by listing Git tags matching that
   module's directory-prefixed tag convention (`entities/<dir>/vX.Y.Z` —
   see [`tag-naming.md`](./tag-naming.md)) — the exact same tag ↔ module
   mapping `go get` itself uses, re-derived independently here rather than
   assumed.
3. **Builds each version's module zip directly from Git history** via
   `golang.org/x/mod/zip.CreateFromVCS`, which shells out to `git archive`
   under the hood and — critically — **automatically excludes nested
   modules**: `shopping-svc`'s zip does not contain `estargz/` or `ipfs/`,
   the same module-boundary rule `go build` enforces (see
   [`module-discovery.md`](./module-discovery.md)). No custom zip-filtering
   logic was needed; this is exactly what the real Go toolchain uses
   internally.
4. **Extracts each version's `go.mod`** straight out of the zip it just
   built (so what's served always matches what's in the zip, by
   construction) and gets each tag's commit time via `git log`.
5. **Writes the standard proxy layout** — `@v/list`, `@v/<version>.info`,
   `@v/<version>.mod`, `@v/<version>.zip`, `@latest` — under
   `module.EscapePath`-escaped module path directories, exactly as
   documented at <https://go.dev/ref/mod#goproxy-protocol>.

## Using the artifact

Download the `goproxy-<sha>` artifact from a workflow run (Actions tab →
the run → Artifacts), unzip it, then:

```bash
export GOPROXY="file:///absolute/path/to/goproxy"
export GOSUMDB=off               # this tree isn't registered with sum.golang.org
export GOPRIVATE="github.com/abhi4u1947/*"   # skip the (external) proxy/sumdb for this path

go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc@v1.1.0
go get github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz@latest
```

Both resolve entirely from the local file tree — no network access needed,
matching every command in [`experiments.md`](./experiments.md) but served
from disk instead of Git. The workflow itself does exactly this as a
sanity check before uploading (see the `sanity check` step in
`goproxy-artifacts.yml`), so a broken artifact fails CI rather than
shipping silently.

## Reproducing it yourself

```bash
cd tools/goproxy-gen
go run . --repo ../.. --out ../../dist/goproxy
```

`--repo` must point at a real (non-bare) Git checkout with the tags you
want included already fetched (`git fetch --tags`). `dist/` is
`.gitignore`d — this is a generated build artifact, not checked-in source.
