package data

import "github.com/sirupsen/logrus"

type UserPostgresDb struct {
	log *logrus.Logger
}

func NewPgDb(logger *logrus.Logger) *UserPostgresDb {
	u := &UserPostgresDb{
		logger,
	}

	return u
}
func (mdb *UserPostgresDb) AddUser(user User) error {
	return nil
}
func (mdb *UserPostgresDb) UpdateUser(user User) error {
	return nil
}

func (mdb *UserPostgresDb) GetUsers() (users []*User) {
	return nil
}

func (mdb *UserPostgresDb) DeleteUser(id int) error {
	return nil
}
