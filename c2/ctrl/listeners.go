package ctrl

import (
	"c2/core/listeners"
	"c2/repos"
	"encoding/json"
	"net/http"
)

func GetListeners(w http.ResponseWriter, r *http.Request) {
	ls := listeners.AsSlice()
	if len(ls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jj(w, &ls)
}

func InsertListener(w http.ResponseWriter, r *http.Request) {
	var l listeners.Listener

	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(l.ProfileID) == 0 {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}
	if len(l.Address) == 0 {
		http.Error(w, "invalid profile address", http.StatusBadRequest)
		return
	}
	if len(l.Port) == 0 {
		http.Error(w, "invalid profile port", http.StatusBadRequest)
		return
	}

	metaProfile, err := repos.GetProfile(l.ProfileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := metaProfile.GetProfile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := listeners.Insert(l, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteListener(w http.ResponseWriter, r *http.Request) {
	ProfileID := r.PathValue("ProfileID")
	Address := r.PathValue("Address")
	Port := r.PathValue("Port")

	listeners.Delete(listeners.Listener{
		ProfileID: ProfileID,
		Address:   Address,
		Port:      Port,
	})
}
