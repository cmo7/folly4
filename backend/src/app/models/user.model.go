package models

import "github.com/cmo7/folly4/src/lib/generics/common"

type UserEntity struct {
	BaseModel `gorm:"embedded"`
	Username  string        `gorm:"unique;not null"`
	Password  string        `gorm:"not null"`
	Email     string        `gorm:"unique;not null"`
	Roles     []*RoleEntity `gorm:"many2many:user_roles;"`
}

// Implement common.Entity interface.
func (u *UserEntity) GetEntityName() common.EntityName {
	return common.EntityName("User")
}

func (u *UserEntity) GetName() string {
	return u.Username
}

// Implement permission.User interface.
func (u *UserEntity) GetRoles() []*RoleEntity {
	return u.Roles
}

func (u *UserEntity) SetRoles(roles []*RoleEntity) {
	u.Roles = roles
}

func (u *UserEntity) GetPermissions() []*PermissionEntity {
	var permissions []*PermissionEntity
	for _, role := range u.Roles {
		permissions = append(permissions, role.Permissions...)
	}
	return permissions
}

func (u *UserEntity) SetPermissions(permissions []*PermissionEntity) {
	// Do nothing.
}
