package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Recipe struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	User      uuid.UUID `json:"user_uuid"`
	TagString string    `json:"tag_string"`
	CreatedAt time.Time `json:"created_at"`
}

type File struct {
	UUID        uuid.UUID `json:"uuid"`
	RecipeUUID  uuid.UUID `json:"recipe_uuid"`
	URL         string    `json:"url"`
	PageNumber  int       `json:"page_number"`
	Image       bool      `json:"image"`
	ContentType string    `json:"-"`
	Data        []byte    `json:"-"`
}

type Tag struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type RecipeTag struct {
	RecipeUUID uuid.UUID `json:"recipe_uuid"`
	TagUUID    uuid.UUID `json:"tag_uuid"`
}

type RecipesResponse struct {
	RecipeData []Recipe `json:"recipeData"`
	FileData   []File   `json:"fileData"`
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Recipe  Recipe `json:"recipe"`
	Tags    []Tag  `json:"tags"`
}
