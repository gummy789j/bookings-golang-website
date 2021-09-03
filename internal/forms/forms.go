package forms

import (
	"net/http"
	"net/url"
	"strings"
)

type Forms struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Forms {
	return &Forms{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Forms) Has(field string, r *http.Request) bool {

	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Valid return true if there are no errors, otherwise false
func (f *Forms) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Forms) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}
