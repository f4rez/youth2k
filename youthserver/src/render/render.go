package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Render struct {
}

func (r *Render) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (r *Render) RespondWithError(w http.ResponseWriter, code int, message string) {
	r.RespondWithJSON(w, code, map[string]string{"error": message})
}

func (r *Render) RespondOK(w http.ResponseWriter) {
	fmt.Fprint(w, http.StatusOK, "OK")
}
