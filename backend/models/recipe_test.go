package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

func TestRecipeJSON(t *testing.T) {
	recipe := Recipe{
		UUID:      uuid.Must(uuid.NewV4()),
		Filename:  "test-recipe",
		User:      uuid.Must(uuid.NewV4()),
		TagString: "dinner, quick",
		CreatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify JSON field names match what frontend expects
	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"uuid", "filename", "user_uuid", "tag_string", "created_at"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}

	// Verify "uuid" is a string (UUID format)
	if _, ok := m["uuid"].(string); !ok {
		t.Error("Expected 'uuid' to be a string (UUID)")
	}
}

func TestFileJSON(t *testing.T) {
	file := File{
		UUID:  uuid.Must(uuid.NewV4()),
		URL:   "https://example.com/image.jpg",
		Image: true,
	}

	data, err := json.Marshal(file)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"uuid", "url", "image"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}
}

func TestFileRecipeJSON(t *testing.T) {
	fileRecipe := FileRecipe{
		FileUUID:   uuid.Must(uuid.NewV4()),
		RecipeUUID: uuid.Must(uuid.NewV4()),
		PageNumber: 1,
	}

	data, err := json.Marshal(fileRecipe)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"file_uuid", "recipe_uuid", "page_number"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}
}

func TestRecipesResponseJSON(t *testing.T) {
	response := RecipesResponse{
		RecipeData:  []Recipe{},
		FileData:    []File{},
		MappingData: []FileRecipe{},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	// These field names MUST match what the frontend expects
	expectedFields := []string{"recipeData", "fileData", "mappingData"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output (frontend depends on this)", field)
		}
	}
}

func TestUploadResponseJSON(t *testing.T) {
	response := UploadResponse{
		Success: true,
		Recipe:  Recipe{UUID: uuid.Must(uuid.NewV4()), Filename: "test"},
		Tags:    []Tag{{UUID: uuid.Must(uuid.NewV4()), Name: "dinner"}},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if m["success"] != true {
		t.Error("Expected success to be true")
	}
}
