package main

import (
	"artchernov.ru/internal/models"
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	db          *sql.DB
	snippets    *models.SnippetModel
	templates   map[string]*template.Template
	formDecoder *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "snippetbox:root@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLogger.Printf("Open db connection with %s", *dsn)
	db, err := initDB(*dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	infoLogger.Println("Db connection established")

	templatesCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog: errorLogger,
		infoLog:  infoLogger,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templates:   templatesCache,
		formDecoder: formDecoder,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLogger,
		Handler:  app.routes(),
	}

	infoLogger.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
}
