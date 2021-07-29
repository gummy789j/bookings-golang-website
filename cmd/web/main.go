package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gummy789j/bookings/pkg/config"
	"github.com/gummy789j/bookings/pkg/handlers"
	"github.com/gummy789j/bookings/pkg/render"
)

const portNum = ":8081"

var app config.AppConfig

var session *scs.SessionManager

func main() {

	//  change this when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

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
