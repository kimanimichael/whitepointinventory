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