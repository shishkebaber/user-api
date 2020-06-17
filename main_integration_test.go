// Integration testing with PostgresDB
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nicholasjackson/env"

	"github.com/shishkebaber/user-api/data"
	"github.com/shishkebaber/user-api/server"
	log "github.com/sirupsen/logrus"
	"net/http/httptest"

	"net/http"
	"testing"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
    id SERIAL,
    first-name TEXT NOT NULL,
    last-name TEXT NOT NULL,
	nickname TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL,
	country TEXT NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
)`

func clearTable(pgdb *data.UserPostgresDb) {
	pgdb.Pool.Exec(context.Background(), "DELETE FROM users")
	pgdb.Pool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func createTable(pgdb *data.UserPostgresDb) {
	if _, err := pgdb.Pool.Exec(context.Background(), tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

var s *server.Server
var bindAddressTest = env.String("BIND_ADDRESS", false, ":9090", "Server bind address")
var postgresBindAddressTest = env.String("PG_BIND_ADDRESS", false, "", "Server bind address")

func init() {
	s = server.NewServer(bindAddressTest, postgresBindAddressTest)
	createTable(s.UserHandlers.Db.(*data.UserPostgresDb))
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Handler.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d. Got: %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable(s.UserHandlers.Db.(*data.UserPostgresDb))

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got: %s", body)
	}
}

func TestCreateUser(t *testing.T) {

	clearTable(s.UserHandlers.Db.(*data.UserPostgresDb))

	var jsonStr = []byte(`{"first-name":"test user", "last-name":"test lname", "nickname":"tester", "email":"test@test.test". "password":"easy", "country":"testia"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetUser(t *testing.T) {
	clearTable(s.UserHandlers.Db.(*data.UserPostgresDb))
	args := make([]interface{}, 6)
	args = append(args, []string{"testname", "testsurn", "testnickname", "testemail", "testpassword", "testcountry"})
	keys := []string{"first-name", "last-name", "nickname", "email", "country"}
	values := []string{"testname", "testsurn", "testnickname", "testemail", "testcountry"}

	_, err := s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Exec(context.Background(), "INSERT INTO users(first-name,last-name,nickname,email,password,country) VALUES($1,$2,$3,$4,$5,$6)",
		args...)
	if err != nil {
		s.Logger.Error("Unable to insert test data")
	}
	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var rM map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &rM)

	if len(values) != len(rM) {
		t.Errorf("Expected an response items count to be equal to fields count. Got: %d", len(rM))
	}

	for i, v := range keys {
		if _, ok := rM[v]; !ok {
			t.Errorf("Expected existance of key %s. ", v)
		}
		if rM[v] != values[i] {
			t.Errorf("Expected value %s. Got: %s", values[i], rM[v])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	clearTable(s.UserHandlers.Db.(*data.UserPostgresDb))
	args := make([]interface{}, 6)
	args = append(args, []string{"testname", "testsurn", "testnickname", "testemail@email.cc", "testpassword", "testcountry"})

	_, err := s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Exec(context.Background(), "INSERT INTO users(first-name,last-name,nickname,email,password,country) VALUES($1,$2,$3,$4,$5,$6)",
		args...)
	if err != nil {
		s.Logger.Error("Unable to insert test data")
	}

	var jsonStr = []byte(`{"first-name":"test user", "last-name":"test lname", "nickname":"testonator", "email":"test@test.test", "country":"testia"}`)
	req, _ := http.NewRequest("PUT", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	reqGet, _ := http.NewRequest("GET", "/users/1", nil)
	reqGet.Header.Set("Content-Type", "application/json")
	responseGet := executeRequest(reqGet)
	var rM map[string]interface{}
	var input map[string]interface{}
	json.Unmarshal(responseGet.Body.Bytes(), &rM)
	json.Unmarshal(jsonStr, &input)

	if len(input) != len(rM) {
		t.Errorf("Expected an response items count to be equal to input fields count. Got: %d", len(rM))
	}

	for k, v := range rM {
		if _, ok := input[k]; !ok {
			t.Errorf("Expected existance of key %s. ", k)
		}
		if v != input[k] {
			t.Errorf("Expected value %s to be equal to input value %s.", v, input[k])
		}
	}
}

func TestDeleteUser(t *testing.T) {
	clearTable(s.UserHandlers.Db.(*data.UserPostgresDb))

	args := make([]interface{}, 6)
	args = append(args, []string{"testname", "testsurn", "testnickname", "testemail@email.cc", "testpassword", "testcountry"})

	_, err := s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Exec(context.Background(), "INSERT INTO users(first-name,last-name,nickname,email,password,country) VALUES($1,$2,$3,$4,$5,$6)",
		args...)
	if err != nil {
		s.Logger.Error("Unable to insert test data")
	}

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	reqGet, _ := http.NewRequest("GET", "/users/1", nil)
	reqGet.Header.Set("Content-Type", "application/json")
	responseGet := executeRequest(reqGet)

	checkResponseCode(t, http.StatusNotFound, responseGet.Code)
}
