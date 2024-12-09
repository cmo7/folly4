package auditservice

import (
	"context"
	"encoding/json"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/repository"
	"github.com/cmo7/folly4/src/lib/generics/service"
	"github.com/cmo7/folly4/src/lib/generics/util/audit"
)

// AuditService is a service that provides audit functionality. It is a wrapper around a CrudService.
// It adds hooks to the CrudService to create an audit log for each action.
// Its generic types are E, which is the entity type, and A, which is the audit type.
// audit.Audit is an interface that represents an audit log. It is expected to be implemented by the user.
type AuditService[E common.Entity, A audit.Audit] struct {
	service.CrudServiceWithHooks[E]                          // Embed the CrudServiceWithHooks to inherit its methods.
	auditRepository                 repository.Repository[A] // The repository to store the audit logs.
}

func NewAuditService[E common.Entity, A audit.Audit](
	crudService service.CrudService[E],
	auditRepository repository.Repository[A],
) *AuditService[E, A] {
	service := &AuditService[E, A]{
		CrudServiceWithHooks: service.NewCrudServiceWithHooks(crudService),
		auditRepository:      auditRepository,
	}

	// Add hooks to create audit logs for each action.
	service.AddBeforeCreateHook(func(ctx context.Context, payload E) error {
		a := audit.GetAudit[A](ctx)
		a.SetAction(audit.AuditActionCreate)
		a.SetEntity(payload.GetEntityName())
		a.SetEntityID(payload.GetID())
		a.SetNewValue(serializeEntity(payload))
		return nil
	})

	service.AddAfterCreateHook(func(ctx context.Context, payload E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultSuccess)
		_, err := service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddOnCreateFailHook(func(ctx context.Context, err error, failedEntity E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultFailure)
		a.SetMessage(err.Error())
		_, err = service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddBeforeUpdateHook(func(ctx context.Context, payload E) error {
		a := audit.GetAudit[A](ctx)
		a.SetAction(audit.AuditActionUpdate)
		a.SetEntity(payload.GetEntityName())
		a.SetEntityID(payload.GetID())
		a.SetNewValue(serializeEntity(payload))
		return nil
	})

	service.AddAfterUpdateHook(func(ctx context.Context, payload E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultSuccess)
		a.SetPrevValue(serializeEntity(payload))
		_, err := service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddOnUpdateFailHook(func(ctx context.Context, err error, failedEntity E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultFailure)
		a.SetMessage(err.Error())
		_, err = service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddBeforeDeleteHook(func(ctx context.Context, id E) error {
		a := audit.GetAudit[A](ctx)
		a.SetAction(audit.AuditActionDelete)
		a.SetEntity(id.GetEntityName())
		a.SetEntityID(id.GetID())
		return nil
	})

	service.AddAfterDeleteHook(func(ctx context.Context, id E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultSuccess)
		_, err := service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddOnDeleteFailHook(func(ctx context.Context, err error, id E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultFailure)
		a.SetMessage(err.Error())
		_, err = service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddBeforeFindHook(func(ctx context.Context, entities ...E) error {
		if len(entities) == 0 {
			return nil
		}
		a := audit.GetAudit[A](ctx)
		a.SetAction(audit.AuditActionRead)
		a.SetEntity(entities[0].GetEntityName())
		return nil
	})

	service.AddAfterFindHook(func(ctx context.Context, entities ...E) error {
		if len(entities) == 0 {
			return nil
		}
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultSuccess)
		return nil
	})

	service.AddOnFindFailHook(func(ctx context.Context, err error, entities ...E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultFailure)
		a.SetMessage(err.Error())
		_, err = service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddBeforeAssocHook(func(ctx context.Context) error {
		a := audit.GetAudit[A](ctx)
		a.SetAction(audit.AuditActionAssociate)
		return nil
	})

	service.AddAfterAssocHook(func(ctx context.Context, entity E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultSuccess)
		_, err := service.auditRepository.Create(ctx, a)
		return err
	})

	service.AddOnAssocFailHook(func(ctx context.Context, err error, entity E) error {
		a := audit.GetAudit[A](ctx)
		a.SetActionResult(audit.AuditActionResultFailure)
		a.SetMessage(err.Error())
		_, err = service.auditRepository.Create(ctx, a)
		return err
	})

	return service
}

// Create creates an entity and an audit log for the creation.

func serializeEntity(entity common.Entity) string {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return ""
	}
	return string(bytes)

}
