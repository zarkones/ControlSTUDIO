package main

import (
	"c2/ctrl"
	"net/http"
)

func initRouting(r *http.ServeMux) {
	r.HandleFunc("GET /v1/agents", ctrl.GetAgents)

	r.HandleFunc("GET /v1/agents/{agentID}", ctrl.GetMessages)
	r.HandleFunc("PUT /v1/agents/{agentID}", ctrl.InsertMessage)

	r.HandleFunc("GET /v1/profiles", ctrl.GetProfiles)
	r.HandleFunc("PUT /v1/profiles", ctrl.InsertProfile)
	r.HandleFunc("POST /v1/profiles", ctrl.UpdateProfile)
	r.HandleFunc("DELETE /v1/profiles", ctrl.DeleteProfile)

	r.HandleFunc("GET /v1/listeners", ctrl.GetListeners)
	r.HandleFunc("PUT /v1/listeners", ctrl.InsertListener)
	r.HandleFunc("DELETE /v1/listeners", ctrl.DeleteListener)
}
