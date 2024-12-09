package repositories

import (
	"github.com/cmo7/folly4/src/app/models"
	gorm_impl "github.com/cmo7/folly4/src/lib/impl/gorm-repository"
	"gorm.io/gorm"
)

type AuditGormRepository struct {
	*gorm_impl.GormGenericRepository[*models.AuditEntity]
}

var auditRepo *AuditGormRepository

func GetAuditRepository(db *gorm.DB) *AuditGormRepository {
	if auditRepo == nil {
		auditRepo = &AuditGormRepository{
			GormGenericRepository: gorm_impl.NewGormGenericRepository[*models.AuditEntity](db),
		}
	}
	return auditRepo
}
