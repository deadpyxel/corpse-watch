package scanner

import "net/url"

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
