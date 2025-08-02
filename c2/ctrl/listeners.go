package ctrl

import (
	"c2/core/listeners"
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

	listeners.Insert(l)

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
