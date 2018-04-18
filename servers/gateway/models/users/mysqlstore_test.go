package users

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"database/sql"
	"fmt"
	"database/sql/driver"
)

func TestNewMySQLStorePanic(t *testing.T) {
	defer func() {
		if r:= recover(); r == nil {
			t.Errorf("NewMySQLStore function did not panic")
		}
	}()
	NewMySQLStore(nil)
}

func TestNewMySQLStore(t *testing.T) {
	db, _ := getNewSQLMock(t)
	defer db.Close()
	sqlStore := NewMySQLStore(db)
	if sqlStore == nil {
		t.Errorf("an error with returning database from NewMySQLStore")
	}
}

func TestMySQLStore_GetByID(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	expectedID := 1
	expectedSQL := "select * from users where id=?"
	row := createUserRows()

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedID).
		WillReturnRows(row)

	if _, err := s.GetByID(int64(expectedID)); err != nil {
		t.Errorf("error was not expected while getting user by id: %s", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedID).
		WillReturnError(fmt.Errorf("scanning error"))

	if _, err := s.GetByID(int64(expectedID)); err == nil {
		t.Errorf("expecting scanning error, but didn't get one")
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	if _, err := s.GetByID(int64(expectedID)); err == nil {
		t.Errorf("expecting an ErrNoRows but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func TestMySQLStore_GetByEmail(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	expectedEmail := "test1.uw.edu"
	expectedSQL := "select * from users where email=?"
	row := createUserRows()

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedEmail).
		WillReturnRows(row)

	if _, err := s.GetByEmail(expectedEmail); err != nil {
		t.Errorf("error was not expected while getting user by email: %s", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedEmail).
		WillReturnError(fmt.Errorf("scanning error"))

	if _, err := s.GetByEmail(expectedEmail); err == nil {
		t.Errorf("expecting scanning error, but didn't get one")
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedEmail).
		WillReturnError(sql.ErrNoRows)

	if _, err := s.GetByEmail(expectedEmail); err == nil {
		t.Errorf("expecting an ErrNoRows but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func TestMySQLStore_GetByUserName(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	expectedUsername := "user1"
	expectedSQL := "select * from users where username=?"
	row := createUserRows()

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedUsername).
		WillReturnRows(row)

	if _, err := s.GetByUserName(expectedUsername); err != nil {
		t.Errorf("error was not expected while getting user by username: %s", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedUsername).
		WillReturnError(sql.ErrNoRows)

	if _, err := s.GetByUserName(expectedUsername); err == nil {
		t.Errorf("expecting an ErrNoRows but didn't get one")
	}

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(expectedUsername).
		WillReturnError(fmt.Errorf("scanning error"))

	if _, err := s.GetByUserName(expectedUsername); err == nil {
		t.Errorf("expecting scanning error, but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func TestMySQLStore_Insert(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	u := &User{
		3, "test3@uw.edu", []byte("somehash"), "user3", "first3", "last3", "photo3",
	}
	sqlInsert := "insert into users (email, passHash, username, firstName, lastName, photoURL) values (?,?,?,?,?,?)"
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := s.Insert(u); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnError(fmt.Errorf("executing insert error"))

	if _, err := s.Insert(u); err == nil {
		t.Errorf("expecting executing insert error, but didn't get one")
	}

	_, ErrNoRowsAffected := driver.ResultNoRows.RowsAffected()
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnResult(sqlmock.NewErrorResult(ErrNoRowsAffected))

	if _, err := s.Insert(u); err == nil {
		t.Errorf("expecting getting new ID error, but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func TestMySQLStore_Update(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	row := createUserRows()
	expectedID := 1
	update := &Updates{"newFirst1", "newFirst2"}
	sqlGetID := "select * from users where id=?"
	sqlUpdate := "update users set firstName=?, lastName=? where id=?"
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(sqlGetID)).
		WithArgs(expectedID).
		WillReturnRows(row)

	if _, err := s.Update(int64(expectedID), update); err != nil {
		t.Errorf("expected no error with update, but got %s", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedID).
		WillReturnError(fmt.Errorf("executing update error"))

	if _, err := s.Update(int64(expectedID), update); err == nil {
		t.Errorf("expecting executing update error, but didn't get one")
	}

	_, ErrNoRowsAffected := driver.ResultNoRows.RowsAffected()
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedID).
		WillReturnResult(sqlmock.NewErrorResult(ErrNoRowsAffected))

	if _, err := s.Update(int64(expectedID), update); err == nil {
		t.Errorf("expecting getting rows affected error, but didn't get one")
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if _, err := s.Update(int64(expectedID), update); err == nil {
		t.Errorf("expecting ErrUserNotFound, but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func TestMySQLStore_Delete(t *testing.T) {
	db, mock := getNewSQLMock(t)
	defer db.Close()

	s := MySQLStore{db}
	expectedID := 1
	sqlDelete := "delete from users where id=?"

	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := s.Delete(int64(expectedID)); err != nil {
		t.Errorf("expected no error with delete, but got %s", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedID).
		WillReturnError(fmt.Errorf("executing delete error"))

	if err := s.Delete(int64(expectedID)); err == nil {
		t.Errorf("expecting executing update error, but didn't get one")
	}

	_, ErrNoRowsAffected := driver.ResultNoRows.RowsAffected()
	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewErrorResult(ErrNoRowsAffected))

	if err := s.Delete(int64(expectedID)); err == nil {
		t.Errorf("expecting getting rows affected error, but didn't get one")
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.Delete(int64(expectedID)); err == nil {
		t.Errorf("expecting ErrUserNotFound, but didn't get one")
	}
	expectationsMetTest(t, mock)
}

func getNewSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error %s was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func expectationsMetTest(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func createUserRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"ID",
		"Email",
		"PassHash",
		"UserName",
		"FirstName",
		"LastName",
		"PhotoURL",
	}).AddRow(1, "test1@uw.edu", "somehash", "user1", "first1", "last1", "photo1").
		AddRow(2, "test2@uw.edu", "somehash", "user2", "first2", "last2", "photo2")

	return rows
}