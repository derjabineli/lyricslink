package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "password123!"
	password2 := "anotherPassword123!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	cases := []struct {
		name string
		password string
		hash string
		wantErr bool
	} {
		{
			name: "Correct Password 1",
			password: password1,
			hash: hash1,
			wantErr: false,
		},
		{
			name: "Correct Password 2",
			password: password2,
			hash: hash2,
			wantErr: false,
		},
		{
			name: "Incorrect Password",
			password: "incorrectPassword",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Incorrect Hash",
			password: password1,
			hash: hash2,
			wantErr: true,
		},
		{
			name: "Invalid Hash",
			password: password1,
			hash: "asd",
			wantErr: true,
		},
		{
			name: "Empty Password",
			password: "",
			hash: hash1,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := CheckPasswordHash(c.password, c.hash)
		if (err != nil) != c.wantErr {
			t.Errorf("%v expecting error %v but received %v", c.name, c.wantErr, err)
		}
		})
	}
}