package scanner

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/Portshift/dockle/config"
	"github.com/Portshift/dockle/pkg/assessor/contentTrust"
	"github.com/Portshift/dockle/pkg/assessor/manifest"
	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/types"
)

func TestScanImage(t *testing.T) {
	log.InitLogger(false, false)
	testcases := map[string]struct {
		config   config.Config
		wantErr  error
		expected []*types.Assessment
	}{
		"Dockerfile.base": {
			// TODO : too large to use github / fileName:  "base.tar",
			// testdata/Dockerfile.base
			config: config.Config{
				ImageName: "goodwithtech/dockle-test:base-test",
			},
			expected: []*types.Assessment{
				{Code: types.AvoidEmptyPassword, Filename: "/etc/shadow"},
				{Code: types.AvoidRootDefault, Filename: manifest.ConfigFileName},
				{Code: types.AvoidCredential, Filename: "/app/credentials.json"},
				{Code: types.CheckSuidGuid, Filename: "/app/gid.txt"},
				{Code: types.CheckSuidGuid, Filename: "/app/suid.txt"},
				{Code: types.CheckSuidGuid, Filename: "/bin/mount"},
				{Code: types.CheckSuidGuid, Filename: "/bin/su"},
				{Code: types.CheckSuidGuid, Filename: "/bin/umount"},
				{Code: types.CheckSuidGuid, Filename: "/usr/lib/openssh/ssh-keysign"},
				{Code: types.UseCOPY, Filename: manifest.ConfigFileName},
				{Code: types.AddHealthcheck, Filename: manifest.ConfigFileName},
				{Code: types.MinimizeAptGet, Filename: manifest.ConfigFileName},
				{Code: types.AvoidCredential, Filename: manifest.ConfigFileName},
				{Code: types.UseContentTrust, Filename: contentTrust.HostEnvironmentFileName},
			},
		},
		"Dockerfile.scratch": {
			config: config.Config{
				FilePath: "./testdata/scratch.tar",
			},
			expected: []*types.Assessment{
				{Code: types.AvoidCredential, Filename: "/credentials.json"},
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
		assesses, err := ScanImage(ctx, v.config)
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
