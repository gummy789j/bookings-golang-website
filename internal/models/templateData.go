package models

import "github.com/gummy789j/bookings/internal/forms"

// 原本是寫在handlers 但因為import cycle的 error
// 必須開一個Package
// golang不像C++允許import cycle

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //cross site request forgery token
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Forms
}
