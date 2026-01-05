package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type IngredientUnit string

const (
	UnitML    IngredientUnit = "ml"
	UnitMG    IngredientUnit = "mg"
	UnitCount IngredientUnit = "count"
)

type IngredientName struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RecipeStep struct {
	UUID         uuid.UUID `json:"uuid"`
	RecipeUUID   uuid.UUID `json:"recipe_uuid"`
	StepNumber   int       `json:"step_number"`
	Instructions string    `json:"instructions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StepIngredient struct {
	UUID               uuid.UUID      `json:"uuid"`
	RecipeStepUUID     uuid.UUID      `json:"recipe_step_uuid"`
	IngredientNameUUID uuid.UUID      `json:"ingredient_name_uuid"`
	IngredientType     IngredientUnit `json:"ingredient_type"`
	Quantity           float64        `json:"quantity"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// StepIngredientWithName includes the ingredient name for convenience
type StepIngredientWithName struct {
	StepIngredient
	IngredientName string `json:"ingredient_name"`
}

// RecipeStepWithIngredients includes all ingredients for a step
type RecipeStepWithIngredients struct {
	RecipeStep
	Ingredients []StepIngredientWithName `json:"ingredients"`
}

// API response types

type RecipeStepResponse struct {
	Instruction string   `json:"instruction"`
	Ingredients []string `json:"ingredients"`
}

type RecipeStepsResponse struct {
	Steps []RecipeStepResponse `json:"steps"`
}
