package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

func CreateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Authenticate via email
	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	authUser, err := storage.GetUserByEmail(email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	var req models.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TargetUserUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "target_user_uuid is required")
		return
	}

	// Use authenticated user's UUID as source (prevents spoofing)
	sourceUserUUID := authUser.UUID

	if sourceUserUUID == req.TargetUserUUID {
		respondWithError(w, http.StatusBadRequest, "cannot connect to yourself")
		return
	}

	// Check if target user exists
	_, err = storage.GetUserByUUID(req.TargetUserUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "target user not found")
		return
	}

	// Check if connection already exists
	exists, err := storage.ConnectionExists(sourceUserUUID, req.TargetUserUUID)
	if err != nil {
		logger.Error(ctx, "failed to check connection", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create connection")
		return
	}
	if exists {
		respondWithError(w, http.StatusConflict, "connection already exists")
		return
	}

	err = storage.CreateConnection(sourceUserUUID, req.TargetUserUUID)
	if err != nil {
		logger.Error(ctx, "failed to create connection", err,
			"source_user_uuid", sourceUserUUID, "target_user_uuid", req.TargetUserUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to create connection")
		return
	}

	logger.Info(ctx, "connection created",
		"source_user_uuid", sourceUserUUID,
		"target_user_uuid", req.TargetUserUUID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userUUIDStr := r.URL.Query().Get("user_uuid")
	if userUUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "user_uuid is required")
		return
	}

	userUUID, err := uuid.FromString(userUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_uuid")
		return
	}

	direction := r.URL.Query().Get("direction")
	if direction == "" {
		direction = "outgoing"
	}

	var users []models.User
	if direction == "incoming" {
		users, err = storage.GetConnectionsByTargetUser(userUUID)
	} else {
		users, err = storage.GetConnectionsBySourceUser(userUUID)
	}

	if err != nil {
		logger.Error(ctx, "failed to get connections", err, "user_uuid", userUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get connections")
		return
	}

	if users == nil {
		users = []models.User{}
	}

	response := models.ConnectionsResponse{
		Success:     true,
		Connections: users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Authenticate via email
	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	authUser, err := storage.GetUserByEmail(email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	targetUserUUIDStr := r.URL.Query().Get("target_user_uuid")
	if targetUserUUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "target_user_uuid is required")
		return
	}

	targetUserUUID, err := uuid.FromString(targetUserUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid target_user_uuid")
		return
	}

	// Use authenticated user's UUID as source (prevents spoofing)
	sourceUserUUID := authUser.UUID

	bidirectional := false
	if bidirectionalStr := r.URL.Query().Get("bidirectional"); bidirectionalStr != "" {
		parsed, parseErr := strconv.ParseBool(bidirectionalStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid bidirectional flag")
			return
		}
		bidirectional = parsed
	}

	if bidirectional {
		err = storage.DeleteConnectionsBidirectional(sourceUserUUID, targetUserUUID)
	} else {
		err = storage.DeleteConnection(sourceUserUUID, targetUserUUID)
	}
	if err != nil {
		logger.Error(ctx, "failed to delete connection", err,
			"source_user_uuid", sourceUserUUID, "target_user_uuid", targetUserUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to delete connection")
		return
	}

	logger.Info(ctx, "connection deleted",
		"source_user_uuid", sourceUserUUID, "target_user_uuid", targetUserUUID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
