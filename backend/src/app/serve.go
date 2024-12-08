package app

import (
	"net/http"

	"github.com/cmo7/folly4/src/lib/generics/router"
)

func Serve() {
	router := router.RouterStack()
	http.ListenAndServe(":8080", router)
}
