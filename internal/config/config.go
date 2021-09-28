package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/gummy789j/bookings/internal/models"
)

type AppConfig struct {
	UseCache      bool                          //是否開啟快取修改的功能
	TemplateCache map[string]*template.Template //以name為Key存放每一個new page Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
