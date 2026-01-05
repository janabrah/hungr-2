package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cobyabrahams/hungr/handlers"
	"github.com/cobyabrahams/hungr/storage"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	if err := storage.Init(dbURL); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/recipes", handleRecipes)
	http.HandleFunc("/api/files/", handleFiles)
	http.HandleFunc("/api/users", handleUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case "GET":
		handlers.GetRecipes(w, r)
	case "POST":
		handlers.CreateRecipe(w, r)
	case "DELETE":
		handlers.DeleteRecipe(w, r)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		handlers.GetFile(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case "GET":
		handlers.GetUser(w, r)
	case "POST":
		handlers.CreateUser(w, r)
	case "PUT":
		handlers.UpdateUser(w, r)
	case "DELETE":
		handlers.DeleteUser(w, r)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
