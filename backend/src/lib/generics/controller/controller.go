package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cmo7/folly4/src/lib/generics"
	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/cmo7/folly4/src/lib/generics/relation"
	"github.com/cmo7/folly4/src/lib/generics/service"
	"github.com/google/uuid"
)

type Controller interface {
	Create() http.HandlerFunc
	Find() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	FindAll() http.HandlerFunc
	Count() http.HandlerFunc
	Associate() http.HandlerFunc
	Dissociate() http.HandlerFunc
	Exists() http.HandlerFunc
	Random() http.HandlerFunc
	First() http.HandlerFunc
	Combo() http.HandlerFunc
}

// CrudController is a generic controller that provides CRUD functionality, compatible with the service.CrudService.
// Every method returns an http.HandlerFunc that can be used to handle HTTP requests.

// The controller uses a generics.Mapper to map entities to DTOs and vice versa.
type CrudController[E common.Entity, D common.Entity] struct {
	service.CrudService[E]
	mapper generics.Mapper[E, D]
}

func NewController[E common.Entity, D common.Entity](crudService service.CrudService[E]) *CrudController[E, D] {
	return &CrudController[E, D]{
		CrudService: crudService,
	}
}

func (c *CrudController[E, D]) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body.
		var entity E
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the entity.
		createdEntity, err := c.CrudService.Create(r.Context(), entity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdEntity)
	}
}

func (c *CrudController[E, D]) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		relations := extractRelationsFromRequest(r)

		entity, err := c.CrudService.FindOne(r.Context(), uid, relations)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

func (c *CrudController[E, D]) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var entity E
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entity.SetID(uid)

		updatedEntity, err := c.CrudService.Update(r.Context(), entity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedEntity)
	}
}

func (c *CrudController[E, D]) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entity, err := c.CrudService.FindOne(r.Context(), uid, nil)

		err = c.CrudService.Delete(r.Context(), entity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *CrudController[E, D]) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageable := extractPageableFromRequest(r)
		filter := extractFilterFromRequest(r)
		relations := extractRelationsFromRequest(r)
		orderBys := extractOrderBysFromRequest(r)

		page, err := c.CrudService.FindAll(r.Context(), pageable, filter, relations, orderBys)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}
}

func (c *CrudController[E, D]) Count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := extractFilterFromRequest(r)

		count, err := c.CrudService.Count(r.Context(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(count)
	}
}

func (c *CrudController[E, D]) Associate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		association := r.URL.Query().Get("association")
		targetID := r.URL.Query().Get("target")
		targetUID, err := uuid.Parse(targetID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entity, err := c.CrudService.Associate(r.Context(), uid, association, targetUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

func (c *CrudController[E, D]) Dissociate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		association := r.URL.Query().Get("association")
		targetID := r.URL.Query().Get("target")
		targetUID, err := uuid.Parse(targetID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entity, err := c.CrudService.Dissociate(r.Context(), uid, association, targetUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

func (c *CrudController[E, D]) Exists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exists, err := c.CrudService.Exists(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exists)
	}
}

func (c *CrudController[E, D]) Random() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, err := c.CrudService.Random(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)

	}
}

func (c *CrudController[E, D]) First() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := extractFilterFromRequest(r)

		entity, err := c.CrudService.First(r.Context(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

func (c *CrudController[E, D]) Combo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entities, err := c.CrudService.ComboBox(r.Context(), extractPageableFromRequest(r), extractFilterFromRequest(r), extractRelationsFromRequest(r), extractOrderBysFromRequest(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

func extractPageableFromRequest(r *http.Request) pagination.Pageable {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	intSize, err := strconv.Atoi(size)
	if err != nil {
		intSize = 10
	}

	return pagination.Pageable{
		Page: intPage,
		Size: intSize,
	}
}

func extractFilterFromRequest(r *http.Request) filter.Filter {
	filterString := r.URL.Query().Get("filter")
	if filterString == "" {
		return nil
	}
	f, err := filter.Parse(filterString)
	if err != nil {
		return nil
	}
	return f
}

func extractOrderBysFromRequest(r *http.Request) []order.OrderBy {
	orderString := r.URL.Query().Get("order")
	if orderString == "" {
		return nil
	}
	order, err := order.Parse(orderString)
	if err != nil {
		return nil
	}
	return order
}

func extractRelationsFromRequest(r *http.Request) []relation.Relation {
	rel := r.URL.Query().Get("relations")
	if rel == "" {
		return nil
	}

	relations, err := relation.Parse(rel)
	if err != nil {
		return nil
	}

	return relations
}
