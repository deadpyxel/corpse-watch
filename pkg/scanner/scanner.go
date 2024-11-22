package scanner

import (
	"net/url"
)

// IsSameDomain checks if two URLs belong to the same domain.
// Returns an error if an invalid URL is passed
func IsSameDomain(baseURL, checkURL string) (bool, error) {
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
