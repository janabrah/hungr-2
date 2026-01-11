package storage

import (
	"context"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const (
	queryCreateConnection = `
		INSERT INTO user_connections (source_user_uuid, target_user_uuid)
		VALUES ($1, $2)`

	queryGetConnectionsBySourceUser = `
		SELECT u.uuid, u.email, u.name, u.created_at
		FROM user_connections uc
		JOIN users u ON u.uuid = uc.target_user_uuid
		WHERE uc.source_user_uuid = $1
		ORDER BY u.name`

	queryGetConnectionsByTargetUser = `
		SELECT u.uuid, u.email, u.name, u.created_at
		FROM user_connections uc
		JOIN users u ON u.uuid = uc.source_user_uuid
		WHERE uc.target_user_uuid = $1
		ORDER BY u.name`

	queryConnectionExists = `
		SELECT 1 FROM user_connections
		WHERE source_user_uuid = $1 AND target_user_uuid = $2`

	queryDeleteConnection = `
		DELETE FROM user_connections
		WHERE source_user_uuid = $1 AND target_user_uuid = $2`
)

func CreateConnection(sourceUserUUID, targetUserUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(), queryCreateConnection, sourceUserUUID, targetUserUUID)
	return err
}

func GetConnectionsBySourceUser(sourceUserUUID uuid.UUID) ([]models.User, error) {
	rows, err := db.Query(context.Background(), queryGetConnectionsBySourceUser, sourceUserUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetConnectionsByTargetUser(targetUserUUID uuid.UUID) ([]models.User, error) {
	rows, err := db.Query(context.Background(), queryGetConnectionsByTargetUser, targetUserUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func ConnectionExists(sourceUserUUID, targetUserUUID uuid.UUID) (bool, error) {
	var exists int
	err := db.QueryRow(context.Background(), queryConnectionExists, sourceUserUUID, targetUserUUID).Scan(&exists)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func DeleteConnection(sourceUserUUID, targetUserUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(), queryDeleteConnection, sourceUserUUID, targetUserUUID)
	return err
}
