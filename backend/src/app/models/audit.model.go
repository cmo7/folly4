package models

import (
	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/util/audit"
	"github.com/google/uuid"
)

type AuditEntity struct {
	BaseModel `gorm:"embedded"`
	Action    audit.AuditAction
	Result    audit.AuditActionResult
	Message   string
	UserID    uuid.UUID
	Entity    common.EntityName
	EntityID  uuid.UUID
	NewValue  string
	PrevValue string
	Location  string
	IP        string
	UserAgent string
}

func (a *AuditEntity) GetEntityName() common.EntityName {
	return common.EntityName("Audit")
}

func (a *AuditEntity) GetName() string {
	return a.ID.String()
}

func (a *AuditEntity) GetAction() audit.AuditAction {
	return a.Action
}

func (a *AuditEntity) GetActionResult() audit.AuditActionResult {
	return a.Result
}

func (a *AuditEntity) GetMessage() string {
	return a.Message
}

func (a *AuditEntity) GetUserID() uuid.UUID {
	return a.UserID
}

func (a *AuditEntity) GetEntity() common.EntityName {
	return a.Entity
}

func (a *AuditEntity) GetEntityID() uuid.UUID {
	return a.EntityID
}

func (a *AuditEntity) GetNewValue() string {
	return a.NewValue
}

func (a *AuditEntity) GetPrevValue() string {
	return a.PrevValue
}

func (a *AuditEntity) GetLocation() string {
	return a.Location
}

func (a *AuditEntity) GetIP() string {
	return a.IP
}

func (a *AuditEntity) GetUserAgent() string {
	return a.UserAgent
}

func (a *AuditEntity) SetAction(action audit.AuditAction) {
	a.Action = action
}

func (a *AuditEntity) SetActionResult(result audit.AuditActionResult) {
	a.Result = result
}

func (a *AuditEntity) SetMessage(message string) {
	a.Message = message
}

func (a *AuditEntity) SetUserID(userID uuid.UUID) {
	a.UserID = userID
}

func (a *AuditEntity) SetEntity(entity common.EntityName) {
	a.Entity = entity
}

func (a *AuditEntity) SetEntityID(entityID uuid.UUID) {
	a.EntityID = entityID
}

func (a *AuditEntity) SetNewValue(newValue string) {
	a.NewValue = newValue
}

func (a *AuditEntity) SetPrevValue(prevValue string) {
	a.PrevValue = prevValue
}

func (a *AuditEntity) SetLocation(location string) {
	a.Location = location
}

func (a *AuditEntity) SetIP(ip string) {
	a.IP = ip
}

func (a *AuditEntity) SetUserAgent(userAgent string) {
	a.UserAgent = userAgent
}
