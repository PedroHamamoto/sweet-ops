package utils

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}
