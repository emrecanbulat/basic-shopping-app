package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	t.Parallel()
	s := httptest.NewServer(http.HandlerFunc(Hello))

	req, err := http.NewRequest(http.MethodGet, s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	want := "Hello From Shopping App"
	if string(body) != want {
		t.Errorf("Unexpected body returned. Want %q, got %q", want, string(body))
	}
}
