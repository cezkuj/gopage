package main

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type Person struct {
	Id    int
	First string
	Last  string
}

func test_db() {
	db, err := sqlx.Open("mysql",
		"django:djangopass@tcp(127.0.0.1:3306)/homepage")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable := `
          CREATE TABLE IF NOT EXISTS people (
          id SERIAL NOT NULL PRIMARY KEY,
          first TEXT NOT NULL,
          last TEXT NOT NULL);
        `
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CREATED")
	insertPerson := "INSERT INTO people(first, last) VALUES (?, ?)"
	_, err = db.Exec(insertPerson, "abc", "ccc")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INSERTED")
        people := []Person{}
	err = db.Select(&people, "SELECT * FROM people")
        if err != nil {
                log.Fatal(err)
        }

        log.Println(people)
	for i, person := range people {
		log.Println(i, person)
	}
	log.Println("SELECTED")

}
