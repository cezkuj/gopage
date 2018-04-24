package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	//"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Env struct {
	db *sqlx.DB
}

func initDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql",
		"django:djangopass@tcp(127.0.0.1:3306)/homepage")

	if err != nil {
		return nil, err
	}
	createTable := `
          CREATE TABLE IF NOT EXISTS users (
          id SERIAL NOT NULL PRIMARY KEY,
          username TEXT NOT NULL,
          hash TEXT NOT NULL,
          token TEXT NOT NULL,
          validity TEXT NOT NULL);
        `
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return db, nil

}

func (env Env) registerUser(username, password string) (string, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users where username=?", username)
	if err != nil {
		return "", err
	}
	if len(users) != 0 {
		return "", errors.New("RegisterError: User is already present")
	} else {
		token, validity := newTokenAndValidity()
		_, err = env.db.Exec("INSERT INTO users (username, hash, token, validity) VALUES (?, ?, ?, ?)", username, getSHA1Hash(password), token, validity)
		return token, err
	}

}

func (env Env) loginUser(username, password string) (string, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users where username=?", username)
	if err != nil {
		return "", err
	}
	if len(users) != 1 {
		return "", errors.New("LoginError: User is not registered")
	}
	user := users[0]
	if getSHA1Hash(password) == user.Hash {

		token, validity := newTokenAndValidity()
		_, err = env.db.Exec("UPDATE users SET token=?, validity=?", token, validity)
		if err != nil {
			return "", nil
		}
		return token, nil
	}
	return "", errors.New("LoginError: Password is invalid")
}
func (env Env) authenticateUser(username, token string) (bool, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users where username=?", username)
	if err != nil {
		return false, err
	}
	if len(users) != 1 {
		return false, errors.New("LoginError: User is not registered")
	}
	user := users[0]
	if user.Token == token {
		return true, nil
	}
	return false, nil

}
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getSHA1Hash(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func newTokenAndValidity() (string, string) {
	now := time.Now()
	token := randSeq(32)
	return token, strings.Join([]string{strconv.Itoa(now.Day()), now.Month().String(), strconv.Itoa(now.Year())}, " ")
}
