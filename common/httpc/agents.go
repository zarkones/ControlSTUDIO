package httpc

import (
	"c2/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	BaseURL = ""
	client  = &http.Client{
		Timeout: time.Minute,
	}

	ErrInvalidBaseURL       = errors.New("invalid base url")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

func GetAgents() (agents []models.Agent, err error) {
	if len(BaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}
	resp, err := client.Get(BaseURL + "/v1/agents")
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
		return []models.Agent{}, nil
	case http.StatusOK:
		return agents, json.NewDecoder(resp.Body).Decode(&agents)
	}
}
