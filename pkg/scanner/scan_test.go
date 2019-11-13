package scanner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goodwithtech/dockle/pkg/assessor/contentTrust"

	"github.com/goodwithtech/dockle/pkg/assessor/manifest"

	"github.com/google/go-cmp/cmp/cmpopts"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/google/go-cmp/cmp"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

func TestScanImage(t *testing.T) {
	log.InitLogger(false)
	testcases := map[string]struct {
		imageName string
		fileName  string
		option    deckodertypes.DockerOption
		wantErr   error
		expected  []*types.Assessment
	}{
		"Dockerfile.base": {
			fileName: "",
			// testdata/Dockerfile.base
			imageName: "goodwithtech/dockle-test:base-test",
			option:    deckodertypes.DockerOption{Timeout: time.Minute},
			expected: []*types.Assessment{
				{Type: types.AvoidEmptyPassword, Filename: "etc/shadow"},
				{Type: types.AvoidRootDefault, Filename: manifest.ConfigFileName},
				{Type: types.AvoidCredentialFile, Filename: "app/credentials.json"},
				{Type: types.CheckSuidGuid, Filename: "app/gid.txt"},
				{Type: types.CheckSuidGuid, Filename: "app/suid.txt"},
				{Type: types.CheckSuidGuid, Filename: "bin/mount"},
				{Type: types.CheckSuidGuid, Filename: "bin/su"},
				{Type: types.CheckSuidGuid, Filename: "bin/umount"},
				{Type: types.CheckSuidGuid, Filename: "usr/lib/openssh/ssh-keysign"},
				{Type: types.UseCOPY, Filename: manifest.ConfigFileName},
				{Type: types.AddHealthcheck, Filename: manifest.ConfigFileName},
				{Type: types.MinimizeAptGet, Filename: manifest.ConfigFileName},
				{Type: types.AvoidEnvKeySecret, Filename: manifest.ConfigFileName},
				{Type: types.UseContentTrust, Filename: contentTrust.HostEnvironmentFileName},
			},
		},
		"emptyArg": {
			wantErr: types.ErrSetImageOrFile,
		},
	}
	for name, v := range testcases {
		ctx := context.Background()
		assesses, err := ScanImage(ctx, v.imageName, v.fileName, v.option)
		if !errors.Is(v.wantErr, err) {
			t.Errorf("%s: error got %v, want %v", name, err, v.wantErr)
		}

		cmpopts := []cmp.Option{
			cmpopts.SortSlices(func(x, y *types.Assessment) bool {
				if x.Type == y.Type {
					return x.Filename < y.Filename
				}
				return x.Type < y.Type
			}),
			cmpopts.IgnoreFields(types.Assessment{}, "Desc"),
		}
		if diff := cmp.Diff(assesses, v.expected, cmpopts...); diff != "" {
			t.Errorf("%s : tasks diff %v", name, diff)
		}
	}
}
