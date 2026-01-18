package storage

import (
	"context"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const (
	queryGetRecipeByUUID = `
		SELECT r.uuid, r.name, r.user_uuid, r.source,
		       COALESCE(STRING_AGG(t.name, ', ' ORDER BY rt.id), '') as tag_string,
		       r.created_at, u.email
		FROM recipes r
		JOIN users u ON r.user_uuid = u.uuid
		LEFT JOIN recipe_tags rt ON r.uuid = rt.recipe_uuid
		LEFT JOIN tags t ON rt.tag_uuid = t.uuid
		WHERE r.uuid = $1
		GROUP BY r.uuid, r.name, r.user_uuid, r.source, r.created_at, u.email`

	queryGetRecipesByUserEmail = `
		WITH viewer AS (
			SELECT uuid FROM users WHERE email = $1
		)
		SELECT r.uuid, r.name, r.user_uuid, r.source,
		       COALESCE(STRING_AGG(t.name, ', ' ORDER BY rt.id), '') as tag_string,
		       r.created_at, u.email
		FROM recipes r
		JOIN users u ON r.user_uuid = u.uuid
		JOIN viewer v ON true
		LEFT JOIN recipe_tags rt ON r.uuid = rt.recipe_uuid
		LEFT JOIN tags t ON rt.tag_uuid = t.uuid
		WHERE r.user_uuid = v.uuid
			OR EXISTS (
				SELECT 1 FROM user_connections uc
				WHERE uc.source_user_uuid = r.user_uuid
					AND uc.target_user_uuid = v.uuid
			)
		GROUP BY r.uuid, r.name, r.user_uuid, r.source, r.created_at, u.email
		ORDER BY r.created_at DESC LIMIT 100`

	queryInsertRecipeByEmail = `
		INSERT INTO recipes (name, user_uuid, source)
		SELECT $1, u.uuid, $3
		FROM users u WHERE u.email = $2
		RETURNING uuid, name, user_uuid, $3 as source, '' as tag_string, created_at, $2 as owner_email`

	queryDeleteRecipeTags   = `DELETE FROM recipe_tags WHERE recipe_uuid = $1`
	queryDeleteRecipeFiles  = `DELETE FROM files WHERE recipe_uuid = $1`
	queryDeleteRecipe       = `DELETE FROM recipes WHERE uuid = $1`
	queryUpdateRecipeSource = `UPDATE recipes SET source = $2 WHERE uuid = $1`
)

func GetRecipeByUUID(recipeUUID uuid.UUID) (*models.Recipe, error) {
	var r models.Recipe
	err := db.QueryRow(context.Background(), queryGetRecipeByUUID, recipeUUID).Scan(
		&r.UUID, &r.Name, &r.User, &r.Source, &r.TagString, &r.CreatedAt, &r.OwnerEmail)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRecipesByUserEmail(email string) ([]models.Recipe, error) {
	rows, err := db.Query(context.Background(), queryGetRecipesByUserEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := []models.Recipe{}
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.UUID, &r.Name, &r.User, &r.Source, &r.TagString, &r.CreatedAt, &r.OwnerEmail); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, rows.Err()
}

func InsertRecipeByEmail(name string, email string, source *string) (*models.Recipe, error) {
	var r models.Recipe
	err := db.QueryRow(context.Background(), queryInsertRecipeByEmail,
		name, email, source).Scan(&r.UUID, &r.Name, &r.User, &r.Source, &r.TagString, &r.CreatedAt, &r.OwnerEmail)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// TxInsertRecipeByEmail inserts a recipe within a transaction
func TxInsertRecipeByEmail(ctx context.Context, tx *Tx, name string, email string, source *string) (*models.Recipe, error) {
	var r models.Recipe
	err := tx.tx.QueryRow(ctx, queryInsertRecipeByEmail,
		name, email, source).Scan(&r.UUID, &r.Name, &r.User, &r.Source, &r.TagString, &r.CreatedAt, &r.OwnerEmail)
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

func TxUpdateRecipeSource(ctx context.Context, tx *Tx, recipeUUID uuid.UUID, source *string) error {
	_, err := tx.tx.Exec(ctx, queryUpdateRecipeSource, recipeUUID, source)
	return err
}
