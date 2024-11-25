package scanner

import (
	"io"
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

// parseLinks parses the HTML content of an HTTP response and extracts all valid URLs found in anchor tags.
// It takes a pointer to an http.Response object as input and returns a slice of strings containing the valid URLs.
// The function uses the html package to tokenize the response body and extract anchor tags.
// It filters out URLs with invalid prefixes such as "#" or "javascript:".
// The function also converts relative URLs to absolute URLs using the base URL from the response.
// If no valid URLs are found, an empty slice is returned.
//
// Note: This function assumes that the response body contains valid HTML content with anchor tags.
// It does not handle malformed HTML or non-HTML content gracefully.
func parseLinks(resp *http.Response, opts *LinkParserOptions) []string {
	if opts == nil {
		opts = DefaultLinkParserOptions()
	}
	foundURLs := make([]string, 0)

	// gracefully exit if nil response or closed body
	if resp == nil || resp.Body == nil {
		return foundURLs
	}

	// TODO: find a way to default reader size to 1MiB without this conditional
	readerLimit := int64(1024 * 1024)
	if opts.ReaderLimit > 0 {
		readerLimit = opts.ReaderLimit
	}
	limitedBody := io.LimitReader(resp.Body, readerLimit)
	tokenizer := html.NewTokenizer(limitedBody)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					link := attr.Val
					if isBrowsableURL(link, opts.InvalidPrefixes) {
						absoluteURL := makeAbsoluteURL(resp.Request.URL, link)
						if absoluteURL != "" {
							foundURLs = append(foundURLs, absoluteURL)
						}
					}
				}
			}
		}
	}

	return foundURLs
}
