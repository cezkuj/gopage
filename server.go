package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
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

type Cfg struct {
	mysqlHost     string
	mysqlUser     string
	mysqlPassword string
	mysqlDatabase string
	production    bool
}

func parseCfg() Cfg {
	return Cfg{mysqlHost: getenv("MYSQL_HOST"), mysqlUser: getenv("MYSQL_USER"), mysqlPassword: getenv("MYSQL_PASSWORD"), mysqlDatabase: getenv("MYSQL_DATABASE"), production: getenv("PRODUCTION") != "0"}

}
func getenv(envVar string) string {
	varName := os.Getenv(envVar)
	if varName == "" {
		log.Fatal(envVar + " is not set")
	}
	return varName
}
func authenticate(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		username, err := r.Cookie("username")
		if err != nil {
			log.Println("Authentication failure due to " + err.Error())
			io.WriteString(w, "Not authenticated")
			return
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Println(err)
			io.WriteString(w, "Not authenticated")
			return
		}
		userPresent, err := env.userIsPresent(username.Value)
		if err != nil {
			log.Println(err)
			io.WriteString(w, "Not authenticated")
			return
		}
		if !userPresent {
			io.WriteString(w, "Not authenticated")
			log.Println("user not present", username.Value)
			return

		}
		authenticated, err := env.authenticateUser(username.Value, Token(token.Value))
		if err != nil {
			log.Println(err)
			io.WriteString(w, "Not authenticated")
			return
		}
		if !authenticated {
			io.WriteString(w, "Not authenticated")
			return
		}
		io.WriteString(w, "Authenticated")

	}
}
func login(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dat, err := parseReaderToJson(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		username := dat["username"]
		password := dat["password"]
		userPresent, err := env.userIsPresent(username)
		if err != nil {
			log.Println(err)
			return
		}
		if !userPresent {
			io.WriteString(w, "User is not present")
			return
		}
		passwordCorrect, err := env.passwordIsCorrect(username, password)
		if err != nil {
			log.Println(err)
			return
		}
		if !passwordCorrect {
			io.WriteString(w, "Password is not correct")
			return
		}
		token, err := env.updateToken(username)
		if err != nil {
			log.Println(err)
			return
		}
		setCookies(w, username, token)
		io.WriteString(w, "Logging in")

	}
}
func register(env Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dat, err := parseReaderToJson(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		username := dat["username"]
		password := dat["password"]
		userPresent, err := env.userIsPresent(username)
		if err != nil {
			log.Println(err)
			return
		}
		if userPresent {
			io.WriteString(w, "User is already present")
			return
		}
		token, err := env.createUser(username, password)
		if err != nil {
			log.Println(err)
			return
		}
		io.WriteString(w, "User created")
		setCookies(w, username, token)

	}

}
func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, indexHTML)
}
func startProdServer() {

}
func startDevServer() {

}
func main() {
	cfg := parseCfg()
	db, err := initDb(cfg.mysqlUser + ":" + cfg.mysqlPassword + "@tcp(" + cfg.mysqlHost + ")/" + cfg.mysqlDatabase)
	if err != nil {
		log.Fatal(err)
	}
	env := Env{db: db}
	CSRF := csrf.Protect([]byte("88D283B4F5882897B13DDE4D422D5"), csrf.Secure(false)) //secure false should be removed while running with TLS
	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("assets")))
	router.HandleFunc("/authenticate", authenticate(env))
	router.HandleFunc("/login", login(env)).Methods("POST")
	router.HandleFunc("/register", register(env)).Methods("POST")
	serveMux := &http.ServeMux{}
	serveMux.Handle("/", CSRF(router))
	srv := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      serveMux,
	}
	log.Println(srv.ListenAndServe())
}
func parseReaderToJson(reader io.Reader) (map[string]string, error) {
	var dat map[string]string
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	err := json.Unmarshal(buf.Bytes(), &dat)
	return dat, err
}
func setCookies(w http.ResponseWriter, username string, token Token) {
	cookies := createCookies(username, token)
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func createCookies(username string, token Token) []http.Cookie {
	nextDay := time.Now().Add(24 * time.Hour)
	midnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location()) //Sets time to midnight of the next day
	return []http.Cookie{http.Cookie{Name: "username", Value: username, Expires: midnight}, http.Cookie{Name: "token", Value: string(token), Expires: midnight}}

}
