package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

const indexHTML = `
<!DOCTYPE HTML>
<html>
  <head>
    <meta charset="utf-8">
    <title>Hi React</title>
  </head>
  <body>
    <div id='root'></div>
    <script src="` + cdnReact + `"></script>
    <script src="` + cdnReactDom + `"></script>
    <script src="` + cdnBabelStandalone + `"></script>
    <script src="` + cdnAxios + `"></script>
    <script src="/js/app.jsx" type="text/babel"></script>
  </body>
</html>
`

func authenticate(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := r.Cookie("username")
		if err != nil {
			io.WriteString(w, "Not authenticated")
			return
		}
		token, err := r.Cookie("token")
		if err != nil {
			io.WriteString(w, "Not authenticated")
			return
		}

		log.Println(username.Value, token.Value)
		authenticated, err := env.authenticateUser(username.Value, token.Value)
		if !authenticated || err != nil {
			io.WriteString(w, "Not authenticated")
		} else {
			io.WriteString(w, "Authenticated")
		}

	}
}
func login(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Println(buf.String())
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		log.Println(r.Form)
	}
}
func register(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Println(buf.String())

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		log.Println(r.Form)
	}
}
func HelloServer(w http.ResponseWriter, r *http.Request) {
	nextDay := time.Now().Add(24 * time.Hour)
	midnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
	cookieUsername := http.Cookie{Name: "username", Value: "cezkuj", Expires: midnight}
	cookieToken := http.Cookie{Name: "token", Value: "secr3t", Expires: midnight}
	http.SetCookie(w, &cookieUsername)
	http.SetCookie(w, &cookieToken)
	for _, cookie := range r.Cookies() {
		log.Println(w, cookie.Name)
	}
	io.WriteString(w, indexHTML)
}
func main() {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	env := Env{db: db}
	/*
		token, err := env.loginUser("admin", "secr3t")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(token)
		auth, err := env.authenticateUser("admin", token)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(auth)
	*/
	router := mux.NewRouter()
	router.HandleFunc("/", HelloServer)
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("assets")))
	router.HandleFunc("/authenticate", authenticate(env))
	router.HandleFunc("/login", login(env)).Methods("POST")
	router.HandleFunc("/register", register(env)).Methods("POST")
	serveMux := &http.ServeMux{}
	serveMux.Handle("/", router)
	srv := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      serveMux,
	}
	log.Println(srv.ListenAndServe())
}
