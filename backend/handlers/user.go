package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "email and name are required")
		return
	}

	// Check if user already exists
	existingUser, _ := storage.GetUserByEmail(req.Email)
	if existingUser != nil {
		respondWithError(w, http.StatusConflict, "user with this email already exists")
		return
	}

	user, err := storage.CreateUser(req.Email, req.Name)
	if err != nil {
		logger.Error(ctx, "failed to create user", err, "email", req.Email)
		respondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	logger.Info(ctx, "user created", "user_uuid", user.UUID, "email", req.Email)
	response := models.UserResponse{
		Success: true,
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUIDStr := r.URL.Query().Get("uuid")
	email := r.URL.Query().Get("email")

	var user *models.User
	var err error

	if userUUIDStr != "" {
		userUUID, parseErr := uuid.FromString(userUUIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid uuid")
			return
		}
		user, err = storage.GetUserByUUID(userUUID)
	} else if email != "" {
		user, err = storage.GetUserByEmail(email)
	} else {
		respondWithError(w, http.StatusBadRequest, "uuid or email is required")
		return
	}

	if err != nil {
		logger.Debug(ctx, "user not found", "uuid", userUUIDStr, "email", email)
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	response := models.UserResponse{
		Success: true,
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUIDStr := r.URL.Query().Get("uuid")
	if userUUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "uuid is required")
		return
	}

	userUUID, err := uuid.FromString(userUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	var updateReq struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if updateReq.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	user, err := storage.UpdateUser(userUUID, updateReq.Name)
	if err != nil {
		logger.Error(ctx, "failed to update user", err, "user_uuid", userUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	logger.Info(ctx, "user updated", "user_uuid", userUUID)
	response := models.UserResponse{
		Success: true,
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUIDStr := r.URL.Query().Get("uuid")
	if userUUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "uuid is required")
		return
	}

	userUUID, err := uuid.FromString(userUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	err = storage.DeleteUser(userUUID)
	if err != nil {
		logger.Error(ctx, "failed to delete user", err, "user_uuid", userUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	logger.Info(ctx, "user deleted", "user_uuid", userUUID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := storage.UpsertUserOnLogin(req.Email)
	if err != nil {
		logger.Error(ctx, "failed to login/upsert user", err, "email", req.Email)
		respondWithError(w, http.StatusInternalServerError, "login failed")
		return
	}

	logger.Info(ctx, "user logged in", "user_uuid", user.UUID, "email", req.Email)
	response := models.UserResponse{
		Success: true,
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
