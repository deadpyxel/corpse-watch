package scanner

import (
	"fmt"
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

type LinkParserOptions struct {
	InvalidPrefixes   []string // invalid prefixed for non browsable links
	ReaderLimit       int64    // limit read size of the response body
	SkipExternalLinks bool     // skip adding external links to results
}

func DefaultLinkParserOptions() *LinkParserOptions {
	return &LinkParserOptions{
		InvalidPrefixes:   []string{"#", "javascript:", "mailto:", "tel:"},
		ReaderLimit:       1024 * 1024,
		SkipExternalLinks: true,
	}
}
