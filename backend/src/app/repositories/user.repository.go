package repositories

import (
	"github.com/cmo7/folly4/src/app/models"
	gorm_impl "github.com/cmo7/folly4/src/lib/impl/gorm-repository"
	"gorm.io/gorm"
)

// UserGormRepository is the repository for the UserEntity model.
// It is a wrapper around the GormGenericRepository. This is so we can easily add custom methods to the repository.
type UserGormRepository struct {
	*gorm_impl.GormGenericRepository[*models.UserEntity]
}

// The singleton instance of the UserGormRepository.
var userRepo *UserGormRepository

// GetUserRepository returns the singleton instance of the UserGormRepository.
func GetUserRepository(db *gorm.DB) *UserGormRepository {
	if userRepo == nil {
		userRepo = &UserGormRepository{
			GormGenericRepository: gorm_impl.NewGormGenericRepository[*models.UserEntity](db),
		}
	}
	return userRepo
}
