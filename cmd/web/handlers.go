package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/serenade010/beatcoin/internal/models"
	"github.com/serenade010/beatcoin/internal/validator"
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
	data.Form = moedlCreateForm{
		Name:   "",
		Crypto: "",
	}
	app.render(w, http.StatusOK, "create.html", data)
}

type moedlCreateForm struct {
	Id             int
	Name           string
	Belongs_to     int
	Ratio_of_train float32
	Look_back      int
	Forecast_days  int
	Crypto         string
	First_layer    int
	Second_layer   sql.NullInt16
	Third_layer    sql.NullInt16
	First_index    sql.NullString
	Second_index   sql.NullString
	Third_index    sql.NullString
	Learning_rate  float32
	Epoch          int
	Batch_size     int
	Modelerr       sql.NullFloat64
	validator.Validator
}

func (app *application) modelCreateModel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

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

	form := moedlCreateForm{
		Name:   r.PostForm.Get("name"),
		Crypto: r.PostForm.Get("crypto"),
	}

	// Initialize a map holding any validations errors from the erroe fields
	//TODO: add a validation for rest of the fields

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 20), "name", "This field cannot be more than 20 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	//Insert all column data into DB
	lastid, err := app.models.Insert(form.Name, 2, float32(ratio_of_train), look_back, forcast_days, form.Crypto, first_layer, second_layer, third_layer, index_one, index_two, index_three, float32(learning_rate), epoch, batch_size, 333)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Model successfully created!")

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

	data := app.newTemplateData(r)
	data.Model = model

	app.render(w, http.StatusOK, "view.html", data)
}
