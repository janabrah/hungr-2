package handlers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/models"
	"golang.org/x/net/html"
)

type ExtractRequest struct {
	URL string `json:"url"`
}

// Shared extraction rules used by all prompts
const extractionRules = `Rules:
1. ALWAYS start with a step that has an empty instruction "" containing ALL ingredients from the recipe. This must be the first element in the steps array.
2. Then add additional steps for the actual cooking instructions (these steps should have empty ingredients arrays since all ingredients are in the first step).
3. Format ingredients as "quantity unit ingredient" (e.g., "2 cups flour", "1 tsp salt", "3 eggs")
4. Use standard cooking units: tsp, tbsp, cup, oz, lb, g, kg, ml, l
5. For countable items without units, just use the number and name (e.g., "2 eggs", "1 onion")
6. Do NOT include temperatures (e.g., "350°F", "180°C") in the ingredients list - temperatures belong in the instruction steps only
7. IMPORTANT: Preserve ALL numbers in instructions including oven temperatures (e.g., "Preheat oven to 350°F"), cooking times (e.g., "bake for 25 minutes"), and quantities. Never omit or round these values.
8. IMPORTANT: Watch for mixed fractions! "3 1/2 cups" means 3.5 cups (three and a half), NOT "3" followed by "1/2 cup". Similarly "2 1/4 tsp" means 2.25 tsp. Convert mixed fractions to decimals.`

const extractImageSystemPrompt = `You are a recipe extraction assistant. Given an image of a recipe (such as a photo from a cookbook, a handwritten recipe card, or a screenshot), extract the recipe steps and ingredients into a structured JSON format.

` + extractionRules + `
9. If the image is unclear or partially visible, extract what you can see`

const extractURLSystemPrompt = `You are a recipe extraction assistant. Given the text content of a recipe webpage, extract the recipe steps and ingredients into a structured JSON format.

` + extractionRules

const extractTextSystemPrompt = `You are a recipe extraction assistant. Given raw text that has been copied and pasted from a recipe website (which may be poorly formatted, contain ads, navigation text, or other noise), extract the recipe steps and ingredients into a structured JSON format.

` + extractionRules + `
9. Ignore any non-recipe content like ads, navigation, comments, ratings, or author bios
10. If the text contains multiple recipes, extract only the main/first recipe`

// Shared JSON schema for recipe steps response
var recipeStepsSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"steps": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"instruction": map[string]interface{}{
						"type":        "string",
						"description": "The step instruction. Empty string if this is just an ingredients list.",
					},
					"ingredients": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":        "string",
							"description": "Ingredient in format 'quantity unit name' (e.g., '2 cups flour', '1 tsp salt', '3 eggs')",
						},
						"description": "Ingredients used in this step",
					},
				},
				"required":             []string{"instruction", "ingredients"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"steps"},
	"additionalProperties": false,
}

type openAIRequest struct {
	Model          string           `json:"model"`
	Messages       []openAIMessage  `json:"messages"`
	ResponseFormat *responseFormat  `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *jsonSchema `json:"json_schema,omitempty"`
}

type jsonSchema struct {
	Name   string      `json:"name"`
	Strict bool        `json:"strict"`
	Schema interface{} `json:"schema"`
}

type openAIMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string for text, []contentPart for vision
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// callOpenAI makes a request to OpenAI's chat completions API and parses the recipe response
func callOpenAI(apiKey, model string, messages []openAIMessage, timeout time.Duration) (*models.RecipeStepsResponse, error) {
	reqBody := openAIRequest{
		Model:    model,
		Messages: messages,
		ResponseFormat: &responseFormat{
			Type: "json_schema",
			JSONSchema: &jsonSchema{
				Name:   "recipe_steps",
				Strict: true,
				Schema: recipeStepsSchema,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var result models.RecipeStepsResponse
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %v", err)
	}

	return &result, nil
}

func ExtractRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		respondWithError(w, http.StatusBadRequest, "url is required")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Error(ctx, "OPENAI_API_KEY not set", fmt.Errorf("missing env var"))
		respondWithError(w, http.StatusInternalServerError, "recipe extraction not configured")
		return
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-5.2"
	}

	// Fetch the webpage
	logger.Info(ctx, "fetching URL", "url", req.URL)
	content, err := fetchAndExtractText(req.URL)
	if err != nil {
		logger.Error(ctx, "failed to fetch URL", err, "url", req.URL)
		respondWithError(w, http.StatusBadRequest, "failed to fetch URL: "+err.Error())
		return
	}

	// Truncate if too long
	if len(content) > 15000 {
		content = content[:15000]
	}
	logger.Info(ctx, "fetched content", "length", len(content))

	// Extract recipe using OpenAI
	logger.Info(ctx, "calling OpenAI API", "model", model)
	result, err := extractRecipeWithOpenAI(apiKey, model, content)
	if err != nil {
		logger.Error(ctx, "failed to extract recipe", err, "url", req.URL)
		respondWithError(w, http.StatusInternalServerError, "failed to extract recipe: "+err.Error())
		return
	}
	logger.Info(ctx, "OpenAI extraction complete", "steps", len(result.Steps))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func ExtractRecipeFromImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse multipart form (max 50MB for multiple images)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "at least one image file is required")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Error(ctx, "OPENAI_API_KEY not set", fmt.Errorf("missing env var"))
		respondWithError(w, http.StatusInternalServerError, "recipe extraction not configured")
		return
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-5.2"
	}

	var imageDataURLs []string
	for _, header := range files {
		contentType := header.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("file %s must be an image", header.Filename))
			return
		}

		file, err := header.Open()
		if err != nil {
			logger.Error(ctx, "failed to open image", err, "filename", header.Filename)
			respondWithError(w, http.StatusInternalServerError, "failed to read image")
			return
		}

		imageData, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			logger.Error(ctx, "failed to read image", err, "filename", header.Filename)
			respondWithError(w, http.StatusInternalServerError, "failed to read image")
			return
		}

		base64Image := base64.StdEncoding.EncodeToString(imageData)
		dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)
		imageDataURLs = append(imageDataURLs, dataURL)

		logger.Info(ctx, "processed image", "filename", header.Filename, "size", len(imageData))
	}

	logger.Info(ctx, "extracting recipe from images", "count", len(imageDataURLs), "model", model)

	// Extract recipe using OpenAI Vision
	result, err := extractRecipeFromImageWithOpenAI(apiKey, model, imageDataURLs)
	if err != nil {
		logger.Error(ctx, "failed to extract recipe from images", err)
		respondWithError(w, http.StatusInternalServerError, "failed to extract recipe: "+err.Error())
		return
	}
	logger.Info(ctx, "OpenAI image extraction complete", "steps", len(result.Steps))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func extractRecipeFromImageWithOpenAI(apiKey, model string, imageDataURLs []string) (*models.RecipeStepsResponse, error) {
	// Build user message with text prompt and images
	promptText := "Extract the recipe from these images. Include all ingredients and cooking steps you can see."
	if len(imageDataURLs) == 1 {
		promptText = "Extract the recipe from this image. Include all ingredients and cooking steps you can see."
	}

	userContent := []contentPart{
		{Type: "text", Text: promptText},
	}
	for _, dataURL := range imageDataURLs {
		userContent = append(userContent, contentPart{Type: "image_url", ImageURL: &imageURL{URL: dataURL}})
	}

	messages := []openAIMessage{
		{Role: "system", Content: extractImageSystemPrompt},
		{Role: "user", Content: userContent},
	}

	return callOpenAI(apiKey, model, messages, 90*time.Second)
}

func fetchAndExtractText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var textContent strings.Builder
	extractText(doc, &textContent)

	return textContent.String(), nil
}

func extractText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript", "iframe", "svg", "nav", "footer", "header":
			return
		}
	}

	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			sb.WriteString(text)
			sb.WriteString(" ")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, sb)
	}
}

func extractRecipeWithOpenAI(apiKey, model, content string) (*models.RecipeStepsResponse, error) {
	messages := []openAIMessage{
		{Role: "system", Content: extractURLSystemPrompt},
		{Role: "user", Content: fmt.Sprintf("Extract the recipe from this webpage content:\n\n%s", content)},
	}

	return callOpenAI(apiKey, model, messages, 60*time.Second)
}

type ExtractTextRequest struct {
	Text string `json:"text"`
}

func ExtractRecipeFromText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req ExtractTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Text == "" {
		respondWithError(w, http.StatusBadRequest, "text is required")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Error(ctx, "OPENAI_API_KEY not set", fmt.Errorf("missing env var"))
		respondWithError(w, http.StatusInternalServerError, "recipe extraction not configured")
		return
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-5.2"
	}

	// Truncate if too long
	text := req.Text
	if len(text) > 15000 {
		text = text[:15000]
	}
	logger.Info(ctx, "extracting recipe from text", "length", len(text), "model", model)

	// Extract recipe using OpenAI
	result, err := extractRecipeFromTextWithOpenAI(apiKey, model, text)
	if err != nil {
		logger.Error(ctx, "failed to extract recipe from text", err)
		respondWithError(w, http.StatusInternalServerError, "failed to extract recipe: "+err.Error())
		return
	}
	logger.Info(ctx, "OpenAI text extraction complete", "steps", len(result.Steps))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func extractRecipeFromTextWithOpenAI(apiKey, model, text string) (*models.RecipeStepsResponse, error) {
	messages := []openAIMessage{
		{Role: "system", Content: extractTextSystemPrompt},
		{Role: "user", Content: fmt.Sprintf("Extract the recipe from this text:\n\n%s", text)},
	}

	return callOpenAI(apiKey, model, messages, 60*time.Second)
}

func loadEnvFile(filename string) {
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

func init() {
	loadEnvFile(".env")
}
