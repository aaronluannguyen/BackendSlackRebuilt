package users

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
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
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlStore := NewMySQLStore(db)
	if sqlStore == nil {
		t.Errorf("an error with returning database from NewMySQLStore")
	}
}

func TestMySQLStore_GetByID(t *testing.T) {

}

func TestMySQLStore_GetByEmail(t *testing.T) {

}

func TestMySQLStore_GetByUserName(t *testing.T) {

}

func TestMySQLStore_Insert(t *testing.T) {

}

func TestMySQLStore_Update(t *testing.T) {

}

func TestMySQLStore_Delete(t *testing.T) {

}