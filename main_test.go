package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)

	if err!= nil {
		t.Fatal(err)
	}

	recorder:= httptest.NewRecorder()

	hf := http.HandlerFunc(handlerHealth)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted%v", status, http.StatusOK)
	}
}

func TestRouter(t *testing.T) {
	r := newTestRouter("/healthz", handlerHealth)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL+"/healthz")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be %v Got %v", http.StatusOK, resp.StatusCode)
	}
}