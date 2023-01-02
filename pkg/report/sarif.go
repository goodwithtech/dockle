package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/owenrumney/go-sarif/v2/sarif"

	"github.com/goodwithtech/dockle/config"
	"github.com/goodwithtech/dockle/pkg/types"
)

type SarifWriter struct {
	Output io.Writer
}

type sarifResult struct {
	ruleID          string
	ruleDescription string
	link            string
	description     string
	severity        string
	locations       []string
}

func (sw SarifWriter) Write(assessMap types.AssessmentMap) (abend bool, err error) {
	var rules []*sarifResult
	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		if _, ok := assessMap[ass.Code]; !ok {
			continue
		}
		assess := assessMap[ass.Code]
		detail := sarifDetail(assess.Code, assess.Level, assess.Assessments)
		if detail != nil {
			rules = append(rules, detail)
		}
		if assess.Level >= config.Conf.ExitLevel {
			abend = true
		}
	}

	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return false, err
	}
	run := sarif.NewRunWithInformationURI("Dockle", "https://github.com/goodwithtech/dockle")
	report.AddRun(run)
	for _, r := range rules {
		result := sarif.NewRuleResult(r.ruleID).
			WithLevel(strings.ToLower(r.severity)).
			WithMessage(sarif.NewTextMessage(r.description))

		for _, uri := range r.locations {
			result.AddLocation(
				sarif.NewLocation().WithPhysicalLocation(
					sarif.NewPhysicalLocation().WithArtifactLocation(
						sarif.NewArtifactLocation().WithUri(
							uri,
						),
					),
				),
			)
		}

		run.AddRule(r.ruleID).
			WithName(r.ruleID).
			WithDescription(r.ruleDescription).
			WithHelpURI(r.link)
		run.AddResult(result)
	}
	if err := report.PrettyWrite(sw.Output); err != nil {
		return false, fmt.Errorf("failed to write sarif: %w", err)
	}
	return abend, nil
}

func sarifDetail(code string, level int, assessments []*types.Assessment) (jsonInfo *sarifResult) {
	if len(assessments) == 0 {
		return nil
	}
	alerts := []string{}
	locations := []string{}
	for _, assessment := range assessments {
		alerts = append(alerts, assessment.Desc)
		locations = append(locations, assessment.Filename)
	}
	return &sarifResult{
		ruleID:          code,
		severity:        sarifAlertLabels[level],
		ruleDescription: types.TitleMap[code],
		link:            fmt.Sprintf("https://github.com/goodwithtech/dockle/blob/master/CHECKPOINT.md#%s", code),
		description:     strings.Join(alerts, ", "),
		locations:       locations,
	}
}
