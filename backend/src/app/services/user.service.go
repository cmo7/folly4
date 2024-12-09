package services

import (
	"github.com/cmo7/folly4/src/app/models"
	"github.com/cmo7/folly4/src/app/repositories"
	"github.com/cmo7/folly4/src/lib/generics/service"
	auditservice "github.com/cmo7/folly4/src/lib/impl/audit-service"
	permissionservice "github.com/cmo7/folly4/src/lib/impl/permission-service"
	"gorm.io/gorm"
)

type UserService struct {
	// Embed a service.CrudService to get all the CRUD methods for free.
	// The crud service will use the user repository to perform the CRUD operations.
	service.CrudService[*models.UserEntity]
}

var userService *UserService

func GetUserService(db *gorm.DB) *UserService {
	if userService == nil {
		userService = &UserService{
			permissionservice.NewPermissionService(
				auditservice.NewAuditService(
					repositories.GetUserRepository(db),
					repositories.GetAuditRepository(db),
				),
				repositories.GetUserRepository(db),
			),
		}
	}
	return userService
}
