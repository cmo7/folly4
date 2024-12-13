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

	// Any other methods that are specific to the user service can be added here.
	// For example, the user service may have a method to get a user by username.
	// This method would not be part of the CRUD operations.
	// The user service can also have a method to get a user by email.
	// This method would not be part of the CRUD operations.
	// The user service can also have a method to get a user by username and password.
	// This method would not be part of the CRUD operations.
	// ...
}

var userService *UserService

// instantiate initializes the user service with the necessary layers and dependencies.
// It composes the user service with the following layers:
// - User Repository: Interacts with the database.
// - User Audit Service: Logs all CRUD operations performed on the user entity.
// - User Permission Service: Checks if the user has the required permissions to perform CRUD operations.
// - User Service: Adds user-specific functionality to the user permission service.
//
// Parameters:
// - db: A pointer to the gorm.DB instance used for database interactions.
func instantiate(db *gorm.DB) {
	userService = nil

	// The user service is a composition of:
	// - A user permission service.
	// - A user audit service.
	// - A user repository.

	// Layer 1: User Repository. The lowest layer in the user service. This layer interacts with the database.
	userRepository := repositories.GetUserRepository(db)

	// Layer 2: User Audit Service. Adds audit functionality to the user repository. The audit service will log all the CRUD operations performed on the user entity.
	userAuditService := auditservice.NewAuditService(
		userRepository,
		repositories.GetAuditRepository(db),
	)

	// Layer 3: User Permission Service. Adds permission functionality to the user audit service. The permission service will check if the user has the required permissions to perform the CRUD operations on the user entity.
	userPermissionService := permissionservice.NewPermissionService(
		userAuditService,
		repositories.GetPermissionRepository(db),
	)

	// Layer 4: User Service. Adds user-specific functionality to the user permission service. The user service will have methods that are specific to the user entity.
	userService = &UserService{
		userPermissionService,
	}
}

// Return the user service singleton.
func GetUserService(db *gorm.DB) *UserService {
	if userService == nil {
		instantiate(db)
	}
	return userService
}
