package firstservice_test

import (
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"testing"
	. "go-to-do/firstservice"
)

func TestFunc_GetTokenHandler(t *testing.T) {
	//check ErrInvalidNameOrPassword
	userBody := &User{Name: "Oleg", Password: "1" ,}
	b, err := json.Marshal(userBody)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBuffer(b)

	r := httptest.NewRequest("GET", "http://127.0.0.1:8080/api/get-token", body)
	w := httptest.NewRecorder()
	GetTokenHandler(w, r)

	expectedBody, _ := NewError(ErrInvalidNameOrPassword).Marshall()

	if w.Body.String() != string (expectedBody) {
		t.Fail()
	}

	userBody = &User{Name: "Oleg", Password: "322",}
	b, err = json.Marshal(userBody)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.NewBuffer(b)

	r = httptest.NewRequest("GET", "http://127.0.0.1:8080/api/get-token", body)
	w = httptest.NewRecorder()
	GetTokenHandler(w, r)

	token := &Token{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Mn0.IJMhcXdqJEqIrXUDP4B7JjAQCGKKxEkAhxU-rLA3djI"}

	expectedBody, _ = token.Marahsll()

	if w.Body.String() != string (expectedBody) {
		t.Fatal(w.Body.String())
		t.Fail()
	}

}

