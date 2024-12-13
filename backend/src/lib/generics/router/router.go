package router

import (
	"fmt"
	"net/http"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/controller"
)

type HTTPMethod string

const (
	HTTPMethodGet    HTTPMethod = "GET"
	HTTPMethodPost   HTTPMethod = "POST"
	HTTPMethodPut    HTTPMethod = "PUT"
	HTTPMethodPatch  HTTPMethod = "PATCH"
	HTTPMethodDelete HTTPMethod = "DELETE"
)

// CrudRouter is a generic router that provides routing functionality.
// It is a wrapper around an http.ServeMux.
type CrudRouter[E common.Entity, D common.Entity] struct {
	*http.ServeMux
	baseRoute  string
	controller *controller.CrudController[E, D]
}

// NewRouter creates a new Router.
func NewRouter[E common.Entity, D common.Entity](controller *controller.CrudController[E, D]) *CrudRouter[E, D] {
	var zero E
	entityName := zero.GetEntityName()
	r := CrudRouter[E, D]{
		ServeMux:   http.NewServeMux(),
		baseRoute:  "/" + string(entityName),
		controller: controller,
	}

	r.Get("/", r.controller.FindAll())
	r.Get("/count", r.controller.Count())
	r.Get("/exists", r.controller.Exists())
	r.Get("/random", r.controller.Random())
	r.Get("/first", r.controller.First())
	r.Get("/combo", r.controller.Combo())
	r.Post("/", r.controller.Create())
	r.Get("/{id}", r.controller.Find())
	r.Put("/{id}", r.controller.Update())
	r.Patch("/{id}", r.controller.Update())
	r.Delete("/{id}", r.controller.Delete())
	r.Post("/{id}/{association}/{target}", r.controller.Associate())
	r.Delete("/{id}/{association}/{target}", r.controller.Dissociate())

	return &r
}

// GetBaseRoute returns the base route of the router.
func (r *CrudRouter[E, D]) GetBaseRoute() string {
	return r.baseRoute
}

func (r *CrudRouter[E, D]) GetRoute(method HTTPMethod, route string) string {
	return fmt.Sprintf("%s %s", method, r.baseRoute+route)
}

// AddRoute adds a route to the router.
func (r *CrudRouter[E, D]) AddRoute(method HTTPMethod, route string, handler http.HandlerFunc) {
	r.HandleFunc(r.GetRoute(method, route), handler)
}

func (r *CrudRouter[E, D]) Get(route string, handler http.HandlerFunc) {
	r.AddRoute(HTTPMethodGet, route, handler)
}

func (r *CrudRouter[E, D]) Post(route string, handler http.HandlerFunc) {
	r.AddRoute(HTTPMethodPost, route, handler)
}

func (r *CrudRouter[E, D]) Put(route string, handler http.HandlerFunc) {
	r.AddRoute(HTTPMethodPut, route, handler)
}

func (r *CrudRouter[E, D]) Patch(route string, handler http.HandlerFunc) {
	r.AddRoute(HTTPMethodPatch, route, handler)
}

func (r *CrudRouter[E, D]) Delete(route string, handler http.HandlerFunc) {
	r.AddRoute(HTTPMethodDelete, route, handler)
}

func (r *CrudRouter[E, D]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.ServeMux.ServeHTTP(w, req)
}

func Stack(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}
	}
}

func RouterStack(routers ...*CrudRouter[common.Entity, common.Entity]) *http.ServeMux {
	mux := http.NewServeMux()
	for _, router := range routers {
		mux.Handle(router.GetBaseRoute(), router)
	}
	return mux
}
