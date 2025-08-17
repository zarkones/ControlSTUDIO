package ctrl

import (
	"c2/repos/operatorsRepo"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SerializedOperator struct {
	Username     string
	PublicKeyHex string
	CreatedAt    time.Time
}

func GetOperators(w http.ResponseWriter, r *http.Request) {
	operators, err := operatorsRepo.GetMultiple()
	if err != nil {
		log.Println("api: error: GetOperators:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(operators) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	serializableOperators := make([]SerializedOperator, len(operators))
	for i, operator := range operators {
		serializableOperators[i] = SerializedOperator{
			Username:     operator.Username,
			PublicKeyHex: operator.PublicKeyHex,
			CreatedAt:    operator.CreatedAt,
		}
	}

	// Just in case someone tries to use this variable which is sensitive,
	// due to it possesing encrypted private key, yes, I know it's encrypted,
	// but still...
	operators = nil

	if err := json.NewEncoder(w).Encode(&serializableOperators); err != nil {
		log.Println("api: error: serializing response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
