package handlers

import (
	"bufio"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func skipUnlessExpensiveTests(t *testing.T) {
	if os.Getenv("RUN_EXPENSIVE_TESTS") == "" {
		t.Skip("Skipping expensive test. Set RUN_EXPENSIVE_TESTS=1 to run.")
	}
}

func loadEnvFileForTest(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func TestExtractRecipeFromImage_Expensive(t *testing.T) {
	skipUnlessExpensiveTests(t)
	loadEnvFileForTest(".env")
	loadEnvFileForTest("../.env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_API_KEY must be set")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-5.2"
	}

	// Fetch a recipe image from Unsplash (reliable public access)
	// Image of chocolate chip cookies
	imageURL := "https://images.unsplash.com/photo-1499636136210-6f4ee915583e?w=800"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", imageURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RecipeTest/1.0)")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to fetch test image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to fetch test image: HTTP %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read image data: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	dataURL := "data:" + contentType + ";base64," + base64Image

	t.Logf("Fetched image: %d bytes, content-type: %s", len(imageData), contentType)

	// Call the extraction function
	result, err := extractRecipeFromImageWithOpenAI(apiKey, model, []string{dataURL})
	if err != nil {
		t.Fatalf("Failed to extract recipe: %v", err)
	}

	// Verify we got some results
	if len(result.Steps) == 0 {
		t.Fatal("Expected at least one step")
	}

	// Count total ingredients
	totalIngredients := 0
	for _, step := range result.Steps {
		totalIngredients += len(step.Ingredients)
	}

	t.Logf("Extracted %d steps with %d total ingredients", len(result.Steps), totalIngredients)

	// Log steps for inspection
	for i, step := range result.Steps {
		if step.Instruction == "" {
			t.Logf("Step %d (ingredients only): %d ingredients", i+1, len(step.Ingredients))
			for _, ing := range step.Ingredients {
				t.Logf("  - %s", ing)
			}
		} else {
			t.Logf("Step %d: %s", i+1, step.Instruction)
		}
	}
}

func TestExtractRecipeFromMultipleImages_Expensive(t *testing.T) {
	skipUnlessExpensiveTests(t)
	loadEnvFileForTest(".env")
	loadEnvFileForTest("../.env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_API_KEY must be set")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-5.2"
	}

	// Use two different food images from Unsplash
	imageURLs := []string{
		"https://images.unsplash.com/photo-1499636136210-6f4ee915583e?w=800", // cookies
		"https://images.unsplash.com/photo-1558961363-fa8fdf82db35?w=800",    // cookies on baking sheet
	}

	client := &http.Client{}
	var dataURLs []string
	for _, url := range imageURLs {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RecipeTest/1.0)")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to fetch image %s: %v", url, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("Failed to fetch image %s: HTTP %d", url, resp.StatusCode)
		}

		imageData, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read image data: %v", err)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		base64Image := base64.StdEncoding.EncodeToString(imageData)
		dataURLs = append(dataURLs, "data:"+contentType+";base64,"+base64Image)

		t.Logf("Fetched image: %d bytes", len(imageData))
	}

	// Call the extraction function with multiple images
	result, err := extractRecipeFromImageWithOpenAI(apiKey, model, dataURLs)
	if err != nil {
		t.Fatalf("Failed to extract recipe: %v", err)
	}

	// Verify we got some results
	if len(result.Steps) == 0 {
		t.Fatal("Expected at least one step")
	}

	// Count total ingredients
	totalIngredients := 0
	for _, step := range result.Steps {
		totalIngredients += len(step.Ingredients)
	}

	t.Logf("Extracted %d steps with %d total ingredients from %d images", len(result.Steps), totalIngredients, len(dataURLs))

	// Log steps for inspection
	for i, step := range result.Steps {
		if step.Instruction == "" {
			t.Logf("Step %d (ingredients only): %d ingredients", i+1, len(step.Ingredients))
		} else {
			t.Logf("Step %d: %s", i+1, step.Instruction)
		}
	}
}
