package users

import "testing"

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		expectError bool
		nu NewUser
	}{
		{
			"Valid Email",
			false,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Invalid Email",
			true,
			NewUser{
				"invalid email",
				"123456",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Valid Password",
			false,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Invalid Password",
			true,
			NewUser{
				"test@uw.edu",
				"12345",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Valid Password Confirmation",
			false,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Invalid Password Confirmation",
			true,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456789",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Valid Username",
			false,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456",
				"Username",
				"FirstName",
				"LastName",
			},
		},
		{
			"Invalid Username",
			true,
			NewUser{
				"test@uw.edu",
				"123456",
				"123456",
				"",
				"FirstName",
				"LastName",
			},
		},
	}

	for _, c := range cases {
		err := c.nu.Validate()
		if err == nil && c.expectError {
			t.Errorf("case: %s: expected error but didn't get one", c.name)
		}
		if err != nil && !c.expectError {
			t.Errorf("case: %s: unexpected error", c.name)
		}
	}
}

func TestToUser(t *testing.T) {

}

func TestFullName(t *testing.T) {
	cases := []struct {
		name string
		u User
		expectedFullName string
	}{
		{
			"Normal First & Last Name",
			User{
				FirstName: "First",
				LastName: "Last",
			},
			"First Last",
		},
		{
			"No First Name",
			User{
				LastName: "Last",
			},
			"Last",
		},
		{
			"No Last Name",
			User{
				FirstName: "First",
			},
			"First",
		},
		{
			"No First or Last Name",
			User{},
			"",
		},
	}

	for _, c := range cases {
		fullName := c.u.FullName()
		if fullName != c.expectedFullName {
			t.Errorf("case: %s: expected full name does not match with produced full name", c.name)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	user := User{}
	user.SetPassword("123456Abc")
	userEmptyString := User{}
	userEmptyString.SetPassword("")

	cases := []struct {
		name string
		user User
		userPassword []byte
		testPassword string
		err bool
	}{
		{
			"Correct Password",
			user,
			user.PassHash,
			"123456Abc",
			false,
		},
		{
			"Incorrect Password",
			user,
			user.PassHash,
			"123456Abc123",
			true,
		},
		{
			"Empty String Password Correct",
			userEmptyString,
			userEmptyString.PassHash,
			"",
			false,
		},
		{
			"Empty String Password Incorrect",
			userEmptyString,
			userEmptyString.PassHash,
			"not empty",
			true,
		},
	}

	for _, c := range cases {
		err := c.user.Authenticate(c.testPassword)
		if err != nil && !c.err {
			t.Errorf("case: %s: unexpected error", c.name)
		}
		if err == nil && c.err {
			t.Errorf("case %s: expected error but didn't get one", c.name)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	user := User{FirstName: "First", LastName: "Last"}
	update := &Updates{"NewFirst", "NewLast"}
	err := user.ApplyUpdates(update)
	if err != nil {
		t.Errorf("unexpected error when updating first and last name")
	}
	if user.FirstName != update.FirstName {
		t.Errorf("error in updating first name")
	}
	if user.LastName != update.LastName {
		t.Errorf("error in updating last name")
	}

	invalidUpdate := &Updates{"", ""}
	err = user.ApplyUpdates(invalidUpdate)
	if err == nil {
		t.Errorf("expected error but didn't get one for an invalid update with empty first and last name strings")
	}
}