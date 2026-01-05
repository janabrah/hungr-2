package storage

import (
	"context"

	"github.com/cobyabrahams/hungr/models"
	"github.com/gofrs/uuid"
)

const (
	queryGetUserByUUID = `
		SELECT uuid, email, name, created_at
		FROM users WHERE uuid = $1`

	queryGetUserByEmail = `
		SELECT uuid, email, name, created_at
		FROM users WHERE email = $1`

	queryCreateUser = `
		INSERT INTO users (email, name, last_seen)
		VALUES ($1, $2, NOW())
		RETURNING uuid, email, name, created_at`

	queryUpsertUserOnLogin = `
		INSERT INTO users (email, name, last_seen)
		VALUES ($1, $1, NOW())
		ON CONFLICT (email) DO UPDATE SET last_seen = NOW()
		RETURNING uuid, email, name, created_at`

	queryUpdateUser = `
		UPDATE users SET name = $1
		WHERE uuid = $2
		RETURNING uuid, email, name, created_at`

	queryDeleteUser = `DELETE FROM users WHERE uuid = $1`
)

func GetUserByUUID(userUUID uuid.UUID) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(), queryGetUserByUUID, userUUID).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(), queryGetUserByEmail, email).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(email, name string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(), queryCreateUser, email, name).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpsertUserOnLogin(email string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(), queryUpsertUserOnLogin, email).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(userUUID uuid.UUID, name string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(context.Background(), queryUpdateUser, name, userUUID).Scan(
		&u.UUID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func DeleteUser(userUUID uuid.UUID) error {
	_, err := db.Exec(context.Background(), queryDeleteUser, userUUID)
	return err
}
