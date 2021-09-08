package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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

func (f *Forms) Has(field string) bool {

	x := f.Get(field)
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

func (f *Forms) MinLength(field string, length int) bool {

	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

func (f *Forms) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
