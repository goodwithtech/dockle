package types

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/knqyf263/fanal/types"
)

type DockerConfig struct {
	AuthURL  string        `env:"DOCKER_GUARD_AUTH_URL"`
	UserName string        `env:"DOCKER_GUARD_USERNAME"`
	Password string        `env:"DOCKER_GUARD_PASSWORD"`
	Timeout  time.Duration `env:"DOCKER_GUARD_TIMEOUT_SEC" envDefault:"60s"`
	Insecure bool          `env:"DOCKER_GUARD_INSECURE" envDefault:"true"`
	NonSSL   bool          `env:"DOCKER_GUARD_NON_SSL" envDefault:"false"`
}

func GetDockerOption() (types.DockerOption, error) {
	cfg := DockerConfig{}
	if err := env.Parse(&cfg); err != nil {
		return types.DockerOption{}, err
	}
	return types.DockerOption{
		AuthURL:  cfg.AuthURL,
		UserName: cfg.UserName,
		Password: cfg.Password,
		Timeout:  cfg.Timeout,
		Insecure: cfg.Insecure,
		NonSSL:   cfg.NonSSL,
	}, nil
}
