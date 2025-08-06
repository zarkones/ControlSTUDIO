package httpc

import (
	"bytes"
	"c2/ctrl"
	"c2/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetMessagesByIDs(messageIDs *[]string) (msgMap map[string]models.Message, err error) {
	if len(BaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}

	jsonIDs, err := json.Marshal(messageIDs)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, BaseURL+"/v1/messages/by-ids", bytes.NewBuffer(jsonIDs))
	if err != nil {
		return nil, err
	}

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
		return nil, nil
	case http.StatusOK:
		return msgMap, json.NewDecoder(resp.Body).Decode(&msgMap)
	}
}

func GetMessages(agentID string, before, after *string, page *int) (messages ctrl.GetMessagesRespCtx, err error) {
	if len(BaseURL) == 0 {
		return ctrl.GetMessagesRespCtx{}, ErrInvalidBaseURL
	}

	path := "/v1/messages/" + agentID
	args := []string{}
	if len(*after) != 0 {
		args = append(args, "after="+*after)
	}
	if len(*before) != 0 {
		args = append(args, "before="+*before)
	}
	if page != nil {
		args = append(args, "page="+strconv.Itoa(*page))
	}
	if len(args) != 0 {
		path += "?" + strings.Join(args, "&")
	}

	resp, err := client.Get(BaseURL + path)
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
