package gopage

import "log"
import "testing"

const (
	db_connection = "django:djangopass@tcp(127.0.0.1:3306)/homepage"
	table_name    = "users"
	test_user     = "test_abc"
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
func TestGetUser(t *testing.T) {
	env := setUp()
	_, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	_, err = env.createUser(test_user+"2", test_user+"2")
	if err != nil {
		t.Error(err)
	}
	user, err := env.getUser(test_user)
	if err != nil {
		t.Error(err)
	}
	want := test_user
	got := user.Username
	test(got, want, t)
}
func TestGetUsers(t *testing.T) {
	env := setUp()
	_, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	_, err = env.createUser(test_user+"2", test_user+"2")
	if err != nil {
		t.Error(err)
	}
	users, err := env.getUsers()
	if err != nil {
		t.Error(err)
	}

	if len(users) != 2 {
		t.Errorf("TestGetUsers failed to invalid amount of users")
	}
}

func TestUserIsPresentAndCreateUser(t *testing.T) {
	env := setUp()
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
	env := setUp()
	_, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}

	want := true
	got, err := env.passwordIsCorrect(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	test(got, want, t)
	want = false
	got, err = env.passwordIsCorrect(test_user, test_user+"1")
	if err != nil {
		t.Error(err)
	}
	test(got, want, t)

}

func TestUpdateToken(t *testing.T) {
	env := setUp()
	_, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	token, err := env.updateToken(test_user)
	if err != nil {
		t.Error(err)
	}
	user, err := env.getUser(test_user)
	if err != nil {
		t.Error(err)
	}
	want := token
	got := user.Token
	test(got, want, t)
}

func TestAuthenticateUser(t *testing.T) {
	env := setUp()
	token, err := env.createUser(test_user, test_user)
	if err != nil {
		t.Error(err)
	}
	want := true
	got, err := env.authenticateUser(test_user, token)
	if err != nil {
		t.Error(err)
	}
	test(got, want, t)
}
