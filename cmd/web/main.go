package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/serenade010/beatcoin/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	models        *models.ModelModel
	templateCache map[string]*template.Template
}

func main() {

	// create info log ande error log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load env variable
	err := loadENV()
	if err != nil {
		errorlog.Fatal(err)
	}

	// add url adress and database flags
	addr := flag.String("addr", ":4000", "HTTP Network adress")
	dsn := flag.String("dsn", os.Getenv("DB_URL"), "SQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorlog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}

	app := &application{
		errorLog:      errorlog,
		infoLog:       infoLog,
		models:        &models.ModelModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func loadENV() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
