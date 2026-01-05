package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	recipes, err := storage.GetRecipesByUserEmail(email)
	if err != nil {
		logger.Error(ctx, "failed to get recipes", err, "email", email)
		respondWithError(w, http.StatusInternalServerError, "failed to load recipes")
		return
	}

	recipeUUIDs := make([]uuid.UUID, len(recipes))
	for i, r := range recipes {
		recipeUUIDs[i] = r.UUID
	}

	files, err := storage.GetFilesByRecipeUUIDs(recipeUUIDs)
	if err != nil {
		logger.Error(ctx, "failed to get files for recipes", err, "recipe_count", len(recipeUUIDs))
		respondWithError(w, http.StatusInternalServerError, "failed to load recipe files")
		return
	}

	response := models.RecipesResponse{
		RecipeData: recipes,
		FileData:   files,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	tagString := r.URL.Query().Get("tagString")

	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Error(ctx, "failed to parse multipart form", err, "email", email)
		respondWithError(w, http.StatusBadRequest, "failed to parse upload")
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "at least one file is required")
		return
	}

	recipe, err := storage.InsertRecipeByEmail(name, email, tagString)
	if err != nil {
		logger.Error(ctx, "failed to insert recipe", err, "email", email, "name", name)
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe - user may not exist")
		return
	}

	logger.Info(ctx, "recipe created", "recipe_uuid", recipe.UUID, "email", email)

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error(ctx, "failed to open uploaded file", err, "recipe_uuid", recipe.UUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to process uploaded file")
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			logger.Error(ctx, "failed to read uploaded file", err, "recipe_uuid", recipe.UUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to read uploaded file")
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		_, err = storage.InsertFile(recipe.UUID, data, contentType, i, true)
		if err != nil {
			logger.Error(ctx, "failed to store file", err, "recipe_uuid", recipe.UUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to store file")
			return
		}
	}

	var insertedTags []models.Tag
	if tagString != "" {
		tags := strings.Split(tagString, ", ")
		for _, tagName := range tags {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			tagUUID := storage.CreateTagUUID(tagName)
			tag, err := storage.UpsertTag(tagUUID, tagName)
			if err != nil {
				logger.Error(ctx, "failed to upsert tag", err, "recipe_uuid", recipe.UUID, "tag", tagName)
				respondWithError(w, http.StatusInternalServerError, "failed to create tag")
				return
			}
			insertedTags = append(insertedTags, *tag)

			if err := storage.InsertRecipeTag(recipe.UUID, tagUUID); err != nil {
				logger.Error(ctx, "failed to link tag to recipe", err, "recipe_uuid", recipe.UUID, "tag_uuid", tagUUID)
				respondWithError(w, http.StatusInternalServerError, "failed to link tag")
				return
			}
		}
	}

	response := models.UploadResponse{
		Success: true,
		Recipe:  *recipe,
		Tags:    insertedTags,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recipeUUIDStr := r.URL.Query().Get("uuid")
	if recipeUUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "uuid is required")
		return
	}

	recipeUUID, err := uuid.FromString(recipeUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	err = storage.DeleteRecipe(recipeUUID)
	if err != nil {
		logger.Error(ctx, "failed to delete recipe", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to delete recipe")
		return
	}

	logger.Info(ctx, "recipe deleted", "recipe_uuid", recipeUUID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		respondWithError(w, http.StatusBadRequest, "invalid file path")
		return
	}
	fileUUIDStr := parts[len(parts)-1]

	fileUUID, err := uuid.FromString(fileUUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid file uuid")
		return
	}

	data, contentType, err := storage.GetFileData(fileUUID)
	if err != nil {
		logger.Error(ctx, "failed to get file data", err, "file_uuid", fileUUID)
		respondWithError(w, http.StatusNotFound, "file not found")
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
