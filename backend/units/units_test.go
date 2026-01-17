package units

import (
	"testing"
)

func TestParseIngredientString(t *testing.T) {
	tests := []struct {
		input    string
		quantity float64
		unit     string
		category UnitCategory
		name     string
	}{
		{"2 cups flour", 2, "cup", CategoryVolume, "flour"},
		{"1 tsp salt", 1, "tsp", CategoryVolume, "salt"},
		{"1/2 cup milk", 0.5, "cup", CategoryVolume, "milk"},
		{"3 tbsp olive oil", 3, "tbsp", CategoryVolume, "olive oil"},
		{"1 lb ground beef", 1, "lb", CategoryMass, "ground beef"},
		{"8 oz cream cheese", 8, "oz", CategoryMass, "cream cheese"},
		{"500 g pasta", 500, "g", CategoryMass, "pasta"},
		{"2 fl oz vanilla extract", 2, "fl_oz", CategoryVolume, "vanilla extract"},
		{"1/4 tsp black pepper", 0.25, "tsp", CategoryVolume, "black pepper"},
		{"3 eggs", 3, "count", CategoryCount, "eggs"},
		{"1 onion", 1, "count", CategoryCount, "onion"},
		{"2 cloves garlic", 2, "count", CategoryCount, "cloves garlic"},
		{"1 cup all-purpose flour", 1, "cup", CategoryVolume, "all-purpose flour"},
		{"2 tablespoons butter", 2, "tbsp", CategoryVolume, "butter"},
		{"1 teaspoon baking powder", 1, "tsp", CategoryVolume, "baking powder"},
		{"100 grams sugar", 100, "g", CategoryMass, "sugar"},
		{"1 pound chicken breast", 1, "lb", CategoryMass, "chicken breast"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseIngredientString(tt.input)
			if err != nil {
				t.Fatalf("ParseIngredientString(%q) returned error: %v", tt.input, err)
			}

			if result.Quantity != tt.quantity {
				t.Errorf("Quantity: got %v, want %v", result.Quantity, tt.quantity)
			}
			if result.Unit != tt.unit {
				t.Errorf("Unit: got %q, want %q", result.Unit, tt.unit)
			}
			if result.Category != tt.category {
				t.Errorf("Category: got %q, want %q", result.Category, tt.category)
			}
			if result.IngredientName != tt.name {
				t.Errorf("IngredientName: got %q, want %q", result.IngredientName, tt.name)
			}
		})
	}
}

func TestParseIngredientString_Fractions(t *testing.T) {
	tests := []struct {
		input    string
		quantity float64
	}{
		{"1/2 cup sugar", 0.5},
		{"1/4 tsp salt", 0.25},
		{"3/4 cup milk", 0.75},
		{"1/3 cup water", 1.0 / 3.0},
		{"2/3 cup broth", 2.0 / 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseIngredientString(tt.input)
			if err != nil {
				t.Fatalf("ParseIngredientString(%q) returned error: %v", tt.input, err)
			}

			delta := 0.0001
			if result.Quantity < tt.quantity-delta || result.Quantity > tt.quantity+delta {
				t.Errorf("Quantity: got %v, want %v", result.Quantity, tt.quantity)
			}
		})
	}
}

func TestParseIngredientString_Errors(t *testing.T) {
	tests := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"1/0 cup flour", "division by zero"},
		{"123", "number only"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := ParseIngredientString(tt.input)
			if err == nil {
				t.Errorf("ParseIngredientString(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func TestParseIngredientString_NoQuantity(t *testing.T) {
	tests := []struct {
		input          string
		expectedName   string
		expectedQty    float64
		expectedUnit   string
	}{
		{"flour", "flour", 1, "count"},
		{"salt to taste", "salt to taste", 1, "count"},
		{"avocado oil, for cooking", "avocado oil, for cooking", 1, "count"},
		{"fresh parsley", "fresh parsley", 1, "count"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseIngredientString(tt.input)
			if err != nil {
				t.Errorf("ParseIngredientString(%q) unexpected error: %v", tt.input, err)
				return
			}
			if result.IngredientName != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, result.IngredientName)
			}
			if result.Quantity != tt.expectedQty {
				t.Errorf("expected quantity %f, got %f", tt.expectedQty, result.Quantity)
			}
			if result.Unit != tt.expectedUnit {
				t.Errorf("expected unit %q, got %q", tt.expectedUnit, result.Unit)
			}
		})
	}
}

func TestParseQuantity(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1", 1},
		{"2", 2},
		{"0.5", 0.5},
		{"1.5", 1.5},
		{"1/2", 0.5},
		{"1/4", 0.25},
		{"3/4", 0.75},
		{"100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseQuantity(tt.input)
			if err != nil {
				t.Fatalf("parseQuantity(%q) returned error: %v", tt.input, err)
			}

			delta := 0.0001
			if result < tt.expected-delta || result > tt.expected+delta {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseQuantity_Errors(t *testing.T) {
	tests := []string{
		"abc",
		"1/2/3",
		"1/0",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := parseQuantity(input)
			if err == nil {
				t.Errorf("parseQuantity(%q) expected error, got nil", input)
			}
		})
	}
}

func TestFindBestIntegerUnit_FractionalUnits(t *testing.T) {
	// Regression test: fractional units should display as fractional units,
	// not convert to smaller units (e.g., 1/2 tsp was previously showing as 49 drops)
	tests := []struct {
		name         string
		baseML       float64
		expectedUnit string
		expectedVal  float64
	}{
		// Fractional teaspoons
		{"half teaspoon", 2.46446, "half_tsp", 1},
		{"quarter teaspoon", 1.23223, "qtr_tsp", 1},
		{"eighth teaspoon", 0.616115, "eighth_tsp", 1},
		{"two half teaspoons", 4.92892, "tsp", 1},
		{"two quarter teaspoons", 2.46446, "half_tsp", 1},
		{"three eighth teaspoons", 1.848345, "eighth_tsp", 3},
		// Fractional cups
		{"half cup", 118.294, "half_cup", 1},
		{"quarter cup", 59.147, "qtr_cup", 1},
		{"two half cups", 236.588, "cup", 1},
		{"two quarter cups", 118.294, "half_cup", 1},
		{"three quarter cups", 177.441, "qtr_cup", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := FindBestIntegerUnit(tt.baseML, CategoryVolume)
			if q.Unit != tt.expectedUnit {
				t.Errorf("Unit: got %q, want %q", q.Unit, tt.expectedUnit)
			}
			if q.Value != tt.expectedVal {
				t.Errorf("Value: got %v, want %v", q.Value, tt.expectedVal)
			}
		})
	}
}

func TestRoundTripIngredientParsing(t *testing.T) {
	// Test that we can parse formatted output back
	tests := []struct {
		quantity float64
		unit     string
		name     string
	}{
		{2, "cup", "flour"},
		{1, "tsp", "salt"},
		{500, "g", "pasta"},
		{3, "tbsp", "oil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the formatted string
			q := Quantity{Value: tt.quantity, Unit: tt.unit, Category: CategoryVolume}
			if _, ok := MassUnits[tt.unit]; ok {
				q.Category = CategoryMass
			}
			formatted := Format(q) + " " + tt.name

			// Parse it back
			result, err := ParseIngredientString(formatted)
			if err != nil {
				t.Fatalf("ParseIngredientString(%q) returned error: %v", formatted, err)
			}

			if result.Quantity != tt.quantity {
				t.Errorf("Quantity: got %v, want %v", result.Quantity, tt.quantity)
			}
			if result.IngredientName != tt.name {
				t.Errorf("Name: got %q, want %q", result.IngredientName, tt.name)
			}
		})
	}
}
