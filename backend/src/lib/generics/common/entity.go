package common

import "github.com/google/uuid"

// Entity is an interface that represents a generic entity.
// It provides a method to retrieve the entity's ID.
type Entity interface {
	GetID() ID
	GetName() string
}

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

type ComboOption struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}
