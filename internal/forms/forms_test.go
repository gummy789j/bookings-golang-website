package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form show invalid when required fields exist")
	}

}

func TestForm_Has(t *testing.T) {

	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("Whatever")
	if has {
		t.Error("form show has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("form show does not has field when it does ")
	}
}

func TestForm_MinLength(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "123")
	form := New(postedData)

	form.MinLength("a", 4)

	if form.Valid() {
		t.Error("form shows the field shorter than required minimum length when it doesn't")
	}

	postedData = url.Values{}
	postedData.Add("b", "12345")
	form = New(postedData)

	form.MinLength("b", 4)

	if !form.Valid() {
		t.Error("form shows the field doesn't shorter than required minimum length when it does")
	}

	isError := form.Errors.Get("another_field")

	if isError != "" {
		t.Error("should not have an error, but get one")
	}
}

func TestForm_IsEmail(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("email1", "gu")
	form := New(postedData)

	form.IsEmail("email1")

	if form.Valid() {
		t.Error("form show valid email when it should not valid")
	}

	postedData = url.Values{}
	postedData.Add("email2", "gu@here.com")
	form = New(postedData)

	form.IsEmail("email2")

	if !form.Valid() {
		t.Error("form show invalid email when it should be valid")
	}

}
