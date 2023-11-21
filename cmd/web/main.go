package main

import (
	"artchernov.ru/internal/models"
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	db             *sql.DB
	snippets       *models.SnippetModel
	templates      map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLogger,
		infoLog:  infoLogger,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templates:      templatesCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		//CipherSuites: []uint16{
		//	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		//	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		//	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		//	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		//},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLogger,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLogger.Fatal(err)
}
