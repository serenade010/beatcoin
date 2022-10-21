package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/serenade010/beatcoin/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) modelRank(w http.ResponseWriter, r *http.Request) {
	models, err := app.models.Best()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "rank.html", &templateData{Models: models})
}

func (app *application) modelPlay(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "play.html", nil)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "login.html", nil)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "signup.html", nil)
}

func (app *application) modelCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) modelCreateModel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.models.Insert("l'amour", 2, 0.8, 1, 1, "BTC", 33, 33, 33, "q", "w", "e", 0.9, 30, 30, 333)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//TODO: FIX this route
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) modelView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	model, err := app.models.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}
	app.render(w, http.StatusOK, "view.html", &templateData{Model: model})
}
