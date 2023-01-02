package report

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/goodwithtech/dockle/pkg/types"
)

func TestSarifWriter_Write(t *testing.T) {
	tests := []struct {
		description string
		assessments types.AssessmentSlice
		sarif       string
	}{
		{
			description: "Should include location when URI",
			assessments: types.AssessmentSlice{
				{
					Code:     "DKL-DI-0006",
					Filename: "alpine:latest",
					Desc:     "Avoid 'latest' tag",
				},
			},
			sarif: "./testdata/DKL-DI-0006.sarif",
		},
		{
			description: "Should include location when file path",
			assessments: types.AssessmentSlice{
				{
					Code:     "CIS-DI-0010",
					Filename: "/some/abs/path",
					Desc:     "Suspicious filename found",
				},
			},
			sarif: "./testdata/CIS-DI-0010.sarif",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Generate the assessment map
			am := types.CreateAssessmentMap(
				tt.assessments,
				map[string]struct{}{},
				false,
			)

			// Write the serif report to a buffer
			output := &bytes.Buffer{}
			writer := &SarifWriter{Output: output}
			_, err := writer.Write(am)
			if err != nil {
				t.Errorf("Write error: %v", err)
			}

			// parse that JSON into a map for easy comparison
			var actual map[string]interface{}
			err = json.NewDecoder(output).Decode(&actual)
			if err != nil {
				t.Errorf("Decode error: %v", err)
			}

			expected := loadSarifFixture(t, tt.sarif)
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("diff: %v", diff)
			}
		})
	}
}

func loadSarifFixture(t testing.TB, path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Fixture read error: %v", err)
	}

	var sarif map[string]interface{}
	err = json.Unmarshal(data, &sarif)
	if err != nil {
		t.Errorf("Fixture decode error: %v", err)
	}

	return sarif
}
