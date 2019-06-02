package manifest

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/image"

	"github.com/goodwithtech/docker-guard/pkg/types"
)

func TestAssess(t *testing.T) {
	var tests = map[string]struct {
		path     string
		assesses []types.Assessment
	}{
		"Valid": {
			path: "./testdata/root_default.json",
			assesses: []types.Assessment{
				{
					Type:     types.AvoidRootDefault,
					Filename: "docker config",
					Desc:     "Avoid default user set root",
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
		var d image.Image
		json.Unmarshal(filebytes, &d)

		actual, err := checkAssessments(d)
		if err != nil {
			t.Errorf("%s : catch the error : %v", testname, err)
		}
		if !reflect.DeepEqual(v.assesses, actual) {
			t.Errorf("[%s]\nexpected : %v\nactual : %v", testname, v.assesses, actual)
		}
	}
}
