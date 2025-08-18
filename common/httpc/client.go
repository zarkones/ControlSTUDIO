package httpc

import (
	"net/http"
	"time"
)

var (
	BaseURL = ""
	client  = &http.Client{
		Timeout: time.Minute,
	}
	AuthHeaderValue = ""
)

func setAuthHeader(r *http.Request) {
	if r == nil || len(AuthHeaderValue) == 0 {
		return
	}
	r.Header.Set("Authorization", AuthHeaderValue)
}
