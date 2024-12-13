package models

import "github.com/cmo7/folly4/src/lib/generics/common"

type RoleEntity struct {
	BaseModel     `gorm:"embedded"`
	Name          string `gorm:"unique;not null"`
	LocalizedName string `gorm:"not null"`

	// A role can have many permissions.
	Permissions []*PermissionEntity `gorm:"many2many:role_permissions;"`

	// A role can have many users.
	Users []*UserEntity `gorm:"many2many:user_roles;"`
}

func (r *RoleEntity) GetEntityName() common.EntityName {
	return common.EntityName("Role")
}

func (r *RoleEntity) GetName() string {
	return r.Name
}
