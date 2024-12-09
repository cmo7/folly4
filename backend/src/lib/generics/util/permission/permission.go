package permission

import (
	"context"
	"fmt"

	"github.com/cmo7/folly4/src/lib/generics/common"
)

type Operation string

const (
	OperationCreate Operation = "CREATE"
	OperationRead   Operation = "READ"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"

	OperationEnable  Operation = "ENABLE"
	OperationDisable Operation = "DISABLE"

	OperationAssociate  Operation = "ASSOCIATE"
	OperationDissociate Operation = "DISSOCIATE"

	OperationLogin  Operation = "LOGIN"
	OperationLogout Operation = "LOGOUT"

	OperationApprove Operation = "APPROVE"
	OperationReject  Operation = "REJECT"
)

func (o Operation) String() string {
	return string(o)
}

// Permission is an interface that represents a permission.
// It provides methods to retrieve the permission's ID, name, and description.
type Permission interface {
	common.Entity
	GetEntity() common.EntityName
	SetEntity(entity common.EntityName)
	GetOperation() Operation
	SetOperation(operation Operation)
	ToString() string
}

type Role interface {
	common.Entity
	GetPermissions() []Permission
	SetPermissions(permissions []Permission)
}

type User interface {
	common.Entity
	GetRoles() []Role
	SetRoles(roles []Role)
	GetPermissions() []Permission
	SetPermissions(permissions []Permission)
}

type UserKey struct{}

type RolesKey struct{}

type PermissionsKey struct{}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, UserKey{}, user)
}

func GetUser(ctx context.Context) User {
	return ctx.Value(UserKey{}).(User)
}

func WithRoles(ctx context.Context, roles []Role) context.Context {
	return context.WithValue(ctx, RolesKey{}, roles)
}

func GetRoles(ctx context.Context) []Role {
	return ctx.Value(RolesKey{}).([]Role)
}

func WithPermissions(ctx context.Context, permissions []Permission) context.Context {
	return context.WithValue(ctx, PermissionsKey{}, permissions)
}

func GetPermissions(ctx context.Context) []Permission {
	return ctx.Value(PermissionsKey{}).([]Permission)
}

func getFullPermissionListFromContext(ctx context.Context) []Permission {
	// Get the permissions from the context.
	permissions := GetPermissions(ctx)

	// Get the roles from the context.
	roles := GetRoles(ctx)
	// Get the permissions from the roles.
	for _, role := range roles {
		permissions = append(permissions, role.GetPermissions()...)
	}

	// Get the user from the context.
	user := GetUser(ctx)
	// Get the permissions from the user.
	permissions = append(permissions, user.GetPermissions()...)

	// Get the permissions from the user's roles.
	for _, role := range user.GetRoles() {
		permissions = append(permissions, role.GetPermissions()...)
	}

	return permissions
}

func HasPermission(ctx context.Context, operation Operation, entity common.EntityName) bool {
	permissions := getFullPermissionListFromContext(ctx)
	for _, p := range permissions {
		if p.GetEntity() == entity && p.GetOperation() == operation {
			return true
		}
	}
	return false
}

func PermissionDenied(ctx context.Context, operation Operation, entity common.EntityName) error {
	return fmt.Errorf("permission denied: %s %s for user %s", operation, entity, GetUser(ctx).GetID())
}
