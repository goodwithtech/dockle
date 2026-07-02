# deckoder

Fetches a container image from a Docker daemon, a remote registry, or a tar
archive, and extracts the files that the caller asks for via a tar filter
function.

This package was originally maintained as a separate library,
[github.com/goodwithtech/deckoder](https://github.com/goodwithtech/deckoder)
(itself a fork of [aquasecurity/fanal](https://github.com/aquasecurity/fanal)),
and was integrated into this repository since Dockle is its only consumer.

Package layout:

- `analyzer/` — walks all image layers and builds the merged `FileMap`
- `extractor/` — the `Extractor` interface; `extractor/docker` implements it
  on top of `containers/image` transports (daemon, registry, archive)
- `extractor/image` — image source resolution and registry authentication
  (`token/ecr`, `token/gcr` provide ECR/GCR credentials)
- `types/` — `FileMap`, `FilterFunc`, `DockerOption`, layer metadata
- `utils/` — small shared helpers
