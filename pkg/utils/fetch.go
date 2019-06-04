package utils

import (
	"net/http"
	"regexp"

	"github.com/goodwithtech/docker-guard/pkg/log"

	"golang.org/x/xerrors"

	"github.com/parnurzeal/gorequest"
)

var versionPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

func fetchURL(url string, cookie *http.Cookie) ([]byte, error) {
	resp, body, err := gorequest.New().AddCookie(cookie).Get(url).Type("text").EndBytes()
	if err != nil {
		return nil, xerrors.Errorf("fail to fetch : %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, xerrors.Errorf("HTTP error code : %d, url : %s", resp.StatusCode, url)
	}
	return body, nil
}

func FetchLatestVersion() (version string, err error) {
	log.Logger.Debug("Fetch latest version from github")
	body, err := fetchURL(
		"https://github.com/goodwithtech/docker-guard/releases/latest",
		&http.Cookie{Name: "user_session", Value: "guard"},
	)
	if err != nil {
		return "", err
	}
	versionMatched := versionPattern.FindString(string(body))
	return versionMatched, nil
}
