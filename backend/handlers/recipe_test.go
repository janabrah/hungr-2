package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

	req := httptest.NewRequest("POST", "/api/recipes?email="+testEmail+"&name=TestRecipe&tagString=dinner,%20quick", body)
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

	// Verify tags are returned correctly
	if len(response.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(response.Tags))
	}
	expectedTags := map[string]bool{"dinner": false, "quick": false}
	for _, tag := range response.Tags {
		if _, ok := expectedTags[tag.Name]; !ok {
			t.Errorf("Unexpected tag %q", tag.Name)
		}
		expectedTags[tag.Name] = true
		if tag.UUID == uuid.Nil {
			t.Errorf("Tag %q has nil UUID", tag.Name)
		}
	}
	for name, found := range expectedTags {
		if !found {
			t.Errorf("Expected tag %q not found", name)
		}
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
	recipe, err := storage.InsertRecipeByEmail("steps-invalid-ing-test", testEmail)
	if err != nil {
		t.Fatalf("Failed to create test recipe: %v", err)
	}
	defer storage.DeleteRecipe(recipe.UUID)

	body := `{"steps": [{"instruction": "Test", "ingredients": ["2"]}]}`
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
	recipe, err := storage.InsertRecipeByEmail("steps-valid-test", testEmail)
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

func TestCreateRecipe_TagStringRoundTrip(t *testing.T) {
	ensureTestUser(t)

	// Create a recipe with tags
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/recipes?email="+testEmail+"&name=TagRoundTripTest&tagString=breakfast,%20easy,%20cereal", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create failed: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var createResp models.UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Fetch recipes via GetRecipes and verify tag_string is computed correctly
	getReq := httptest.NewRequest("GET", "/api/recipes?email="+testEmail, nil)
	getW := httptest.NewRecorder()

	GetRecipes(getW, getReq)

	getResp := getW.Result()
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(getResp.Body)
		t.Fatalf("GetRecipes failed: status %d: %s", getResp.StatusCode, string(bodyBytes))
	}

	var recipesResp models.RecipesResponse
	if err := json.NewDecoder(getResp.Body).Decode(&recipesResp); err != nil {
		t.Fatalf("Failed to decode recipes response: %v", err)
	}

	// Find our recipe in the response
	var foundRecipe *models.Recipe
	for i := range recipesResp.RecipeData {
		if recipesResp.RecipeData[i].UUID == createResp.Recipe.UUID {
			foundRecipe = &recipesResp.RecipeData[i]
			break
		}
	}

	if foundRecipe == nil {
		t.Fatal("Created recipe not found in GetRecipes response")
	}

	// Verify tag_string matches the order tags were provided
	if foundRecipe.TagString != "breakfast, easy, cereal" {
		t.Errorf("Expected tag_string 'breakfast, easy', got %q", foundRecipe.TagString)
	}
}

func TestUpdateRecipeSteps_RoundTrip(t *testing.T) {
	ensureTestUser(t)

	// Create a recipe
	recipe, err := storage.InsertRecipeByEmail("steps-roundtrip-test", testEmail)
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

func TestPatchRecipe_MissingUUID(t *testing.T) {
	body := `{"tagString": "test"}`
	req := httptest.NewRequest("PATCH", "/api/recipes/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PatchRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestPatchRecipe_InvalidUUID(t *testing.T) {
	body := `{"tagString": "test"}`
	req := httptest.NewRequest("PATCH", "/api/recipes/not-a-uuid", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PatchRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestPatchRecipe_RecipeNotFound(t *testing.T) {
	body := `{"tagString": "test"}`
	req := httptest.NewRequest("PATCH", "/api/recipes/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PatchRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestPatchRecipe_InvalidBody(t *testing.T) {
	ensureTestUser(t)

	recipe, err := storage.InsertRecipeByEmail("patch-invalid-body-test", testEmail)
	if err != nil {
		t.Fatalf("Failed to create test recipe: %v", err)
	}
	defer storage.DeleteRecipe(recipe.UUID)

	body := `{invalid json`
	req := httptest.NewRequest("PATCH", "/api/recipes/"+recipe.UUID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PatchRecipe(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

// Helper to create a recipe with tags for patch tests
func createRecipeWithTags(t *testing.T, name, tagString string) models.UploadResponse {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	url := "/api/recipes?email=" + testEmail + "&name=" + name
	if tagString != "" {
		url += "&tagString=" + strings.ReplaceAll(tagString, " ", "%20")
	}
	createReq := httptest.NewRequest("POST", url, body)
	createReq.Header.Set("Content-Type", writer.FormDataContentType())
	createW := httptest.NewRecorder()

	CreateRecipe(createW, createReq)

	createResp := createW.Result()
	defer createResp.Body.Close()

	var response models.UploadResponse
	if err := json.NewDecoder(createResp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	return response
}

func TestPatchRecipe_SameTags(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with tags
	createResp := createRecipeWithTags(t, "PatchSameTagsTest", "alpha, beta, gamma")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with identical tags
	patchBody := `{"tagString": "alpha, beta, gamma"}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify tags unchanged
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "alpha, beta, gamma" {
		t.Errorf("Expected tag_string 'alpha, beta, gamma', got %q", updatedRecipe.TagString)
	}
}

func TestPatchRecipe_SubsetDifferentOrder(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with tags
	createResp := createRecipeWithTags(t, "PatchSubsetTest", "alpha, beta, gamma")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with subset in different order
	patchBody := `{"tagString": "gamma, alpha"}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify only subset tags in new order
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "gamma, alpha" {
		t.Errorf("Expected tag_string 'gamma, alpha', got %q", updatedRecipe.TagString)
	}
}

func TestPatchRecipe_Superset(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with tags
	createResp := createRecipeWithTags(t, "PatchSupersetTest", "alpha, beta")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with superset
	patchBody := `{"tagString": "alpha, beta, gamma, delta"}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify all tags present
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "alpha, beta, gamma, delta" {
		t.Errorf("Expected tag_string 'alpha, beta, gamma, delta', got %q", updatedRecipe.TagString)
	}
}

func TestPatchRecipe_MixedNewAndOld(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with tags
	createResp := createRecipeWithTags(t, "PatchMixedTest", "alpha, beta, gamma")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with mix of old and new tags
	patchBody := `{"tagString": "beta, delta, epsilon"}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify mixed tags
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "beta, delta, epsilon" {
		t.Errorf("Expected tag_string 'beta, delta, epsilon', got %q", updatedRecipe.TagString)
	}
}

func TestPatchRecipe_NewTagNotInTable(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with existing tag
	createResp := createRecipeWithTags(t, "PatchNewTagTest", "existing-tag")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with a completely new tag that doesn't exist in tags table
	patchBody := `{"tagString": "brand-new-unique-tag-12345"}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify new tag was created and assigned
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "brand-new-unique-tag-12345" {
		t.Errorf("Expected tag_string 'brand-new-unique-tag-12345', got %q", updatedRecipe.TagString)
	}
}

func TestPatchRecipe_ClearTags(t *testing.T) {
	ensureTestUser(t)

	// Create recipe with tags
	createResp := createRecipeWithTags(t, "PatchClearTagsTest", "breakfast, lunch")
	defer storage.DeleteRecipe(createResp.Recipe.UUID)

	// Patch with empty tag string to clear tags
	patchBody := `{"tagString": ""}`
	patchReq := httptest.NewRequest("PATCH", "/api/recipes/"+createResp.Recipe.UUID.String(), strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()

	PatchRecipe(patchW, patchReq)

	patchResp := patchW.Result()
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("Expected status 200, got %d: %s", patchResp.StatusCode, string(bodyBytes))
	}

	// Verify tags were cleared
	updatedRecipe, err := storage.GetRecipeByUUID(createResp.Recipe.UUID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	if updatedRecipe.TagString != "" {
		t.Errorf("Expected empty tag_string, got %q", updatedRecipe.TagString)
	}
}
