package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cobyabrahams/hungr/storage"
)

const testUserEmail = "coby@hungr.com"

var recipes = []struct {
	name      string
	tags      string
	imageURLs []string
}{
	{
		name:      "Grandma's Chocolate Chip Cookies",
		tags:      "dessert, cookies, baking",
		imageURLs: []string{"https://picsum.photos/seed/cookies1/800/600", "https://picsum.photos/seed/cookies2/800/600"},
	},
	{
		name:      "Quick Weeknight Pasta",
		tags:      "dinner, pasta, quick",
		imageURLs: []string{"https://picsum.photos/seed/pasta/800/600"},
	},
	{
		name:      "Sunday Morning Pancakes",
		tags:      "breakfast, pancakes",
		imageURLs: []string{"https://picsum.photos/seed/pancakes1/800/600", "https://picsum.photos/seed/pancakes2/800/600", "https://picsum.photos/seed/pancakes3/800/600"},
	},
	{
		name:      "Spicy Thai Curry",
		tags:      "dinner, thai, spicy",
		imageURLs: []string{"https://picsum.photos/seed/curry/800/600"},
	},
	{
		name:      "Garden Salad",
		tags:      "lunch, salad, healthy, quick",
		imageURLs: []string{"https://picsum.photos/seed/salad/800/600"},
	},
	{
		name:      "Homemade Pizza Dough",
		tags:      "dinner, pizza, baking",
		imageURLs: []string{"https://picsum.photos/seed/pizza1/800/600", "https://picsum.photos/seed/pizza2/800/600"},
	},
	{
		name:      "Beef Stew",
		tags:      "dinner, soup, comfort food",
		imageURLs: []string{"https://picsum.photos/seed/stew/800/600"},
	},
	{
		name:      "Avocado Toast",
		tags:      "breakfast, quick, healthy",
		imageURLs: []string{"https://picsum.photos/seed/avocado/800/600"},
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
