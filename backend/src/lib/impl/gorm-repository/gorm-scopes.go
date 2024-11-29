package gorm_impl

import (
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/cmo7/folly4/src/lib/generics/relation"
	"gorm.io/gorm"
)

/**
* The following functions are scopes that can be used to build queries.
 */

// scopePage recieves a Pageable and applies the limit and offset to the query.
func scopePage(pageable pagination.Pageable) func(*gorm.DB) *gorm.DB {
	limit := pageable.Size
	offset := pageable.Size * (pageable.Page - 1)
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}

// scopePreload recieves a list of relations and applies the preload to the query.
func scopePreload(relations []relation.Relation) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(string(relation))
		}
		return db
	}
}

// scopeOrder recieves a list of OrderBys and applies the order to the query.
func scopeOrder(orderBys []order.OrderBy) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, orderBy := range orderBys {
			if orderBy.Direction == order.Asc {
				db = db.Order(orderBy.Field)
			} else {
				db = db.Order(orderBy.Field + " DESC")
			}
		}
		return db
	}
}

// scopeFilter recieves a Filter and applies the filter to the query.
func scopeFilter(f filter.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.IsComposite() {
			return filterComposite(f.(filter.Composite))(db)
		}
		return db
	}
}

func filterComposite(f filter.Composite) func(*gorm.DB) *gorm.DB {
	switch f.Operator {
	case filter.And:
		return andCompose(f.Filters)
	case filter.Or:
		return orCompose(f.Filters)
	case filter.Not:
		return notCompose(f.Filters)
	}
	return nil
}

func andCompose(filters []filter.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, f := range filters {
			if f.IsComposite() {
				db = db.Where(filterComposite(f.(filter.Composite))(db))
			} else {
				db = db.Where(compare(f.(filter.Leaf)), f.(filter.Leaf).Value)
			}
		}
		return db
	}
}

func orCompose(filters []filter.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, f := range filters {
			if f.IsComposite() {
				db = db.Or(filterComposite(f.(filter.Composite))(db))
			} else {
				db = db.Or(compare(f.(filter.Leaf)), f.(filter.Leaf).Value)
			}
		}
		return db
	}
}

func notCompose(filters []filter.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, f := range filters {
			if f.IsComposite() {
				db = db.Not(filterComposite(f.(filter.Composite))(db))
			} else {
				db = db.Not(compare(f.(filter.Leaf)), f.(filter.Leaf).Value)
			}
		}
		return db
	}
}

func compare(Filter filter.Leaf) string {
	switch Filter.Comparator {
	case filter.Equal:
		return Filter.Field + " = ?"
	case filter.NotEqual:
		return Filter.Field + " != ?"
	case filter.GreaterThan:
		return Filter.Field + " > ?"
	case filter.GreaterThanOrEqual:
		return Filter.Field + " >= ?"
	case filter.LessThan:
		return Filter.Field + " < ?"
	case filter.LessThanOrEqual:
		return Filter.Field + " <= ?"
	case filter.Like:
		return Filter.Field + " LIKE ?"
	case filter.NotLike:
		return Filter.Field + " NOT LIKE ?"
	case filter.In:
		return Filter.Field + " IN (?)"
	case filter.NotIn:
		return Filter.Field + " NOT IN (?)"
	case filter.IsNull:
		return Filter.Field + " IS NULL"
	case filter.IsNotNull:
		return Filter.Field + " IS NOT NULL"
	default:
		return " = ?"
	}
}
