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
	data.User = app.users.UserInfo(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) modelPlay(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.User = app.users.UserInfo(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	app.render(w, http.StatusOK, "play.html", data)
}

func (app *application) modelRank(w http.ResponseWriter, r *http.Request) {
	models, err := app.models.Best()
	if err != nil {
		app.serverError(w, err)
		return
	}
	idMapping := make(map[int]string, 10)
	for _, model := range models {
		user := app.users.UserInfo(model.Belongs_to)
		idMapping[model.Belongs_to] = user.Name
	}

	data := app.newTemplateData(r)
	data.RankModels = models
	data.UseridMatch = idMapping
	data.User = app.users.UserInfo(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	fmt.Println(data.UseridMatch)
	app.render(w, http.StatusOK, "rank.html", data)
}

func (app *application) myModelsView(w http.ResponseWriter, r *http.Request) {
	models, err := app.models.MyModels(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.MyModels = models
	fmt.Println(data.MyModels)
	data.User = app.users.UserInfo(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	app.render(w, http.StatusOK, "mymodel.html", data)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) modelTrain(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "train.html", data)
}

func (app *application) modelTrainPost(w http.ResponseWriter, r *http.Request) {
	// data := app.newTemplateData(r)
	// app.render(w, http.StatusOK, "train.html", data)
}




func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm

	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {

			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/model/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) modelCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = modelCreateForm{}
	data.User = app.users.UserInfo(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	app.render(w, http.StatusOK, "create.html", data)
}

type modelCreateForm struct {
	Id                  int
	Name                string         `form:"name"`
	Belongs_to          int            `form:"-"`
	Ratio_of_train      float32        `form:"ratio_of_train"`
	Look_back           int            `form:"look_back"`
	Forecast_days       int            `form:"forcast_days"`
	Crypto              string         `form:"crypto"`
	First_layer         int            `form:"first_layer"`
	Second_layer        int            `form:"second_layer"`
	Third_layer         int            `form:"third_layer"`
	First_index         sql.NullString `form:"-"`
	Second_index        sql.NullString `form:"-"`
	Third_index         sql.NullString `form:"-"`
	Learning_rate       float32        `form:"learning_rate"`
	Epoch               int            `form:"epoch"`
	Batch_size          int            `form:"batch_size"`
	Modelerr            sql.NullFloat64
	validator.Validator `form:"-"`
}

func (app *application) modelCreatePost(w http.ResponseWriter, r *http.Request) {

	var form modelCreateForm

	err := app.decodePostForm(r, &form)
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
	lastid, err := app.models.Insert(form.Name, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), form.Ratio_of_train, form.Look_back, form.Forecast_days, form.Crypto, form.First_layer, form.Second_layer, form.Third_layer, index_one, index_two, index_three, form.Learning_rate, form.Epoch, form.Batch_size, 333)
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
	if !app.models.Belong(model.Id, app.sessionManager.GetInt(r.Context(), "authenticatedUserID")) {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)
	data.Model = model

	app.render(w, http.StatusOK, "view.html", data)
}

// func myModelsView(w http.ResponseWriter, r *http.Request) {

// }
