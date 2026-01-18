package storage

import (
	"math"
	"testing"

	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/units"
)

func TestCreateStepIngredientWithUnit(t *testing.T) {
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("ingredient-unit-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
	}
	defer DeleteRecipe(recipe.UUID)

	step, err := CreateRecipeStep(recipe.UUID, 1, "Mix ingredients")
	if err != nil {
		t.Fatalf("CreateRecipeStep failed: %v", err)
	}

	ingredientName, err := UpsertIngredientName("flour")
	if err != nil {
		t.Fatalf("UpsertIngredientName failed: %v", err)
	}

	tests := []struct {
		name         string
		unit         string
		quantity     float64
		wantType     models.IngredientUnit
		wantQuantity float64
		tolerance    float64
	}{
		// Volume units
		{"teaspoon", "tsp", 2, models.UnitML, 9.85784, 0.001},
		{"tablespoon", "tbsp", 1, models.UnitML, 14.7868, 0.001},
		{"cup", "cup", 1, models.UnitML, 236.588, 0.001},
		{"cups plural", "cups", 2, models.UnitML, 473.176, 0.001},
		{"quart", "qt", 1, models.UnitML, 946.353, 0.001},
		{"gallon", "gal", 1, models.UnitML, 3785.41, 0.01},
		{"fluid ounce", "fl oz", 1, models.UnitML, 29.5735, 0.001},
		{"milliliter", "ml", 100, models.UnitML, 100, 0.001},
		{"liter", "l", 1, models.UnitML, 1000, 0.001},
		{"liter spelled out", "liter", 2, models.UnitML, 2000, 0.001},

		// Mass units
		{"gram", "g", 100, models.UnitMG, 100000, 0.001},
		{"grams plural", "grams", 50, models.UnitMG, 50000, 0.001},
		{"kilogram", "kg", 1, models.UnitMG, 1000000, 0.001},
		{"ounce", "oz", 1, models.UnitMG, 28349.5, 0.1},
		{"ounces plural", "ounces", 8, models.UnitMG, 226796, 1},
		{"pound", "lb", 1, models.UnitMG, 453592, 1},
		{"pounds plural", "lbs", 2, models.UnitMG, 907184, 1},
		{"milligram", "mg", 500, models.UnitMG, 500, 0.001},

		// Count units
		{"count", "count", 3, models.UnitCount, 3, 0.001},
		{"piece", "piece", 1, models.UnitCount, 1, 0.001},
		{"pieces", "pieces", 4, models.UnitCount, 4, 0.001},
		{"each", "each", 2, models.UnitCount, 2, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si, err := CreateStepIngredientWithUnit(step.UUID, ingredientName.UUID, tt.unit, tt.quantity)
			if err != nil {
				t.Fatalf("CreateStepIngredientWithUnit(%q, %f) failed: %v", tt.unit, tt.quantity, err)
			}

			if si.IngredientType != tt.wantType {
				t.Errorf("got type %v, want %v", si.IngredientType, tt.wantType)
			}

			if math.Abs(si.Quantity-tt.wantQuantity) > tt.tolerance {
				t.Errorf("got quantity %f, want %f (tolerance %f)", si.Quantity, tt.wantQuantity, tt.tolerance)
			}

			DeleteStepIngredient(si.UUID)
		})
	}
}

func TestCreateStepIngredientWithUnit_InvalidUnit(t *testing.T) {
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("invalid-unit-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
	}
	defer DeleteRecipe(recipe.UUID)

	step, err := CreateRecipeStep(recipe.UUID, 1, "Mix")
	if err != nil {
		t.Fatalf("CreateRecipeStep failed: %v", err)
	}

	ingredientName, err := UpsertIngredientName("sugar")
	if err != nil {
		t.Fatalf("UpsertIngredientName failed: %v", err)
	}

	_, err = CreateStepIngredientWithUnit(step.UUID, ingredientName.UUID, "bogusunit", 1)
	if err == nil {
		t.Error("expected error for invalid unit, got nil")
	}
}

func TestComplexRecipeWithIngredients(t *testing.T) {
	ensureTestUser(t)

	recipe, err := InsertRecipeByEmail("complex-recipe-test", testEmail, nil)
	if err != nil {
		t.Fatalf("InsertRecipeByEmail failed: %v", err)
	}
	defer DeleteRecipe(recipe.UUID)

	flour, _ := UpsertIngredientName("all-purpose flour")
	sugar, _ := UpsertIngredientName("granulated sugar")
	butter, _ := UpsertIngredientName("unsalted butter")
	eggs, _ := UpsertIngredientName("large eggs")
	vanilla, _ := UpsertIngredientName("vanilla extract")
	milk, _ := UpsertIngredientName("whole milk")
	salt, _ := UpsertIngredientName("salt")

	// Step 1: Dry ingredients
	step1, err := CreateRecipeStep(recipe.UUID, 1, "Combine dry ingredients in a large bowl")
	if err != nil {
		t.Fatalf("CreateRecipeStep failed: %v", err)
	}

	CreateStepIngredientWithUnit(step1.UUID, flour.UUID, "cups", 2)
	CreateStepIngredientWithUnit(step1.UUID, sugar.UUID, "cup", 1)
	CreateStepIngredientWithUnit(step1.UUID, salt.UUID, "tsp", 0.5)

	// Step 2: Wet ingredients
	step2, err := CreateRecipeStep(recipe.UUID, 2, "Cream butter and mix wet ingredients")
	if err != nil {
		t.Fatalf("CreateRecipeStep failed: %v", err)
	}

	CreateStepIngredientWithUnit(step2.UUID, butter.UUID, "oz", 4)
	CreateStepIngredientWithUnit(step2.UUID, eggs.UUID, "pieces", 2)
	CreateStepIngredientWithUnit(step2.UUID, vanilla.UUID, "tsp", 1)
	CreateStepIngredientWithUnit(step2.UUID, milk.UUID, "cup", 0.5)

	// Step 3: Combine and bake (no ingredients)
	_, err = CreateRecipeStep(recipe.UUID, 3, "Fold wet into dry and bake at 350F for 30 minutes")
	if err != nil {
		t.Fatalf("CreateRecipeStep failed: %v", err)
	}

	// Query back the full recipe
	stepsWithIngredients, err := GetRecipeStepsWithIngredients(recipe.UUID)
	if err != nil {
		t.Fatalf("GetRecipeStepsWithIngredients failed: %v", err)
	}

	if len(stepsWithIngredients) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(stepsWithIngredients))
	}

	// Verify step 1
	if stepsWithIngredients[0].StepNumber != 1 {
		t.Errorf("expected step 1, got %d", stepsWithIngredients[0].StepNumber)
	}
	if len(stepsWithIngredients[0].Ingredients) != 3 {
		t.Errorf("expected 3 ingredients in step 1, got %d", len(stepsWithIngredients[0].Ingredients))
	}

	// Verify step 2
	if stepsWithIngredients[1].StepNumber != 2 {
		t.Errorf("expected step 2, got %d", stepsWithIngredients[1].StepNumber)
	}
	if len(stepsWithIngredients[1].Ingredients) != 4 {
		t.Errorf("expected 4 ingredients in step 2, got %d", len(stepsWithIngredients[1].Ingredients))
	}

	// Verify step 3 has no ingredients
	if stepsWithIngredients[2].StepNumber != 3 {
		t.Errorf("expected step 3, got %d", stepsWithIngredients[2].StepNumber)
	}
	if len(stepsWithIngredients[2].Ingredients) != 0 {
		t.Errorf("expected 0 ingredients in step 3, got %d", len(stepsWithIngredients[2].Ingredients))
	}

	// Verify ingredient names are populated and convert back to display units
	foundFlour := false
	for _, ing := range stepsWithIngredients[0].Ingredients {
		if ing.IngredientName == "all-purpose flour" {
			foundFlour = true
			// Stored as ml (base unit for volume)
			if ing.IngredientType != models.UnitML {
				t.Errorf("flour should be stored as ml, got %v", ing.IngredientType)
			}
			q := units.FindBestIntegerUnit(ing.Quantity, units.CategoryVolume)
			if q.Unit != "cup" || q.Value != 2 {
				t.Errorf("flour should convert to 2 cups, got %v %s", q.Value, q.Unit)
			}
		}
	}
	if !foundFlour {
		t.Error("flour not found in step 1 ingredients")
	}

	// Test GetAllIngredientsForRecipe
	allIngredients, err := GetAllIngredientsForRecipe(recipe.UUID)
	if err != nil {
		t.Fatalf("GetAllIngredientsForRecipe failed: %v", err)
	}

	if len(allIngredients) != 7 {
		t.Errorf("expected 7 total ingredients, got %d", len(allIngredients))
	}

	// Verify eggs are stored as count
	for _, ing := range allIngredients {
		if ing.IngredientName == "large eggs" {
			if ing.IngredientType != models.UnitCount {
				t.Errorf("eggs should be count, got %v", ing.IngredientType)
			}
			if ing.Quantity != 2 {
				t.Errorf("expected 2 eggs, got %f", ing.Quantity)
			}
		}
	}

	// Verify ingredient_names persist (they're shared across recipes)
	flourCheck, err := GetIngredientNameByName("all-purpose flour")
	if err != nil || flourCheck == nil {
		t.Error("ingredient names should persist")
	}
}

func TestUnitConversionRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		quantity float64
	}{
		{"1 cup to ml and back", "cup", 1},
		{"2.5 tbsp to ml and back", "tbsp", 2.5},
		{"100g to mg and back", "g", 100},
		{"1 lb to mg and back", "lb", 1},
		{"3 pieces", "pieces", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitKey, category, err := units.ParseUnit(tt.unit)
			if err != nil {
				t.Fatalf("ParseUnit failed: %v", err)
			}

			baseValue, _, err := units.ToBaseUnit(tt.quantity, unitKey)
			if err != nil {
				t.Fatalf("ToBaseUnit failed: %v", err)
			}

			convertedBack, err := units.FromBaseUnit(baseValue, category, unitKey)
			if err != nil {
				t.Fatalf("FromBaseUnit failed: %v", err)
			}

			if math.Abs(convertedBack-tt.quantity) > 0.0001 {
				t.Errorf("round trip failed: started with %f %s, got back %f", tt.quantity, tt.unit, convertedBack)
			}
		})
	}
}

func TestFindBestIntegerUnit(t *testing.T) {
	tests := []struct {
		name      string
		baseValue float64
		category  units.UnitCategory
		wantUnit  string
		wantValue float64
	}{
		{"236ml is 1 cup", 236.588, units.CategoryVolume, "cup", 1},
		{"1000ml is 1 liter", 1000, units.CategoryVolume, "l", 1},
		{"453592mg is 1 lb", 453592, units.CategoryMass, "lb", 1},
		{"28349.5mg is 1 oz", 28349.5, units.CategoryMass, "oz", 1},
		{"1000mg is 1 gram", 1000, units.CategoryMass, "g", 1},
		{"5 count stays 5 count", 5, units.CategoryCount, "count", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := units.FindBestIntegerUnit(tt.baseValue, tt.category)

			if q.Unit != tt.wantUnit {
				t.Errorf("got unit %q, want %q", q.Unit, tt.wantUnit)
			}

			if q.Value != tt.wantValue {
				t.Errorf("got value %f, want %f", q.Value, tt.wantValue)
			}
		})
	}
}

func TestFindBestIntegerUnit_NoMatch(t *testing.T) {
	// 0.123ml is too small for any unit to produce value >= 1
	q := units.FindBestIntegerUnit(0.123, units.CategoryVolume)

	// Should return base unit when no good match
	if q.Unit != "ml" {
		t.Errorf("expected ml for non-matching value, got %q", q.Unit)
	}
	if q.Value != 0.123 {
		t.Errorf("expected 0.123, got %f", q.Value)
	}
}
