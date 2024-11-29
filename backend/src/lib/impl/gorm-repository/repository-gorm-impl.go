package gorm_impl

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/cmo7/folly4/src/lib/generics/relation"
	"gorm.io/gorm"
)

type GormGenericRepository[E common.Entity] struct {
	db *gorm.DB
}

func NewGormGenericRepository[E common.Entity](db *gorm.DB) *GormGenericRepository[E] {
	return &GormGenericRepository[E]{db: db}
}

func (r *GormGenericRepository[E]) Create(ctx context.Context, payload E) (E, error) {
	result := r.db.WithContext(ctx).Create(&payload)
	return payload, result.Error
}

func (r *GormGenericRepository[E]) Update(ctx context.Context, payload E) (E, error) {
	result := r.db.WithContext(ctx).Save(&payload)
	return payload, result.Error
}

func (r *GormGenericRepository[E]) UpdateField(ctx context.Context, payload E, field string, value interface{}) (E, error) {
	result := r.db.WithContext(ctx).Model(&payload).Update(field, value)
	return payload, result.Error
}

func (r *GormGenericRepository[E]) Delete(ctx context.Context, payload E) error {
	result := r.db.WithContext(ctx).Delete(&payload)
	return result.Error
}

func (r *GormGenericRepository[E]) FindOne(ctx context.Context, id common.ID, relations []relation.Relation) (E, error) {
	var entity E
	result := r.db.WithContext(ctx).Scopes(scopePreload(relations)).First(&entity, id)
	return entity, result.Error
}

func (r *GormGenericRepository[E]) FindAll(ctx context.Context, pageable pagination.Pageable, f filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[E], error) {
	var entities []E

	result := r.db.WithContext(ctx).Scopes(
		scopePage(pageable),
		scopePreload(relations),
		scopeOrder(orderBys),
		scopeFilter(f),
	).Find(&entities)
	if result.Error != nil {
		return pagination.Page[E]{}, result.Error
	}

	filteredCount, err := r.Count(ctx, f)
	if err != nil {
		return pagination.Page[E]{}, err
	}

	count, err := r.Count(ctx, filter.Composite{Operator: filter.And, Filters: []filter.Filter{}})
	if err != nil {
		return pagination.Page[E]{}, err
	}

	return pagination.Page[E]{Content: entities, Page: pageable.Page, Size: pageable.Size, Total: count, Filtered: filteredCount}, result.Error
}

func (r *GormGenericRepository[E]) Count(ctx context.Context, filter filter.Filter) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(new(E)).Scopes(scopeFilter(filter)).Count(&count)
	return count, result.Error
}

func (r *GormGenericRepository[E]) Associate(ctx context.Context, id common.ID, association string, targetId common.ID) (E, error) {
	var entity E
	r.db.WithContext(ctx).Model(&entity).Association(association).Append(&targetId)
	return entity, nil
}

func (r *GormGenericRepository[E]) Dissociate(ctx context.Context, id common.ID, association string, targetId common.ID) (E, error) {
	var entity E
	r.db.WithContext(ctx).Model(&entity).Association(association).Delete(&targetId)
	return entity, nil
}

func (r *GormGenericRepository[E]) Exists(ctx context.Context, id common.ID) (bool, error) {
	var entity E
	result := r.db.WithContext(ctx).First(&entity, id)
	return result.RowsAffected > 0, result.Error
}

func (r *GormGenericRepository[E]) Random(ctx context.Context) (E, error) {
	var entity E
	result := r.db.WithContext(ctx).Order("RANDOM()").First(&entity)
	return entity, result.Error
}

func (r *GormGenericRepository[E]) First(ctx context.Context, filter filter.Filter) (E, error) {
	var entity E
	result := r.db.WithContext(ctx).Scopes(scopeFilter(filter)).First(&entity)
	return entity, result.Error
}

func (r *GormGenericRepository[E]) ComboBox(ctx context.Context, pageable pagination.Pageable, f filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[common.ComboOption], error) {
	var entities []E
	result := r.db.WithContext(ctx).Scopes(
		scopePage(pageable),
		scopePreload(relations),
		scopeOrder(orderBys),
		scopeFilter(f),
	).Find(&entities)
	if result.Error != nil {
		return pagination.Page[common.ComboOption]{}, result.Error
	}

	filteredCount, err := r.Count(ctx, f)
	if err != nil {
		return pagination.Page[common.ComboOption]{}, err
	}

	count, err := r.Count(ctx, filter.Composite{Operator: filter.And, Filters: []filter.Filter{}})
	if err != nil {
		return pagination.Page[common.ComboOption]{}, err
	}

	var options []common.ComboOption
	for _, entity := range entities {
		options = append(options, common.ComboOption{ID: entity.GetID(), Name: entity.GetName()})
	}

	return pagination.Page[common.ComboOption]{Content: options, Page: pageable.Page, Size: pageable.Size, Total: count, Filtered: filteredCount}, result.Error
}
