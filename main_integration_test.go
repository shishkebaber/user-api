// Integration testing with PostgresDB
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/shishkebaber/user-api/data"
	"github.com/shishkebaber/user-api/server"
	"net/http/httptest"
	"strings"

	"net/http"
	"testing"
)

var s *server.Server

func init() {
	s = server.NewServer()
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

func createData() {
	_, err := s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Exec(context.Background(), "INSERT INTO users(first_name,last_name,nickname,email,password,country) VALUES($1,$2,$3,$4,$5,$6)",
		"testname", "testsurn", "testnickname", "testemail", "testpassword", "testcountry")
	if err != nil {
		s.Logger.Error("Unable to insert test data ", err)
	}
}

func TestEmptyTable(t *testing.T) {
	data.ClearTable(s.UserHandlers.Db.(*data.UserPostgresDb), s.Logger)

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); strings.TrimSpace(body) != "null" {
		t.Errorf("Expected an empty array. Got: %s", body)
	}
}

func TestCreateUser(t *testing.T) {

	data.ClearTable(s.UserHandlers.Db.(*data.UserPostgresDb), s.Logger)

	var jsonStr = []byte(`{"first_name":"test user", "last_name":"test lname", "nickname":"tester", "email":"test@test.test", "password":"easy", "country":"testia"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetUser(t *testing.T) {
	data.ClearTable(s.UserHandlers.Db.(*data.UserPostgresDb), s.Logger)
	keys := []string{"first_name", "last_name", "nickname", "email", "country"}
	values := []string{"testname", "testsurn", "testnickname", "testemail", "testcountry"}

	createData()

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var rM []map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &rM)
	if err != nil {
		s.Logger.Error("Unable to unmarshal get result: ", err)
	}

	for i, v := range keys {
		if _, ok := rM[0][v]; !ok {
			t.Errorf("Expected existance of key %s. ", v)
		}
		if rM[0][v] != values[i] {
			t.Errorf("Expected value %s. Got: %s", values[i], rM[0][v])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	data.ClearTable(s.UserHandlers.Db.(*data.UserPostgresDb), s.Logger)
	createData()

	var jsonStr = []byte(`{"id":1, "first_name":"test user", "last_name":"test lname", "nickname":"testonator", "email":"test@test.test", "country":"testia"}`)
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
	data.ClearTable(s.UserHandlers.Db.(*data.UserPostgresDb), s.Logger)

	createData()

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	reqGet, _ := http.NewRequest("GET", "/users?id=1", nil)
	reqGet.Header.Set("Content-Type", "application/json")
	responseGet := executeRequest(reqGet)

	body := responseGet.Body.String()
	if strings.TrimSpace(body) != "null" {
		t.Errorf("Expected an empty array. Got: %s", body)
	}
}
