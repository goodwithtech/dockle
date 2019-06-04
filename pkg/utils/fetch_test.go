package utils

import (
	"net/http"
	"strings"
	"testing"
)

func TestFetchURL(t *testing.T) {
	var tests = map[string]struct {
		url     string
		cookie  *http.Cookie
		expect  string
		wantErr error
	}{
		"Github": {
			url: "https://github.com/goodwithtech/docker-guard/releases/latest",
			cookie: &http.Cookie{
				Name:  "user_session",
				Value: "guard",
			},
			expect: "hoge",
		},
	}

	for testname, v := range tests {

		result, err := fetchURL(v.url, v.cookie)
		if err != nil {
			t.Errorf("%s : fail to fetch %s", testname, v.url)
		}
		if string(result) != v.expect {
			t.Errorf("%s : want %s, actual %s", testname, v.expect, string(result))
		}
	}
}

func TestFetchVersion(t *testing.T) {
	result, err := FetchLatestVersion()
	if err != nil {
		t.Errorf("fail to fetch version : %s", err)
	}
	if !strings.Contains("v", result) {
		t.Errorf("fail to fetch version : %s", result)
	}
}
