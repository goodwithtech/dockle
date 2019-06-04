package utils

import (
	"strings"
	"testing"
)

func TestFetchVersion(t *testing.T) {
	result, err := FetchLatestVersion()
	if err != nil {
		t.Errorf("fail to fetch version : %s", err)
	}
	if !strings.Contains("v", result) {
		t.Errorf("fail to fetch version : %s", result)
	}
}
