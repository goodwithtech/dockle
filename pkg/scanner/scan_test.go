package scanner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/google/go-cmp/cmp"

	"github.com/goodwithtech/dockle/pkg/assessor/contentTrust"
	"github.com/goodwithtech/dockle/pkg/assessor/manifest"
	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

func TestScanImage(t *testing.T) {
	log.InitLogger(false, false)
	testcases := map[string]struct {
		imageName string
		fileName  string
		option    deckodertypes.DockerOption
		wantErr   error
		expected  []*types.Assessment
	}{
		"Dockerfile.base": {
			// TODO : too large to use github / fileName:  "base.tar",
			// testdata/Dockerfile.base
			imageName: "goodwithtech/dockle-test:base-test",
			option:    deckodertypes.DockerOption{Timeout: time.Minute},
			expected: []*types.Assessment{
				{Code: types.AvoidEmptyPassword, Filename: "etc/shadow"},
				{Code: types.AvoidRootDefault, Filename: manifest.ConfigFileName},
				{Code: types.AvoidCredential, Filename: "app/credentials.json"},
				{Code: types.CheckSuidGuid, Filename: "app/gid.txt"},
				{Code: types.CheckSuidGuid, Filename: "app/suid.txt"},
				{Code: types.CheckSuidGuid, Filename: "bin/mount"},
				{Code: types.CheckSuidGuid, Filename: "bin/su"},
				{Code: types.CheckSuidGuid, Filename: "bin/umount"},
				{Code: types.CheckSuidGuid, Filename: "usr/lib/openssh/ssh-keysign"},
				{Code: types.UseCOPY, Filename: manifest.ConfigFileName},
				{Code: types.AddHealthcheck, Filename: manifest.ConfigFileName},
				{Code: types.MinimizeAptGet, Filename: manifest.ConfigFileName},
				{Code: types.AvoidCredential, Filename: manifest.ConfigFileName},
				{Code: types.UseContentTrust, Filename: contentTrust.HostEnvironmentFileName},
			},
		},
		"Dockerfile.scratch": {
			fileName: "./testdata/scratch.tar",
			expected: []*types.Assessment{
				{Code: types.AvoidCredential, Filename: "credentials.json"},
				{Code: types.AddHealthcheck, Filename: manifest.ConfigFileName},
				{Code: types.UseContentTrust, Filename: contentTrust.HostEnvironmentFileName},
				{Code: types.AvoidEmptyPassword, Level: types.SkipLevel},
				{Code: types.AvoidDuplicateUserGroup, Level: types.SkipLevel},
				{Code: types.AvoidDuplicateUserGroup, Level: types.SkipLevel},
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
				if x.Code == y.Code {
					return x.Filename < y.Filename
				}
				return x.Code < y.Code
			}),
			cmpopts.IgnoreFields(types.Assessment{}, "Desc"),
		}
		if diff := cmp.Diff(assesses, v.expected, cmpopts...); diff != "" {
			t.Errorf("%s : tasks diff %v", name, diff)
		}
	}
}
