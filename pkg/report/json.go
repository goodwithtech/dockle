package report

import (
	"encoding/json"
	"fmt"
	"io"

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

func (jw JsonWriter) Write(assessMap types.AssessmentMap) (bool, error) {
	abend := types.AssessmentSlice{}
	abendAssessments := &abend
	jsonSummary := JsonSummary{}
	jsonDetails := []*JsonDetail{}
	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		assesses, ok := assessMap[ass.Code]
		if !ok {
			jsonSummary.Pass++
			continue
		}
		level, detail := jsonDetail(ass.Code, assesses)
		if detail != nil {
			jsonDetails = append(jsonDetails, detail)
		}

		// increment summary
		switch level {
		case types.FatalLevel:
			jsonSummary.Fatal++
		case types.WarnLevel:
			jsonSummary.Warn++
		case types.InfoLevel:
			jsonSummary.Info++
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
	return len(*abendAssessments) > 0, nil
}
func jsonDetail(code string, assessments []*types.Assessment) (level int, jsonInfo *JsonDetail) {
	if len(assessments) == 0 {
		return types.PassLevel, nil
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
	return level, jsonInfo
}
