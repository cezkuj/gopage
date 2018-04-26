package main

import "log"
import "testing"

const table_name = "users_test"

func dropDb(table_name string) {
	db, err := initDb(table_name)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("DROP TABLE " + table_name)
}
func TestInitDb(t *testing.T) {
	dropDb(table_name)
	db, err := initDb(table_name)
	if err != nil {
		log.Fatal(err)
	}
	tables := []string{}
	db.Select(&tables, "SHOW TABLES")
	for _, table := range tables {
		if table == table_name {
			return
		}
	}
	log.Fatal("TestInitDB FAILED - Table not found")

}
func TestUserIsPresent(t *testing.T) {

}
func TestCreateUser(t *testing.T) {

}
func TestPasswordIsCorrect(t *testing.T) {

}
func TestUpdateToken(t *testing.T) {

}
func TestAuthenticate(t *testing.T) {

}
