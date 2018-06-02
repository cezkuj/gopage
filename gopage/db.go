package gopage

import (
	"crypto/sha1"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	//"log"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Token string

type Env struct {
	db *sqlx.DB
}

func initDb(db_connection string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql",
		db_connection)

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
func (env Env) getUsers() ([]User, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users")
	return users, err
}

func (env Env) getUser(username string) (User, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users where username=?", username)
	if len(users) != 0 {
		return users[0], err
	}
	return User{}, errors.New("There is no user with such username")
}

func (env Env) userIsPresent(username string) (bool, error) {
	users := []User{}
	err := env.db.Select(&users, "SELECT * FROM users where username=?", username)
	if len(users) != 0 {
		return true, err
	}
	return false, err
}

func (env Env) createUser(username, password string) (Token, error) {
	token, validity := newTokenAndValidity()
	_, err := env.db.Exec("INSERT INTO users (username, hash, token, validity) VALUES (?, ?, ?, ?)", username, getSHA1Hash(password), token, validity)
	return token, err

}
func (env Env) passwordIsCorrect(username, password string) (bool, error) {
	user, err := env.getUser(username)
	if err != nil {
		return false, err
	}
	if getSHA1Hash(password) == user.Hash {
		return true, nil
	}
	return false, nil

}
func (env Env) updateToken(username string) (Token, error) {
	token, validity := newTokenAndValidity()
	_, err := env.db.Exec("UPDATE users SET token=?, validity=? where username=?", token, validity, username)
	return token, err

}
func (env Env) authenticateUser(username string, token Token) (bool, error) {
	user, err := env.getUser(username)
	if err != nil {
		return false, err
	}
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

func newTokenAndValidity() (Token, string) {
	now := time.Now()
	token := Token(randSeq(32))
	return token, strings.Join([]string{strconv.Itoa(now.Day()), now.Month().String(), strconv.Itoa(now.Year())}, " ")
}
