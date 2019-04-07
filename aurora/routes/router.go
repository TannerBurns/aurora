package routes

import (
	"net/http"

	"../controllers"
	"../models"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() (*mux.Router, *controllers.Controller) {
	controller := &controllers.Controller{Name: "API.Controller"}
	controller.Logger = models.NewLogger()
	lc, _ := models.NewConfig("./config/settings.conf")
	controller.Session = models.NewSession(lc)

	Routes := []Routes{
		StatusRoutes(controller),
		UserRoutes(controller),
		AuthenticationRoutes(controller),
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, routes := range Routes {
		for _, route := range routes {
			var handler http.Handler
			handler = route.HandlerFunc

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	return router, controller
}
