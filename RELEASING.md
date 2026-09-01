# Releasing CharVomit

Maintainer-only notes for cutting a release. There's no label-based release
automation in this repo — [release.yml](.github/workflows/release.yml)
triggers solely on pushing a `v*` tag.

## Steps

1. Merge the PR to `main` and confirm [ci.yml](.github/workflows/ci.yml) is
   green on `main` (lint + build on all platforms, tests on Linux/Windows;
   see README "Known CI limitations" for why macOS skips test execution).
2. Decide the version bump using [semver](https://semver.org/):
   - **patch** (`vX.Y.Z+1`): bug fixes, no API changes.
   - **minor** (`vX.Y+1.0`): backwards-compatible features/additions.
   - **major** (`vX+1.0.0`): breaking changes to exported APIs in `pkg/`, or
     breaking CLI behavior changes.
3. **If bumping the major version**, first update the Go module path to
   match, since this repo doesn't use a `/vN` subdirectory layout:
   - `go.mod`: `module github.com/mojotx/CharVomit/vN`
   - All internal import paths (`cmd/CharVomit/*.go`, `pkg/CharVomit/*.go`)
   - `.goreleaser.yaml`'s `ldflags` target path
   - README's Go Reference badge, `go install` command, and any other
     `github.com/mojotx/CharVomit/...` import references
   - Add a "Breaking changes in vN.0.0" note to the README for library
     consumers (see the "Breaking changes in v2.0.0" section as an example)
4. Tag `main` locally and push the tag:
   ```shell
   git checkout main && git pull
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```
5. This triggers `release.yml`, which re-verifies the build/tests/lint, then
   runs GoReleaser (`.goreleaser.yaml`) to cross-compile
   linux/windows/darwin binaries and publish them as a GitHub release
   attached to the tag.
6. Confirm the [release](https://github.com/mojotx/CharVomit/releases)
   published with the expected artifacts, and that `CharVomit --version`
   from a downloaded binary reports the new tag.
