package users

import (
	"database/sql"
	"fmt"
)

const sqlSelectAll = "select * from users"
const sqlSelectBy =  sqlSelectAll + " where "
const sqlSelectByID = sqlSelectBy + "id=?"
const sqlSelectByEmail = sqlSelectBy + "email=?"
const sqlSelectByUsername = sqlSelectBy + "username=?"
const sqlInsert = "insert into users (email, passHash, username, firstName, lastName, photoURL) values (?,?,?,?,?,?)"
const sqlUpdate = "update users set firstName=?, lastName=? where id=?"
const sqlDelete = "delete from users where id=?"

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil database pointer")
	}
	return &MySQLStore{db}
}

func (s *MySQLStore) GetByID(id int64) (*User, error) {
	row := s.db.QueryRow(sqlSelectByID, id)
	return getByX(row)
}

func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	row := s.db.QueryRow(sqlSelectByEmail, email)
	return getByX(row)
}

func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	row := s.db.QueryRow(sqlSelectByUsername, username)
	return getByX(row)
}

func (s *MySQLStore) Insert(user *User) (*User, error) {
	result, err := s.db.Exec(sqlInsert, user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("executing insert: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting new ID: %v", err)
	}
	user.ID = id
	return user, nil
}

func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	result, err := s.db.Exec(sqlUpdate, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("updating: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting rows affected: %v", err)
	}
	if affected == 0 {
		return nil, ErrUserNotFound
	}
	return s.GetByID(id)
}

func (s *MySQLStore) Delete(id int64) error {
	result, err := s.db.Exec(sqlDelete, id)
	if err != nil {
		return fmt.Errorf("deleting: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %v", err)
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func getByX(r *sql.Row) (*User, error) {
	user := &User{}
	if err := r.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}