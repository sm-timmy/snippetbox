package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLogger,
		infoLog:  infoLogger,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLogger,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLogger.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
}
