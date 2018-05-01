package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticateWithoutCookie(t *testing.T) {
	env := setUp()
	req, err := http.NewRequest("GET", "/authenticate", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate(env))
	handler.ServeHTTP(rr, req)
	expected := "Not authenticated"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
func TestAuthenticateWithCookie(t *testing.T) {
	env := setUp()
	token, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/authenticate", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookies := createCookies(test_user, token)
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate(env))
	handler.ServeHTTP(rr, req)
	expected := "Authenticated"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestLogin(t *testing.T) {
	env := setUp()
	token, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Fatal(err)
	}
	dat := make(map[string]string)
	dat["username"] = test_user
	dat["token"] = string(token)
	body, err := json.Marshal(dat)
	if err != nil {
		t.Fatal(err)
	}
	body_reader := bytes.NewReader(body)
	req, err := http.NewRequest("POST", "/login", body_reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(login(env))
	handler.ServeHTTP(rr, req)
	expected := "Logging in"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestRegister(t *testing.T) {

}
