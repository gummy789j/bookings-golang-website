package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gummy789j/bookings/internal/config"
	"github.com/gummy789j/bookings/internal/driver"
	"github.com/gummy789j/bookings/internal/handlers"
	"github.com/gummy789j/bookings/internal/helpers"
	"github.com/gummy789j/bookings/internal/models"
	"github.com/gummy789j/bookings/internal/render"
)

const portNum = ":8081"

var app config.AppConfig

var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	ListenForMail()

	// from := "me@here.com"
	// auth := smtp.PlainAuth("", from, "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@here.com"}, []byte("Hello, world"))
	// if err != nil {
	// 	log.Println(err)
	// }

	//http.HandleFunc("/", handlers.Repo.Home)

	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port %s\n", portNum)

	//_ = http.ListenAndServe(portNum, nil)

	srv := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	// what am I going to put in the session
	// use to initialize
	// only be used on the initialization
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})

	// Build a new mail channal
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//  change this when in production
	app.InProduction = false

	// Build a new info logger for later
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// store the new info logger
	app.InfoLog = infoLog

	// Build a new error logger for later
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// store the new error logger
	app.ErrorLog = errorLog

	// Build a new Session manager and set some parameters
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// store seesion manager
	app.Session = session

	// connect with database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=gummy789j password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to database!")

	// CreateTemplateCache to help the development faster (do not need to re-execute when modified templates)
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	// store template cache
	app.TemplateCache = tc

	// whether using template cache or not
	app.UseCache = false

	// Build a new repository which is when you get the request and call handler, it can store the data and function that you need
	repo := handlers.NewRepo(&app, db)

	// Build a new handlers
	handlers.NewHandlers(repo)

	// After your handler end you need to render the template on the browser
	render.NewRenderer(&app)

	// Helper can help you to handle the error message
	helpers.NewHelpers(&app)

	return db, nil
}
