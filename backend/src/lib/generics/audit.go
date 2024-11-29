package generics

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
)

type AuditAction string

const (
	AuditActionNone AuditAction = "NONE"

	AuditActionCreate AuditAction = "CREATE"
	AuditActionRead   AuditAction = "READ"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"

	AuditActionEnable  AuditAction = "ENABLE"
	AuditActionDisable AuditAction = "DISABLE"

	AuditActionAssociate  AuditAction = "ASSOCIATE"
	AuditActionDissociate AuditAction = "DISSOCIATE"

	AuditActionLogin  AuditAction = "LOGIN"
	AuditActionLogout AuditAction = "LOGOUT"

	AuditActionApprove AuditAction = "APPROVE"
	AuditActionReject  AuditAction = "REJECT"
)

type Audit interface {
	Audit(ctx context.Context, audit AuditContext)
}

// Represents something like "User with ID 123 created a new entity with ID 456 of type 'EntityName' with the following fields: {field1: value1, field2: value2}"
type AuditContext struct {
	Action   AuditAction
	UserID   common.ID
	Entity   string
	EntityID common.ID
	Fields   map[string]interface{}

	// Origin of the request
	IP        string
	UserAgent string

	// Additional fields
	OldValues map[string]interface{}
}

func WithAudit(ctx context.Context, audit AuditContext) context.Context {
	return context.WithValue(ctx, "audit", audit)
}

func GetAudit(ctx context.Context) AuditContext {
	audit, ok := ctx.Value("audit").(AuditContext)
	if !ok {
		return AuditContext{}
	}
	return audit
}
