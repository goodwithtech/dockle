# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Dockle is a container image linter for security, written in Go. It scans built Docker images (not Dockerfiles) against CIS Benchmark checkpoints (`CIS-DI-*`) and Dockle original checkpoints (`DKL-DI-*` for Docker best practices, `DKL-LI-*` for Linux best practices). All checkpoints are documented in `CHECKPOINT.md`.

## Commands

```bash
# Build
go build -o dockle cmd/dockle/main.go

# Run all tests (CI runs with CGO_ENABLED=0)
go test ./...

# Run tests for a single package
go test ./pkg/assessor/manifest/

# Run a single test
go test ./pkg/assessor/manifest/ -run TestAssess

# Run locally against an image
./dockle [IMAGE_NAME]
./dockle --input image.tar   # scan a saved image file
```

Releases are built by GoReleaser (`goreleaser.yaml`) via GitHub Actions when a `v*` tag is pushed (`.github/workflows/releasebuild.yml`).

## Architecture

The scan pipeline flows: **CLI → scanner → extractor (deckoder) → assessors → assessment map → report writer**.

1. **Entry point**: `cmd/dockle/main.go` calls `pkg.NewApp()` (`pkg/app.go`), which defines all CLI flags using `urfave/cli` v1. The action is `pkg.Run` (`pkg/run.go`), which wires everything together. `config.CreateFromCli` (`config/config.go`) populates the global `config.Conf` (ignore rules from flags/`DOCKLE_IGNORES`/`.dockleignore`, exit code, etc.).

2. **Image extraction**: `pkg/scanner/scan.go` uses `github.com/goodwithtech/deckoder` (a companion library maintained in a separate repo) to fetch the image from a Docker daemon, remote registry, or tar archive. Only files that assessors declare they need — via `RequiredFiles()` / `RequiredExtensions()` / `RequiredPermissions()` — are extracted, using a tar filter function. Acceptance flags (`--accept-file`, `--accept-file-extension`) remove files from that filter.

3. **Assessors** (`pkg/assessor/`): each subpackage implements the `Assessor` interface (`assessor.go`) and is registered in its `init()`-style list in `pkg/assessor/assessor.go`. Each assessor inspects the extracted `FileMap` and returns `[]*types.Assessment` tagged with a checkpoint code. The `manifest` assessor is the largest — it parses image config/history to lint Dockerfile-derived instructions. Assessor-level allow/deny lists (sensitive words, credential file names) are injected from CLI flags in `pkg/run.go`.

4. **Checkpoint definitions**: `pkg/types/checkpoint.go` holds the code constants, `DefaultLevelMap` (FATAL/WARN/INFO/SKIP/PASS levels), and `TitleMap`. `types.CreateAssessmentMap` (`pkg/types/assessment.go`) groups assessments and applies ignores/levels.

5. **Output**: `pkg/report/` has three `Writer` implementations — list (default, colored), JSON, SARIF — selected by `--format`. The writer returns `abend`, which combined with `--exit-code`/`--exit-level` determines the process exit status.

### Adding a new checkpoint

Touch all of these: code constant + level + title in `pkg/types/checkpoint.go`, detection logic in the relevant assessor under `pkg/assessor/` (or a new assessor registered in `pkg/assessor/assessor.go`), documentation in `CHECKPOINT.md`, and the summary table in `README.md`.

### Test conventions

Tests are table-driven and live next to the code. The `manifest` assessor tests use JSON image-config fixtures in `pkg/assessor/manifest/testdata/`; SARIF writer tests compare against golden files in `pkg/report/testdata/`.
