package permissionservice

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/repository"
	"github.com/cmo7/folly4/src/lib/generics/service"
	"github.com/cmo7/folly4/src/lib/generics/util/permission"
)

type PermissionService[E common.Entity, P permission.Permission] struct {
	service.CrudServiceWithHooks[E]
	permissionRepository repository.Repository[P]
}

func NewPermissionService[E common.Entity, P permission.Permission](
	crudService service.CrudService[E],
	permissionRepository repository.Repository[P],
) *PermissionService[E, P] {
	service := &PermissionService[E, P]{
		CrudServiceWithHooks: service.NewCrudServiceWithHooks(
			crudService,
		),
		permissionRepository: permissionRepository,
	}

	service.AddBeforeCreateHook(func(ctx context.Context, payload E) error {
		if permission.HasPermission(ctx, permission.OperationCreate, payload.GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationCreate, payload.GetEntityName())
	})

	service.AddBeforeFindHook(func(ctx context.Context, entities ...E) error {
		if permission.HasPermission(ctx, permission.OperationRead, entities[0].GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationRead, entities[0].GetEntityName())
	})

	service.AddBeforeUpdateHook(func(ctx context.Context, payload E) error {
		if permission.HasPermission(ctx, permission.OperationUpdate, payload.GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationUpdate, payload.GetEntityName())
	})

	service.AddBeforeDeleteHook(func(ctx context.Context, payload E) error {
		if permission.HasPermission(ctx, permission.OperationDelete, payload.GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationDelete, payload.GetEntityName())
	})

	service.AddBeforeAssocHook(func(ctx context.Context) error {
		var entity E
		if permission.HasPermission(ctx, permission.OperationAssociate, entity.GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationAssociate, entity.GetEntityName())
	})

	service.AddBeforeDissocHook(func(ctx context.Context) error {
		var entity E
		if permission.HasPermission(ctx, permission.OperationDissociate, entity.GetEntityName()) {
			return nil
		}
		return permission.PermissionDenied(ctx, permission.OperationDissociate, entity.GetEntityName())
	})

	return service
}
