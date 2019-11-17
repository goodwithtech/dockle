package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goodwithtech/dockle/config"

	"github.com/goodwithtech/dockle/pkg/types"
)

type JsonWriter struct {
	Output io.Writer
}

type JsonOutputFormat struct {
	Summary JsonSummary   `json:"summary"`
	Details []*JsonDetail `json:"details"`
}
type JsonSummary struct {
	Fatal int `json:"fatal"`
	Warn  int `json:"warn"`
	Info  int `json:"info"`
	Pass  int `json:"pass"`
}
type JsonDetail struct {
	Code   string   `json:"code"`
	Title  string   `json:"title"`
	Level  string   `json:"level"`
	Alerts []string `json:"alerts"`
}

func (jw JsonWriter) Write(assessMap types.AssessmentMap) (abend bool, err error) {
	jsonSummary := JsonSummary{}
	jsonDetails := []*JsonDetail{}
	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		if _, ok := assessMap[ass.Code]; !ok {
			jsonSummary.Pass++
			continue
		}
		assesses := assessMap[ass.Code].Assessments
		detail := jsonDetail(ass.Code, ass.Level, assesses)
		if detail != nil {
			jsonDetails = append(jsonDetails, detail)
		}

		// increment summary
		switch ass.Level {
		case types.FatalLevel:
			jsonSummary.Fatal++
		case types.WarnLevel:
			jsonSummary.Warn++
		case types.InfoLevel:
			jsonSummary.Info++
		}
		if ass.Level >= config.Conf.ExitLevel {
			abend = true
		}
	}
	result := JsonOutputFormat{
		Summary: jsonSummary,
		Details: jsonDetails,
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return false, fmt.Errorf("failed to marshal json: %w", err)
	}

	if _, err = fmt.Fprint(jw.Output, string(output)); err != nil {
		return false, fmt.Errorf("failed to write json: %w", err)
	}
	return abend, nil
}
func jsonDetail(code string, level int, assessments []*types.Assessment) (jsonInfo *JsonDetail) {
	if len(assessments) == 0 {
		return nil
	}
	alerts := []string{}
	for _, assessment := range assessments {
		alerts = append(alerts, assessment.Desc)
	}
	jsonInfo = &JsonDetail{
		Code:   code,
		Title:  types.TitleMap[code],
		Level:  AlertLabels[level],
		Alerts: alerts,
	}
	return jsonInfo
}
