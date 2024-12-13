package models

import (
	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/util/permission"
)

type PermissionEntity struct {
	BaseModel `gorm:"embedded"`
	Entity    common.EntityName
	Operation permission.Operation
}

// Implement common.Entity interface.
func (p *PermissionEntity) GetEntityName() common.EntityName {
	return common.EntityName("Permission")
}

func (p *PermissionEntity) GetName() string {
	return p.ToString()
}

// Implement permission.Permission interface.
func (p *PermissionEntity) GetEntity() common.EntityName {
	return p.Entity
}

func (p *PermissionEntity) SetEntity(entity common.EntityName) {
	p.Entity = entity
}

func (p *PermissionEntity) GetOperation() permission.Operation {
	return p.Operation
}

func (p *PermissionEntity) SetOperation(operation permission.Operation) {
	p.Operation = operation
}

func (p *PermissionEntity) ToString() string {
	return p.Entity.String() + ":" + p.Operation.String()
}
