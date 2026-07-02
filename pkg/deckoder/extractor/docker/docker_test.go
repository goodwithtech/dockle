package docker

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencontainers/go-digest"

	"github.com/stretchr/testify/require"

	"github.com/goodwithtech/dockle/pkg/deckoder/extractor/image"
	"github.com/goodwithtech/dockle/pkg/deckoder/types"
	"github.com/goodwithtech/dockle/pkg/deckoder/utils"
)

const (
	NormalFileMode os.FileMode = 0644
	SuFileMode     os.FileMode = 0600
	SetSuidNormal  os.FileMode = 040000644
)

func TestExtractor_ExtractLayerFiles(t *testing.T) {
	type fields struct {
		option types.DockerOption
		image  image.RealImage
	}
	type args struct {
		ctx    context.Context
		dig    digest.Digest
		filter types.FilterFunc
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		imagePath       string
		expectedDigest  digest.Digest
		expectedFileMap types.FileMap
		expectedOpqDirs []string
		expectedWhFiles []string
		wantErr         string
	}{
		{
			name:      "happy path",
			imagePath: "testdata/image1.tar",
			args: args{
				ctx:    nil,
				dig:    "sha256:d9ff549177a94a413c425ffe14ae1cc0aa254bc9c7df781add08e7d2fba25d27",
				filter: utils.CreateFilterPathFunc([]string{"etc/hostname"}),
			},
			expectedDigest: "sha256:d9ff549177a94a413c425ffe14ae1cc0aa254bc9c7df781add08e7d2fba25d27",
			expectedFileMap: types.FileMap{
				"etc/hostname": types.FileData{
					Body:     []byte("localhost\n"),
					FileMode: NormalFileMode,
				},
			},
		},
		{
			name:      "symbolic path",
			imagePath: "testdata/symbolic.tar",
			args: args{
				ctx:    nil,
				dig:    "sha256:98d172aa39eb52759aa79fda88452c3d78528ea21b170332cf45759c902c519a",
				filter: utils.CreateFilterPathFunc([]string{"app/once-suid.txt"}),
			},
			expectedDigest: "sha256:677c191235a22a3125057f747c7f44b606e17701b7a273c1d2ff1d8dc825deea",
			expectedFileMap: types.FileMap{
				"app/once-suid.txt": {Body: []byte(""), FileMode: NormalFileMode},
			},
		},
		{
			name:      "symbolic path",
			imagePath: "testdata/symbolic.tar",
			args: args{
				ctx:    nil,
				dig:    "sha256:434ba219e3907e89fe29f9b7de597fdf2305c615356f6a760e880570486fb4bb",
				filter: utils.CreateFilterPathFunc([]string{"app/once-suid.txt"}),
			},
			expectedDigest: "sha256:677c191235a22a3125057f747c7f44b606e17701b7a273c1d2ff1d8dc825deea",
			expectedFileMap: types.FileMap{
				"app/once-suid.txt": {Body: []byte(""), FileMode: SetSuidNormal},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, cleanup, err := NewDockerArchiveExtractor(context.Background(), tt.imagePath, types.DockerOption{})
			require.NoError(t, err)
			defer cleanup()

			actualFileMap, actualOpqDirs, actualWhFiles, err := d.ExtractLayerFiles(tt.args.ctx, tt.args.dig, tt.args.filter)
			if tt.wantErr != "" {
				require.NotNil(t, err, tt.name)
				assert.Contains(t, err.Error(), tt.wantErr, tt.name)
				return
			} else {
				require.NoError(t, err, tt.name)
			}

			assert.Equal(t, tt.expectedFileMap, actualFileMap)
			assert.Equal(t, tt.expectedOpqDirs, actualOpqDirs)
			assert.Equal(t, tt.expectedWhFiles, actualWhFiles)
		})
	}
}
