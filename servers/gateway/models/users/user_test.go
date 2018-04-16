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

}

func TestAuthenticate(t *testing.T) {

}

func TestApplyUpdates(t *testing.T) {

}