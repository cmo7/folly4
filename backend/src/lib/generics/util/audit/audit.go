package audit

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/google/uuid"
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

type AuditActionResult string

const (
	AuditActionResultNone    AuditActionResult = "NONE"
	AuditActionResultSuccess AuditActionResult = "SUCCESS"
	AuditActionResultFailure AuditActionResult = "FAILURE"
)

type Audit interface {
	common.Entity
	GetAction() AuditAction
	SetAction(action AuditAction)
	GetActionResult() AuditActionResult
	SetActionResult(actionResult AuditActionResult)
	SetMessage(message string)
	GetMessage() string
	GetUserID() uuid.UUID
	SetUserID(userID uuid.UUID)
	GetEntity() common.EntityName
	SetEntity(entity common.EntityName)
	GetEntityID() uuid.UUID
	SetEntityID(entityID uuid.UUID)
	GetNewValue() string
	SetNewValue(newValue string)
	GetPrevValue() string
	SetPrevValue(prevValue string)
	GetLocation() string
	SetLocation(location string)
	GetIP() string
	SetIP(ip string)
	GetUserAgent() string
	SetUserAgent(userAgent string)
}

type AuditContextKey struct{}

func WithAudit(ctx context.Context, audit Audit) context.Context {
	return context.WithValue(ctx, AuditContextKey{}, audit)
}

func GetAudit[A Audit](ctx context.Context) A {
	var zero A
	audit, ok := ctx.Value(AuditContextKey{}).(A)
	if !ok {
		return zero
	}
	return audit
}

func SetAction[A Audit](ctx context.Context, action AuditAction) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetAction(action)
	return WithAudit(ctx, audit)
}

func SetUserID[A Audit](ctx context.Context, userID uuid.UUID) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetUserID(userID)
	return WithAudit(ctx, audit)
}

func SetEntity[A Audit](ctx context.Context, entity common.EntityName) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetEntity(entity)
	return WithAudit(ctx, audit)
}

func SetEntityID[A Audit](ctx context.Context, entityID uuid.UUID) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetEntityID(entityID)
	return WithAudit(ctx, audit)
}

func SetNewValue[A Audit](ctx context.Context, newValue string) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetNewValue(newValue)
	return WithAudit(ctx, audit)
}

func SetLocation[A Audit](ctx context.Context, location string) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetLocation(location)
	return WithAudit(ctx, audit)
}

func SetIP[A Audit](ctx context.Context, ip string) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetIP(ip)
	return WithAudit(ctx, audit)
}

func SetUserAgent[A Audit](ctx context.Context, userAgent string) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetUserAgent(userAgent)
	return WithAudit(ctx, audit)
}

func SetPrevValue[A Audit](ctx context.Context, prevValue string) context.Context {
	audit := GetAudit[A](ctx)
	audit.SetPrevValue(prevValue)
	return WithAudit(ctx, audit)
}
