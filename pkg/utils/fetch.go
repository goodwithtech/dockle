package utils

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/goodwithtech/dockle/pkg/log"
)

var versionPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

func fetchLocation(ctx context.Context, url string, cookie *http.Cookie) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.AddCookie(cookie)
	resp, err := (&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return nil, fmt.Errorf("HTTP error code : %d, url : %s", resp.StatusCode, url)
	}
	location, err := resp.Location()
	if err != nil {
		return nil, err
	}
	locationString := location.String()
	return &locationString, nil
}

func FetchLatestVersion(ctx context.Context) (version string, err error) {
	log.Logger.Debug("Fetch latest version from github")
	body, err := fetchLocation(
		ctx,
		"https://github.com/goodwithtech/dockle/releases/latest",
		&http.Cookie{Name: "user_session", Value: "guard"},
	)
	if err != nil {
		return "", err
	}
	if versionMatched := versionPattern.FindString(*body); versionMatched != "" {
		return versionMatched, nil
	}
	return "", fmt.Errorf("not found version patterns parsing GH response: %s", *body)
}
