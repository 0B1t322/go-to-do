package firstservice_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	. "go-to-do/firstservice"
)

func TestFunc_ParseUserJSON(t *testing.T) {
	// checking for ErrNameTaken
	userBody := &User{Name: "Oleg",Password: "322"}
	b, err := json.Marshal(userBody)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBuffer(b)

	r := httptest.NewRequest("POST", "http://127.0.0.1:8080/create-user", body)
	w := httptest.NewRecorder()
	ParseJSONUser(w, r)
	
	e := NewError(ErrNameTaken)
	b, err = e.Marshall()
	if err != nil {
		t.Fatal(err)
	}

	excpectedBody := bytes.NewBuffer(b)
	if w.Body.String() != excpectedBody.String() {
		t.Log("Server ans:", w.Body.String())
		t.Log("Excpected:", excpectedBody.String())
		t.Fatal("Error with POST when name is taken")
	}

	// Checking for ErrNameFieldEmpty
	userBody = &User{Name: "", Password: "322"}
	b, err = json.Marshal(userBody)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.NewBuffer(b)

	r = httptest.NewRequest("POST", "http://127.0.0.1:8080/create-user", body)
	w = httptest.NewRecorder()
	ParseJSONUser(w, r)

	e = NewError(ErrNameFieldEmpty)
	b, err = e.Marshall()
	if err != nil {
		t.Fatal(err)
	}

	excpectedBody = bytes.NewBuffer(b)
	if w.Body.String() != excpectedBody.String() {
		t.Log("Server ans:", w.Body.String())
		t.Log("Excpected:", excpectedBody.String())
		t.Fatal("Error with POST when name field is empty")
	}

	// Checking for ErrPassFieldEmpty
	userBody = &User{Name: "Oleg", Password: ""}
	b, err = json.Marshal(userBody)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.NewBuffer(b)

	r = httptest.NewRequest("POST", "http://127.0.0.1:8080/create-user", body)
	w = httptest.NewRecorder()
	ParseJSONUser(w, r)

	e = NewError(ErrPassFieldEmpty)
	b, err = e.Marshall()
	if err != nil {
		t.Fatal(err)
	}

	excpectedBody = bytes.NewBuffer(b)
	if w.Body.String() != excpectedBody.String() {
		t.Log("Server ans:", w.Body.String())
		t.Log("Excpected:", excpectedBody.String())
		t.Fatal("Error with POST when pass field is empty")
	}
}