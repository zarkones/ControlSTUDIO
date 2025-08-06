package ctrl

import (
	"c2/models"
	"c2/repos"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type GetMessagesRespCtx struct {
	Messages []models.Message
	Before   string
	After    string
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agentID")

	q := r.URL.Query()
	before, errBefore := strconv.ParseInt(q.Get("before"), 10, 64)
	after, errAfter := strconv.ParseInt(q.Get("after"), 10, 64)
	page, _ := strconv.Atoi(q.Get("page"))
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = repos.DEFAULT_LIMIT
	}
	offset := page * limit

	var messages []models.Message

	if errAfter == nil {
		messages, err = repos.GetMessagesAfter(agentID, after, limit)
	} else {
		if errBefore == nil {
			messages, err = repos.GetMessagesBefore(agentID, before, limit)
		} else {
			messages, err = repos.GetMessages(agentID, offset, limit)
		}
	}

	if err != nil {
		log.Println("api: error: GetMessages:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := GetMessagesRespCtx{
		Messages: messages,
		Before:   fmt.Sprint(messages[len(messages)-1].CreatedAt),
		After:    fmt.Sprint(messages[0].CreatedAt),
	}

	jj(w, &resp)
}

type GetMessageUpdatesReqCtx []string
type GetMessageUpdatesRespCtx map[string]models.Message

func GetMessageByIDs(w http.ResponseWriter, r *http.Request) {
	var messageIDs GetMessageUpdatesReqCtx

	if err := json.NewDecoder(r.Body).Decode(&messageIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if messageIDs == nil || len(messageIDs) == 0 {
		http.Error(w, "no ids supplied", http.StatusUnprocessableEntity)
		return
	}

	messages, err := repos.GetMessagesByIDs(messageIDs)
	if err != nil {
		log.Println("api: error: GetMessages:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	msgMap := make(GetMessageUpdatesRespCtx, len(messages))
	for _, msg := range messages {
		msgMap[msg.ID] = msg
	}

	jj(w, &msgMap)
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agentID")

	msgRequest, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(msgRequest) == 0 {
		http.Error(w, "invalid message request", http.StatusBadRequest)
		return
	}

	if err := repos.InsertMessage(&models.Message{
		AgentID: agentID,
		Request: string(msgRequest),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
