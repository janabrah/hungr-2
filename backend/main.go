package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cobyabrahams/hungr/handlers"
	"github.com/cobyabrahams/hungr/logger"
	"github.com/cobyabrahams/hungr/middleware"
	"github.com/cobyabrahams/hungr/storage"
)

func main() {
	logger.Init()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	if err := storage.Init(dbURL); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/recipes", middleware.RequestLogger(middleware.CORS(handleRecipes, "GET, POST, DELETE, OPTIONS")))
	http.HandleFunc("/api/recipes/", middleware.RequestLogger(middleware.CORS(handleRecipeSubresources, "GET, PUT, PATCH, POST, OPTIONS")))
	http.HandleFunc("/api/files/", middleware.RequestLogger(middleware.CORS(handleFiles, "GET")))
	http.HandleFunc("/api/users", middleware.RequestLogger(middleware.CORS(handleUsers, "GET, POST, PUT, DELETE, OPTIONS")))
	http.HandleFunc("/api/auth/login", middleware.RequestLogger(middleware.CORS(handleLogin, "POST, OPTIONS")))
	http.HandleFunc("/api/extract-recipe", middleware.RequestLogger(middleware.CORS(handleExtractRecipe, "POST, OPTIONS")))
	http.HandleFunc("/api/extract-recipe-image", middleware.RequestLogger(middleware.CORS(handleExtractRecipeImage, "POST, OPTIONS")))
	http.HandleFunc("/api/extract-recipe-text", middleware.RequestLogger(middleware.CORS(handleExtractRecipeText, "POST, OPTIONS")))
	http.HandleFunc("/api/tags", middleware.RequestLogger(middleware.CORS(handleTags, "GET, OPTIONS")))
	http.HandleFunc("/api/connections", middleware.RequestLogger(middleware.CORS(handleConnections, "GET, POST, DELETE, OPTIONS")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info("server starting", "port", port)
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
	case "DELETE":
		handlers.DeleteRecipe(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRecipeSubresources(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/steps") {
		switch r.Method {
		case "GET":
			handlers.GetRecipeSteps(w, r)
		case "PUT":
			handlers.UpdateRecipeSteps(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasSuffix(r.URL.Path, "/files") {
		if r.Method == "POST" {
			handlers.AddRecipeFiles(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasSuffix(r.URL.Path, "/public") {
		switch r.Method {
		case "GET":
			handlers.GetPublicRecipe(w, r)
		case "POST":
			handlers.SetRecipePublic(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if r.Method == "PATCH" {
		handlers.PatchRecipe(w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handlers.GetFile(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetUser(w, r)
	case "POST":
		handlers.CreateUser(w, r)
	case "PUT":
		handlers.UpdateUser(w, r)
	case "DELETE":
		handlers.DeleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.Login(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleExtractRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.ExtractRecipe(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleExtractRecipeImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.ExtractRecipeFromImage(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleExtractRecipeText(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.ExtractRecipeFromText(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTags(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handlers.GetTags(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetConnections(w, r)
	case "POST":
		handlers.CreateConnection(w, r)
	case "DELETE":
		handlers.DeleteConnection(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
