package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

func TestRecipeJSON(t *testing.T) {
	recipe := Recipe{
		UUID:       uuid.Must(uuid.NewV4()),
		Name:       "test-recipe",
		User:       uuid.Must(uuid.NewV4()),
		OwnerEmail: "owner@example.com",
		TagString:  "dinner, quick",
		CreatedAt:  time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"uuid", "name", "user_uuid", "owner_email", "tag_string", "created_at"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}

	if _, ok := m["uuid"].(string); !ok {
		t.Error("Expected 'uuid' to be a string (UUID)")
	}
}

func TestFileJSON(t *testing.T) {
	file := File{
		UUID:       uuid.Must(uuid.NewV4()),
		RecipeUUID: uuid.Must(uuid.NewV4()),
		URL:        "https://example.com/image.jpg",
		PageNumber: 0,
		Image:      true,
	}

	data, err := json.Marshal(file)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"uuid", "recipe_uuid", "url", "page_number", "image"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}
}

func TestRecipesResponseJSON(t *testing.T) {
	response := RecipesResponse{
		RecipeData: []Recipe{},
		FileData:   []File{},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	expectedFields := []string{"recipeData", "fileData"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("Expected field %q in JSON output", field)
		}
	}
}

func TestUploadResponseJSON(t *testing.T) {
	response := UploadResponse{
		Success: true,
		Recipe:  Recipe{UUID: uuid.Must(uuid.NewV4()), Name: "test"},
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
