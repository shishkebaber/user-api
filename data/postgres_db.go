package data

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"strconv"
)

type UserPostgresDb struct {
	log  *logrus.Logger
	Pool *pgxpool.Pool
}

// Creates new postgresDB object with logger and pool
func NewPgDb(logger *logrus.Logger, dbUrl string) *UserPostgresDb {
	pool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		logrus.Error("Unable to connect to PostgresDB \n ", err)
	}
	u := &UserPostgresDb{
		logger,
		pool,
	}

	return u
}

//Creates a User
func (pdb *UserPostgresDb) AddUser(user User) error {
	conn, err := pdb.Pool.Acquire(context.Background())
	if err != nil {
		pdb.log.Error("Unable to acquire DB connection. \n", err)
		return err
	}
	defer conn.Release()

	hashedPassword, err := hashSaltPassword(user.Password)
	if err != nil {
		pdb.log.Error("Failed to hash password \n", err)
		return err
	}

	row := conn.QueryRow(context.Background(),
		"INSERT INTO users (first_name, last_name, nickname, password, email, country) VALUES ($1, $2, $3, $4, $5, %6) RETURNING id",
		user.FirstName, user.LastName, user.Nickname, hashedPassword, user.Email, user.Country)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		pdb.log.Error("Failed to insert user \n", err)
		return err
	}
	pdb.log.Printf("User with ID:%d created \n", id)
	return nil
}

// Updates User
func (pdb *UserPostgresDb) UpdateUser(user User) error {
	conn, err := pdb.Pool.Acquire(context.Background())
	if err != nil {
		pdb.log.Error("Unable to acquire DB connection. \n", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"UPDATE users SET first_name = %2, last_name = %3, nickname = %4, email = %5, country = %6 WHERE id = %1",
		user.Id, user.FirstName, user.LastName, user.Nickname, user.Email, user.Country)

	if err != nil {
		pdb.log.Errorf("Failed to update user with ID:%d \n %v", user.Id, err)
		return err
	}
	pdb.log.Printf("User with ID:%d updated \n", user.Id)
	return nil
}

// List Users from DB
func (pdb *UserPostgresDb) GetUsers(filters map[string][]string) ([]*User, error) {
	conn, err := pdb.Pool.Acquire(context.Background())
	if err != nil {
		pdb.log.Error("Unable to acquire DB connection. \n", err)
		return nil, err
	}
	defer conn.Release()

	sql, args := getSQLSelect(filters)

	rows, err := conn.Query(context.Background(),
		sql, args...)
	if err != nil {
		pdb.log.Error("Failed to select users from db. \n", err)
		return nil, err
	}
	defer rows.Close()

	var result []*User
	for rows.Next() {
		rowUser := &User{}
		err = rows.Scan(&rowUser.Id, &rowUser.FirstName, &rowUser.LastName, &rowUser.Nickname, &rowUser.Email, &rowUser.Country)
		if err != nil {
			pdb.log.Error("Unable to scan user. \n", err)
			return nil, err
		}
		result = append(result, rowUser)
	}

	return result, nil
}

// Remove User from DB
func (pdb *UserPostgresDb) DeleteUser(id int64) (int64, error) {
	conn, err := pdb.Pool.Acquire(context.Background())
	if err != nil {
		pdb.log.Error("Unable to acquire DB connection. \n", err)
		return -1, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		"DELETE FROM users  WHERE id = %1",
		id)

	if err != nil {
		pdb.log.Errorf("Failed to delete user with ID:%d \n %v", id, err)
		return -1, err
	}
	pdb.log.Printf("User with ID:%d deleted \n", id)

	return ct.RowsAffected(), nil
}

// Function generates SQL and arguments for ListAll request
func getSQLSelect(filters map[string][]string) (string, []interface{}) {
	if len(filters) == 0 {
		return "SELECT * FROM users", []interface{}{}
	}
	resultSQL := "SELECT * FROM users WHERE "
	args := make([]interface{}, len(filters))
	counter := 1
	lenCounter := 1
	for k, v := range filters {
		var clause string
		if lenCounter == 1 {
			clause = k + " "
		} else {
			clause = " AND " + k + " "
		}
		if len(v) > 1 {
			subClause := "IN ("
			for i, sV := range v {
				if i == len(v)-1 {
					subClause += "%" + strconv.Itoa(counter) + ")"
					args = append(args, sV)
					counter += 1
					continue
				}
				subClause += "%" + strconv.Itoa(counter) + ", "
				args = append(args, sV)
				counter += 1
			}
			clause += subClause
		} else {
			clause += "= " + "%" + strconv.Itoa(counter)
			args = append(args, v[0])
			counter += 1
		}
		lenCounter += 1
		resultSQL += clause

	}

	return resultSQL + ";", args
}
