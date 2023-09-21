package credential

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/Portshift/dockle/pkg/log"

	"github.com/Portshift/dockle/pkg/types"
)

type CredentialAssessor struct{}

func (a CredentialAssessor) Assess(imageData *types.ImageData) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : credential files")
	assesses := []*types.Assessment{}
	fmap := makeMaps(a.RequiredFiles())
	fexts := makeMaps(a.RequiredExtensions())
	for filename := range imageData.FileMap {
		basename := filepath.Base(filename)
		// check exist target files
		if _, ok := fmap[basename]; ok {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.AvoidCredential,
					Filename: filename,
					Desc:     fmt.Sprintf("Suspicious filename found : %s (You can suppress it with \"-af %s\")", filename, basename),
				})
		} else if _, ok := fexts[filepath.Ext(basename)]; ok {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.AvoidCredential,
					Filename: filename,
					Desc:     fmt.Sprintf("Suspicious file extension found : %s (You can suppress it with \"-ae %s\")", filename, trimFirstRune(filepath.Ext(basename))),
				})
		}
	}
	return assesses, nil
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
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
		// TODO: Only check .docker/config.json
		// "config.json",
		"credentials",
		"credential",
		"password.txt",
		"id_rsa",
		"id_dsa",
		"id_ecdsa",
		"id_ed25519",
		"secret_token.rb",
		"carrierwave.rb",
		"omniauth.rb",
		"settings.py",
		"database.yml",
		"credentials.xml",
	}
}

func (a CredentialAssessor) RequiredExtensions() []string {
	return []string{
		// reference: https://github.com/eth0izzle/shhgit/blob/master/config.yaml
		// TODO: potential sensitive data but they have many false-positives.
		//       Dockle need to analyze each file.
		//".pem",
		//".key",
		//".p12",
		//".pkcs12",
		//".pfx",
		//".asc",

		".secret",
		".ovpn",
		".private_key",
		".cscfg",
		".rdp",
		".mdf",
		".sdf",
		".bek",
		".tpm",
		".fve",
		".jks",
		".psafe3",
		".agilekeychain",
		".keychain",
		".pcap",
		".gnucache",
	}
}

func (a CredentialAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
