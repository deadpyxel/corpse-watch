package scanner

type Result struct {
  URL string `json:"url"`
  Status int `json:"status"`
  Error error `json:"error"`
}
