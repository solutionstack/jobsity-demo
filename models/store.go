package models

import "github.com/google/uuid"

type UserRecord struct {
	ID uuid.UUID `json:"id"`
	Signup
}
