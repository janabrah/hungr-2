package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

const tagNamespace = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

func Init(connString string) error {
	var err error
	db, err = pgx.Connect(context.Background(), connString)
	return err
}

func GetRecipesByUserUUID(userUUID uuid.UUID) ([]models.Recipe, error) {
	rows, err := db.Query(context.Background(),
		`SELECT uuid, filename, user_uuid, tag_string, created_at
		 FROM recipes WHERE user_uuid = $1
		 ORDER BY created_at DESC LIMIT 100`, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.UUID, &r.Filename, &r.User, &r.TagString, &r.CreatedAt); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, rows.Err()
}

func GetFileMappingsByRecipeUUIDs(recipeUUIDs []uuid.UUID) ([]models.FileRecipe, error) {
	if len(recipeUUIDs) == 0 {
		return []models.FileRecipe{}, nil
	}

	rows, err := db.Query(context.Background(),
		`SELECT file_uuid, recipe_uuid, page_number
		 FROM file_recipes WHERE recipe_uuid = ANY($1)`, recipeUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []models.FileRecipe
	for rows.Next() {
		var m models.FileRecipe
		if err := rows.Scan(&m.FileUUID, &m.RecipeUUID, &m.PageNumber); err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}
	return mappings, rows.Err()
}

func GetFilesByUUIDs(fileUUIDs []uuid.UUID) ([]models.File, error) {
	if len(fileUUIDs) == 0 {
		return []models.File{}, nil
	}

	rows, err := db.Query(context.Background(),
		`SELECT uuid, url, image FROM files WHERE uuid = ANY($1)`, fileUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(&f.UUID, &f.URL, &f.Image); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func InsertRecipe(filename string, user uuid.UUID, tagString string) (*models.Recipe, error) {
	var r models.Recipe
	err := db.QueryRow(context.Background(),
		`INSERT INTO recipes (filename, user_uuid, tag_string)
		 VALUES ($1, $2, $3)
		 RETURNING uuid, filename, user_uuid, tag_string, created_at`,
		filename, user, tagString).Scan(&r.UUID, &r.Filename, &r.User, &r.TagString, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func InsertFile(url string, isImage bool) (*models.File, error) {
	var f models.File
	err := db.QueryRow(context.Background(),
		`INSERT INTO files (url, image) VALUES ($1, $2) RETURNING uuid, url, image`,
		url, isImage).Scan(&f.UUID, &f.URL, &f.Image)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func InsertFileRecipe(fileUUID, recipeUUID uuid.UUID, pageNumber int) error {
	_, err := db.Exec(context.Background(),
		`INSERT INTO file_recipes (file_uuid, recipe_uuid, page_number) VALUES ($1, $2, $3)`,
		fileUUID, recipeUUID, pageNumber)
	return err
}

func UpsertTag(tagUUID uuid.UUID, name string) (*models.Tag, error) {
	var t models.Tag
	err := db.QueryRow(context.Background(),
		`INSERT INTO tags (uuid, name) VALUES ($1, $2)
		 ON CONFLICT (uuid) DO UPDATE SET name = EXCLUDED.name
		 RETURNING uuid, name`,
		tagUUID, name).Scan(&t.UUID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func InsertRecipeTag(recipeUUID, tagUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`INSERT INTO recipe_tags (recipe_uuid, tag_uuid) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		recipeUUID, tagUUID)
	return err
}

func UploadFile(filename string, file io.Reader) (string, error) {
	// For now, return a placeholder URL
	// TODO: Implement actual file storage (S3, local disk, etc.)
	return fmt.Sprintf("/uploads/%s", filename), nil
}

func CreateTagUUID(tag string) uuid.UUID {
	namespace := uuid.Must(uuid.FromString(tagNamespace))
	return uuid.NewV5(namespace, tag)
}
