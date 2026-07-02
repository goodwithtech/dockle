package gcr

import (
	"errors"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/docker-credential-gcr/v2/store"

	"github.com/goodwithtech/dockle/pkg/types"
)

func TestCheckOptions(t *testing.T) {
	var tests = map[string]struct {
		domain  string
		opt     types.DockerOption
		gcr     *GCR
		wantErr error
	}{
		"InvalidURL": {
			domain:  "alpine:3.9",
			opt:     types.DockerOption{},
			wantErr: types.InvalidURLPattern,
		},
		"NoOption": {
			domain: "gcr.io",
			opt:    types.DockerOption{},
			gcr:    &GCR{domain: "gcr.io"},
		},
		"CredOption": {
			domain: "gcr.io",
			opt:    types.DockerOption{GcpCredPath: "/path/to/file.json"},
			gcr:    &GCR{domain: "gcr.io", Store: store.NewGCRCredStore("/path/to/file.json")},
		},
	}

	for testname, v := range tests {
		g := &GCR{}
		err := g.CheckOptions(v.domain, v.opt)
		if v.wantErr != nil {
			if err == nil {
				t.Errorf("%s : expected error but no error", testname)
				continue
			}
			if !errors.Is(err, v.wantErr) {
				t.Errorf("[%s]\nexpected error based on %v\nactual : %v", testname, v.wantErr, err)
			}
			continue
		}
		if !reflect.DeepEqual(v.gcr, g) {
			t.Errorf("[%s]\nexpected : %v\nactual : %v", testname, v.gcr, g)
		}
	}
}
