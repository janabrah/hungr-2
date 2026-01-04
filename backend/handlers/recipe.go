package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {
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

	recipes, err := storage.GetRecipesByUserUUID(userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	recipeUUIDs := make([]uuid.UUID, len(recipes))
	for i, r := range recipes {
		recipeUUIDs[i] = r.UUID
	}

	files, err := storage.GetFilesByRecipeUUIDs(recipeUUIDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
	name := r.URL.Query().Get("name")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	tagString := r.URL.Query().Get("tagString")

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

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "at least one file is required")
		return
	}

	recipe, err := storage.InsertRecipe(name, userUUID, tagString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe: "+err.Error())
		return
	}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to open file: "+err.Error())
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to read file: "+err.Error())
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		_, err = storage.InsertFile(recipe.UUID, data, contentType, i, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to store file: "+err.Error())
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
				respondWithError(w, http.StatusInternalServerError, "failed to create tag: "+err.Error())
				return
			}
			insertedTags = append(insertedTags, *tag)

			if err := storage.InsertRecipeTag(recipe.UUID, tagUUID); err != nil {
				respondWithError(w, http.StatusInternalServerError, "failed to link tag to recipe: "+err.Error())
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

func GetFile(w http.ResponseWriter, r *http.Request) {
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
