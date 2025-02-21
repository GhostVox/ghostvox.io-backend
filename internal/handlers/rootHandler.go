package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type RootHandler struct {
	queries *database.Queries
}

func NewRootHandler(queries *database.Queries) *RootHandler {
	return &RootHandler{
		queries: queries,
	}
}

func (db *RootHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {

	response := struct {
		Message string `json:"message"`
	}{
		Message: "welcome to ghostvox.io-backend",
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Could not encode response error:%s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(encoded)
	if err != nil {
		fmt.Printf("Could not write response error:%s", err)
	}
	return
}
