package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

// Silence unused import error
var _ = uuid.Nil

func TestGetRecipes_ReturnsJSON(t *testing.T) {
	t.Skip("Enable when GetRecipes is implemented")

	testUserUUID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest("GET", "/api/recipes?user_uuid="+testUserUUID.String(), nil)
	w := httptest.NewRecorder()

	GetRecipes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %q", contentType)
	}

	var response models.RecipesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Response should have all three fields (even if empty)
	if response.RecipeData == nil {
		t.Error("Expected recipeData to be non-nil")
	}
	if response.FileData == nil {
		t.Error("Expected fileData to be non-nil")
	}
	if response.MappingData == nil {
		t.Error("Expected mappingData to be non-nil")
	}
}

func TestGetRecipes_MissingUserUUID(t *testing.T) {
	t.Skip("Enable when GetRecipes is implemented")

	req := httptest.NewRequest("GET", "/api/recipes", nil)
	w := httptest.NewRecorder()

	GetRecipes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Should return error for missing user_uuid
	if resp.StatusCode == http.StatusOK {
		t.Error("Expected error status for missing user_uuid")
	}
}

func TestCreateRecipe_MissingFilename(t *testing.T) {
	t.Skip("Enable when CreateRecipe is implemented")

	// Create empty multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?tagString=test", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Should return error for missing filename
	if resp.StatusCode == http.StatusOK {
		t.Error("Expected error status for missing filename")
	}
}

func TestCreateRecipe_WithFiles(t *testing.T) {
	t.Skip("Enable when CreateRecipe is implemented and storage is mocked")

	// Create multipart form with a test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add a fake file
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake image data"))

	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?filename=TestRecipe&tagString=dinner,quick", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Recipe.Filename != "TestRecipe" {
		t.Errorf("Expected filename 'TestRecipe', got %q", response.Recipe.Filename)
	}
	if response.Recipe.UUID == uuid.Nil {
		t.Error("Expected non-nil UUID for recipe")
	}
}

func TestCreateRecipe_MultipleFiles(t *testing.T) {
	t.Skip("Enable when CreateRecipe is implemented and storage is mocked")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add multiple files
	for i := 0; i < 3; i++ {
		part, _ := writer.CreateFormFile("file", "test.jpg")
		part.Write([]byte("fake image data"))
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?filename=MultiPageRecipe&tagString=cookbook", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// TODO: Verify that 3 files were created and linked
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, http.StatusBadRequest, "test error")

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}
