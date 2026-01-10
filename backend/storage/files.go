package storage

import (
	"context"
	"fmt"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const (
	queryGetFilesByRecipeUUIDs = `
		SELECT uuid, recipe_uuid, url, page_number, image
		FROM files WHERE recipe_uuid = ANY($1)
		ORDER BY page_number`

	queryInsertFile = `
		INSERT INTO files (recipe_uuid, data, content_type, url, page_number, image)
		VALUES ($1, $2, $3, '', $4, $5)
		RETURNING uuid, recipe_uuid, page_number, image`

	queryUpdateFileURL = `UPDATE files SET url = $1 WHERE uuid = $2`

	queryGetFileData = `SELECT data, content_type FROM files WHERE uuid = $1`
)

func GetFilesByRecipeUUIDs(recipeUUIDs []uuid.UUID) ([]models.File, error) {
	if len(recipeUUIDs) == 0 {
		return []models.File{}, nil
	}

	rows, err := db.Query(context.Background(), queryGetFilesByRecipeUUIDs, recipeUUIDs)
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

func InsertFile(recipeUUID uuid.UUID, data []byte, contentType string, pageNumber int, isImage bool) (*models.File, error) {
	var f models.File
	var fileUUID uuid.UUID
	err := db.QueryRow(context.Background(), queryInsertFile,
		recipeUUID, data, contentType, pageNumber, isImage).Scan(&fileUUID, &f.RecipeUUID, &f.PageNumber, &f.Image)
	if err != nil {
		return nil, err
	}
	f.UUID = fileUUID
	f.URL = fmt.Sprintf("/api/files/%s", fileUUID.String())
	f.Image = isImage

	_, err = db.Exec(context.Background(), queryUpdateFileURL, f.URL, fileUUID)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// TxInsertFile inserts a file within a transaction
func TxInsertFile(ctx context.Context, tx *Tx, recipeUUID uuid.UUID, data []byte, contentType string, pageNumber int, isImage bool) (*models.File, error) {
	var f models.File
	var fileUUID uuid.UUID
	err := tx.tx.QueryRow(ctx, queryInsertFile,
		recipeUUID, data, contentType, pageNumber, isImage).Scan(&fileUUID, &f.RecipeUUID, &f.PageNumber, &f.Image)
	if err != nil {
		return nil, err
	}
	f.UUID = fileUUID
	f.URL = fmt.Sprintf("/api/files/%s", fileUUID.String())
	f.Image = isImage

	_, err = tx.tx.Exec(ctx, queryUpdateFileURL, f.URL, fileUUID)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func GetFileData(fileUUID uuid.UUID) ([]byte, string, error) {
	var data []byte
	var contentType string
	err := db.QueryRow(context.Background(), queryGetFileData, fileUUID).Scan(&data, &contentType)
	if err != nil {
		return nil, "", err
	}
	return data, contentType, nil
}
