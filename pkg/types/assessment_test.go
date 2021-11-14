package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateAssessmentMap(t *testing.T) {
	testcases := map[string]struct {
		as       AssessmentSlice
		ig       map[string]struct{}
		debug    bool
		expected AssessmentMap
	}{
		"OK": {
			as: AssessmentSlice{
				{Code: "a", Filename: "a"},
				{Code: "b", Filename: "b"},
				{Code: "a", Filename: "c"},
				{Code: "a", Filename: "b"},
			},
			ig: map[string]struct{}{},
			expected: map[string]CodeInfo{
				"a": {
					Code:  "a",
					Level: 0,
					Assessments: []*Assessment{
						{Code: "a", Filename: "a"},
						{Code: "a", Filename: "c"},
						{Code: "a", Filename: "b"},
					},
				},
				"b": {
					Code:  "b",
					Level: 0,
					Assessments: []*Assessment{
						{Code: "b", Filename: "b"},
					},
				},
			},
		},
		"IgnoreB": {
			as: AssessmentSlice{
				{Code: "a", Filename: "a"},
				{Code: "b", Filename: "b"},
				{Code: "a", Filename: "c"},
				{Code: "a", Filename: "b"},
			},
			ig: map[string]struct{}{"b": {}},
			expected: map[string]CodeInfo{
				"a": {
					Code:  "a",
					Level: 0,
					Assessments: []*Assessment{
						{Code: "a", Filename: "a"},
						{Code: "a", Filename: "c"},
						{Code: "a", Filename: "b"},
					},
				},
			},
		},
		"IgnoreBwithDebug": {
			as: AssessmentSlice{
				{Code: "a", Filename: "a"},
				{Code: "b", Filename: "b"},
				{Code: "a", Filename: "c"},
				{Code: "a", Filename: "b"},
			},
			ig:    map[string]struct{}{"b": {}},
			debug: true,
			expected: map[string]CodeInfo{
				"a": {
					Code:  "a",
					Level: 0,
					Assessments: []*Assessment{
						{Code: "a", Filename: "a"},
						{Code: "a", Filename: "c"},
						{Code: "a", Filename: "b"},
					},
				},
				"b": {
					Code:  "b",
					Level: IgnoreLevel,
					Assessments: []*Assessment{
						{Code: "b", Filename: "b"},
					},
				},
			},
		},
	}

	for name, v := range testcases {
		actual := CreateAssessmentMap(v.as, v.ig, v.debug)
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
