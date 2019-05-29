module github.com/tomoyamachi/docker-guard

go 1.12

require (
	github.com/docker/docker v0.0.0-20180924202107-a9c061deec0f
	github.com/genuinetools/reg v0.16.0
	github.com/knqyf263/fanal v0.0.0-20190528042547-07e27879b658
	github.com/mattn/go-shellwords v1.0.5 // indirect
	github.com/moby/moby v1.13.1
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/urfave/cli v1.20.0
	github.com/vbatts/tar-split v0.11.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
)

replace github.com/genuinetools/reg => github.com/tomoyamachi/reg v0.16.2-0.20190418055600-c6010b917a55
