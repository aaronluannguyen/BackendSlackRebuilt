package users

import (
	"database/sql"
	"fmt"
	"github.com/challenges-aaronluannguyen/servers/gateway/indexes"
	"strings"
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

//getByX scans in the given row and returns a user struct with the scanned information
// from the row
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

//LoadExistingUsersToTrie loads and returns a trie with all current users from
//the existing database
func (s *MySQLStore) LoadExistingUsersToTrie() (*indexes.Trie, error) {
	trie := indexes.NewTrie()
	rows, err := s.db.Query(sqlSelectAll)
	if err != nil {
		return nil, fmt.Errorf("selecting: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrUserNotFound
			}
			return nil, fmt.Errorf("scanning: %v", err)
		}
		AddUserToTrie(trie, user)
	}
	return trie, nil
}

//AddUserToTrie adds the user's username, first name, and last name to the trie
//accounting for multiple-word names and accidental spaces before and after names
func AddUserToTrie(trie *indexes.Trie, user *User) {
	username := strings.Split(strings.ToLower(user.UserName), " ")
	firstName := strings.Split(strings.ToLower(user.FirstName), " ")
	lastName := strings.Split(strings.ToLower(user.LastName), " ")

	for _, name := range username {
		trimName := strings.TrimSpace(name)
		trie.Add(trimName, user.ID)
	}

	for _, name := range firstName {
		trimName := strings.TrimSpace(name)
		trie.Add(trimName, user.ID)
	}

	for _, name := range lastName {
		trimName := strings.TrimSpace(name)
		trie.Add(trimName, user.ID)
	}
}

//SortTopTwentyUsersByUsername orders the top twenty users by username and returns
//the correct order of users
func (s *MySQLStore) SortTopTwentyUsersByUsername(users []int64) (*[]*User, error) {
	var sortedUsers []*User
	if len(users) < 1 {
		return sortedUsers, nil
	}
	var queryQMarks string
	var ids string
	for _, id := range users {
		queryQMarks += ", ?"
		addIdString := fmt.Sprintf(", %d", id)
		ids += addIdString
	}
	queryQMarks = strings.TrimPrefix(queryQMarks, ", ")
	ids = strings.TrimPrefix(ids, ", ")
	query := fmt.Sprintf("select * from users where id in (%s) order by username", queryQMarks)
	idArg := GetIdInterface(users)
	rows, err := s.db.Query(query, idArg...)
	if err != nil {
		return nil, fmt.Errorf("error querying top twenty users from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrUserNotFound
			}
			return nil, fmt.Errorf("scanning: %v", err)
		}
		sortedUsers = append(sortedUsers, user)
	}
	return &sortedUsers, nil
}

//GetIdInterface retrieves the ids and inputs them into
//an interface for a query call
func GetIdInterface(users []int64) []interface{} {
	args := make([]interface{}, len(users))
	for i, id := range users {
		args[i] = id
	}
	return args
}