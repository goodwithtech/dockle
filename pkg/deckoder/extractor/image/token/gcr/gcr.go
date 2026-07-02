package gcr

import (
	"context"
	"fmt"
	"strings"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/GoogleCloudPlatform/docker-credential-gcr/v2/config"
	"github.com/GoogleCloudPlatform/docker-credential-gcr/v2/credhelper"
	"github.com/GoogleCloudPlatform/docker-credential-gcr/v2/store"
)

type GCR struct {
	Store  store.GCRCredStore
	domain string
}

const gcrURL = "gcr.io"

func (g *GCR) CheckOptions(domain string, d types.DockerOption) error {
	if !strings.HasSuffix(domain, gcrURL) {
		return fmt.Errorf("GCR : %w", types.InvalidURLPattern)
	}
	g.domain = domain
	if d.GcpCredPath != "" {
		g.Store = store.NewGCRCredStore(d.GcpCredPath)
	}
	return nil
}

func (g *GCR) GetCredential(ctx context.Context) (username, password string, err error) {
	var credStore store.GCRCredStore
	if g.Store == nil {
		credStore, err = store.DefaultGCRCredStore()
		if err != nil {
			return "", "", fmt.Errorf("failed to get GCRCredStore: %w", err)
		}
	} else {
		credStore = g.Store
	}
	userCfg, err := config.LoadUserConfig()
	if err != nil {
		return "", "", fmt.Errorf("failed to load user config: %w", err)
	}
	helper := credhelper.NewGCRCredentialHelper(credStore, userCfg)
	return helper.Get(g.domain)
}
