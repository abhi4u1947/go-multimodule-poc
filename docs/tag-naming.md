# How Go derives tag names from module paths

## The rule

For a module whose `module` directive is a subpath of the repository's
root import path, Go requires the release tag to be prefixed with that
module's directory path **relative to the repository root**, followed by
`/vMAJOR.MINOR.PATCH`:

```
tag = <dir-of-go.mod-relative-to-repo-root>/vX.Y.Z
```

For the *root* module of a repository (the module whose `go.mod` sits at
the repository root), no prefix is used — just `vX.Y.Z`. None of this
repo's four modules is a root module (there is no `go.mod` at
`monorepo/`), so all four use the prefixed form.

## Applying it to this repository

The repository root import path is `github.com/abhi4u1947/go-multimodule-poc`.
Module directories and their correct tags:

| Module (`module` directive)                                                    | Directory (relative to repo root)         | Correct tag                                    |
|----------------------------------------------------------------------------------|--------------------------------------------|-------------------------------------------------|
| `.../go-multimodule-poc/entities/shared-lib`                                    | `entities/shared-lib`                      | `entities/shared-lib/v1.0.0`, `v1.1.0`           |
| `.../go-multimodule-poc/entities/shopping-svc`                                  | `entities/shopping-svc`                    | `entities/shopping-svc/v1.0.0`, `v1.1.0`         |
| `.../go-multimodule-poc/entities/shopping-svc/estargz`                         | `entities/shopping-svc/estargz`            | `entities/shopping-svc/estargz/v0.18.1`, `v0.18.2` |
| `.../go-multimodule-poc/entities/shopping-svc/ipfs`                            | `entities/shopping-svc/ipfs`               | `entities/shopping-svc/ipfs/v0.18.1`, `v0.18.2`  |

This is why this PoC's tags carry the `entities/` prefix, even though the
original task brief sketched shorter tags like `shopping-svc/v1.0.0`. Those
shorter names are exactly what you'd get if the module lived at
`monorepo/shopping-svc` (no `entities/` layer) — or, as demonstrated in
[experiment 4](./experiments.md#4-incorrect-git-tags), what happens when
someone tags at the *wrong* prefix: `go get` reports
`invalid version: unknown revision ...` because no ref named
`entities/shopping-svc/estargz/vX.Y.Z` exists, even though a same-named tag
exists elsewhere in the repo.

## Why the prefix has to be there at all

A single Git repository can contain many modules, so `vX.Y.Z` alone is
ambiguous — `v1.0.0` could mean "shared-lib 1.0.0" or "shopping-svc 1.0.0."
The directory-prefixed tag disambiguates *which* module's release history a
given tag belongs to, and simultaneously tells `go` exactly where inside
the repo to read `go.mod`/packages from for that version — it does not need
a separate "module → subdirectory" lookup table; the tag *is* that lookup.

## Mapping a request to a tag, concretely

```
requested module:  github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz @ v0.18.2
                    └──────────────── repo root ────────────────┘└──────── subdir ────────┘
                    github.com/abhi4u1947/go-multimodule-poc      entities/shopping-svc/estargz

resolution steps:
  1. go tries successively shorter prefixes of the module path as
     candidate repository roots (github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/estargz,
     then .../entities/shopping-svc, then .../entities, then
     github.com/abhi4u1947/go-multimodule-poc) until one resolves to a
     real VCS repository. Here that's the last one.
  2. subdir = module path with the repo root prefix removed
            = "entities/shopping-svc/estargz"
  3. expected tag = subdir + "/v0.18.2"
            = "entities/shopping-svc/estargz/v0.18.2"
  4. go fetches that tag's tree and reads entities/shopping-svc/estargz/go.mod
     from it.
```

See this run for real, unedited `go get` output:
[experiment 9](./experiments.md#9-mapping-a-requested-module-to-its-git-tag).

## Module identity vs. tag: two different jobs

- **Module identity** comes *entirely* from the `module` directive inside
  `go.mod`. It's what other code imports, what appears in `go.mod`
  `require` lines, and what the Go toolchain uses to detect a
  [mismatch](./experiments.md#3-module-path-mismatch) (`module declares its
  path as X but was required as Y`).
- **The Git tag** is only a *label on a commit* used to locate the
  requested *version* of that identity inside the repository's history. It
  plays no role once the module is fetched — nothing in `go.mod`, `go.sum`,
  or the build ever again refers to the tag name itself, only to the
  resolved semantic version and content hash.

Concretely: renaming a tag doesn't change a module's identity (the
`module` line is untouched), and two completely different repositories
could tag their releases identically (`v1.0.0`) without any conflict,
because module identity is namespaced by import path, not by tag string.
See [experiment 8](./experiments.md#8-module-identity-vs-git-tags).
