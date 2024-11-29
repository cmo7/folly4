package generics

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/cmo7/folly4/src/lib/generics/relation"
)

type CrudService[E common.Entity] interface {
	Create(ctx context.Context, payload E) (E, error)
	Update(ctx context.Context, payload E) (E, error)
	UpdateField(ctx context.Context, payload E, field string, value interface{}) (E, error)
	Delete(ctx context.Context, payload E) error
	FindOne(ctx context.Context, id common.ID, relations []relation.Relation) (E, error)
	FindAll(ctx context.Context, pageable pagination.Pageable, filter filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[E], error)
	Count(ctx context.Context, filter filter.Filter) (int64, error)
	Associate(ctx context.Context, id common.ID, association string, targetId common.ID) (E, error)
	Dissociate(ctx context.Context, id common.ID, association string, targetId common.ID) (E, error)
	Exists(ctx context.Context, id common.ID) (bool, error)
	Random(ctx context.Context) (E, error)
	First(ctx context.Context, filter filter.Filter) (E, error)
	ComboBox(ctx context.Context, pageable pagination.Pageable, filter filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[common.ComboOption], error)
}