package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

// Silence unused import error
var _ = uuid.Nil

func init() {
	logger.Init()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL must be set to run tests")
	}
	if err := storage.Init(dbURL); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

const testEmail = "test@example.com"

func ensureTestUser(t *testing.T) {
	_, err := storage.GetUserByEmail(testEmail)
	if err != nil {
		_, err = storage.CreateUser(testEmail, "Test User")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
}

func TestGetRecipes_ReturnsJSON(t *testing.T) {
	ensureTestUser(t)

	req := httptest.NewRequest("GET", "/api/recipes?email="+testEmail, nil)
	w := httptest.NewRecorder()

	GetRecipes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %q", contentType)
	}

	var response models.RecipesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.RecipeData == nil {
		t.Error("Expected recipeData to be non-nil")
	}
	if response.FileData == nil {
		t.Error("Expected fileData to be non-nil")
	}
}

func TestGetRecipes_MissingEmail(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/recipes", nil)
	w := httptest.NewRecorder()

	GetRecipes(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateRecipe_MissingName(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?email=test@example.com&tagString=test", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateRecipe_MissingEmail(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?name=Test&tagString=test", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateRecipe_WithFiles(t *testing.T) {
	ensureTestUser(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake image data"))

	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?email="+testEmail+"&name=TestRecipe&tagString=dinner,quick", body)
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
	if response.Recipe.Name != "TestRecipe" {
		t.Errorf("Expected name 'TestRecipe', got %q", response.Recipe.Name)
	}
	if response.Recipe.UUID == uuid.Nil {
		t.Error("Expected non-nil UUID for recipe")
	}

	// Cleanup
	storage.DeleteRecipe(response.Recipe.UUID)
}

func TestCreateRecipe_MultipleFiles(t *testing.T) {
	ensureTestUser(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i := 0; i < 3; i++ {
		part, _ := writer.CreateFormFile("file", "test.jpg")
		part.Write([]byte("fake image data"))
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?email="+testEmail+"&name=MultiPageRecipe&tagString=cookbook", body)
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

	// Cleanup
	storage.DeleteRecipe(response.Recipe.UUID)
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

func TestGetRecipeSteps_MissingUUID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/recipes//steps", nil)
	w := httptest.NewRecorder()

	GetRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetRecipeSteps_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/recipes/not-a-uuid/steps", nil)
	w := httptest.NewRecorder()

	GetRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateRecipeSteps_MissingUUID(t *testing.T) {
	body := `{"steps": []}`
	req := httptest.NewRequest("PUT", "/api/recipes//steps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateRecipeSteps_InvalidUUID(t *testing.T) {
	body := `{"steps": []}`
	req := httptest.NewRequest("PUT", "/api/recipes/not-a-uuid/steps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateRecipeSteps_RecipeNotFound(t *testing.T) {

	body := `{"steps": []}`
	req := httptest.NewRequest("PUT", "/api/recipes/00000000-0000-0000-0000-000000000001/steps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestUpdateRecipeSteps_InvalidIngredient(t *testing.T) {
	ensureTestUser(t)

	// Create a recipe first
	recipe, err := storage.InsertRecipeByEmail("steps-invalid-ing-test", testEmail, "test")
	if err != nil {
		t.Fatalf("Failed to create test recipe: %v", err)
	}
	defer storage.DeleteRecipe(recipe.UUID)

	body := `{"steps": [{"instruction": "Test", "ingredients": ["invalid"]}]}`
	req := httptest.NewRequest("PUT", "/api/recipes/"+recipe.UUID.String()+"/steps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, string(bodyBytes))
	}
}

func TestUpdateRecipeSteps_ValidRequest(t *testing.T) {
	ensureTestUser(t)

	// Create a recipe first
	recipe, err := storage.InsertRecipeByEmail("steps-valid-test", testEmail, "test")
	if err != nil {
		t.Fatalf("Failed to create test recipe: %v", err)
	}
	defer storage.DeleteRecipe(recipe.UUID)

	body := `{
		"steps": [
			{
				"instruction": "Mix dry ingredients",
				"ingredients": ["2 cups flour", "1 tsp salt"]
			},
			{
				"instruction": "Add wet ingredients",
				"ingredients": ["1 cup milk", "2 eggs"]
			}
		]
	}`
	req := httptest.NewRequest("PUT", "/api/recipes/"+recipe.UUID.String()+"/steps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateRecipeSteps(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response["success"] {
		t.Error("Expected success to be true")
	}
}

func TestUpdateRecipeSteps_RoundTrip(t *testing.T) {
	ensureTestUser(t)

	// Create a recipe
	recipe, err := storage.InsertRecipeByEmail("steps-roundtrip-test", testEmail, "test")
	if err != nil {
		t.Fatalf("Failed to create test recipe: %v", err)
	}
	defer storage.DeleteRecipe(recipe.UUID)

	// PUT some steps
	putBody := `{
		"steps": [
			{
				"instruction": "Preheat oven to 350F",
				"ingredients": []
			},
			{
				"instruction": "Mix dry ingredients",
				"ingredients": ["2 cups flour", "1 tsp baking powder", "1/2 tsp salt"]
			},
			{
				"instruction": "Add wet ingredients and stir",
				"ingredients": ["1 cup milk", "2 eggs", "3 tbsp butter"]
			}
		]
	}`
	putReq := httptest.NewRequest("PUT", "/api/recipes/"+recipe.UUID.String()+"/steps", bytes.NewBufferString(putBody))
	putReq.Header.Set("Content-Type", "application/json")
	putW := httptest.NewRecorder()

	UpdateRecipeSteps(putW, putReq)

	putResp := putW.Result()
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(putResp.Body)
		t.Fatalf("PUT failed: status %d: %s", putResp.StatusCode, string(bodyBytes))
	}

	// GET the steps back
	getReq := httptest.NewRequest("GET", "/api/recipes/"+recipe.UUID.String()+"/steps", nil)
	getW := httptest.NewRecorder()

	GetRecipeSteps(getW, getReq)

	getResp := getW.Result()
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(getResp.Body)
		t.Fatalf("GET failed: status %d: %s", getResp.StatusCode, string(bodyBytes))
	}

	var stepsResp models.RecipeStepsResponse
	if err := json.NewDecoder(getResp.Body).Decode(&stepsResp); err != nil {
		t.Fatalf("Failed to decode GET response: %v", err)
	}

	// Verify we got 3 steps
	if len(stepsResp.Steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(stepsResp.Steps))
	}

	// Check first step (no ingredients)
	if stepsResp.Steps[0].Instruction != "Preheat oven to 350F" {
		t.Errorf("Step 1 instruction: got %q", stepsResp.Steps[0].Instruction)
	}
	if len(stepsResp.Steps[0].Ingredients) != 0 {
		t.Errorf("Step 1 should have 0 ingredients, got %d", len(stepsResp.Steps[0].Ingredients))
	}

	// Check second step
	if stepsResp.Steps[1].Instruction != "Mix dry ingredients" {
		t.Errorf("Step 2 instruction: got %q", stepsResp.Steps[1].Instruction)
	}
	if len(stepsResp.Steps[1].Ingredients) != 3 {
		t.Errorf("Step 2 should have 3 ingredients, got %d", len(stepsResp.Steps[1].Ingredients))
	}

	// Check third step
	if stepsResp.Steps[2].Instruction != "Add wet ingredients and stir" {
		t.Errorf("Step 3 instruction: got %q", stepsResp.Steps[2].Instruction)
	}
	if len(stepsResp.Steps[2].Ingredients) != 3 {
		t.Errorf("Step 3 should have 3 ingredients, got %d", len(stepsResp.Steps[2].Ingredients))
	}
}
