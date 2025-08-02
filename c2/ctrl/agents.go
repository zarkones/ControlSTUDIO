package ctrl

import (
	"c2/repos"
	"net/http"
)

func GetAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := repos.GetAgents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(agents) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jj(w, &agents)
}
