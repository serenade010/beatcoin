package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create middleware for loading and saving sessions for dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/model/rank", dynamic.ThenFunc(app.modelRank))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodGet, "/model/view/:id", dynamic.ThenFunc(app.modelView))
	router.Handler(http.MethodGet, "/model/play", dynamic.ThenFunc(app.modelPlay))
	router.Handler(http.MethodGet, "/model/create", dynamic.ThenFunc(app.modelCreate))
	router.Handler(http.MethodPost, "/model/create", dynamic.ThenFunc(app.modelCreateModel))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)

}
