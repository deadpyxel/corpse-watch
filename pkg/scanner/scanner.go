package scanner

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
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

func isBrowsableURL(url string, invalidPrefixes []string) bool {
	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(url, prefix) {
			return false
		}
	}
	return true
}

func makeAbsoluteURL(baseURL *url.URL, href string) string {
	relURL, err := url.Parse(href)
	if err != nil {
		return ""
	}
	absURL := baseURL.ResolveReference(relURL)
	return absURL.String()
}

func parseLinks(resp *http.Response) []string {
	foundURLs := make([]string, 0)
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				link := attr.Val
				invalidPrefixes := []string{"#", "javascript:", "mailto:", "tel:"}
				if isBrowsableURL(link, invalidPrefixes) {
					absoluteURL := makeAbsoluteURL(resp.Request.URL, link)
					if absoluteURL != "" {
						foundURLs = append(foundURLs, absoluteURL)
					}
				}
			}
		}
	}

	return foundURLs
}
