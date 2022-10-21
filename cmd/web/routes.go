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

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/model/rank", app.modelRank)
	router.HandlerFunc(http.MethodGet, "/user/signup", app.userSignup)
	router.HandlerFunc(http.MethodGet, "/user/login", app.userLogin)
	router.HandlerFunc(http.MethodGet, "/model/view/:id", app.modelView)
	router.HandlerFunc(http.MethodGet, "/model/play", app.modelPlay)
	router.HandlerFunc(http.MethodGet, "/model/create", app.modelCreate)
	router.HandlerFunc(http.MethodPost, "/model/create", app.modelCreateModel)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)

}
