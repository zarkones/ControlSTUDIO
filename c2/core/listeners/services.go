package listeners

import (
	"common/profiles"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ListenerService struct {
	Server *http.Server
}

func errorHandler(err error, r *http.Request) (abort bool) {
	if err == nil {
		return false
	}
	// TODO: Do something if needed.
	log.Println("endpoint error handler:", err)
	return true
}

func getLatestMessage(profile *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		obfuscatedAgentID, err := extractPlacement(&profile.ID.Placement, r)
		errorHandler(err, r)
		agentID, err := profiles.OperateData(&profile.ID.Operations, &obfuscatedAgentID)
		errorHandler(err, r)

		obfuscatedData, err := extractPlacement(&profile.Payload.Placement, r)
		errorHandler(err, r)
		payload, err := profiles.OperateData(&profile.Payload.Operations, &obfuscatedData)

		fmt.Println(agentID, "\n", payload)
	}

}

func respondToMessage(profile *profiles.ProfileRequest) (handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

func handlersFromProfile(profile *profiles.Profile) (handler *http.ServeMux, err error) {
	handler = http.NewServeMux()

	for _, method := range profile.Incall.Methods {
		for _, route := range profile.Incall.Routes {
			if !strings.HasPrefix(route, "/") {
				route = "/" + route
			}
			handler.HandleFunc(method+" "+route, getLatestMessage(&profile.Incall))
		}
	}

	for _, method := range profile.Outcall.Methods {
		for _, route := range profile.Outcall.Routes {
			if !strings.HasPrefix(route, "/") {
				route = "/" + route
			}
			handler.HandleFunc(method+" "+route, respondToMessage(&profile.Outcall))
		}
	}

	return handler, nil
}
