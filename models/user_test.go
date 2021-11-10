package models

import "testing"

func TestUserPasswordHash(t *testing.T) {
	tables := []struct {
		x string
		y string
		n bool
	}{
		{"", "", true},
		{"password", "password", true},
		{"password", "Password", false},
	}

	for _, table := range tables {
		u := User{}

		err := u.SetPassword(table.x)
		if err != nil {
			t.Errorf("got error setting password %s: %s", table.x, err.Error())

		}

		valid := u.CheckPasswordHash(table.y)
		if valid != table.n {
			t.Errorf("check password failed matching %s to %s, got: %v, want: %v,", table.x, table.y, valid, table.n)
		}
	}
}
