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
	req, err := http.NewRequest("GET", "/api/authenticate", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate(env))
	handler.ServeHTTP(rr, req)
	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("handler returned unexpected http code: got %v want %v",
			rr.Code, expected)
	}

}
func TestAuthenticateWithCookie(t *testing.T) {
	env := setUp()
	token, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/api/authenticate", nil)
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
	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("handler returned unexpected http code: got %v want %v",
			rr.Code, expected)
	}

}

func TestLogin(t *testing.T) {
	env := setUp()
	_, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Fatal(err)
	}
	dat := make(map[string]string)
	dat["username"] = test_user
	dat["password"] = test_user
	body, err := json.Marshal(dat)
	if err != nil {
		t.Fatal(err)
	}
	body_reader := bytes.NewReader(body)
	req, err := http.NewRequest("POST", "/api/login", body_reader)
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
	env := setUp()
	dat := make(map[string]string)
	dat["username"] = test_user
	dat["password"] = test_user
	body, err := json.Marshal(dat)
	if err != nil {
		t.Fatal(err)
	}
	body_reader := bytes.NewReader(body)
	req, err := http.NewRequest("POST", "/api/register", body_reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(register(env))
	handler.ServeHTTP(rr, req)
	expected := "User created"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
