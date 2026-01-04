package storage

import (
	"os"
	"testing"

	"github.com/gofrs/uuid"
)

func TestCreateTagUUID(t *testing.T) {
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

			if got1 != got2 {
				t.Errorf("CreateTagUUID(%q) not deterministic: got %v and %v", tt.tag, got1, got2)
			}

			if tt.tag != "" && got1 == uuid.Nil {
				t.Errorf("CreateTagUUID(%q) returned nil UUID", tt.tag)
			}
		})
	}

	uuid1 := CreateTagUUID("dinner")
	uuid2 := CreateTagUUID("breakfast")
	if uuid1 == uuid2 {
		t.Error("Different tags should produce different UUIDs")
	}
}

func skipIfNoDatabase(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}
	if db == nil {
		if err := Init(dbURL); err != nil {
			t.Skipf("Could not connect to database: %v", err)
		}
	}
}

func TestGetRecipesByUserUUID(t *testing.T) {
	skipIfNoDatabase(t)

	testUserUUID := uuid.Must(uuid.NewV4())
	recipes, err := GetRecipesByUserUUID(testUserUUID)
	if err != nil {
		t.Fatalf("GetRecipesByUserUUID failed: %v", err)
	}

	if recipes == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestInsertRecipe(t *testing.T) {
	skipIfNoDatabase(t)

	testUserUUID := uuid.Must(uuid.NewV4())
	recipe, err := InsertRecipe("test-recipe", testUserUUID, "test, tags")
	if err != nil {
		t.Fatalf("InsertRecipe failed: %v", err)
	}

	if recipe == nil {
		t.Fatal("Expected recipe, got nil")
	}
	if recipe.Name != "test-recipe" {
		t.Errorf("Expected name 'test-recipe', got %q", recipe.Name)
	}
	if recipe.UUID == uuid.Nil {
		t.Error("Expected non-nil UUID")
	}
}

func TestInsertAndGetFile(t *testing.T) {
	skipIfNoDatabase(t)

	userUUID := uuid.Must(uuid.NewV4())
	recipe, err := InsertRecipe("file-test", userUUID, "test")
	if err != nil {
		t.Fatalf("InsertRecipe failed: %v", err)
	}

	testData := []byte("fake image data")
	file, err := InsertFile(recipe.UUID, testData, "image/jpeg", 0, true)
	if err != nil {
		t.Fatalf("InsertFile failed: %v", err)
	}

	if file == nil {
		t.Fatal("Expected file, got nil")
	}
	if file.RecipeUUID != recipe.UUID {
		t.Errorf("Expected recipe UUID %v, got %v", recipe.UUID, file.RecipeUUID)
	}

	data, contentType, err := GetFileData(file.UUID)
	if err != nil {
		t.Fatalf("GetFileData failed: %v", err)
	}
	if string(data) != string(testData) {
		t.Errorf("Expected data %q, got %q", testData, data)
	}
	if contentType != "image/jpeg" {
		t.Errorf("Expected content type 'image/jpeg', got %q", contentType)
	}

	files, err := GetFilesByRecipeUUIDs([]uuid.UUID{recipe.UUID})
	if err != nil {
		t.Fatalf("GetFilesByRecipeUUIDs failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestUpsertTag(t *testing.T) {
	skipIfNoDatabase(t)

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

	tag2, err := UpsertTag(tagUUID, "test-tag")
	if err != nil {
		t.Fatalf("UpsertTag (second call) failed: %v", err)
	}
	if tag2.UUID != tag.UUID {
		t.Error("Expected same UUID on upsert")
	}
}

func TestMultipleFilesPerRecipe(t *testing.T) {
	skipIfNoDatabase(t)

	userUUID := uuid.Must(uuid.NewV4())
	recipe, err := InsertRecipe("multi-file-test", userUUID, "test")
	if err != nil {
		t.Fatalf("InsertRecipe failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		_, err := InsertFile(recipe.UUID, []byte("page data"), "image/jpeg", i, true)
		if err != nil {
			t.Fatalf("InsertFile failed: %v", err)
		}
	}

	files, err := GetFilesByRecipeUUIDs([]uuid.UUID{recipe.UUID})
	if err != nil {
		t.Fatalf("GetFilesByRecipeUUIDs failed: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
	if files[0].PageNumber != 0 || files[1].PageNumber != 1 || files[2].PageNumber != 2 {
		t.Error("Expected files to be ordered by page_number")
	}
}

func TestDeleteRecipe(t *testing.T) {
	skipIfNoDatabase(t)

	userUUID := uuid.Must(uuid.NewV4())
	recipe, err := InsertRecipe("delete-test", userUUID, "test")
	if err != nil {
		t.Fatalf("InsertRecipe failed: %v", err)
	}

	_, err = InsertFile(recipe.UUID, []byte("file data"), "image/jpeg", 0, true)
	if err != nil {
		t.Fatalf("InsertFile failed: %v", err)
	}

	err = DeleteRecipe(recipe.UUID)
	if err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	recipes, err := GetRecipesByUserUUID(userUUID)
	if err != nil {
		t.Fatalf("GetRecipesByUserUUID failed: %v", err)
	}
	for _, r := range recipes {
		if r.UUID == recipe.UUID {
			t.Error("Recipe should have been deleted")
		}
	}

	files, err := GetFilesByRecipeUUIDs([]uuid.UUID{recipe.UUID})
	if err != nil {
		t.Fatalf("GetFilesByRecipeUUIDs failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected 0 files after recipe deletion, got %d", len(files))
	}
}
