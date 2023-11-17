package main

import (
	"artchernov.ru/internal/models"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	db        *sql.DB
	snippets  *models.SnippetModel
	templates map[string]*template.Template
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

	app := &application{
		errorLog: errorLogger,
		infoLog:  infoLogger,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templates: templatesCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLogger,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLogger.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
}
