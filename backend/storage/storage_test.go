package storage

import (
	"context"
	"os"
	"testing"

	"github.com/gofrs/uuid"
)

const testEmail = "test@example.com"

func ensureTestUser(t *testing.T) {
	_, err := GetUserByEmail(testEmail)
	if err != nil {
		_, err = CreateUser(testEmail, "Test User")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
}

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

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL must be set to run tests")
	}
	if err := Init(dbURL); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

func TestGetRecipesByUserEmail(t *testing.T) {
	ensureTestUser(t)

	recipes, err := GetRecipesByUserEmail(testEmail)
	if err != nil {
		t.Fatalf("GetRecipesByUserEmail failed: %v", err)
	}

	if recipes == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestInsertRecipeByEmail(t *testing.T) {
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("test-recipe", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
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
	if recipe.Source != nil {
		t.Errorf("Expected nil source, got %v", recipe.Source)
	}
}

func TestInsertRecipeByEmail_WithSource(t *testing.T) {
	ensureTestUser(t)

	source := "website"
	recipe, err := InsertRecipeByEmail("source-recipe", testEmail, &source)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
	}

	if recipe.Source == nil || *recipe.Source != source {
		t.Fatalf("Expected source %q, got %v", source, recipe.Source)
	}

	found, err := GetRecipeByUUID(recipe.UUID)
	if err != nil {
		t.Fatalf("GetRecipeByUUID failed: %v", err)
	}

	if found.Source == nil || *found.Source != source {
		t.Fatalf("Expected source %q from GetRecipeByUUID, got %v", source, found.Source)
	}
}

func TestInsertAndGetFile(t *testing.T) {
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("file-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
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
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("multi-file-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
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
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("delete-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
	}

	_, err = InsertFile(recipe.UUID, []byte("file data"), "image/jpeg", 0, true)
	if err != nil {
		t.Fatalf("InsertFile failed: %v", err)
	}

	err = DeleteRecipe(recipe.UUID)
	if err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	recipes, err := GetRecipesByUserEmail(testEmail)
	if err != nil {
		t.Fatalf("GetRecipesByUserEmail failed: %v", err)
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

func TestTransactionCommit(t *testing.T) {
	ensureTestUser(t)

	ctx := context.Background()

	tx, err := BeginTx(ctx)
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}

	// Insert recipe in transaction
	recipe, err := TxInsertRecipeByEmail(ctx, tx, "tx-commit-test", testEmail, nil)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxInsertRecipeByEmail failed: %v", err)
	}

	// Insert file in transaction
	_, err = TxInsertFile(ctx, tx, recipe.UUID, []byte("tx file data"), "image/png", 0, true)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxInsertFile failed: %v", err)
	}

	// Insert tag in transaction
	tagUUID := CreateTagUUID("tx-test-tag")
	_, err = TxUpsertTag(ctx, tx, tagUUID, "tx-test-tag")
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxUpsertTag failed: %v", err)
	}

	err = TxInsertRecipeTag(ctx, tx, recipe.UUID, tagUUID)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxInsertRecipeTag failed: %v", err)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Verify recipe exists after commit
	foundRecipe, err := GetRecipeByUUID(recipe.UUID)
	if err != nil {
		t.Fatalf("GetRecipeByUUID failed: %v", err)
	}
	if foundRecipe.Name != "tx-commit-test" {
		t.Errorf("Expected name 'tx-commit-test', got %q", foundRecipe.Name)
	}

	// Verify file exists after commit
	files, err := GetFilesByRecipeUUIDs([]uuid.UUID{recipe.UUID})
	if err != nil {
		t.Fatalf("GetFilesByRecipeUUIDs failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file after commit, got %d", len(files))
	}

	// Cleanup
	DeleteRecipe(recipe.UUID)
}

func TestTransactionRollback(t *testing.T) {
	ensureTestUser(t)

	ctx := context.Background()

	tx, err := BeginTx(ctx)
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}

	// Insert recipe in transaction
	recipe, err := TxInsertRecipeByEmail(ctx, tx, "tx-rollback-test", testEmail, nil)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxInsertRecipeByEmail failed: %v", err)
	}

	// Insert file in transaction
	_, err = TxInsertFile(ctx, tx, recipe.UUID, []byte("rollback file data"), "image/png", 0, true)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatalf("TxInsertFile failed: %v", err)
	}

	// Rollback instead of commit
	if err := tx.Rollback(ctx); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify recipe does NOT exist after rollback
	_, err = GetRecipeByUUID(recipe.UUID)
	if err == nil {
		t.Error("Expected error getting recipe after rollback, but got nil")
		// Cleanup if it somehow exists
		DeleteRecipe(recipe.UUID)
	}

	// Verify no files exist for the rolled-back recipe
	files, err := GetFilesByRecipeUUIDs([]uuid.UUID{recipe.UUID})
	if err != nil {
		t.Fatalf("GetFilesByRecipeUUIDs failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected 0 files after rollback, got %d", len(files))
	}
}
