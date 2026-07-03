# go-multimodule-poc

A runnable proof of concept demonstrating Go module resolution in a
realistic multi-module monorepo: nested modules, module identity, Git tag
resolution, Minimal Version Selection, pseudo-versions, `replace`
directives, and module path mismatches.

Companion consumer repo:
[`go-multimodule-poc-consumer`](https://github.com/abhi4u1947/go-multimodule-poc-consumer).

## CI and releases

- [`.github/workflows/ci.yml`](.github/workflows/ci.yml) builds, vets, and
  tests every module standalone (`GOWORK=off`, mirroring how an external
  consumer resolves them) and via the local `go.work`, on every push/PR to
  `main`.
- [`.github/workflows/release.yml`](.github/workflows/release.yml) tags new
  releases automatically: pushing a change under `entities/<module>/` on
  `main` bumps that module's patch version and creates the correctly-prefixed
  tag (`docs/tag-naming.md`) plus a GitHub Release. `workflow_dispatch` lets
  you pick the module and bump level (patch/minor/major) manually. It also
  pings the consumer repo's `auto-update` workflow immediately if a
  `CONSUMER_DISPATCH_TOKEN` repo secret is configured (a PAT with `repo`
  scope on `go-multimodule-poc-consumer`) — without it, the consumer still
  picks up new tags via its own Dependabot config and scheduled check.

> **One deliberate deviation from the original brief:** module paths use
> this repository's real import path
> (`github.com/abhi4u1947/go-multimodule-poc`) instead of the placeholder
> `github.com/example/monorepo`, so every `go get` in this README and in
> [`docs/experiments.md`](docs/experiments.md) actually resolves against
> real Git history rather than a path nobody can fetch.
>
> **A second, environment-driven deviation:** release tags exist in this
> repo's local history but are not yet pushed to GitHub — the sandbox that
> built this PoC can push branches but not tags. See
> [`docs/NOTE-on-tags.md`](docs/NOTE-on-tags.md) for exactly why, and the
> one command that finishes publishing them.

## Layout

```
go-multimodule-poc/
├── README.md
├── go.work                          # local multi-module dev only (see docs/module-discovery.md)
├── docs/
│   ├── module-discovery.md          # how `go` finds module boundaries; nested-module traversal
│   ├── tag-naming.md                # how tag names map to module paths
│   ├── version-resolution.md        # MVS, pseudo-versions, replace
│   ├── experiments.md               # all 9 experiments, real captured `go` output
│   └── NOTE-on-tags.md              # why tags aren't pushed to GitHub yet, and how to finish
├── scripts/
│   ├── tag-releases.sh              # (re)creates all 8 release tags locally
│   └── run-experiments.sh           # reproduces experiment 1-2 output for yourself
└── entities/
    ├── shared-lib/                  # module: .../entities/shared-lib
    │   ├── go.mod
    │   ├── logger/  config/  utils/
    └── shopping-svc/                # module: .../entities/shopping-svc
        ├── go.mod
        ├── cmd/shopping/  api/  internal/  pkg/  version/
        ├── estargz/                 # independent module: .../entities/shopping-svc/estargz
        │   ├── go.mod
        │   ├── pkg/  version/
        └── ipfs/                    # independent module: .../entities/shopping-svc/ipfs
            ├── go.mod
            ├── pkg/  version/
```

Four independent Go modules, one Git repository. `estargz` and `ipfs` are
nested *directories* under `shopping-svc/` but are **not** part of the
`shopping-svc` module — each has its own `go.mod`, its own version history,
and its own dependency on `shared-lib`. See
[`docs/module-discovery.md`](docs/module-discovery.md) for why that nesting
doesn't collapse them into one module.

## Modules and versions

| Module path | Directory | Tags |
|---|---|---|
| `.../entities/shared-lib` | `entities/shared-lib` | `v1.0.0`, `v1.1.0` |
| `.../entities/shopping-svc` | `entities/shopping-svc` | `v1.0.0`, `v1.1.0` |
| `.../entities/shopping-svc/estargz` | `entities/shopping-svc/estargz` | `v0.18.1`, `v0.18.2` |
| `.../entities/shopping-svc/ipfs` | `entities/shopping-svc/ipfs` | `v0.18.1`, `v0.18.2` |

Full paths and the exact Git tag naming rule that derives from them:
[`docs/tag-naming.md`](docs/tag-naming.md).

## Requirements

- Go 1.24+
- Git

## Running it from scratch

```bash
git clone https://github.com/abhi4u1947/go-multimodule-poc.git
cd go-multimodule-poc

# Local multi-module development, via the checked-in go.work:
go run ./entities/shopping-svc/cmd/shopping
# 2026-07-03T10:13:27Z [INFO] (shopping-svc) starting shopping-svc v1.1.0
# 2026-07-03T10:13:27Z [INFO] (shopping-svc) api self-check: pong
# shopping-svc ready

# Each module also builds standalone (no go.work), resolving shared-lib
# from Git the same way an external consumer would, once tags are
# published (see docs/NOTE-on-tags.md):
cd entities/shopping-svc/estargz && GOWORK=off go build ./...
```

## Experiments

All 9 experiments from the brief, each with real, unedited `go get` /
`go mod tidy` / `go list` / `go mod graph` / `go mod why` output:

1. [Correct dependency resolution](docs/experiments.md#1-correct-dependency-resolution)
2. [Requiring different versions (MVS)](docs/experiments.md#2-requiring-different-versions-minimal-version-selection)
3. [Module path mismatch](docs/experiments.md#3-module-path-mismatch)
4. [Incorrect Git tags](docs/experiments.md#4-incorrect-git-tags)
5. [Missing tags → pseudo-versions](docs/experiments.md#5-missing-git-tags--pseudo-versions)
6. [Local development with `replace`, then removing it](docs/experiments.md#6-local-development-with-replace-then-removing-it)
7. [Nested modules stop parent traversal](docs/experiments.md#7-nested-modules-stop-parent-traversal)
8. [Module identity vs. Git tags](docs/experiments.md#8-module-identity-vs-git-tags)
9. [Mapping a requested module to its Git tag](docs/experiments.md#9-mapping-a-requested-module-to-its-git-tag)

## Further reading

- [`docs/module-discovery.md`](docs/module-discovery.md) — how `go` walks
  directories to find module boundaries, and why nested `go.mod` files stop
  that walk.
- [`docs/tag-naming.md`](docs/tag-naming.md) — the exact algorithm mapping
  a module path + version to a Git tag name, and back.
- [`docs/version-resolution.md`](docs/version-resolution.md) — Minimal
  Version Selection, pseudo-versions, and `replace` semantics.
