package service

/**
 * @api {interface} CrudService[E] CrudService[E]
 * @apiName CrudService[E]
 * @apiGroup Interface
 * @apiVersion 1.0.0
 *
 */

import (
	"context"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/cmo7/folly4/src/lib/generics/relation"
	"github.com/google/uuid"
)

type CrudService[E common.Entity] interface {
	Create(ctx context.Context, payload E) (E, error)
	Update(ctx context.Context, payload E) (E, error)
	UpdateField(ctx context.Context, payload E, field string, value interface{}) (E, error)
	Delete(ctx context.Context, payload E) error
	FindOne(ctx context.Context, id uuid.UUID, relations []relation.Relation) (E, error)
	FindAll(ctx context.Context, pageable pagination.Pageable, filter filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[E], error)
	Count(ctx context.Context, filter filter.Filter) (int64, error)
	Associate(ctx context.Context, id uuid.UUID, association string, targetId uuid.UUID) (E, error)
	Dissociate(ctx context.Context, id uuid.UUID, association string, targetId uuid.UUID) (E, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Random(ctx context.Context) (E, error)
	First(ctx context.Context, filter filter.Filter) (E, error)
	ComboBox(ctx context.Context, pageable pagination.Pageable, filter filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[common.ComboOption], error)
}

type ServiceHook string // ServiceHook is a type that represents a hook that can be executed before or after a service method.

const (
	// Create hooks
	BeforeCreate ServiceHook = "BeforeCreate"
	AfterCreate  ServiceHook = "AfterCreate"
	OnCreateFail ServiceHook = "OnCreateFail"
	// Update hooks
	BeforeUpdate ServiceHook = "BeforeUpdate"
	AfterUpdate  ServiceHook = "AfterUpdate"
	OnUpdateFail ServiceHook = "OnUpdateFail"
	// Delete hooks
	BeforeDelete ServiceHook = "BeforeDelete"
	AfterDelete  ServiceHook = "AfterDelete"
	OnDeleteFail ServiceHook = "OnDeleteFail"
	// Find hooks
	BeforeFind ServiceHook = "BeforeFind"
	AfterFind  ServiceHook = "AfterFind"
	OnFindFail ServiceHook = "OnFindFail"
	// Count hooks
	BeforeCount ServiceHook = "BeforeCount"
	AfterCount  ServiceHook = "AfterCount"
	OnCountFail ServiceHook = "OnCountFail"
	// Association hooks
	BeforeAssoc ServiceHook = "BeforeAssoc"
	AfterAssoc  ServiceHook = "AfterAssoc"
	OnAssocFail ServiceHook = "OnAssocFail"
	// Dissociation hooks
	BeforeDissoc ServiceHook = "BeforeDissoc"
	AfterDissoc  ServiceHook = "AfterDissoc"
	OnDissocFail ServiceHook = "OnDissocFail"
	// Exists hooks
	BeforeExists ServiceHook = "BeforeExists"
	AfterExists  ServiceHook = "AfterExists"
	OnExistsFail ServiceHook = "OnExistsFail"
	// Random hooks
	BeforeRandom ServiceHook = "BeforeRandom"
	AfterRandom  ServiceHook = "AfterRandom"
	OnRandomFail ServiceHook = "OnRandomFail"
	// First hooks
	BeforeFirst ServiceHook = "BeforeFirst"
	AfterFirst  ServiceHook = "AfterFirst"
	OnFirstFail ServiceHook = "OnFirstFail"
	// ComboBox hooks
	BeforeCombo ServiceHook = "BeforeCombo"
	AfterCombo  ServiceHook = "AfterCombo"
	OnComboFail ServiceHook = "OnComboFail"
)

type ServiceHookFunc[E common.Entity] func(ctx context.Context, payload ...E) error // ServiceHookFunc is a function that represents a hook that can be executed before or after a service method.

// CrudServiceWithHooks is a struct that represents a service that allows simple interactions with a single entity.
// Wraps a service and allows hooks to be executed before or after each method.
// As repository fulfills the CrudService interface, it can be used as a service.
// This allows multiple services to be combined into a single service with hooks, which can be useful for auditing, logging, or other purposes.
type CrudServiceWithHooks[E common.Entity] struct {
	repo CrudService[E]
	// Hooks for each service method

	// Create hooks
	BeforeCreate func(ctx context.Context, payload E) error
	AfterCreate  func(ctx context.Context, createdEntity E) error
	OnCreateFail func(ctx context.Context, err error, failedEntity E) error

	// Update hooks
	BeforeUpdate func(ctx context.Context, payload E) error
	AfterUpdate  func(ctx context.Context, updatedEntity E) error
	OnUpdateFail func(ctx context.Context, err error, failedEntity E) error

	// Delete hooks
	BeforeDelete func(ctx context.Context, payload E) error
	AfterDelete  func(ctx context.Context, deletedEntity E) error
	OnDeleteFail func(ctx context.Context, err error, failedEntity E) error

	// Find hooks
	BeforeFind func(ctx context.Context, entities ...E) error
	AfterFind  func(ctx context.Context, entities ...E) error
	OnFindFail func(ctx context.Context, err error, entities ...E) error

	// Count hooks
	BeforeCount func(ctx context.Context) error
	AfterCount  func(ctx context.Context) error
	OnCountFail func(ctx context.Context, err error) error

	// Association hooks
	BeforeAssoc func(ctx context.Context) error
	AfterAssoc  func(ctx context.Context, entity E) error
	OnAssocFail func(ctx context.Context, err error, entity E) error

	// Dissociation hooks
	BeforeDissoc func(ctx context.Context) error
	AfterDissoc  func(ctx context.Context, entity E) error
	OnDissocFail func(ctx context.Context, err error, entity E) error

	// Exists hooks
	BeforeExists func(ctx context.Context) error
	AfterExists  func(ctx context.Context) error
	OnExistsFail func(ctx context.Context, err error) error

	// Random hooks
	BeforeRandom func(ctx context.Context) error
	AfterRandom  func(ctx context.Context, entity E) error
	OnRandomFail func(ctx context.Context, err error, entity E) error

	// First hooks
	BeforeFirst func(ctx context.Context) error
	AfterFirst  func(ctx context.Context, entity E) error
	OnFirstFail func(ctx context.Context, err error, entity E) error

	// ComboBox hooks
	BeforeCombo func(ctx context.Context) error
	AfterCombo  func(ctx context.Context) error
	OnComboFail func(ctx context.Context, err error) error
}

func NewCrudServiceWithHooks[E common.Entity](repo CrudService[E]) CrudServiceWithHooks[E] {
	return CrudServiceWithHooks[E]{repo: repo}
}

func (s *CrudServiceWithHooks[E]) Create(ctx context.Context, payload E) (E, error) {
	if s.BeforeCreate != nil {
		if err := s.BeforeCreate(ctx, payload); err != nil {
			return payload, err
		}
	}

	entity, err := s.repo.Create(ctx, payload)

	if err != nil {
		if s.OnCreateFail != nil {
			if err := s.OnCreateFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterCreate != nil {
		if err := s.AfterCreate(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) Update(ctx context.Context, payload E) (E, error) {
	if s.BeforeUpdate != nil {
		if err := s.BeforeUpdate(ctx, payload); err != nil {
			return payload, err
		}
	}

	entity, err := s.repo.Update(ctx, payload)

	if err != nil {
		if s.OnUpdateFail != nil {
			if err := s.OnUpdateFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterUpdate != nil {
		if err := s.AfterUpdate(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) UpdateField(ctx context.Context, payload E, field string, value interface{}) (E, error) {
	if s.BeforeUpdate != nil {
		if err := s.BeforeUpdate(ctx, payload); err != nil {
			return payload, err
		}
	}

	entity, err := s.repo.UpdateField(ctx, payload, field, value)

	if err != nil {
		if s.OnUpdateFail != nil {
			if err := s.OnUpdateFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterUpdate != nil {
		if err := s.AfterUpdate(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) Delete(ctx context.Context, payload E) error {
	if s.BeforeDelete != nil {
		if err := s.BeforeDelete(ctx, payload); err != nil {
			return err
		}
	}

	err := s.repo.Delete(ctx, payload)

	if err != nil {
		if s.OnDeleteFail != nil {
			if err := s.OnDeleteFail(ctx, err, payload); err != nil {
				return err
			}
		}
	}

	if s.AfterDelete != nil {
		if err := s.AfterDelete(ctx, payload); err != nil {
			return err
		}
	}

	return err
}

func (s *CrudServiceWithHooks[E]) FindOne(ctx context.Context, id uuid.UUID, relations []relation.Relation) (E, error) {
	var entity E
	if s.BeforeFind != nil {
		if err := s.BeforeFind(ctx, entity); err != nil {
			return entity, err
		}
	}

	entity, err := s.repo.FindOne(ctx, id, relations)

	if err != nil {
		if s.OnFindFail != nil {
			if err := s.OnFindFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterFind != nil {
		if err := s.AfterFind(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) FindAll(ctx context.Context, pageable pagination.Pageable, f filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[E], error) {
	if s.BeforeFind != nil {
		if err := s.BeforeFind(ctx); err != nil {
			return pagination.NewPage[E]([]E{}, 0, 0, 0, 0), err
		}
	}

	page, err := s.repo.FindAll(ctx, pageable, f, relations, orderBys)

	if err != nil {
		if s.OnFindFail != nil {
			if err := s.OnFindFail(ctx, err, page.Content...); err != nil {
				return page, err
			}
		}
	}

	if s.AfterFind != nil {
		if err := s.AfterFind(ctx, page.Content...); err != nil {
			return page, err
		}
	}

	return page, err
}

func (s *CrudServiceWithHooks[E]) Count(ctx context.Context, f filter.Filter) (int64, error) {
	if s.BeforeCount != nil {
		if err := s.BeforeCount(ctx); err != nil {
			return 0, err
		}
	}

	count, err := s.repo.Count(ctx, f)

	if err != nil {
		if s.OnCountFail != nil {
			if err := s.OnCountFail(ctx, err); err != nil {
				return 0, err
			}
		}
	}

	if s.AfterCount != nil {
		if err := s.AfterCount(ctx); err != nil {
			return count, err
		}
	}

	return count, err
}

func (s *CrudServiceWithHooks[E]) Associate(ctx context.Context, id uuid.UUID, association string, targetId uuid.UUID) (E, error) {
	var entity E
	if s.BeforeAssoc != nil {
		if err := s.BeforeAssoc(ctx); err != nil {
			return entity, err
		}
	}

	entity, err := s.repo.Associate(ctx, id, association, targetId)

	if err != nil {
		if s.OnAssocFail != nil {
			if err := s.OnAssocFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterAssoc != nil {
		if err := s.AfterAssoc(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) Dissociate(ctx context.Context, id uuid.UUID, association string, targetId uuid.UUID) (E, error) {
	var entity E
	if s.BeforeDissoc != nil {
		if err := s.BeforeDissoc(ctx); err != nil {
			return entity, err
		}
	}

	entity, err := s.repo.Dissociate(ctx, id, association, targetId)

	if err != nil {
		if s.OnDissocFail != nil {
			if err := s.OnDissocFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterDissoc != nil {
		if err := s.AfterDissoc(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err
}

func (s *CrudServiceWithHooks[E]) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if s.BeforeExists != nil {
		if err := s.BeforeExists(ctx); err != nil {
			return false, err
		}
	}

	exists, err := s.repo.Exists(ctx, id)

	if err != nil {
		if s.OnExistsFail != nil {
			if err := s.OnExistsFail(ctx, err); err != nil {
				return false, err
			}
		}
	}

	if s.AfterExists != nil {
		if err := s.AfterExists(ctx); err != nil {
			return exists, err
		}
	}

	return exists, err
}

func (s *CrudServiceWithHooks[E]) Random(ctx context.Context) (E, error) {
	var entity E
	if s.BeforeRandom != nil {
		if err := s.BeforeRandom(ctx); err != nil {
			return entity, err
		}
	}

	entity, err := s.repo.Random(ctx)

	if err != nil {
		if s.OnRandomFail != nil {
			if err := s.OnRandomFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterRandom != nil {
		if err := s.AfterRandom(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err

}

func (s *CrudServiceWithHooks[E]) First(ctx context.Context, f filter.Filter) (E, error) {
	var entity E
	if s.BeforeFirst != nil {
		if err := s.BeforeFirst(ctx); err != nil {
			return entity, err
		}
	}

	entity, err := s.repo.First(ctx, f)

	if err != nil {
		if s.OnFirstFail != nil {
			if err := s.OnFirstFail(ctx, err, entity); err != nil {
				return entity, err
			}
		}
	}

	if s.AfterFirst != nil {
		if err := s.AfterFirst(ctx, entity); err != nil {
			return entity, err
		}
	}

	return entity, err

}

func (s *CrudServiceWithHooks[E]) ComboBox(ctx context.Context, pageable pagination.Pageable, f filter.Filter, relations []relation.Relation, orderBys []order.OrderBy) (pagination.Page[common.ComboOption], error) {
	if s.BeforeCombo != nil {
		if err := s.BeforeCombo(ctx); err != nil {
			return pagination.NewPage([]common.ComboOption{}, 0, 0, 0, 0), err
		}
	}

	page, err := s.repo.ComboBox(ctx, pageable, f, relations, orderBys)

	if err != nil {
		if s.OnComboFail != nil {
			if err := s.OnComboFail(ctx, err); err != nil {
				return page, err
			}
		}
	}

	if s.AfterCombo != nil {
		if err := s.AfterCombo(ctx); err != nil {
			return page, err
		}
	}

	return page, err
}

func (s *CrudServiceWithHooks[E]) AddHook(hook ServiceHook, f interface{}) {
	switch hook {
	case BeforeCreate:
		s.BeforeCreate = f.(func(ctx context.Context, payload E) error)
	case AfterCreate:
		s.AfterCreate = f.(func(ctx context.Context, createdEntity E) error)
	case OnCreateFail:
		s.OnCreateFail = f.(func(ctx context.Context, err error, failedEntity E) error)
	case BeforeUpdate:
		s.BeforeUpdate = f.(func(ctx context.Context, payload E) error)
	case AfterUpdate:
		s.AfterUpdate = f.(func(ctx context.Context, updatedEntity E) error)
	case OnUpdateFail:
		s.OnUpdateFail = f.(func(ctx context.Context, err error, failedEntity E) error)
	case BeforeDelete:
		s.BeforeDelete = f.(func(ctx context.Context, payload E) error)
	case AfterDelete:
		s.AfterDelete = f.(func(ctx context.Context, deletedEntity E) error)
	case OnDeleteFail:
		s.OnDeleteFail = f.(func(ctx context.Context, err error, failedEntity E) error)
	case BeforeFind:
		s.BeforeFind = f.(func(ctx context.Context, entities ...E) error)
	case AfterFind:
		s.AfterFind = f.(func(ctx context.Context, entities ...E) error)
	case OnFindFail:
		s.OnFindFail = f.(func(ctx context.Context, err error, entities ...E) error)
	case BeforeCount:
		s.BeforeCount = f.(func(ctx context.Context) error)
	case AfterCount:
		s.AfterCount = f.(func(ctx context.Context) error)
	case OnCountFail:
		s.OnCountFail = f.(func(ctx context.Context, err error) error)
	case BeforeAssoc:
		s.BeforeAssoc = f.(func(ctx context.Context) error)
	case AfterAssoc:
		s.AfterAssoc = f.(func(ctx context.Context, entity E) error)
	case OnAssocFail:
		s.OnAssocFail = f.(func(ctx context.Context, err error, entity E) error)
	case BeforeDissoc:
		s.BeforeDissoc = f.(func(ctx context.Context) error)
	case AfterDissoc:
		s.AfterDissoc = f.(func(ctx context.Context, entity E) error)
	case OnDissocFail:
		s.OnDissocFail = f.(func(ctx context.Context, err error, entity E) error)
	case BeforeExists:
		s.BeforeExists = f.(func(ctx context.Context) error)
	case AfterExists:
		s.AfterExists = f.(func(ctx context.Context) error)
	case OnExistsFail:
		s.OnExistsFail = f.(func(ctx context.Context, err error) error)
	case BeforeRandom:
		s.BeforeRandom = f.(func(ctx context.Context) error)
	case AfterRandom:
		s.AfterRandom = f.(func(ctx context.Context, entity E) error)
	case OnRandomFail:
		s.OnRandomFail = f.(func(ctx context.Context, err error, entity E) error)
	case BeforeFirst:
		s.BeforeFirst = f.(func(ctx context.Context) error)
	case AfterFirst:
		s.AfterFirst = f.(func(ctx context.Context, entity E) error)
	case OnFirstFail:
		s.OnFirstFail = f.(func(ctx context.Context, err error, entity E) error)
	case BeforeCombo:
		s.BeforeCombo = f.(func(ctx context.Context) error)
	case AfterCombo:
		s.AfterCombo = f.(func(ctx context.Context) error)
	case OnComboFail:
		s.OnComboFail = f.(func(ctx context.Context, err error) error)
	}
}

func (s *CrudServiceWithHooks[E]) RemoveHook(hook ServiceHook) {
	switch hook {
	case BeforeCreate:
		s.BeforeCreate = nil
	case AfterCreate:
		s.AfterCreate = nil
	case OnCreateFail:
		s.OnCreateFail = nil
	case BeforeUpdate:
		s.BeforeUpdate = nil
	case AfterUpdate:
		s.AfterUpdate = nil
	case OnUpdateFail:
		s.OnUpdateFail = nil
	case BeforeDelete:
		s.BeforeDelete = nil
	case AfterDelete:
		s.AfterDelete = nil
	case OnDeleteFail:
		s.OnDeleteFail = nil
	case BeforeFind:
		s.BeforeFind = nil
	case AfterFind:
		s.AfterFind = nil
	case OnFindFail:
		s.OnFindFail = nil
	case BeforeCount:
		s.BeforeCount = nil
	case AfterCount:
		s.AfterCount = nil
	case OnCountFail:
		s.OnCountFail = nil
	case BeforeAssoc:
		s.BeforeAssoc = nil
	case AfterAssoc:
		s.AfterAssoc = nil
	case OnAssocFail:
		s.OnAssocFail = nil
	case BeforeDissoc:
		s.BeforeDissoc = nil
	case AfterDissoc:
		s.AfterDissoc = nil
	case OnDissocFail:
		s.OnDissocFail = nil
	case BeforeExists:
		s.BeforeExists = nil
	case AfterExists:
		s.AfterExists = nil
	case OnExistsFail:
		s.OnExistsFail = nil
	case BeforeRandom:
		s.BeforeRandom = nil
	case AfterRandom:
		s.AfterRandom = nil
	case OnRandomFail:
		s.OnRandomFail = nil
	case BeforeFirst:
		s.BeforeFirst = nil
	case AfterFirst:
		s.AfterFirst = nil
	case OnFirstFail:
		s.OnFirstFail = nil
	case BeforeCombo:
		s.BeforeCombo = nil
	case AfterCombo:
		s.AfterCombo = nil
	case OnComboFail:
		s.OnComboFail = nil
	}
}

func (s *CrudServiceWithHooks[E]) GetHook(hook ServiceHook) (interface{}, bool) {
	switch hook {
	case BeforeCreate:
		return s.BeforeCreate, s.BeforeCreate != nil
	case AfterCreate:
		return s.AfterCreate, s.AfterCreate != nil
	case OnCreateFail:
		return s.OnCreateFail, s.OnCreateFail != nil
	case BeforeUpdate:
		return s.BeforeUpdate, s.BeforeUpdate != nil
	case AfterUpdate:
		return s.AfterUpdate, s.AfterUpdate != nil
	case OnUpdateFail:
		return s.OnUpdateFail, s.OnUpdateFail != nil
	case BeforeDelete:
		return s.BeforeDelete, s.BeforeDelete != nil
	case AfterDelete:
		return s.AfterDelete, s.AfterDelete != nil
	case OnDeleteFail:
		return s.OnDeleteFail, s.OnDeleteFail != nil
	case BeforeFind:
		return s.BeforeFind, s.BeforeFind != nil
	case AfterFind:
		return s.AfterFind, s.AfterFind != nil
	case OnFindFail:
		return s.OnFindFail, s.OnFindFail != nil
	case BeforeCount:
		return s.BeforeCount, s.BeforeCount != nil
	case AfterCount:
		return s.AfterCount, s.AfterCount != nil
	case OnCountFail:
		return s.OnCountFail, s.OnCountFail != nil
	case BeforeAssoc:
		return s.BeforeAssoc, s.BeforeAssoc != nil
	case AfterAssoc:
		return s.AfterAssoc, s.AfterAssoc != nil
	case OnAssocFail:
		return s.OnAssocFail, s.OnAssocFail != nil
	case BeforeDissoc:
		return s.BeforeDissoc, s.BeforeDissoc != nil
	case AfterDissoc:
		return s.AfterDissoc, s.AfterDissoc != nil
	case OnDissocFail:
		return s.OnDissocFail, s.OnDissocFail != nil
	case BeforeExists:
		return s.BeforeExists, s.BeforeExists != nil
	case AfterExists:
		return s.AfterExists, s.AfterExists != nil
	case OnExistsFail:
		return s.OnExistsFail, s.OnExistsFail != nil
	case BeforeRandom:
		return s.BeforeRandom, s.BeforeRandom != nil
	case AfterRandom:
		return s.AfterRandom, s.AfterRandom != nil
	case OnRandomFail:
		return s.OnRandomFail, s.OnRandomFail != nil
	case BeforeFirst:
		return s.BeforeFirst, s.BeforeFirst != nil
	case AfterFirst:
		return s.AfterFirst, s.AfterFirst != nil
	case OnFirstFail:
		return s.OnFirstFail, s.OnFirstFail != nil
	case BeforeCombo:
		return s.BeforeCombo, s.BeforeCombo != nil
	case AfterCombo:
		return s.AfterCombo, s.AfterCombo != nil
	case OnComboFail:
		return s.OnComboFail, s.OnComboFail != nil
	default:
		return nil, false
	}
}

func (s *CrudServiceWithHooks[E]) GetRepo() CrudService[E] {
	return s.repo
}

func (s *CrudServiceWithHooks[E]) SetRepo(repo CrudService[E]) {
	s.repo = repo
}

func (s *CrudServiceWithHooks[E]) AddBeforeCreateHook(f func(ctx context.Context, payload E) error) {
	s.BeforeCreate = f
}

func (s *CrudServiceWithHooks[E]) AddAfterCreateHook(f func(ctx context.Context, createdEntity E) error) {
	s.AfterCreate = f
}

func (s *CrudServiceWithHooks[E]) AddOnCreateFailHook(f func(ctx context.Context, err error, failedEntity E) error) {
	s.OnCreateFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeUpdateHook(f func(ctx context.Context, payload E) error) {
	s.BeforeUpdate = f
}

func (s *CrudServiceWithHooks[E]) AddAfterUpdateHook(f func(ctx context.Context, updatedEntity E) error) {
	s.AfterUpdate = f
}

func (s *CrudServiceWithHooks[E]) AddOnUpdateFailHook(f func(ctx context.Context, err error, failedEntity E) error) {
	s.OnUpdateFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeDeleteHook(f func(ctx context.Context, payload E) error) {
	s.BeforeDelete = f
}

func (s *CrudServiceWithHooks[E]) AddAfterDeleteHook(f func(ctx context.Context, deletedEntity E) error) {
	s.AfterDelete = f
}

func (s *CrudServiceWithHooks[E]) AddOnDeleteFailHook(f func(ctx context.Context, err error, failedEntity E) error) {
	s.OnDeleteFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeFindHook(f func(ctx context.Context, entities ...E) error) {
	s.BeforeFind = f
}

func (s *CrudServiceWithHooks[E]) AddAfterFindHook(f func(ctx context.Context, entities ...E) error) {
	s.AfterFind = f
}

func (s *CrudServiceWithHooks[E]) AddOnFindFailHook(f func(ctx context.Context, err error, entities ...E) error) {
	s.OnFindFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeCountHook(f func(ctx context.Context) error) {
	s.BeforeCount = f
}

func (s *CrudServiceWithHooks[E]) AddAfterCountHook(f func(ctx context.Context) error) {
	s.AfterCount = f
}

func (s *CrudServiceWithHooks[E]) AddOnCountFailHook(f func(ctx context.Context, err error) error) {
	s.OnCountFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeAssocHook(f func(ctx context.Context) error) {
	s.BeforeAssoc = f
}

func (s *CrudServiceWithHooks[E]) AddAfterAssocHook(f func(ctx context.Context, entity E) error) {
	s.AfterAssoc = f
}

func (s *CrudServiceWithHooks[E]) AddOnAssocFailHook(f func(ctx context.Context, err error, entity E) error) {
	s.OnAssocFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeDissocHook(f func(ctx context.Context) error) {
	s.BeforeDissoc = f
}

func (s *CrudServiceWithHooks[E]) AddAfterDissocHook(f func(ctx context.Context, entity E) error) {
	s.AfterDissoc = f
}

func (s *CrudServiceWithHooks[E]) AddOnDissocFailHook(f func(ctx context.Context, err error, entity E) error) {
	s.OnDissocFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeExistsHook(f func(ctx context.Context) error) {
	s.BeforeExists = f
}

func (s *CrudServiceWithHooks[E]) AddAfterExistsHook(f func(ctx context.Context) error) {
	s.AfterExists = f
}

func (s *CrudServiceWithHooks[E]) AddOnExistsFailHook(f func(ctx context.Context, err error) error) {
	s.OnExistsFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeRandomHook(f func(ctx context.Context) error) {
	s.BeforeRandom = f
}

func (s *CrudServiceWithHooks[E]) AddAfterRandomHook(f func(ctx context.Context, entity E) error) {
	s.AfterRandom = f
}

func (s *CrudServiceWithHooks[E]) AddOnRandomFailHook(f func(ctx context.Context, err error, entity E) error) {
	s.OnRandomFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeFirstHook(f func(ctx context.Context) error) {
	s.BeforeFirst = f
}

func (s *CrudServiceWithHooks[E]) AddAfterFirstHook(f func(ctx context.Context, entity E) error) {
	s.AfterFirst = f
}

func (s *CrudServiceWithHooks[E]) AddOnFirstFailHook(f func(ctx context.Context, err error, entity E) error) {
	s.OnFirstFail = f
}

func (s *CrudServiceWithHooks[E]) AddBeforeComboHook(f func(ctx context.Context) error) {
	s.BeforeCombo = f
}

func (s *CrudServiceWithHooks[E]) AddAfterComboHook(f func(ctx context.Context) error) {
	s.AfterCombo = f
}

func (s *CrudServiceWithHooks[E]) AddOnComboFailHook(f func(ctx context.Context, err error) error) {
	s.OnComboFail = f
}
