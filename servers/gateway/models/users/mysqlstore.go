package users

import "database/sql"

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil database pointer")
	}
	return &MySQLStore{db}
}

func GetByID(id int64) (*User, error) {

}

func GetByEmail(email string) (*User, error) {

}

func GetByUserName(username string) (*User, error) {

}

func Insert(user *User) (*User, error) {

}

func Update(id int64, updates *Updates) (*User, error) {

}

func Delete(id int64) error {

}