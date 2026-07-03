# Note: how tags were created and validated for this PoC

This repository's source, `go.mod`/`go.sum` files, and full commit history
were pushed to GitHub on branch `claude/go-module-resolution-poc-1780jc`.
**Git tags were not pushed to GitHub** — the automation environment that
built this PoC can push branches but its outbound Git proxy rejects
`refs/tags/*` pushes (`HTTP 403`) as a hard, backend-level restriction, not
a client-side approval setting.

To keep every experiment in [`experiments.md`](./experiments.md) a **real**
`go get`/`go mod tidy`/`go list`/`go mod graph`/`go mod why` transcript
(not a hand-typed prediction), tags were created locally with plain `git
tag` (see [`scripts/tag-releases.sh`](../scripts/tag-releases.sh) for the
exact commands) and pushed to a byte-identical local bare Git mirror of
this same repository. `go`'s module-resolution code has no special-case
for "local vs. remote" transport — it drives `git` the same way regardless
of whether the remote is `https://github.com/...` or a local path — so
every resolution rule, error message, and pseudo-version in this PoC's docs
is exactly what you'll see once the tags below exist on GitHub.

## To finish publishing the real tags

Run this once, from a clone with permission to push tags to
`github.com/abhi4u1947/go-multimodule-poc` (your own machine, or any
environment/session without this sandbox's tag-push restriction):

```bash
git fetch origin claude/go-module-resolution-poc-1780jc
git checkout claude/go-module-resolution-poc-1780jc
bash scripts/tag-releases.sh   # creates all 8 tags locally
git push origin --tags
```

After that, every command in `docs/experiments.md` reproduces verbatim
against `https://github.com/abhi4u1947/go-multimodule-poc` directly — no
mirror, no special `GOPRIVATE`/git config needed beyond what any private
repo normally requires (see [`README.md`](../README.md#requirements) for
that one-time setup).

## Update: tags are now published, and releases are automated

The 8 tags above were pushed for real from a machine without this sandbox's
restriction, and every `docs/experiments.md` transcript was re-verified
against the genuine GitHub remote (see that file's header for the
verification note). `scripts/tag-releases.sh` remains for reference/bootstrap
only — going forward, new releases are cut automatically by
[`.github/workflows/release.yml`](../.github/workflows/release.yml): pushing
a change under `entities/<module>/` on `main` bumps that module's patch
version and tags+releases it, and `workflow_dispatch` allows a manual
module/bump-level choice. See that workflow file for the exact tag-naming
and version-bump logic, and
[`go-multimodule-poc-consumer`](https://github.com/abhi4u1947/go-multimodule-poc-consumer)'s
`.github/workflows/auto-update.yml` + `.github/dependabot.yml` for how the
consumer picks up new tags automatically.
