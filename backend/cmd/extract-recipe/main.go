package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type RecipeStepsResponse struct {
	Steps []RecipeStepResponse `json:"steps"`
	Tags  []string             `json:"tags"`
}

type RecipeStepResponse struct {
	Instruction string   `json:"instruction"`
	Ingredients []string `json:"ingredients"`
}

type OpenAIRequest struct {
	Model          string          `json:"model"`
	Messages       []OpenAIMessage `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

type JSONSchema struct {
	Name   string      `json:"name"`
	Strict bool        `json:"strict"`
	Schema interface{} `json:"schema"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // .env file is optional
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

const systemPrompt = `You are a recipe extraction assistant. Given the text content of a recipe webpage, extract the recipe steps and ingredients into the JSON schema provided.

Rules:
1. ALWAYS start with a step that has an empty instruction "" containing ALL ingredients from the recipe. This must be the first element in the steps array.
2. Then add additional steps for the actual cooking instructions (these steps should have empty ingredients arrays since all ingredients are in the first step).
3. Format ingredients as "quantity unit ingredient" (e.g., "2 cups flour", "1 tsp salt", "3 eggs")
4. Use standard cooking units: tsp, tbsp, cup, oz, lb, g, kg, ml, l
5. For countable items without units, just use the number and name (e.g., "2 eggs", "1 onion")
6. Do NOT include temperatures (e.g., "350°F", "180°C") - these are not ingredients
7. Also return a "tags" array with 3-8 concise, lowercase tags (1-3 words each). Include at least one role tag (e.g., "dinner", "appetizer", "breakfast", "dessert") and 1-2 tags for major ingredients (e.g., "chicken", "salmon", "mushroom")
8. Return ONLY valid JSON, no markdown formatting or explanation`

func main() {
	loadEnvFile(".env")

	if len(os.Args) < 2 {
		fmt.Println("Usage: extract-recipe <url>")
		fmt.Println("Environment: OPENAI_API_KEY must be set (or in .env file)")
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable must be set")
	}

	url := os.Args[1]

	// Fetch the webpage
	fmt.Fprintf(os.Stderr, "Fetching %s...\n", url)
	content, err := fetchAndExtractText(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	// Truncate if too long (OpenAI has token limits)
	if len(content) > 15000 {
		content = content[:15000]
	}

	fmt.Fprintf(os.Stderr, "Extracted %d characters of text\n", len(content))

	// Send to OpenAI
	fmt.Fprintf(os.Stderr, "Sending to OpenAI...\n")
	result, err := extractRecipeWithOpenAI(apiKey, content)
	if err != nil {
		log.Fatalf("Failed to extract recipe: %v", err)
	}

	// Output the JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}

	fmt.Println(string(output))
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
	// Skip script, style, and other non-content tags
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

func extractRecipeWithOpenAI(apiKey, content string) (*RecipeStepsResponse, error) {
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
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":        "string",
					"description": "Concise tag describing the recipe (e.g., 'dessert', 'gluten-free', 'quick')",
				},
				"description": "Suggested tags for the recipe",
			},
		},
		"required":             []string{"steps", "tags"},
		"additionalProperties": false,
	}

	reqBody := OpenAIRequest{
		Model: "gpt-5.2",
		Messages: []OpenAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Extract the recipe from this webpage content:\n\n%s", content)},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_schema",
			JSONSchema: &JSONSchema{
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the JSON from OpenAI's response
	responseContent := openAIResp.Choices[0].Message.Content

	// Strip markdown code blocks if present
	responseContent = strings.TrimPrefix(responseContent, "```json")
	responseContent = strings.TrimPrefix(responseContent, "```")
	responseContent = strings.TrimSuffix(responseContent, "```")
	responseContent = strings.TrimSpace(responseContent)

	var result RecipeStepsResponse
	if err := json.Unmarshal([]byte(responseContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %v\nResponse was: %s", err, responseContent)
	}

	return &result, nil
}
