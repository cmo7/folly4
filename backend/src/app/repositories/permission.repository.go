package repositories

import (
	"github.com/cmo7/folly4/src/app/models"
	gorm_impl "github.com/cmo7/folly4/src/lib/impl/gorm-repository"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	*gorm_impl.GormGenericRepository[*models.PermissionEntity]
}

var permissionRepo *PermissionRepository

func GetPermissionRepository(db *gorm.DB) *PermissionRepository {
	if permissionRepo == nil {
		permissionRepo = &PermissionRepository{
			GormGenericRepository: gorm_impl.NewGormGenericRepository[*models.PermissionEntity](db),
		}
	}
	return permissionRepo
}
