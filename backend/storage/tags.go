package storage

import (
	"context"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const tagNamespace = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

const (
	queryUpsertTag = `
		INSERT INTO tags (uuid, name) VALUES ($1, $2)
		ON CONFLICT (uuid) DO UPDATE SET name = EXCLUDED.name
		RETURNING uuid, name`

	queryInsertRecipeTag = `
		INSERT INTO recipe_tags (recipe_uuid, tag_uuid) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	queryGetAllTags = `
		SELECT uuid, name FROM tags ORDER BY name`

	queryGetTagsByRecipeUUID = `
		SELECT t.uuid, t.name FROM tags t
		JOIN recipe_tags rt ON rt.tag_uuid = t.uuid
		WHERE rt.recipe_uuid = $1
		ORDER BY t.name`
)

func UpsertTag(tagUUID uuid.UUID, name string) (*models.Tag, error) {
	var t models.Tag
	err := db.QueryRow(context.Background(), queryUpsertTag, tagUUID, name).Scan(&t.UUID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func InsertRecipeTag(recipeUUID, tagUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(), queryInsertRecipeTag, recipeUUID, tagUUID)
	return err
}

// TxUpsertTag upserts a tag within a transaction
func TxUpsertTag(ctx context.Context, tx *Tx, tagUUID uuid.UUID, name string) (*models.Tag, error) {
	var t models.Tag
	err := tx.tx.QueryRow(ctx, queryUpsertTag, tagUUID, name).Scan(&t.UUID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TxInsertRecipeTag inserts a recipe-tag link within a transaction
func TxInsertRecipeTag(ctx context.Context, tx *Tx, recipeUUID, tagUUID uuid.UUID) error {
	_, err := tx.tx.Exec(ctx, queryInsertRecipeTag, recipeUUID, tagUUID)
	return err
}

// TxDeleteRecipeTags deletes all tags for a recipe within a transaction
func TxDeleteRecipeTags(ctx context.Context, tx *Tx, recipeUUID uuid.UUID) error {
	_, err := tx.tx.Exec(ctx, "DELETE FROM recipe_tags WHERE recipe_uuid = $1", recipeUUID)
	return err
}

func CreateTagUUID(tag string) uuid.UUID {
	namespace := uuid.Must(uuid.FromString(tagNamespace))
	return uuid.NewV5(namespace, tag)
}

func GetAllTags() ([]models.Tag, error) {
	rows, err := db.Query(context.Background(), queryGetAllTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.UUID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func GetTagsByRecipeUUID(recipeUUID uuid.UUID) ([]models.Tag, error) {
	rows, err := db.Query(context.Background(), queryGetTagsByRecipeUUID, recipeUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.UUID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}
