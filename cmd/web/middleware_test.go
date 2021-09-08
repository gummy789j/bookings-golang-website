package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var myh http.HandlerFunc

	h := NoSurf(&myh)

	switch h.(type) {
	case http.Handler:

	default:
		t.Error("type is not http.Handler\n")
	}
}
