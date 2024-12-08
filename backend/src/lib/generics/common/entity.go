package common

import "github.com/google/uuid"

// Entity is an interface that represents a generic entity.
// It provides a method to retrieve the entity's ID.
type Entity interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
	GetName() string
	GetEntityName() EntityName
}

type ComboOption struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EntityName string

func (e EntityName) String() string {
	return string(e)
}
