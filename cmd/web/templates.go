package main

import (
	"html/template"
	"path/filepath"

	"github.com/serenade010/beatcoin/internal/models"
)

type templateData struct {
	CurrentYear     int
	Model           *models.Model
	RankModels      []*models.Model
	MyModels        []*models.Model
	User            *models.User
	Form            any
	Flash           string
	IsAuthenticated bool
	UseridMatch     map[int]string
	Result          Response
}

func inc(i int) int {
	return i + 1
}

// func splitTimeArray[v int | float64](arrs [][]v) []v {
// 	var newarr []v

// 	for _, arr := range arrs {
// 		newarr = append(newarr, arr[0])
// 	}

// 	return newarr
// }

// func splitPRiceArray([][]float64) []int {

// }

var functions = template.FuncMap{
	"inc": inc,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil

}
