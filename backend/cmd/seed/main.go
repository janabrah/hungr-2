package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cobyabrahams/hungr/storage"
	"github.com/cobyabrahams/hungr/units"
)

const testUserEmail = "coby@hungr.com"

type recipeStep struct {
	instruction string
	ingredients []string
}

var recipes = []struct {
	name      string
	tags      string
	imageURLs []string
	steps     []recipeStep
}{
	{
		name:      "Grandma's Chocolate Chip Cookies",
		tags:      "dessert, cookies, baking",
		imageURLs: []string{"https://picsum.photos/seed/cookies1/800/600", "https://picsum.photos/seed/cookies2/800/600"},
		steps: []recipeStep{
			{instruction: "Preheat oven to 375°F (190°C)", ingredients: []string{}},
			{instruction: "Cream together butter and sugars until fluffy", ingredients: []string{"1 cup butter", "3/4 cup sugar", "3/4 cup brown sugar"}},
			{instruction: "Beat in eggs and vanilla", ingredients: []string{"2 eggs", "1 tsp vanilla extract"}},
			{instruction: "Mix in flour, baking soda, and salt", ingredients: []string{"2 cups flour", "1 tsp baking soda", "1 tsp salt"}},
			{instruction: "Fold in chocolate chips", ingredients: []string{"2 cups chocolate chips"}},
			{instruction: "Drop rounded tablespoons onto baking sheets and bake for 9-11 minutes", ingredients: []string{}},
		},
	},
	{
		name:      "Quick Weeknight Pasta",
		tags:      "dinner, pasta, quick",
		imageURLs: []string{"https://picsum.photos/seed/pasta/800/600"},
		steps: []recipeStep{
			{instruction: "Bring a large pot of salted water to boil", ingredients: []string{"1 tbsp salt"}},
			{instruction: "Cook pasta according to package directions", ingredients: []string{"1 lb pasta"}},
			{instruction: "While pasta cooks, sauté garlic in olive oil", ingredients: []string{"4 cloves garlic", "3 tbsp olive oil"}},
			{instruction: "Add tomatoes and simmer for 5 minutes", ingredients: []string{"1 can crushed tomatoes"}},
			{instruction: "Toss pasta with sauce and top with parmesan", ingredients: []string{"1/2 cup parmesan cheese", "1/4 tsp red pepper flakes"}},
		},
	},
	{
		name:      "Sunday Morning Pancakes",
		tags:      "breakfast, pancakes",
		imageURLs: []string{"https://picsum.photos/seed/pancakes1/800/600", "https://picsum.photos/seed/pancakes2/800/600", "https://picsum.photos/seed/pancakes3/800/600"},
		steps: []recipeStep{
			{instruction: "Mix dry ingredients in a large bowl", ingredients: []string{"2 cups flour", "2 tbsp sugar", "2 tsp baking powder", "1 tsp salt"}},
			{instruction: "Whisk wet ingredients in another bowl", ingredients: []string{"2 cups milk", "2 eggs", "1/4 cup melted butter"}},
			{instruction: "Combine wet and dry ingredients, mixing until just combined", ingredients: []string{}},
			{instruction: "Heat griddle to 350°F and grease lightly", ingredients: []string{"1 tbsp butter"}},
			{instruction: "Pour 1/4 cup batter per pancake, flip when bubbles form", ingredients: []string{}},
			{instruction: "Serve with maple syrup and fresh berries", ingredients: []string{"1/2 cup maple syrup", "1 cup mixed berries"}},
		},
	},
	{
		name:      "Spicy Thai Curry",
		tags:      "dinner, thai, spicy",
		imageURLs: []string{"https://picsum.photos/seed/curry/800/600"},
		steps: []recipeStep{
			{instruction: "Heat oil in a wok over high heat", ingredients: []string{"2 tbsp coconut oil"}},
			{instruction: "Sauté curry paste until fragrant", ingredients: []string{"3 tbsp red curry paste"}},
			{instruction: "Add coconut milk and bring to simmer", ingredients: []string{"1 can coconut milk"}},
			{instruction: "Add chicken and cook through", ingredients: []string{"1 lb chicken breast"}},
			{instruction: "Add vegetables and simmer until tender", ingredients: []string{"1 cup bell peppers", "1 cup bamboo shoots", "1/2 cup thai basil"}},
			{instruction: "Season with fish sauce and serve over rice", ingredients: []string{"2 tbsp fish sauce", "1 tbsp sugar", "2 cups jasmine rice"}},
		},
	},
	{
		name:      "Garden Salad",
		tags:      "lunch, salad, healthy, quick",
		imageURLs: []string{"https://picsum.photos/seed/salad/800/600"},
		steps: []recipeStep{
			{instruction: "Wash and dry lettuce, tear into pieces", ingredients: []string{"1 head romaine lettuce"}},
			{instruction: "Slice vegetables", ingredients: []string{"1 cucumber", "1 cup cherry tomatoes", "1/2 red onion"}},
			{instruction: "Make dressing by whisking ingredients together", ingredients: []string{"3 tbsp olive oil", "1 tbsp red wine vinegar", "1 tsp dijon mustard"}},
			{instruction: "Toss salad with dressing and top with cheese", ingredients: []string{"1/4 cup feta cheese", "1/4 cup croutons"}},
		},
	},
	{
		name:      "Homemade Pizza Dough",
		tags:      "dinner, pizza, baking",
		imageURLs: []string{"https://picsum.photos/seed/pizza1/800/600", "https://picsum.photos/seed/pizza2/800/600"},
		steps: []recipeStep{
			{instruction: "Activate yeast in warm water with sugar", ingredients: []string{"1 cup warm water", "1 tsp sugar", "1 packet active dry yeast"}},
			{instruction: "Mix flour and salt in a large bowl", ingredients: []string{"3 cups flour", "1 tsp salt"}},
			{instruction: "Add yeast mixture and olive oil, knead for 10 minutes", ingredients: []string{"2 tbsp olive oil"}},
			{instruction: "Let dough rise in oiled bowl for 1 hour", ingredients: []string{}},
			{instruction: "Punch down, divide in two, and shape into rounds", ingredients: []string{}},
			{instruction: "Add toppings and bake at 450°F for 12-15 minutes", ingredients: []string{"1 cup pizza sauce", "2 cups mozzarella cheese"}},
		},
	},
	{
		name:      "Beef Stew",
		tags:      "dinner, soup, comfort food",
		imageURLs: []string{"https://picsum.photos/seed/stew/800/600"},
		steps: []recipeStep{
			{instruction: "Cut beef into 1-inch cubes and season", ingredients: []string{"2 lb beef chuck", "1 tsp salt", "1/2 tsp black pepper"}},
			{instruction: "Brown beef in batches in a dutch oven", ingredients: []string{"2 tbsp olive oil"}},
			{instruction: "Sauté onions, carrots, and celery", ingredients: []string{"1 onion", "3 carrots", "3 celery stalks"}},
			{instruction: "Add tomato paste and flour, stir for 1 minute", ingredients: []string{"2 tbsp tomato paste", "2 tbsp flour"}},
			{instruction: "Add beef broth, potatoes, and herbs", ingredients: []string{"4 cups beef broth", "1 lb potatoes", "2 bay leaves", "1 tsp thyme"}},
			{instruction: "Simmer for 2 hours until beef is tender", ingredients: []string{}},
		},
	},
	{
		name:      "Avocado Toast",
		tags:      "breakfast, quick, healthy",
		imageURLs: []string{"https://picsum.photos/seed/avocado/800/600"},
		steps: []recipeStep{
			{instruction: "Toast bread until golden", ingredients: []string{"2 slices sourdough bread"}},
			{instruction: "Mash avocado with lime juice and salt", ingredients: []string{"1 avocado", "1 tbsp lime juice", "1/4 tsp salt"}},
			{instruction: "Spread avocado on toast", ingredients: []string{}},
			{instruction: "Top with eggs, red pepper flakes, and microgreens", ingredients: []string{"2 eggs", "1/4 tsp red pepper flakes", "1/4 cup microgreens"}},
		},
	},
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	if err := storage.Init(dbURL); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Seeding database...")

	// Create test user if not exists
	_, err := storage.GetUserByEmail(testUserEmail)
	if err != nil {
		_, err = storage.CreateUser(testUserEmail, "Test User")
		if err != nil {
			log.Fatal("Failed to create test user: ", err)
		}
		fmt.Println("Created test user:", testUserEmail)
	}

	for _, r := range recipes {
		recipe, err := storage.InsertRecipeByEmail(r.name, testUserEmail, r.tags)
		if err != nil {
			log.Printf("Failed to insert recipe %q: %v", r.name, err)
			continue
		}
		fmt.Printf("  Created recipe: %s\n", r.name)

		for i, url := range r.imageURLs {
			data, contentType, err := fetchImage(url)
			if err != nil {
				log.Printf("    Failed to fetch image %s: %v", url, err)
				continue
			}

			_, err = storage.InsertFile(recipe.UUID, data, contentType, i, true)
			if err != nil {
				log.Printf("    Failed to insert file: %v", err)
				continue
			}
		}
		fmt.Printf("    Added %d file(s)\n", len(r.imageURLs))

		tags := splitTags(r.tags)
		for _, tagName := range tags {
			tagUUID := storage.CreateTagUUID(tagName)
			_, err := storage.UpsertTag(tagUUID, tagName)
			if err != nil {
				log.Printf("    Failed to upsert tag %q: %v", tagName, err)
				continue
			}
			if err := storage.InsertRecipeTag(recipe.UUID, tagUUID); err != nil {
				log.Printf("    Failed to link tag %q: %v", tagName, err)
			}
		}

		// Add recipe steps and ingredients
		if len(r.steps) > 0 {
			steps := make([]storage.StepInput, len(r.steps))
			for i, s := range r.steps {
				ingredients := make([]storage.IngredientInput, len(s.ingredients))
				for j, ingStr := range s.ingredients {
					parsed, err := parseIngredient(ingStr)
					if err != nil {
						log.Printf("    Failed to parse ingredient %q: %v", ingStr, err)
						continue
					}
					ingredients[j] = parsed
				}
				steps[i] = storage.StepInput{
					Instruction: s.instruction,
					Ingredients: ingredients,
				}
			}

			if err := storage.ReplaceRecipeSteps(recipe.UUID, steps); err != nil {
				log.Printf("    Failed to add steps: %v", err)
			} else {
				fmt.Printf("    Added %d step(s)\n", len(steps))
			}
		}
	}

	fmt.Printf("\nDone! Test user email: %s\n", testUserEmail)
}

func fetchImage(url string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	return data, contentType, nil
}

func splitTags(tagString string) []string {
	var tags []string
	current := ""
	for _, c := range tagString {
		if c == ',' {
			if t := trim(current); t != "" {
				tags = append(tags, t)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if t := trim(current); t != "" {
		tags = append(tags, t)
	}
	return tags
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

func parseIngredient(s string) (storage.IngredientInput, error) {
	parsed, err := units.ParseIngredientString(s)
	if err != nil {
		return storage.IngredientInput{}, err
	}
	return storage.IngredientInput{
		Name:     parsed.IngredientName,
		Unit:     parsed.Unit,
		Quantity: parsed.Quantity,
	}, nil
}
