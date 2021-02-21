package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/goodwithtech/dockle/pkg/log"
)

var versionPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

// Dockle just want to check latest version string. No need to readall.
const enoughLength = 8000

func fetchURL(ctx context.Context, url string, cookie *http.Cookie, dataLen int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.AddCookie(cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error code : %d, url : %s", resp.StatusCode, url)
	}
	data := make([]byte, dataLen)
	if _, err := io.ReadFull(resp.Body, data); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return data, nil
}

func FetchLatestVersion(ctx context.Context) (version string, err error) {
	log.Logger.Debug("Fetch latest version from github")
	body, err := fetchURL(
		ctx,
		"https://github.com/goodwithtech/dockle/releases/latest",
		&http.Cookie{Name: "user_session", Value: "guard"},
		enoughLength,
	)
	if err != nil {
		return "", err
	}
	versionMatched := versionPattern.FindString(string(body))
	return versionMatched, nil
}
