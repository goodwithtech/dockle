package credential

import (
	"fmt"
	"os"
	"path/filepath"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/log"

	"github.com/goodwithtech/dockle/pkg/types"
)

type CredentialAssessor struct{}

func (a CredentialAssessor) Assess(fileMap deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : credential files")
	assesses := []*types.Assessment{}
	fmap := makeMaps(a.RequiredFiles())
	fexts := makeMaps(a.RequiredExtensions())
	for filename := range fileMap {
		basename := filepath.Base(filename)
		// check exist target files
		_, ok1 := fmap[basename]
		_, ok2 := fexts[filepath.Ext(basename)]
		if ok1 || ok2 {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.AvoidCredential,
					Filename: filename,
					Desc:     fmt.Sprintf("Suspicious filename found : %s (You can suppress it with \"-af %s\")", filename, basename),
				})
		}
	}
	return assesses, nil
}

func makeMaps(keys []string) map[string]struct{} {
	maps := make(map[string]struct{})
	for i := 0; i < len(keys); i++ {
		maps[keys[i]] = struct{}{}
	}
	return maps
}

func (a CredentialAssessor) RequiredFiles() []string {
	return []string{
		"credentials.json",
		"credential.json",
		"config.json",
		"credentials",
		"credential",
		"password.txt",
		"id_rsa",
		"id_dsa",
		"id_ecdsa",
		"id_ed25519",
	}
}

func (a CredentialAssessor) RequiredExtensions() []string {
	return []string{".key", ".secret"}
}

func (a CredentialAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
