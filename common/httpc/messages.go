package httpc

import (
	"c2/models"
	"encoding/json"
	"net/http"
)

func GetMessages(agentID string, before, after *string, page *int) (messages []models.Message, err error) {
	if len(BaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}
	resp, err := client.Get(BaseURL + "/v1/messages")
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
		return []models.Message{}, nil
	case http.StatusOK:
		return messages, json.NewDecoder(resp.Body).Decode(&messages)
	}
}

func InsertMessage(agentID, request string) (err error) {
	req, err := http.NewRequest(http.MethodPut, BaseURL+"/v1/messages/"+agentID, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return ErrUnexpectedStatusCode
	}
	return nil
}
