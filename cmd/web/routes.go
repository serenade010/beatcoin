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

	// Unprotected
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/model/rank", dynamic.ThenFunc(app.modelRank))
	router.Handler(http.MethodGet, "/model/play", dynamic.ThenFunc(app.modelPlay))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Model Related Route
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/model/mymodel", protected.ThenFunc(app.myModelsView))
	router.Handler(http.MethodGet, "/model/view/:id", protected.ThenFunc(app.modelView))
	router.Handler(http.MethodGet, "/model/create", protected.ThenFunc(app.modelCreate))
	router.Handler(http.MethodPost, "/model/create", protected.ThenFunc(app.modelCreatePost))
	router.Handler(http.MethodPost, "/model/modify", protected.ThenFunc(app.modelModify))
	router.Handler(http.MethodGet, "/model/delete", protected.ThenFunc(app.modelDelete))
	router.Handler(http.MethodGet, "/model/train", protected.ThenFunc(app.modelTrain))
	router.Handler(http.MethodPost, "/model/train", protected.ThenFunc(app.modelTrainPost))
	router.Handler(http.MethodGet, "/model/result", protected.ThenFunc(app.modelTrainResult))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)

}
