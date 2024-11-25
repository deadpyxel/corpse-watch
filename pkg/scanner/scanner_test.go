package scanner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"testing"
)

func TestIsSameDomain(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		checkURL  string
		expected  bool
		expectErr bool
	}{
		{
			name:      "when same domain urls returns true with no errors",
			baseURL:   "https://example.com",
			checkURL:  "https://example.com/page",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "when different domain urls returns false with no errors",
			baseURL:   "https://example.com",
			checkURL:  "https://other.com/page",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "when one invalid domain urls returns false with errors",
			baseURL:   "https://example.com",
			checkURL:  ":invalid-url",
			expected:  false,
			expectErr: true,
		},
		{
			name:      "when both invalid domain urls returns false with errors",
			baseURL:   ":another-invalid-url",
			checkURL:  ":invalid-url",
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsSameDomain(tt.baseURL, tt.checkURL)

			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got error: %v", tt.expectErr, err)
			}

			if result != tt.expected {
				t.Errorf("For baseURL: %s and checkURL: %s, expected: %v, got: %v", tt.baseURL, tt.checkURL, tt.expected, result)
			}
		})
	}
}

func ExampleIsSameDomain() {
	url1 := "https://example.com"
	url2 := "https://example.com/page"

	isSame, err := IsSameDomain(url1, url2)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("url1: %s and url2: %s belong to the same domain? %v", url1, url2, isSame)
	// Output:
	// url1: https://example.com and url2: https://example.com/page belong to the same domain? true
}

func TestIsBrowsableURL(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		invalidPrefixes []string
		expected        bool
	}{
		{
			name:            "when given a url without any invalid prefix returns true",
			url:             "https://www.example.com",
			invalidPrefixes: []string{"http://", "ftp://"},
			expected:        true,
		},
		{
			name:            "when given a url with an invalid prefix returns false",
			url:             "ftp://example.com",
			invalidPrefixes: []string{"http://", "ftp://"},
			expected:        false,
		},
		{
			name:            "when given a empty slice of invalid prefixes returns true",
			url:             "http://example.com",
			invalidPrefixes: []string{},
			expected:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isBrowsableURL(test.url, test.invalidPrefixes)
			if result != test.expected {
				t.Errorf("For URL %s and invalid prefixes %v, expected %t but got %t", test.url, test.invalidPrefixes, test.expected, result)
			}
		})
	}
}

func TestMakeAbsoluteURL(t *testing.T) {
	baseURL, _ := url.Parse("https://www.example.com")
	tests := []struct {
		name    string
		baseURL *url.URL
		href    string
		want    string
	}{
		{"relative path", baseURL, "/page1", "https://www.example.com/page1"},
		{"absolute URL", baseURL, "https://www.example.com/page2", "https://www.example.com/page2"},
		{"relative path not formatted", baseURL, "page3", "https://www.example.com/page3"},
		{"invalid path", baseURL, ":invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeAbsoluteURL(tt.baseURL, tt.href)
			if got != tt.want {
				t.Errorf("makeAbsoluteURL(%s, %s) results in %s; want %s,", tt.baseURL, tt.href, got, tt.want)
			}
		})
	}
}

// createMockHttpResponse creates a mock http.Response object with the given HTML content and base URL.
func createMockHttpResponse(htmlContent string, baseURL string) *http.Response {
	parsedURL, _ := url.Parse(baseURL)
	resp := &http.Response{
		Body:    io.NopCloser(bytes.NewBufferString(htmlContent)),
		Request: &http.Request{URL: parsedURL},
	}
	return resp
}

func TestParseLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		baseURL  string
		expected []string
	}{
		{
			name:     "When no links in response return empty slice",
			html:     "<html><body><p>Dummy</p></body></html>",
			baseURL:  "http://example.com",
			expected: []string{},
		},
		{
			name:     "When single relative link is in response returns absolute URL",
			html:     `<html><body><a href="/about">About</a></body></html>`,
			baseURL:  "http://example.com",
			expected: []string{"http://example.com/about"},
		},
		{
			name:     "When multiple links are in response returns all present URLs",
			html:     `<html><body><a href="/about">About</a><a href="/contact">Contact</a></body></html>`,
			baseURL:  "http://example.com",
			expected: []string{"http://example.com/about", "http://example.com/contact"},
		},
		{
			name: "When links with invalid prefixes are present only returns the valid links",
			html: `<html><body>
				<a href="#">Skip</a>
				<a href="javascript:void(0)">JS Link</a>
				<a href="mailto:test@example.com">Email</a>
				<a href="tel:1234567890">Phone</a>
				<a href="/valid">Valid Link</a>
			</body></html>`,
			baseURL:  "http://example.com",
			expected: []string{"http://example.com/valid"},
		},
		{
			name:     "When external links are present they are returned",
			html:     `<html><body><a href="http://another.com/page">External Link</a></body></html>`,
			baseURL:  "http://example.com",
			expected: []string{"http://another.com/page"},
		},
		{
			name:     "When missing href returns an empty slice",
			html:     `<html><body><a target="_blank">External Link</a></body></html>`,
			baseURL:  "http://example.com",
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := createMockHttpResponse(tc.html, tc.baseURL)
			result := parseLinks(resp, nil)

			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d links, got %d instead", len(tc.expected), len(result))
			}

			for _, expectedURL := range tc.expected {
				if !slices.Contains(result, expectedURL) {
					t.Errorf("Missing expected URL: %s", expectedURL)
				}
			}

		})
	}
}

func TestParseLinksErrors(t *testing.T) {
	t.Run("When an empty response returns empty slice", func(t *testing.T) {
		result := parseLinks(nil, nil)

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %d elements instead", len(result))
		}
	})
}
