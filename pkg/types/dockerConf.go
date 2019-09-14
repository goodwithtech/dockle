package types

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/goodwithtech/deckoder/types"
)

type DockerConfig struct {
	AuthURL  string        `env:"DOCKLE_AUTH_URL"`
	UserName string        `env:"DOCKLE_USERNAME"`
	Password string        `env:"DOCKLE_PASSWORD"`
	Timeout  time.Duration `env:"DOCKLE_TIMEOUT_SEC" envDefault:"60s"`
	Insecure bool          `env:"DOCKLE_INSECURE" envDefault:"true"`
	NonSSL   bool          `env:"DOCKLE_NON_SSL" envDefault:"false"`
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
