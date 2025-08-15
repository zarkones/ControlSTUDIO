package ctrl

import (
	"c2/core/listeners"
	"c2/models"
	"c2/repos"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := repos.GetProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(profiles) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jj(w, profiles)
}

func InsertProfile(w http.ResponseWriter, r *http.Request) {
	var profileWithMeta models.MetaProfile

	if err := json.NewDecoder(r.Body).Decode(&profileWithMeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := profileWithMeta.GetProfile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(profileWithMeta.Name) == 0 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	if _, err := profileWithMeta.GetProfile(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repos.InsertProfile(&profileWithMeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		for _, listenerService := range profile.GetHosts() {
			if err := listeners.Insert(
				listeners.Listener{
					ProfileID: profileWithMeta.ID,
					Address:   listenerService.Address,
					Port:      listenerService.Port,
				},
				&profile,
			); err != nil {
				fmt.Println(
					"failed to start a listener service",
					profileWithMeta.ID,
					listenerService.Address,
					listenerService.Port,
					"error:",
					err,
				)
				continue
			}
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profileWithMeta models.MetaProfile

	if err := json.NewDecoder(r.Body).Decode(&profileWithMeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(profileWithMeta.Name) == 0 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	if _, err := profileWithMeta.GetProfile(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repos.UpsertProfile(&profileWithMeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileID")

	if err := repos.DeleteProfile(profileID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
