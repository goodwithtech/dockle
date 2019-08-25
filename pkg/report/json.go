package report

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/xerrors"

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

func (jw JsonWriter) Write(assessments AssessmentSlice) (bool, error) {
	abend := AssessmentSlice{}
	abendAssessments := &abend
	jsonSummary := JsonSummary{}
	jsonDetails := []*JsonDetail{}
	targetType := types.MinTypeNumber
	for targetType <= types.MaxTypeNumber {
		filtered := assessments.FilteredByTargetCode(targetType)
		level, detail := jsonDetail(targetType, filtered)
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
		default:
			jsonSummary.Pass++
		}

		for _, assessment := range filtered {
			abendAssessments.AddAbend(assessment)
		}
		targetType++
	}
	result := JsonOutputFormat{
		Summary: jsonSummary,
		Details: jsonDetails,
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return false, xerrors.Errorf("failed to marshal json: %w", err)
	}

	if _, err = fmt.Fprint(jw.Output, string(output)); err != nil {
		return false, xerrors.Errorf("failed to write json: %w", err)
	}
	return len(*abendAssessments) > 0, nil
}
func jsonDetail(assessmentType int, assessments []*types.Assessment) (level int, jsonInfo *JsonDetail) {
	if len(assessments) == 0 {
		return types.PassLevel, nil
	}
	if assessments[0].Level == types.SkipLevel {
		return types.SkipLevel, nil
	}

	detail := types.AlertDetails[assessmentType]
	level = detail.DefaultLevel
	if assessments[0].Level == types.IgnoreLevel {
		level = types.IgnoreLevel
	}

	alerts := []string{}
	for _, assessment := range assessments {
		alerts = append(alerts, assessment.Desc)
	}
	jsonInfo = &JsonDetail{
		Code:   detail.Code,
		Title:  detail.Title,
		Level:  AlertLabels[level],
		Alerts: alerts,
	}
	return level, jsonInfo
}
