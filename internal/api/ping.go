package api

import (
	"encoding/json"
	"net/http"
)

func (h Handler) Ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	json.NewEncoder(w).Encode("pong")
}
