package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Recipe struct {
	UUID      uuid.UUID `json:"uuid"`
	Filename  string    `json:"filename"`
	User      uuid.UUID `json:"user_uuid"`
	TagString string    `json:"tag_string"`
	CreatedAt time.Time `json:"created_at"`
}

type File struct {
	UUID  uuid.UUID `json:"uuid"`
	URL   string    `json:"url"`
	Image bool      `json:"image"`
}

type FileRecipe struct {
	FileUUID   uuid.UUID `json:"file_uuid"`
	RecipeUUID uuid.UUID `json:"recipe_uuid"`
	PageNumber int       `json:"page_number"`
}

type Tag struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type RecipeTag struct {
	RecipeUUID uuid.UUID `json:"recipe_uuid"`
	TagUUID    uuid.UUID `json:"tag_uuid"`
}

// Response for GET /api/recipes
type RecipesResponse struct {
	RecipeData  []Recipe     `json:"recipeData"`
	FileData    []File       `json:"fileData"`
	MappingData []FileRecipe `json:"mappingData"`
}

// Response for POST /api/recipes
type UploadResponse struct {
	Success bool   `json:"success"`
	Recipe  Recipe `json:"recipe"`
	Tags    []Tag  `json:"tags"`
}
