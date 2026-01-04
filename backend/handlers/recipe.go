package handlers

import (
	"net/http"

	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
)

// Silence unused import errors until implemented
var _ = storage.Init
var _ = models.Recipe{}

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Parse user_uuid from r.URL.Query().Get("user_uuid")
	// 2. Call storage.GetRecipesByUserUUID(userUUID)
	// 3. Extract recipe UUIDs
	// 4. Call storage.GetFileMappingsByRecipeUUIDs(recipeUUIDs)
	// 5. Extract file UUIDs from mappings
	// 6. Call storage.GetFilesByUUIDs(fileUUIDs)
	// 7. Build models.RecipesResponse
	// 8. json.NewEncoder(w).Encode(response)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Parse filename and tagString from r.URL.Query()
	// 2. Parse multipart form: r.ParseMultipartForm(32 << 20)
	// 3. Get files: r.MultipartForm.File["file"]
	// 4. For each file:
	//    a. Open file: fileHeader.Open()
	//    b. Upload: storage.UploadFile(filename+pageNum, file)
	//    c. Collect URLs
	// 5. Insert recipe: storage.InsertRecipe(filename, userUUID, tagString)
	// 6. For each URL:
	//    a. Insert file: storage.InsertFile(url, true)
	//    b. Insert mapping: storage.InsertFileRecipe(fileUUID, recipeUUID, pageNum)
	// 7. Parse tags: strings.Split(tagString, ", ")
	// 8. For each tag:
	//    a. Generate UUID: storage.CreateTagUUID(tag)
	//    b. Upsert tag: storage.UpsertTag(uuid, name)
	//    c. Insert mapping: storage.InsertRecipeTag(recipeUUID, tagUUID)
	// 9. Build models.UploadResponse
	// 10. json.NewEncoder(w).Encode(response)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	// TODO: Helper to return JSON error responses
	http.Error(w, message, code)
}
