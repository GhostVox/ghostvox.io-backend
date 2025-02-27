package handlers

import (
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
)

type RootHandler struct {
	cfg *config.APIConfig
}

func NewRootHandler(cfg *config.APIConfig) *RootHandler {
	return &RootHandler{
		cfg: cfg,
	}
}

func (rh *RootHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {

	type route struct {
		Path              string `json:"path"`
		Method            string `json:"method"`
		Description       string `json:"description"`
		User_token_Needed bool   `json:"User_token_Needed"`
		Public            bool   `json:"public"`
	}
	type welcome struct {
		Message string  `json:"message"`
		Routes  []route `json:"routes"`
	}

	respondWithJSON(w, http.StatusOK, welcome{
		Message: "Welcome to GhostVox.io-backend",
		Routes: []route{
			{
				Path:              "/api/v1/users",
				Method:            "GET",
				Description:       "Get all users",
				User_token_Needed: false,
				Public:            true,
			},
			{
				Path:              "/api/v1/users/:id",
				Method:            "GET",
				Description:       "Get a user by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/users",
				Method:            "POST",
				Description:       "Create a new user",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/users/:id",
				Method:            "PUT",
				Description:       "Update a user by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/users/:id",
				Method:            "DELETE",
				Description:       "Delete a user by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll",
				Method:            "GET",
				Description:       "Get all polls",
				User_token_Needed: false,
				Public:            true,
			},
			{
				Path:              "/api/v1/poll/:id",
				Method:            "GET",
				Description:       "Get a poll by ID",
				User_token_Needed: false,
				Public:            true,
			},
			{
				Path:              "/api/v1/poll",
				Method:            "POST",
				Description:       "Create a new poll",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id",
				Method:            "PUT",
				Description:       "Update a poll by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id",
				Method:            "DELETE",
				Description:       "Delete a poll by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id/vote",
				Method:            "POST",
				Description:       "Vote on a poll",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id/vote",
				Method:            "DELETE",
				Description:       "Unvote on a poll",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id/options",
				Method:            "GET",
				Description:       "Get all options",
				User_token_Needed: false,
				Public:            true,
			},
			{
				Path:              "/api/v1/poll/:id/options/:option_id",
				Method:            "GET",
				Description:       "Get an option by ID",
				User_token_Needed: false,
				Public:            true,
			},
			{
				Path:              "/api/v1/poll/:id/options",
				Method:            "POST",
				Description:       "Create a new option",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id/options/:option_id",
				Method:            "PUT",
				Description:       "Update an option by ID",
				User_token_Needed: true,
				Public:            false,
			},
			{
				Path:              "/api/v1/poll/:id/options/:option_id",
				Method:            "DELETE",
				Description:       "Delete an option by ID",
				User_token_Needed: true,
				Public:            false,
			},
		},
	})
	return
}
