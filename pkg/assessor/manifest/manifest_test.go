package manifest

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/goodwithtech/dockle/pkg/types"
)

func TestAssess(t *testing.T) {
	var tests = map[string]struct {
		path     string
		assesses []*types.Assessment
	}{
		"RootDefault": {
			path: "./testdata/root_default.json",
			assesses: []*types.Assessment{
				{
					Type:     types.AvoidRootDefault,
					Filename: "docker config",
				},
				{
					Type:     types.AddHealthcheck,
					Filename: "docker config",
				},
			},
		},
		"ApkCached": {
			path: "./testdata/apk_cache.json",

			assesses: []*types.Assessment{
				{
					Type:     types.AvoidRootDefault,
					Filename: "docker config",
				},
				{
					Type:     types.AddHealthcheck,
					Filename: "docker config",
				},
				{
					Type:     types.UseApkAddNoCache,
					Filename: "docker config",
				},
				{
					Type:     types.UseCOPY,
					Filename: "docker config",
				},
			},
		},
	}

	for testname, v := range tests {
		read, err := os.Open(v.path)
		if err != nil {
			t.Errorf("%s : can't open file %s", testname, v.path)
		}
		filebytes, err := ioutil.ReadAll(read)
		if err != nil {
			t.Errorf("%s : can't open file %s", testname, v.path)
		}
		var d types.Image
		err = json.Unmarshal(filebytes, &d)
		if err != nil {
			t.Errorf("%s : failed to unmarshal : %s", testname, v.path)
		}

		actual, err := checkAssessments(d)
		if err != nil {
			t.Errorf("%s : catch the error : %v", testname, err)
		}

		diff, equal := messagediff.PrettyDiff(
			v.assesses,
			actual,
			messagediff.IgnoreStructField("Desc"),
		)
		if !equal {
			t.Errorf("%s diff : %v", testname, diff)
		}
	}
}
