package scanner

import (
	"fmt"
	"net/url"
)

type Result struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
	Error  error  `json:"error"`
}

func (r *Result) String() string {
	errorStr := "no errors"
	if r.Error != nil {
		errorStr = fmt.Sprintf("error [%v]", r.Error)
	}
	return fmt.Sprintf("Got status %d for URL %s with %s", r.Status, r.URL, errorStr)
}

func isSameDomain(baseURL, checkURL string) (bool, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false, err
	}

	check, err := url.Parse(checkURL)
	if err != nil {
		return false, err
	}

	return base.Hostname() == check.Hostname(), nil
}
