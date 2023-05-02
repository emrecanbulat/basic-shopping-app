package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Hello(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	Hello(res, req)

	exp := "Hello From Shopping App"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s gog %s", exp, act)
	}
}
