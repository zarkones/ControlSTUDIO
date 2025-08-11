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
)
