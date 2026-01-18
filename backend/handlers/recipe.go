package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/cobyabrahams/hungr/units"
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

	// Read all file data before starting transaction
	type fileData struct {
		data        []byte
		contentType string
	}
	filesData := make([]fileData, 0, len(files))
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error(ctx, "failed to open uploaded file", err, "email", email, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to process uploaded file")
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			logger.Error(ctx, "failed to read uploaded file", err, "email", email, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to read uploaded file")
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}
		filesData = append(filesData, fileData{data: data, contentType: contentType})
	}

	// Start transaction
	tx, err := storage.BeginTx(ctx)
	if err != nil {
		logger.Error(ctx, "failed to begin transaction", err, "email", email)
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe")
		return
	}
	defer tx.Rollback(ctx)

	recipe, err := storage.TxInsertRecipeByEmail(ctx, tx, name, email)
	if err != nil {
		logger.Error(ctx, "failed to insert recipe", err, "email", email, "name", name)
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe - user may not exist")
		return
	}
	recipe.TagString = tagString

	logger.Info(ctx, "recipe created", "recipe_uuid", recipe.UUID, "email", email)

	for i, fd := range filesData {
		_, err = storage.TxInsertFile(ctx, tx, recipe.UUID, fd.data, fd.contentType, i, true)
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
			tag, err := storage.TxUpsertTag(ctx, tx, tagUUID, tagName)
			if err != nil {
				logger.Error(ctx, "failed to upsert tag", err, "recipe_uuid", recipe.UUID, "tag", tagName)
				respondWithError(w, http.StatusInternalServerError, "failed to create tag")
				return
			}
			insertedTags = append(insertedTags, *tag)

			if err := storage.TxInsertRecipeTag(ctx, tx, recipe.UUID, tagUUID); err != nil {
				logger.Error(ctx, "failed to link tag to recipe", err, "recipe_uuid", recipe.UUID, "tag_uuid", tagUUID)
				respondWithError(w, http.StatusInternalServerError, "failed to link tag")
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		logger.Error(ctx, "failed to commit transaction", err, "recipe_uuid", recipe.UUID)
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe")
		return
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

func AddRecipeFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse recipe UUID from path: /api/recipes/{uuid}/files
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/recipes/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		respondWithError(w, http.StatusBadRequest, "recipe uuid is required")
		return
	}

	recipeUUID, err := uuid.FromString(parts[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid recipe uuid")
		return
	}

	// Check if recipe exists
	_, err = storage.GetRecipeByUUID(recipeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "recipe not found")
			return
		}
		logger.Error(ctx, "failed to get recipe", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Error(ctx, "failed to parse multipart form", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusBadRequest, "failed to parse upload")
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "at least one file is required")
		return
	}

	type fileData struct {
		data        []byte
		contentType string
	}

	filesData := make([]fileData, 0, len(files))
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error(ctx, "failed to open uploaded file", err, "recipe_uuid", recipeUUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to process uploaded file")
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			logger.Error(ctx, "failed to read uploaded file", err, "recipe_uuid", recipeUUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to read uploaded file")
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		filesData = append(filesData, fileData{data: data, contentType: contentType})
	}

	existingFiles, err := storage.GetFilesByRecipeUUIDs([]uuid.UUID{recipeUUID})
	if err != nil {
		logger.Error(ctx, "failed to load recipe files", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to load recipe files")
		return
	}

	maxPage := -1
	for _, f := range existingFiles {
		if f.PageNumber > maxPage {
			maxPage = f.PageNumber
		}
	}

	tx, err := storage.BeginTx(ctx)
	if err != nil {
		logger.Error(ctx, "failed to begin transaction", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to store files")
		return
	}
	defer tx.Rollback(ctx)

	insertedFiles := make([]models.File, 0, len(filesData))
	for i, fd := range filesData {
		file, err := storage.TxInsertFile(ctx, tx, recipeUUID, fd.data, fd.contentType, maxPage+i+1, true)
		if err != nil {
			logger.Error(ctx, "failed to store file", err, "recipe_uuid", recipeUUID, "file_index", i)
			respondWithError(w, http.StatusInternalServerError, "failed to store file")
			return
		}
		insertedFiles = append(insertedFiles, *file)
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error(ctx, "failed to commit transaction", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to store files")
		return
	}

	response := models.FileUploadResponse{
		Success: true,
		Files:   insertedFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func GetTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := storage.GetAllTags()
	if err != nil {
		logger.Error(ctx, "failed to get tags", err)
		respondWithError(w, http.StatusInternalServerError, "failed to load tags")
		return
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	response := models.TagsResponse{
		Tags: tags,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetRecipeSteps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse recipe UUID from path: /api/recipes/{uuid}/steps
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/recipes/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		respondWithError(w, http.StatusBadRequest, "recipe uuid is required")
		return
	}

	recipeUUID, err := uuid.FromString(parts[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid recipe uuid")
		return
	}

	// Check if recipe exists
	_, err = storage.GetRecipeByUUID(recipeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "recipe not found")
			return
		}
		logger.Error(ctx, "failed to get recipe", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	// Get steps with ingredients
	stepsWithIngredients, err := storage.GetRecipeStepsWithIngredients(recipeUUID)
	if err != nil {
		logger.Error(ctx, "failed to get recipe steps", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe steps")
		return
	}

	// Build response
	response := models.RecipeStepsResponse{
		Steps: make([]models.RecipeStepResponse, len(stepsWithIngredients)),
	}

	for i, step := range stepsWithIngredients {
		ingredients := make([]string, len(step.Ingredients))
		for j, ing := range step.Ingredients {
			category := units.GetCategoryForIngredientUnit(ing.IngredientType)
			formatted := units.FormatBest(ing.Quantity, category)
			ingredients[j] = fmt.Sprintf("%s %s", formatted, ing.IngredientName)
		}

		response.Steps[i] = models.RecipeStepResponse{
			Instruction: step.Instructions,
			Ingredients: ingredients,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateRecipeSteps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse recipe UUID from path: /api/recipes/{uuid}/steps
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/recipes/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		respondWithError(w, http.StatusBadRequest, "recipe uuid is required")
		return
	}

	recipeUUID, err := uuid.FromString(parts[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid recipe uuid")
		return
	}

	// Check if recipe exists
	_, err = storage.GetRecipeByUUID(recipeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "recipe not found")
			return
		}
		logger.Error(ctx, "failed to get recipe", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	// Parse request body
	var request models.RecipeStepsResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Convert to storage input format
	steps := make([]storage.StepInput, len(request.Steps))
	for i, step := range request.Steps {
		ingredients := make([]storage.IngredientInput, len(step.Ingredients))
		for j, ingStr := range step.Ingredients {
			parsed, err := units.ParseIngredientString(ingStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid ingredient %q: %v", ingStr, err))
				return
			}
			ingredients[j] = storage.IngredientInput{
				Name:     parsed.IngredientName,
				Unit:     parsed.Unit,
				Quantity: parsed.Quantity,
			}
		}
		steps[i] = storage.StepInput{
			Instruction: step.Instruction,
			Ingredients: ingredients,
		}
	}

	// Replace all steps
	if err := storage.ReplaceRecipeSteps(recipeUUID, steps); err != nil {
		logger.Error(ctx, "failed to update recipe steps", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to update recipe steps")
		return
	}

	logger.Info(ctx, "recipe steps updated", "recipe_uuid", recipeUUID, "step_count", len(steps))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func PatchRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse recipe UUID from path: /api/recipes/{uuid}
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/recipes/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		respondWithError(w, http.StatusBadRequest, "recipe uuid is required")
		return
	}

	recipeUUID, err := uuid.FromString(parts[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid recipe uuid")
		return
	}

	// Check if recipe exists
	_, err = storage.GetRecipeByUUID(recipeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "recipe not found")
			return
		}
		logger.Error(ctx, "failed to get recipe", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	// Parse request body
	var request models.PatchRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Start transaction
	tx, err := storage.BeginTx(ctx)
	if err != nil {
		logger.Error(ctx, "failed to begin transaction", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to update recipe")
		return
	}
	defer tx.Rollback(ctx)

	// Delete existing recipe tags
	if err := storage.TxDeleteRecipeTags(ctx, tx, recipeUUID); err != nil {
		logger.Error(ctx, "failed to delete existing tags", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to update tags")
		return
	}

	// Insert new tags
	if request.TagString != "" {
		tags := strings.Split(request.TagString, ", ")
		for _, tagName := range tags {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			tagUUID := storage.CreateTagUUID(tagName)
			_, err := storage.TxUpsertTag(ctx, tx, tagUUID, tagName)
			if err != nil {
				logger.Error(ctx, "failed to upsert tag", err, "recipe_uuid", recipeUUID, "tag", tagName)
				respondWithError(w, http.StatusInternalServerError, "failed to create tag")
				return
			}

			if err := storage.TxInsertRecipeTag(ctx, tx, recipeUUID, tagUUID); err != nil {
				logger.Error(ctx, "failed to link tag to recipe", err, "recipe_uuid", recipeUUID, "tag_uuid", tagUUID)
				respondWithError(w, http.StatusInternalServerError, "failed to link tag")
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		logger.Error(ctx, "failed to commit transaction", err, "recipe_uuid", recipeUUID)
		respondWithError(w, http.StatusInternalServerError, "failed to update recipe")
		return
	}

	logger.Info(ctx, "recipe tags updated", "recipe_uuid", recipeUUID, "tag_string", request.TagString)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
