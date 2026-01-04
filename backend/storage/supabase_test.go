package storage

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestCreateTagUUID(t *testing.T) {
	// Deterministic UUIDs should return the same value for the same input
	tests := []struct {
		tag string
	}{
		{"dinner"},
		{"quick"},
		{"breakfast"},
		{"vegetarian"},
		{"dessert"},
		{""},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got1 := CreateTagUUID(tt.tag)
			got2 := CreateTagUUID(tt.tag)

			// Same input should produce same UUID
			if got1 != got2 {
				t.Errorf("CreateTagUUID(%q) not deterministic: got %v and %v", tt.tag, got1, got2)
			}

			// Should not be nil UUID (unless that's intentional for empty string)
			if tt.tag != "" && got1 == uuid.Nil {
				t.Errorf("CreateTagUUID(%q) returned nil UUID", tt.tag)
			}
		})
	}

	// Different inputs should produce different UUIDs
	uuid1 := CreateTagUUID("dinner")
	uuid2 := CreateTagUUID("breakfast")
	if uuid1 == uuid2 {
		t.Error("Different tags should produce different UUIDs")
	}
}

// Integration tests - these require a real Supabase connection
// Skip them if SUPABASE_URL is not set

func TestGetRecipesByUserUUID(t *testing.T) {
	t.Skip("Integration test - implement when Supabase client is ready")

	testUserUUID := uuid.Must(uuid.NewV4())
	recipes, err := GetRecipesByUserUUID(testUserUUID)
	if err != nil {
		t.Fatalf("GetRecipesByUserUUID failed: %v", err)
	}

	if recipes == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestGetFileMappingsByRecipeUUIDs(t *testing.T) {
	t.Skip("Integration test - implement when Supabase client is ready")

	testUUIDs := []uuid.UUID{
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
	}

	mappings, err := GetFileMappingsByRecipeUUIDs(testUUIDs)
	if err != nil {
		t.Fatalf("GetFileMappingsByRecipeUUIDs failed: %v", err)
	}

	if mappings == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestGetFilesByUUIDs(t *testing.T) {
	t.Skip("Integration test - implement when Supabase client is ready")

	testUUIDs := []uuid.UUID{
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
	}

	files, err := GetFilesByUUIDs(testUUIDs)
	if err != nil {
		t.Fatalf("GetFilesByUUIDs failed: %v", err)
	}

	if files == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestInsertRecipe(t *testing.T) {
	t.Skip("Integration test - implement when Supabase client is ready")

	testUserUUID := uuid.Must(uuid.NewV4())
	recipe, err := InsertRecipe("test-recipe", testUserUUID, "test, tags")
	if err != nil {
		t.Fatalf("InsertRecipe failed: %v", err)
	}

	if recipe == nil {
		t.Fatal("Expected recipe, got nil")
	}
	if recipe.Filename != "test-recipe" {
		t.Errorf("Expected filename 'test-recipe', got %q", recipe.Filename)
	}
	if recipe.UUID == uuid.Nil {
		t.Error("Expected non-nil UUID")
	}
}

func TestUpsertTag(t *testing.T) {
	t.Skip("Integration test - implement when Supabase client is ready")

	tagUUID := CreateTagUUID("test-tag")
	tag, err := UpsertTag(tagUUID, "test-tag")
	if err != nil {
		t.Fatalf("UpsertTag failed: %v", err)
	}

	if tag == nil {
		t.Fatal("Expected tag, got nil")
	}
	if tag.Name != "test-tag" {
		t.Errorf("Expected name 'test-tag', got %q", tag.Name)
	}

	// Upsert again - should not error
	tag2, err := UpsertTag(tagUUID, "test-tag")
	if err != nil {
		t.Fatalf("UpsertTag (second call) failed: %v", err)
	}
	if tag2.UUID != tag.UUID {
		t.Error("Expected same UUID on upsert")
	}
}
