package storage

import (
	"io"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

// TODO: Add client variable to hold Supabase connection

func Init(url, key string) error {
	// TODO: Initialize Supabase client
	return nil
}

// --- Database Reads ---

func GetRecipesByUserUUID(userUUID uuid.UUID) ([]models.Recipe, error) {
	// TODO: Query "recipes" table
	// SELECT uuid, filename, tag_string, created_at
	// WHERE user_uuid = userUUID
	// ORDER BY created_at DESC
	// LIMIT 100
	return nil, nil
}

func GetFileMappingsByRecipeUUIDs(recipeUUIDs []uuid.UUID) ([]models.FileRecipe, error) {
	// TODO: Query "file_recipes" table
	// WHERE recipe_id IN recipeUUIDs
	return nil, nil
}

func GetFilesByUUIDs(fileUUIDs []uuid.UUID) ([]models.File, error) {
	// TODO: Query "files" table
	// WHERE id IN fileUUIDs
	return nil, nil
}

// --- Database Writes ---

func InsertRecipe(filename string, user uuid.UUID, tagString string) (*models.Recipe, error) {
	// TODO: Insert into "recipes" table, return created record
	// Generate UUID with uuid.NewV4() or let database generate it
	return nil, nil
}

func InsertFile(url string, isImage bool) (*models.File, error) {
	// TODO: Insert into "files" table, return created record
	// Generate UUID with uuid.NewV4() or let database generate it
	return nil, nil
}

func InsertFileRecipe(fileUUID, recipeUUID uuid.UUID, pageNumber int) error {
	// TODO: Insert into "file_recipes" table
	return nil
}

func UpsertTag(tagUUID uuid.UUID, name string) (*models.Tag, error) {
	// TODO: Insert into "tags" table, ON CONFLICT update
	return nil, nil
}

func InsertRecipeTag(recipeUUID, tagUUID uuid.UUID) error {
	// TODO: Insert into "recipe_tags" table
	return nil
}

// --- File Storage ---

func UploadFile(filename string, file io.Reader) (string, error) {
	// TODO: Upload to Supabase Storage bucket "recipe-images"
	// Return public URL
	return "", nil
}

// --- Helpers ---

func CreateTagUUID(tag string) uuid.UUID {
	// TODO: Generate deterministic UUID from tag name
	// Use uuid.NewV5 with a namespace for deterministic UUIDs:
	//   namespace := uuid.Must(uuid.FromString("your-namespace-uuid"))
	//   return uuid.NewV5(namespace, tag)
	return uuid.UUID{}
}
