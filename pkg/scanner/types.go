package scanner

import "fmt"

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
