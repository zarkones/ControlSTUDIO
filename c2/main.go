package main

import (
	"encoding/hex"
	"net/http"

	"github.com/zarkones/netescape"
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /something", func(w http.ResponseWriter, r *http.Request) {
		payload := "ls && uname -a"
		h := hex.EncodeToString([]byte(payload))
		csv, _ := netescape.ToCsv(&h)
		w.Write([]byte(csv))
	})

	http.ListenAndServe("0.0.0.0:8000", r)
}
