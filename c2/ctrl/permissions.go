package ctrl

import (
	"c2/core/permissions"
	"c2/repos/permissionsRepo"
	"encoding/json"
	"net/http"
	"strconv"

	access "github.com/zarkones/ControlACCESS"
)

type ExtendedPermission struct {
	access.Permission
	Acquired bool `json:"acquired"`
}

type GetPermissionsRespCtx map[access.PermissionKey]ExtendedPermission

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	userPermissions, err := permissionsRepo.GetMultipleByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userPermissions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	operatorPermissionMap := make(map[access.PermissionKey]access.Permission, len(userPermissions))
	for _, permission := range userPermissions {
		operatorPermissionMap[permission.Key] = permission
	}

	permissionMap := make(GetPermissionsRespCtx, len(permissions.AllPermissions))
	for _, permissionKey := range permissions.AllPermissions {
		if _, ok := operatorPermissionMap[permissionKey]; ok {
			permissionMap[permissionKey] = ExtendedPermission{
				Permission: operatorPermissionMap[permissionKey],
				Acquired:   true,
			}
			continue
		}
		permissionMap[permissionKey] = ExtendedPermission{
			Permission: access.Permission{
				Key:    permissionKey,
				UserID: userID,
			},
			Acquired: false,
		}
	}

	if err := json.NewEncoder(w).Encode(&permissionMap); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InsertPermission(w http.ResponseWriter, r *http.Request) {
	var newPermission access.Permission

	if err := json.NewDecoder(r.Body).Decode(&newPermission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newPermission.UserID) == 0 {
		http.Error(w, "username must be provided", http.StatusUnprocessableEntity)
		return
	}
	if newPermission.Key == access.PERMISSION_NOT_SPECIFIED {
		http.Error(w, "permission key cannot be PERMISSION_NOT_SPECIFIED", http.StatusUnprocessableEntity)
		return
	}

	if err := permissionsRepo.Insert(&newPermission); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeletePermission(w http.ResponseWriter, r *http.Request) {
	rawPermissionKey := r.PathValue("permissionKey")
	userID := r.PathValue("userID")

	permissionKey, err := strconv.Atoi(rawPermissionKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := permissionsRepo.Delete(userID, access.PermissionKey(permissionKey)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
