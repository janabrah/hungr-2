package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cobyabrahams/hungr/storage"
	"github.com/gofrs/uuid"
)

var testUserUUID = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))

var recipes = []struct {
	name     string
	tags     string
	fileURLs []string
}{
	{
		name:     "Grandma's Chocolate Chip Cookies",
		tags:     "dessert, cookies, baking",
		fileURLs: []string{"https://picsum.photos/seed/cookies1/800/600", "https://picsum.photos/seed/cookies2/800/600"},
	},
	{
		name:     "Quick Weeknight Pasta",
		tags:     "dinner, pasta, quick",
		fileURLs: []string{"https://picsum.photos/seed/pasta/800/600"},
	},
	{
		name:     "Sunday Morning Pancakes",
		tags:     "breakfast, pancakes",
		fileURLs: []string{"https://picsum.photos/seed/pancakes1/800/600", "https://picsum.photos/seed/pancakes2/800/600", "https://picsum.photos/seed/pancakes3/800/600"},
	},
	{
		name:     "Spicy Thai Curry",
		tags:     "dinner, thai, spicy",
		fileURLs: []string{"https://picsum.photos/seed/curry/800/600"},
	},
	{
		name:     "Garden Salad",
		tags:     "lunch, salad, healthy, quick",
		fileURLs: []string{"https://picsum.photos/seed/salad/800/600"},
	},
	{
		name:     "Homemade Pizza Dough",
		tags:     "dinner, pizza, baking",
		fileURLs: []string{"https://picsum.photos/seed/pizza1/800/600", "https://picsum.photos/seed/pizza2/800/600"},
	},
	{
		name:     "Beef Stew",
		tags:     "dinner, soup, comfort food",
		fileURLs: []string{"https://picsum.photos/seed/stew/800/600"},
	},
	{
		name:     "Avocado Toast",
		tags:     "breakfast, quick, healthy",
		fileURLs: []string{"https://picsum.photos/seed/avocado/800/600"},
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

	for _, r := range recipes {
		recipe, err := storage.InsertRecipe(r.name, testUserUUID, r.tags)
		if err != nil {
			log.Printf("Failed to insert recipe %q: %v", r.name, err)
			continue
		}
		fmt.Printf("  Created recipe: %s\n", r.name)

		for i, url := range r.fileURLs {
			_, err := storage.InsertFile(recipe.UUID, url, i, true)
			if err != nil {
				log.Printf("    Failed to insert file: %v", err)
				continue
			}
		}
		fmt.Printf("    Added %d file(s)\n", len(r.fileURLs))

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
	}

	fmt.Printf("\nDone! Test user UUID: %s\n", testUserUUID)
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
