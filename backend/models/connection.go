package models

import "github.com/gofrs/uuid"

type CreateConnectionRequest struct {
	TargetUserUUID uuid.UUID `json:"target_user_uuid"`
}

type ConnectionsResponse struct {
	Success     bool   `json:"success"`
	Connections []User `json:"connections"`
}
