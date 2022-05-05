package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/techarm/go-stripe/internal/driver"
	"github.com/techarm/go-stripe/internal/models"
)

const version = "1.0.0"
const cssVersion = "1"

var session *scs.Session

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Start http server in %s mode on port %d\n", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "dev", "Application enviornment {dev|prod}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "Backend API URL")
	flag.StringVar(&cfg.db.dsn, "dsn", "techarm:secret@tcp(localhost:3306)/widgets?parseTime=True&tls=false", "Data source name")
	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "[INFO ] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatalln(err)
	}
	defer conn.Close()

	// setup session
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Session:       session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
