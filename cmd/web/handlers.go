package main

import (
	"errors"
	"fmt"
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

	name := r.PostForm.Get("name")
	crypto := r.PostForm.Get("crypto")

	//Parse every column in the form
	ratio_of_train, err := strconv.ParseFloat(r.PostForm.Get("ratio_of_train"), 32)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	look_back, err := strconv.Atoi(r.PostForm.Get("look_back"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	forcast_days, err := strconv.Atoi(r.PostForm.Get("forcast_days"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	first_layer, err := strconv.Atoi(r.PostForm.Get("first_layer"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	second_layer, err := strconv.Atoi(r.PostForm.Get("second_layer"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	third_layer, err := strconv.Atoi(r.PostForm.Get("third_layer"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	epoch, err := strconv.Atoi(r.PostForm.Get("epoch"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	learning_rate, err := strconv.ParseFloat(r.PostForm.Get("learning_rate"), 32)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	batch_size, err := strconv.Atoi(r.PostForm.Get("batch_size"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var index_one, index_two, index_three string
	for i, item := range r.PostForm["index"] {
		switch i {
		case 0:
			index_one = item
		case 1:
			index_two = item
		case 2:
			index_three = item
		}

	}

	//Insert all column data into DB
	lastid, err := app.models.Insert(name, 2, float32(ratio_of_train), look_back, forcast_days, crypto, first_layer, second_layer, third_layer, index_one, index_two, index_three, float32(learning_rate), epoch, batch_size, 333)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/model/view/%d", lastid), http.StatusSeeOther)
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
