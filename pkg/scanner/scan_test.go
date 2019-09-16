package scanner

import (
	"errors"
	"testing"
	"time"

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
		"test-image": {
			fileName:  "",
			imageName: "goodwithtech/test-image:v1",
			option:    deckodertypes.DockerOption{Timeout: time.Minute},
			expected: []*types.Assessment{
				{Type: types.AvoidEmptyPassword},
				{Type: types.AvoidRootDefault},
				{Type: types.AvoidCredentialFile},
				{Type: types.UseCOPY},
				{Type: types.AddHealthcheck},
				{Type: types.MinimizeAptGet},
				{Type: types.AvoidEnvKeySecret},
				{Type: types.UseContentTrust},
			},
		},
		"emptyArg": {
			wantErr: types.ErrSetImageOrFile,
		},
	}
	for name, v := range testcases {
		assesses, err := ScanImage(v.imageName, v.fileName, v.option)
		if !errors.Is(v.wantErr, err) {
			t.Errorf("%s: error got %v, want %v", name, err, v.wantErr)
		}

		cmpopts := []cmp.Option{
			cmpopts.SortSlices(func(x, y *types.Assessment) bool { return x.Type < y.Type }),
			cmpopts.IgnoreTypes(""),
		}
		if diff := cmp.Diff(assesses, v.expected, cmpopts...); diff != "" {
			t.Errorf("%s : tasks diff %v", name, diff)
		}
	}
}
