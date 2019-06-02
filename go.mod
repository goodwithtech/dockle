module github.com/goodwithtech/docker-guard

go 1.12

require (
	github.com/Sirupsen/logrus v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/docker v0.7.3-0.20190602164837-acdbaaa3ed04
	github.com/fatih/color v1.7.0
	github.com/genuinetools/reg v0.16.0
	github.com/knqyf263/fanal v0.0.0-20190528042547-07e27879b658
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-shellwords v1.0.5 // indirect
	github.com/moby/moby v0.7.3-0.20190602164837-acdbaaa3ed04
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
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

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
