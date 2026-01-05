package main

import (
	"os"
	"testing"
)

func skipUnlessExpensiveTests(t *testing.T) {
	if os.Getenv("RUN_EXPENSIVE_TESTS") == "" {
		t.Skip("Skipping expensive test. Set RUN_EXPENSIVE_TESTS=1 to run.")
	}
}

func TestExtractRecipeFromURL(t *testing.T) {
	skipUnlessExpensiveTests(t)
	loadEnvFile(".env")
	loadEnvFile("../../.env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_API_KEY must be set")
	}

	url := "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/"

	content, err := fetchAndExtractText(url)
	if err != nil {
		t.Fatalf("Failed to fetch URL: %v", err)
	}

	if len(content) < 100 {
		t.Fatalf("Content too short: %d chars", len(content))
	}

	result, err := extractRecipeWithOpenAI(apiKey, content)
	if err != nil {
		t.Fatalf("Failed to extract recipe: %v", err)
	}

	if len(result.Steps) == 0 {
		t.Fatal("Expected at least one step")
	}

	// Check that we got some ingredients
	totalIngredients := 0
	for _, step := range result.Steps {
		totalIngredients += len(step.Ingredients)
	}

	if totalIngredients < 3 {
		t.Errorf("Expected at least 3 ingredients, got %d", totalIngredients)
	}

	// Log the result for manual inspection
	t.Logf("Extracted %d steps with %d total ingredients", len(result.Steps), totalIngredients)
	for i, step := range result.Steps {
		t.Logf("Step %d: %s (%d ingredients)", i+1, step.Instruction, len(step.Ingredients))
	}
}

func TestFetchAndExtractText(t *testing.T) {
	// This test is cheap - just fetches a webpage
	url := "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/"

	content, err := fetchAndExtractText(url)
	if err != nil {
		t.Fatalf("Failed to fetch URL: %v", err)
	}

	if len(content) < 500 {
		t.Errorf("Content seems too short: %d chars", len(content))
	}

	// Should contain some recipe-related words
	keywords := []string{"cookie", "chocolate", "butter", "sugar", "flour"}
	found := 0
	for _, kw := range keywords {
		if containsIgnoreCase(content, kw) {
			found++
		}
	}

	if found < 3 {
		t.Errorf("Expected to find at least 3 recipe keywords, found %d", found)
	}
}

func containsIgnoreCase(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			sc := s[i+j]
			pc := substr[j]
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if pc >= 'A' && pc <= 'Z' {
				pc += 32
			}
			if sc != pc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func TestLoadEnvFile(t *testing.T) {
	// Just test that it doesn't crash on missing file
	loadEnvFile("nonexistent.env")
}
