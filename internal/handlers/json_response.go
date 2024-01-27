package handlers

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	out, _ := json.MarshalIndent(body, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}
