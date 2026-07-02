package image

import (
	"context"
	"github.com/containers/image/v5/pkg/docker/config"
	imageTypes "github.com/containers/image/v5/types"
	"github.com/goodwithtech/dockle/pkg/deckoder/types"
)

var (
	registries []Registry
)

type Registry interface {
	CheckOptions(domain string, option types.DockerOption) error
	GetCredential(ctx context.Context) (string, string, error)
}

func RegisterRegistry(registry Registry) {
	registries = append(registries, registry)
}

func GetToken(ctx context.Context, domain string, opt types.DockerOption) (auth *imageTypes.DockerAuthConfig) {
	if opt.UserName != "" || opt.Password != "" {
		return &imageTypes.DockerAuthConfig{Username: opt.UserName, Password: opt.Password}
	}

	if config, _ := config.GetCredentials(nil, domain); config.Username != "" {
		return &config
	}

	// check registry which particular to get credential
	for _, registry := range registries {
		err := registry.CheckOptions(domain, opt)
		if err != nil {
			continue
		}
		username, password, err := registry.GetCredential(ctx)
		if err != nil {
			// only skip check registry if error occurred
			break
		}
		return &imageTypes.DockerAuthConfig{Username: username, Password: password}
	}
	return nil
}
