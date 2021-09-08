package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Before exit, run the whole testing in this package
	os.Exit(m.Run())
}

type myHandler struct{}

func (this *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
