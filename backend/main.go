package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cobyabrahams/hungr/handlers"
)

func main() {
	// TODO: Load environment variables (SUPABASE_URL, SUPABASE_SERVICE_KEY)

	// TODO: Initialize Supabase client

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/recipes", handleRecipes)

	port := "8080"
	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleRecipes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetRecipes(w, r)
	case "POST":
		handlers.CreateRecipe(w, r)
	case "OPTIONS":
		// TODO: Handle CORS preflight
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
