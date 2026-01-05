package storage

import (
	"context"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const (
	queryGetRecipesByUserEmail = `
		SELECT r.uuid, r.name, r.user_uuid, r.tag_string, r.created_at
		FROM recipes r
		JOIN users u ON r.user_uuid = u.uuid
		WHERE u.email = $1
		ORDER BY r.created_at DESC LIMIT 100`

	queryInsertRecipeByEmail = `
		INSERT INTO recipes (name, user_uuid, tag_string)
		SELECT $1, u.uuid, $2
		FROM users u WHERE u.email = $3
		RETURNING uuid, name, user_uuid, tag_string, created_at`

	queryDeleteRecipeTags = `DELETE FROM recipe_tags WHERE recipe_uuid = $1`
	queryDeleteRecipeFiles = `DELETE FROM files WHERE recipe_uuid = $1`
	queryDeleteRecipe = `DELETE FROM recipes WHERE uuid = $1`
)

func GetRecipesByUserEmail(email string) ([]models.Recipe, error) {
	rows, err := db.Query(context.Background(), queryGetRecipesByUserEmail, email)
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

func InsertRecipeByEmail(name string, email string, tagString string) (*models.Recipe, error) {
	var r models.Recipe
	err := db.QueryRow(context.Background(), queryInsertRecipeByEmail,
		name, tagString, email).Scan(&r.UUID, &r.Name, &r.User, &r.TagString, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func DeleteRecipe(recipeUUID uuid.UUID) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if _, err = tx.Exec(context.Background(), queryDeleteRecipeTags, recipeUUID); err != nil {
		return err
	}

	if _, err = tx.Exec(context.Background(), queryDeleteRecipeFiles, recipeUUID); err != nil {
		return err
	}

	if _, err = tx.Exec(context.Background(), queryDeleteRecipe, recipeUUID); err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
