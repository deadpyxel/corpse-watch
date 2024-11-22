package scanner

import (
	"fmt"
	"log"
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
			name:      "same domain urls returns true with no errors",
			baseURL:   "https://example.com",
			checkURL:  "https://example.com/page",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "different domain urls returns false with no errors",
			baseURL:   "https://example.com",
			checkURL:  "https://other.com/page",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "invalid domain urls returns false with errors",
			baseURL:   "https://example.com",
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
