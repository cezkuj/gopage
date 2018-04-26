package main

import "log"
import "testing"

const (
	db_connection = "django:djangopass@tcp(127.0.0.1:3306)/homepage"
	table_name    = "users"
)

func dropTable(table_name string) {
	db, err := initDb(db_connection)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("DROP TABLE " + table_name)
}
func setUp() Env {
	dropTable(table_name)
	db, err := initDb(db_connection)
	if err != nil {
		log.Fatal(err)
	}
	return Env{db: db}

}
func test(got interface{}, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
func TestInitDb(t *testing.T) {
	dropTable(table_name)
	db, err := initDb(db_connection)
	if err != nil {
		t.Error(err)
	}
	tables := []string{}
	db.Select(&tables, "SHOW TABLES")
	for _, table := range tables {
		if table == table_name {
			return
		}
	}
	t.Errorf("TestInitDB FAILED - Table not found")

}
func TestUserIsPresentAndCreateUser(t *testing.T) {
	env := setUp()
	test_user := "test_abc"
	got, err := env.userIsPresent(test_user)
	if err != nil {
		t.Error(err)
	}
	want := false
	test(got, want, t)
	_, err = env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	got, err = env.userIsPresent(test_user)
	if err != nil {
		t.Error(err)
	}
	want = true
	test(got, want, t)
	got, err = env.userIsPresent("abc")
	if err != nil {
		t.Error(err)
	}
	want = false
	test(got, want, t)

}
func TestPasswordIsCorrect(t *testing.T) {

}
func TestUpdateToken(t *testing.T) {

}
func TestAuthenticate(t *testing.T) {

}
