package data

import (
	"context"
	"github.com/sirupsen/logrus"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
    id SERIAL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
	nickname TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL,
	country TEXT NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
)`

func CreateTable(pgdb *UserPostgresDb, l *logrus.Logger) {
	l.Info("Creating table users")
	if _, err := pgdb.Pool.Exec(context.Background(), tableCreationQuery); err != nil {
		l.Fatal(err)
	}
}

func ClearTable(pgdb *UserPostgresDb, l *logrus.Logger) {
	l.Info("Clearing table users")
	if _, err := pgdb.Pool.Exec(context.Background(), "DELETE FROM users"); err != nil {
		l.Fatal(err)
	}

	if _, err := pgdb.Pool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1"); err != nil {
		l.Fatal(err)
	}
}
