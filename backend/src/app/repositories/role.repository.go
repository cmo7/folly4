package repositories

import (
	"github.com/cmo7/folly4/src/app/models"
	gorm_impl "github.com/cmo7/folly4/src/lib/impl/gorm-repository"

	"gorm.io/gorm"
)

type RoleRepository struct {
	*gorm_impl.GormGenericRepository[*models.RoleEntity]
}

var roleRepo *RoleRepository

func GetRoleRepository(db *gorm.DB) *RoleRepository {
	if roleRepo == nil {
		roleRepo = &RoleRepository{
			GormGenericRepository: gorm_impl.NewGormGenericRepository[*models.RoleEntity](db),
		}
	}
	return roleRepo
}
