package storage

import (
	"context"
	"fmt"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

const tagNamespace = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

func Init(connString string) error {
	var err error
	db, err = pgxpool.New(context.Background(), connString)
	return err
}

func GetRecipesByUserEmail(email string) ([]models.Recipe, error) {
	rows, err := db.Query(context.Background(),
		`SELECT r.uuid, r.name, r.user_uuid, r.tag_string, r.created_at
		 FROM recipes r
		 JOIN users u ON r.user_uuid = u.uuid
		 WHERE u.email = $1
		 ORDER BY r.created_at DESC LIMIT 100`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := []models.Recipe{}
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.UUID, &r.Name, &r.User, &r.TagString, &r.CreatedAt); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, rows.Err()
}

func GetFilesByRecipeUUIDs(recipeUUIDs []uuid.UUID) ([]models.File, error) {
	if len(recipeUUIDs) == 0 {
		return []models.File{}, nil
	}

	rows, err := db.Query(context.Background(),
		`SELECT uuid, recipe_uuid, url, page_number, image
		 FROM files WHERE recipe_uuid = ANY($1)
		 ORDER BY page_number`, recipeUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(&f.UUID, &f.RecipeUUID, &f.URL, &f.PageNumber, &f.Image); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func InsertRecipeByEmail(name string, email string, tagString string) (*models.Recipe, error) {
	var r models.Recipe
	err := db.QueryRow(context.Background(),
		`INSERT INTO recipes (name, user_uuid, tag_string)
		 SELECT $1, u.uuid, $2
		 FROM users u WHERE u.email = $3
		 RETURNING uuid, name, user_uuid, tag_string, created_at`,
		name, tagString, email).Scan(&r.UUID, &r.Name, &r.User, &r.TagString, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func InsertFile(recipeUUID uuid.UUID, data []byte, contentType string, pageNumber int, isImage bool) (*models.File, error) {
	var f models.File
	var fileUUID uuid.UUID
	err := db.QueryRow(context.Background(),
		`INSERT INTO files (recipe_uuid, data, content_type, url, page_number, image)
		 VALUES ($1, $2, $3, '', $4, $5)
		 RETURNING uuid, recipe_uuid, page_number, image`,
		recipeUUID, data, contentType, pageNumber, isImage).Scan(&fileUUID, &f.RecipeUUID, &f.PageNumber, &f.Image)
	if err != nil {
		return nil, err
	}
	f.UUID = fileUUID
	f.URL = fmt.Sprintf("/api/files/%s", fileUUID.String())
	f.Image = isImage

	_, err = db.Exec(context.Background(),
		`UPDATE files SET url = $1 WHERE uuid = $2`, f.URL, fileUUID)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func GetFileData(fileUUID uuid.UUID) ([]byte, string, error) {
	var data []byte
	var contentType string
	err := db.QueryRow(context.Background(),
		`SELECT data, content_type FROM files WHERE uuid = $1`, fileUUID).Scan(&data, &contentType)
	if err != nil {
		return nil, "", err
	}
	return data, contentType, nil
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

func DeleteRecipe(recipeUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM recipe_tags WHERE recipe_uuid = $1`, recipeUUID)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(),
		`DELETE FROM files WHERE recipe_uuid = $1`, recipeUUID)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(),
		`DELETE FROM recipes WHERE uuid = $1`, recipeUUID)
	return err
}

func CreateTagUUID(tag string) uuid.UUID {
	namespace := uuid.Must(uuid.FromString(tagNamespace))
	return uuid.NewV5(namespace, tag)
}

// User storage functions

func GetUserByUUID(userUUID uuid.UUID) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(),
		`SELECT uuid, email, name, created_at
		 FROM users WHERE uuid = $1`, userUUID).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(),
		`SELECT uuid, email, name, created_at
		 FROM users WHERE email = $1`, email).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(email, name string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(),
		`INSERT INTO users (email, name)
		 VALUES ($1, $2)
		 RETURNING uuid, email, name, created_at`,
		email, name).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(userUUID uuid.UUID, name string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(),
		`UPDATE users SET name = $1
		 WHERE uuid = $2
		 RETURNING uuid, email, name, created_at`,
		name, userUUID).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func DeleteUser(userUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM users WHERE uuid = $1`, userUUID)
	return err
}
