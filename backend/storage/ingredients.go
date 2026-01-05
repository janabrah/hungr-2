package storage

import (
	"context"
	"fmt"

	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/units"
	"github.com/gofrs/uuid"
)

// IngredientName operations

func GetIngredientNameByUUID(ingredientUUID uuid.UUID) (*models.IngredientName, error) {
	var i models.IngredientName
	err := db.QueryRow(context.Background(),
		`SELECT uuid, name, created_at, updated_at
		 FROM ingredient_names WHERE uuid = $1`, ingredientUUID).Scan(
		&i.UUID, &i.Name, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func GetIngredientNameByName(name string) (*models.IngredientName, error) {
	var i models.IngredientName
	err := db.QueryRow(context.Background(),
		`SELECT uuid, name, created_at, updated_at
		 FROM ingredient_names WHERE name = $1`, name).Scan(
		&i.UUID, &i.Name, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func CreateIngredientName(name string) (*models.IngredientName, error) {
	var i models.IngredientName
	err := db.QueryRow(context.Background(),
		`INSERT INTO ingredient_names (name)
		 VALUES ($1)
		 RETURNING uuid, name, created_at, updated_at`, name).Scan(
		&i.UUID, &i.Name, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func UpsertIngredientName(name string) (*models.IngredientName, error) {
	var i models.IngredientName
	err := db.QueryRow(context.Background(),
		`INSERT INTO ingredient_names (name)
		 VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET updated_at = NOW()
		 RETURNING uuid, name, created_at, updated_at`, name).Scan(
		&i.UUID, &i.Name, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func ListIngredientNames() ([]models.IngredientName, error) {
	rows, err := db.Query(context.Background(),
		`SELECT uuid, name, created_at, updated_at
		 FROM ingredient_names ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.IngredientName
	for rows.Next() {
		var i models.IngredientName
		if err := rows.Scan(&i.UUID, &i.Name, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, rows.Err()
}

func DeleteIngredientName(ingredientUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM ingredient_names WHERE uuid = $1`, ingredientUUID)
	return err
}

// RecipeStep operations

func GetRecipeStepByUUID(stepUUID uuid.UUID) (*models.RecipeStep, error) {
	var s models.RecipeStep
	err := db.QueryRow(context.Background(),
		`SELECT uuid, recipe_uuid, step_number, instructions, created_at, updated_at
		 FROM recipe_steps WHERE uuid = $1`, stepUUID).Scan(
		&s.UUID, &s.RecipeUUID, &s.StepNumber, &s.Instructions, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetRecipeStepsByRecipeUUID(recipeUUID uuid.UUID) ([]models.RecipeStep, error) {
	rows, err := db.Query(context.Background(),
		`SELECT uuid, recipe_uuid, step_number, instructions, created_at, updated_at
		 FROM recipe_steps WHERE recipe_uuid = $1
		 ORDER BY step_number`, recipeUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []models.RecipeStep
	for rows.Next() {
		var s models.RecipeStep
		if err := rows.Scan(&s.UUID, &s.RecipeUUID, &s.StepNumber, &s.Instructions, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}

func CreateRecipeStep(recipeUUID uuid.UUID, stepNumber int, instructions string) (*models.RecipeStep, error) {
	var s models.RecipeStep
	err := db.QueryRow(context.Background(),
		`INSERT INTO recipe_steps (recipe_uuid, step_number, instructions)
		 VALUES ($1, $2, $3)
		 RETURNING uuid, recipe_uuid, step_number, instructions, created_at, updated_at`,
		recipeUUID, stepNumber, instructions).Scan(
		&s.UUID, &s.RecipeUUID, &s.StepNumber, &s.Instructions, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateRecipeStep(stepUUID uuid.UUID, stepNumber int, instructions string) (*models.RecipeStep, error) {
	var s models.RecipeStep
	err := db.QueryRow(context.Background(),
		`UPDATE recipe_steps SET step_number = $1, instructions = $2, updated_at = NOW()
		 WHERE uuid = $3
		 RETURNING uuid, recipe_uuid, step_number, instructions, created_at, updated_at`,
		stepNumber, instructions, stepUUID).Scan(
		&s.UUID, &s.RecipeUUID, &s.StepNumber, &s.Instructions, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteRecipeStep(stepUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM recipe_steps WHERE uuid = $1`, stepUUID)
	return err
}

func DeleteRecipeStepsByRecipeUUID(recipeUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM recipe_steps WHERE recipe_uuid = $1`, recipeUUID)
	return err
}

// StepIngredient operations

func GetStepIngredientByUUID(ingredientUUID uuid.UUID) (*models.StepIngredient, error) {
	var si models.StepIngredient
	err := db.QueryRow(context.Background(),
		`SELECT uuid, recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity, created_at, updated_at
		 FROM step_ingredients WHERE uuid = $1`, ingredientUUID).Scan(
		&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func GetStepIngredientsByStepUUID(stepUUID uuid.UUID) ([]models.StepIngredient, error) {
	rows, err := db.Query(context.Background(),
		`SELECT uuid, recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity, created_at, updated_at
		 FROM step_ingredients WHERE recipe_step_uuid = $1`, stepUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.StepIngredient
	for rows.Next() {
		var si models.StepIngredient
		if err := rows.Scan(&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, si)
	}
	return ingredients, rows.Err()
}

func GetStepIngredientsWithNamesByStepUUID(stepUUID uuid.UUID) ([]models.StepIngredientWithName, error) {
	rows, err := db.Query(context.Background(),
		`SELECT si.uuid, si.recipe_step_uuid, si.ingredient_name_uuid, si.ingredient_type, si.quantity, si.created_at, si.updated_at, n.name
		 FROM step_ingredients si
		 JOIN ingredient_names n ON si.ingredient_name_uuid = n.uuid
		 WHERE si.recipe_step_uuid = $1`, stepUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.StepIngredientWithName
	for rows.Next() {
		var si models.StepIngredientWithName
		if err := rows.Scan(&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt, &si.IngredientName); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, si)
	}
	return ingredients, rows.Err()
}

func CreateStepIngredient(stepUUID, ingredientNameUUID uuid.UUID, ingredientType models.IngredientUnit, quantity float64) (*models.StepIngredient, error) {
	var si models.StepIngredient
	err := db.QueryRow(context.Background(),
		`INSERT INTO step_ingredients (recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity)
		 VALUES ($1, $2, $3, $4)
		 RETURNING uuid, recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity, created_at, updated_at`,
		stepUUID, ingredientNameUUID, ingredientType, quantity).Scan(
		&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func CreateStepIngredientWithUnit(stepUUID, ingredientNameUUID uuid.UUID, unit string, quantity float64) (*models.StepIngredient, error) {
	unitKey, category, err := units.ParseUnit(unit)
	if err != nil {
		return nil, fmt.Errorf("invalid unit %q: %w", unit, err)
	}

	baseValue, _, err := units.ToBaseUnit(quantity, unitKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %f %s to base unit: %w", quantity, unit, err)
	}

	var ingredientType models.IngredientUnit
	switch category {
	case units.CategoryVolume:
		ingredientType = models.UnitML
	case units.CategoryMass:
		ingredientType = models.UnitMG
	case units.CategoryCount:
		ingredientType = models.UnitCount
	default:
		return nil, fmt.Errorf("unknown unit category: %s", category)
	}

	return CreateStepIngredient(stepUUID, ingredientNameUUID, ingredientType, baseValue)
}

func UpdateStepIngredient(ingredientUUID uuid.UUID, ingredientType models.IngredientUnit, quantity float64) (*models.StepIngredient, error) {
	var si models.StepIngredient
	err := db.QueryRow(context.Background(),
		`UPDATE step_ingredients SET ingredient_type = $1, quantity = $2, updated_at = NOW()
		 WHERE uuid = $3
		 RETURNING uuid, recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity, created_at, updated_at`,
		ingredientType, quantity, ingredientUUID).Scan(
		&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func DeleteStepIngredient(ingredientUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM step_ingredients WHERE uuid = $1`, ingredientUUID)
	return err
}

func DeleteStepIngredientsByStepUUID(stepUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM step_ingredients WHERE recipe_step_uuid = $1`, stepUUID)
	return err
}

// Aggregate queries

func GetRecipeStepsWithIngredients(recipeUUID uuid.UUID) ([]models.RecipeStepWithIngredients, error) {
	steps, err := GetRecipeStepsByRecipeUUID(recipeUUID)
	if err != nil {
		return nil, err
	}

	result := make([]models.RecipeStepWithIngredients, len(steps))
	for i, step := range steps {
		ingredients, err := GetStepIngredientsWithNamesByStepUUID(step.UUID)
		if err != nil {
			return nil, err
		}
		result[i] = models.RecipeStepWithIngredients{
			RecipeStep:  step,
			Ingredients: ingredients,
		}
	}
	return result, nil
}

// GetAllIngredientsForRecipe returns all ingredients across all steps for a recipe
func GetAllIngredientsForRecipe(recipeUUID uuid.UUID) ([]models.StepIngredientWithName, error) {
	rows, err := db.Query(context.Background(),
		`SELECT si.uuid, si.recipe_step_uuid, si.ingredient_name_uuid, si.ingredient_type, si.quantity, si.created_at, si.updated_at, n.name
		 FROM step_ingredients si
		 JOIN ingredient_names n ON si.ingredient_name_uuid = n.uuid
		 JOIN recipe_steps rs ON si.recipe_step_uuid = rs.uuid
		 WHERE rs.recipe_uuid = $1
		 ORDER BY rs.step_number, n.name`, recipeUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.StepIngredientWithName
	for rows.Next() {
		var si models.StepIngredientWithName
		if err := rows.Scan(&si.UUID, &si.RecipeStepUUID, &si.IngredientNameUUID, &si.IngredientType, &si.Quantity, &si.CreatedAt, &si.UpdatedAt, &si.IngredientName); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, si)
	}
	return ingredients, rows.Err()
}

type StepInput struct {
	Instruction string
	Ingredients []IngredientInput
}

type IngredientInput struct {
	Name     string
	Unit     string
	Quantity float64
}

// ReplaceRecipeSteps deletes all existing steps and creates new ones
func ReplaceRecipeSteps(recipeUUID uuid.UUID, steps []StepInput) error {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete existing steps (cascades to ingredients)
	_, err = tx.Exec(ctx, `DELETE FROM recipe_steps WHERE recipe_uuid = $1`, recipeUUID)
	if err != nil {
		return fmt.Errorf("failed to delete existing steps: %w", err)
	}

	// Create new steps
	for i, step := range steps {
		var stepUUID uuid.UUID
		err := tx.QueryRow(ctx,
			`INSERT INTO recipe_steps (recipe_uuid, step_number, instructions)
			 VALUES ($1, $2, $3) RETURNING uuid`,
			recipeUUID, i+1, step.Instruction).Scan(&stepUUID)
		if err != nil {
			return fmt.Errorf("failed to create step %d: %w", i+1, err)
		}

		// Create ingredients for this step
		for _, ing := range step.Ingredients {
			// Upsert ingredient name
			var ingredientNameUUID uuid.UUID
			err := tx.QueryRow(ctx,
				`INSERT INTO ingredient_names (name)
				 VALUES ($1)
				 ON CONFLICT (name) DO UPDATE SET updated_at = NOW()
				 RETURNING uuid`, ing.Name).Scan(&ingredientNameUUID)
			if err != nil {
				return fmt.Errorf("failed to upsert ingredient name %q: %w", ing.Name, err)
			}

			// Convert to base unit
			baseValue, category, err := units.ToBaseUnit(ing.Quantity, ing.Unit)
			if err != nil {
				return fmt.Errorf("failed to convert unit %q: %w", ing.Unit, err)
			}

			var ingredientType models.IngredientUnit
			switch category {
			case units.CategoryVolume:
				ingredientType = models.UnitML
			case units.CategoryMass:
				ingredientType = models.UnitMG
			case units.CategoryCount:
				ingredientType = models.UnitCount
			}

			_, err = tx.Exec(ctx,
				`INSERT INTO step_ingredients (recipe_step_uuid, ingredient_name_uuid, ingredient_type, quantity)
				 VALUES ($1, $2, $3, $4)`,
				stepUUID, ingredientNameUUID, ingredientType, baseValue)
			if err != nil {
				return fmt.Errorf("failed to create ingredient: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}
