package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateAssessmentMap(t *testing.T) {
	testcases := map[string]struct {
		as       AssessmentSlice
		expected AssessmentMap
	}{

		"OK": {
			as: AssessmentSlice{
				{Code: "a", Filename: "a"},
				{Code: "b", Filename: "b"},
				{Code: "a", Filename: "c"},
				{Code: "a", Filename: "b"},
			},
			expected: map[string][]*Assessment{
				"a": {
					{Code: "a", Filename: "a"},
					{Code: "a", Filename: "c"},
					{Code: "a", Filename: "b"},
				},
				"b": {
					{Code: "b", Filename: "b"},
				},
			},
		},
	}

	for name, v := range testcases {
		actual := CreateAssessmentMap(v.as)
		cmpopts := []cmp.Option{
			cmpopts.SortSlices(func(x, y Assessment) bool {
				if x.Code == y.Code {
					return x.Filename < y.Filename
				}
				return x.Code < y.Code
			}),
		}
		if diff := cmp.Diff(actual, v.expected, cmpopts...); diff != "" {
			t.Errorf("%s : diff %v", name, diff)
		}
	}
}
