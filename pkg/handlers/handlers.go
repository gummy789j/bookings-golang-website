package handlers

import (
	"net/http"

	"github.com/gummy789j/bookings/pkg/config"
	"github.com/gummy789j/bookings/pkg/models"
	"github.com/gummy789j/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (this *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	this.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (this *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := map[string]string{
		"test": "Hello, again",
	}

	remoteIP := this.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
