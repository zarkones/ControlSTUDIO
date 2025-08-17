package main

import (
	"c2/core/permissions"
	"c2/ctrl"
	"net/http"

	access "github.com/zarkones/ControlACCESS"
)

func initRouting(r *http.ServeMux) {
	access.Route(r, "GET", "/v1/agents", permissions.AGENTS_LIST, access.AuthMiddleware(ctrl.GetAgents))

	access.Route(r, "POST", "/v1/messages/by-ids", permissions.MESSAGES_LIST, access.AuthMiddleware(ctrl.GetMessageByIDs))
	access.Route(r, "GET", "/v1/messages/{agentID}", permissions.MESSAGES_LIST, access.AuthMiddleware(ctrl.GetMessages))
	access.Route(r, "PUT", "/v1/messages/{agentID}", permissions.MESSAGES_INSERT, access.AuthMiddleware(ctrl.InsertMessage))

	access.Route(r, "GET", "/v1/profiles", permissions.PROFILES_LIST, access.AuthMiddleware(ctrl.GetProfiles))
	access.Route(r, "PUT", "/v1/profiles", permissions.PROFILES_INSERT, access.AuthMiddleware(ctrl.InsertProfile))
	access.Route(r, "POST", "/v1/profiles", permissions.PROFILES_UPDATE, access.AuthMiddleware(ctrl.UpdateProfile))
	access.Route(r, "DELETE", "/v1/profiles", permissions.PROFILES_DELETE, access.AuthMiddleware(ctrl.DeleteProfile))

	access.Route(r, "GET", "/v1/listeners", permissions.LISTENERS_LIST, access.AuthMiddleware(ctrl.GetListeners))

	access.Route(r, "GET", "/v1/permissions/{operatorUsername}", permissions.PERMISSIONS_LIST, access.AuthMiddleware(ctrl.GetPermissions))
	access.Route(r, "PUT", "/v1/permissions", permissions.PERMISSIONS_INSERT, access.AuthMiddleware(ctrl.InsertPermission))
	access.Route(r, "DELETE", "/v1/permissions/{permissionID}", permissions.PERMISSIONS_DELETE, access.AuthMiddleware(ctrl.DeletePermission))

	access.Route(r, "GET", "/v1/operators", permissions.OPERATORS_LIST, access.AuthMiddleware(ctrl.GetOperators))

	// r.HandleFunc("PUT /v1/listeners", ctrl.InsertListener)
	// r.HandleFunc("DELETE /v1/listeners", ctrl.DeleteListener)
}
