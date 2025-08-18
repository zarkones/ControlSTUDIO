package httpc

import (
	"c2/core/listeners"
	"encoding/json"
	"net/http"
)

func GetListeners() (ls []listeners.Listener, err error) {
	if len(BaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}
	req, err := http.NewRequest(http.MethodGet, BaseURL+"/v1/listeners", nil)
	if err != nil {
		return nil, err
	}
	setAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	switch resp.StatusCode {
	default:
		return nil, ErrUnexpectedStatusCode
	case http.StatusNoContent:
		return []listeners.Listener{}, nil
	case http.StatusOK:
		return ls, json.NewDecoder(resp.Body).Decode(&ls)
	}
}
