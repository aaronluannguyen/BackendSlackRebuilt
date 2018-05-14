package users

import "errors"

type MockStore struct {
	Error bool
	User *User
}

func NewMockStore(err bool, user *User) *MockStore {
	return &MockStore{
		err,
		user,
	}
}

func (m *MockStore) GetByID(id int64) (*User, error) {
	if m.Error {
		return nil, errors.New("triggered error")
	}
	return m.User, nil
}

func (m *MockStore) GetByEmail(email string) (*User, error) {
	if m.Error {
		return nil, errors.New("triggered error")
	}
	return m.User, nil
}

func (m *MockStore) GetByUserName(username string) (*User, error) {
	if m.Error {
		return nil, errors.New("triggered error")
	}
	return m.User, nil
}

func (m *MockStore) Insert(user *User) (*User, error) {
	if m.Error {
		return nil, errors.New("triggered error")
	}
	return m.User, nil
}

func (m *MockStore) Update(id int64, updates *Updates) (*User, error) {
	if m.Error {
		return nil, errors.New("triggered error")
	}
	return m.User, nil
}

func (m *MockStore) Delete(id int64) error {
	if m.Error {
		return errors.New("triggered error")
	}
	return nil
}