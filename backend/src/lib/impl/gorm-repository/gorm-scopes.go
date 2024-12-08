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
	// Validate the page and size
	if pageable.Page < 1 {
		pageable.Page = 1
	}
	if pageable.Size < 1 {
		pageable.Size = 1
	}
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
	case filter.LogicalAnd:
		return andCompose(f.Filters)
	case filter.LogicalOr:
		return orCompose(f.Filters)
	case filter.LogicalNot:
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
				db = db.Where(compare(f.(filter.Leaf))(db))
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
				db = db.Or(compare(f.(filter.Leaf))(db))
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
				db = db.Not(compare(f.(filter.Leaf))(db))
			}
		}
		return db
	}
}

func compare(f filter.Leaf) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Value == nil {
			return db
		}
		switch f.Comparator {
		case filter.ComparatorEqual:
			return db.Where(f.Field + " = ?")
		case filter.ComparatorNotEqual:
			return db.Where(f.Field + " != ?")
		case filter.ComparatorGreaterThan:
			return db.Where(f.Field + " > ?")
		case filter.ComparatorGreaterThanOrEqual:
			return db.Where(f.Field + " >= ?")
		case filter.ComparatorLessThan:
			return db.Where(f.Field + " < ?")
		case filter.ComparatorLessThanOrEqual:
			return db.Where(f.Field + " <= ?")
		case filter.ComparatorLike:
			return db.Where(f.Field + " LIKE ?")
		case filter.ComparatorNotLike:
			return db.Where(f.Field + " NOT LIKE ?")
		case filter.ComparatorIn:
			return db.Where(f.Field + " IN (?)")
		case filter.ComparatorNotIn:
			return db.Where(f.Field + " NOT IN (?)")
		case filter.ComparatorIsNull:
			return db.Where(f.Field + " IS NULL")
		case filter.ComparatorIsNotNull:
			return db.Where(f.Field + " IS NOT NULL")
		default:
			return db.Where(f.Field + " = ?")
		}
	}
}
