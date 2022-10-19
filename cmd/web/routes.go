package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/model/rank", app.modelRank)
	mux.HandleFunc("/model/create", app.modelCreate)
	mux.HandleFunc("/user/signup", app.userSignup)
	mux.HandleFunc("/user/login", app.userLogin)

	return mux
}
