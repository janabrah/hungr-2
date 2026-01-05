package handlers

import (
	"bufio"
	"bytes"
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
	Role    string `json:"role"`
	Content string `json:"content"`
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

const extractSystemPrompt = `You are a recipe extraction assistant. Given the text content of a recipe webpage, extract the recipe steps and ingredients into a structured JSON format.

Rules:
1. ALWAYS start with a step that has an empty instruction "" containing ALL ingredients from the recipe. This must be the first element in the steps array.
2. Then add additional steps for the actual cooking instructions (these steps should have empty ingredients arrays since all ingredients are in the first step).
3. Format ingredients as "quantity unit ingredient" (e.g., "2 cups flour", "1 tsp salt", "3 eggs")
4. Use standard cooking units: tsp, tbsp, cup, oz, lb, g, kg, ml, l
5. For countable items without units, just use the number and name (e.g., "2 eggs", "1 onion")
6. Do NOT include temperatures (e.g., "350°F", "180°C") - these are not ingredients`

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
	schema := map[string]interface{}{
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

	reqBody := openAIRequest{
		Model: model,
		Messages: []openAIMessage{
			{Role: "system", Content: extractSystemPrompt},
			{Role: "user", Content: fmt.Sprintf("Extract the recipe from this webpage content:\n\n%s", content)},
		},
		ResponseFormat: &responseFormat{
			Type: "json_schema",
			JSONSchema: &jsonSchema{
				Name:   "recipe_steps",
				Strict: true,
				Schema: schema,
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

	client := &http.Client{Timeout: 60 * time.Second}
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

	responseContent := openAIResp.Choices[0].Message.Content

	var result models.RecipeStepsResponse
	if err := json.Unmarshal([]byte(responseContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %v", err)
	}

	return &result, nil
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
