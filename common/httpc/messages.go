package httpc

import (
	"bytes"
	"c2/ctrl"
	"encoding/json"
	"net/http"
)

func GetMessages(agentID string, before, after *string, page *int) (messages ctrl.GetMessagesRespCtx, err error) {
	if len(BaseURL) == 0 {
		return ctrl.GetMessagesRespCtx{}, ErrInvalidBaseURL
	}
	resp, err := client.Get(BaseURL + "/v1/messages/" + agentID)
	if err != nil {
		return ctrl.GetMessagesRespCtx{}, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	switch resp.StatusCode {
	default:
		return ctrl.GetMessagesRespCtx{}, ErrUnexpectedStatusCode
	case http.StatusNoContent:
		return ctrl.GetMessagesRespCtx{}, nil
	case http.StatusOK:
		return messages, json.NewDecoder(resp.Body).Decode(&messages)
	}
}

func InsertMessage(agentID, request string) (err error) {
	req, err := http.NewRequest(http.MethodPut, BaseURL+"/v1/messages/"+agentID, bytes.NewBuffer([]byte(request)))
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
