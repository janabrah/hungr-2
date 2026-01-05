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

func CreateTagUUID(tag string) uuid.UUID {
	namespace := uuid.Must(uuid.FromString(tagNamespace))
	return uuid.NewV5(namespace, tag)
}
