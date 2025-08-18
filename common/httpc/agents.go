package httpc

import (
	"c2/models"
	"encoding/json"
	"net/http"
)

func GetAgents() (agents []models.Agent, err error) {
	if len(BaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}
	req, err := http.NewRequest(http.MethodGet, BaseURL+"/v1/agents", nil)
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
		return []models.Agent{}, nil
	case http.StatusOK:
		return agents, json.NewDecoder(resp.Body).Decode(&agents)
	}
}
