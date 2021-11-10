package models

import (
	"github.com/google/uuid"
	"testing"
)

func TestUser_IsMemberOfGroup(t *testing.T) {
	tables := []struct {
		x []uuid.UUID
		y []uuid.UUID
		n bool
	}{
		// Check GroupsAll
		{[]uuid.UUID{GroupDNSAdmin}, GroupsAll, true},
		{[]uuid.UUID{GroupUserAdmin}, GroupsAll, true},
		{[]uuid.UUID{GroupSuperAdmin}, GroupsAll, true},

		// Check GroupsAllAdmins
		{[]uuid.UUID{GroupDNSAdmin}, GroupsAllAdmins, true},
		{[]uuid.UUID{GroupUserAdmin}, GroupsAllAdmins, true},
		{[]uuid.UUID{GroupSuperAdmin}, GroupsAllAdmins, true},

		// Check GroupsDNSAdmin
		{[]uuid.UUID{GroupDNSAdmin}, GroupsDNSAdmin, true},
		{[]uuid.UUID{GroupUserAdmin}, GroupsDNSAdmin, false},
		{[]uuid.UUID{GroupSuperAdmin}, GroupsDNSAdmin, true},

		// Check GroupsUserAdmin
		{[]uuid.UUID{GroupDNSAdmin}, GroupsUserAdmin, false},
		{[]uuid.UUID{GroupUserAdmin}, GroupsUserAdmin, true},
		{[]uuid.UUID{GroupSuperAdmin}, GroupsUserAdmin, true},

		// Check multi groups
		{[]uuid.UUID{
			GroupUserAdmin,
			uuid.Must(uuid.Parse("6bb2ffc1-a060-4230-aae7-e7a16eb860b7")),
		}, GroupsUserAdmin, true},
		{[]uuid.UUID{
			uuid.Must(uuid.Parse("b7f67f64-aa55-4852-ba79-ebcb9507cfe7")),
			GroupDNSAdmin,
			uuid.Must(uuid.Parse("59640372-d188-4b28-8b07-d01e8c031fad")),
		}, GroupsAllAdmins, true},
		{[]uuid.UUID{
			uuid.Must(uuid.Parse("dce50784-b1cd-4966-83b6-a92dc9a8c4d8")),
			uuid.Must(uuid.Parse("9ba68d8f-586b-4c4e-90d0-013710949b8c")),
			uuid.Must(uuid.Parse("6a754597-9b87-4bfe-851a-d05b8b0fe2d6")),
			GroupSuperAdmin,
		}, GroupsDNSAdmin, true},
		{[]uuid.UUID{
			uuid.Must(uuid.Parse("e0e08f64-b07f-44b6-bfdc-b3b18a098e70")),
			uuid.Must(uuid.Parse("89ab8121-c140-4057-a237-5186dadf8978")),
			uuid.Must(uuid.Parse("2fb1b9d7-800c-47b4-8210-b717f1b44af1")),
		}, GroupsAll, false},
	}

	for _, table := range tables {
		u := User{
			Groups: table.x,
		}

		valid := u.IsMemberOfGroup(table.y...)
		if valid != table.n {
			t.Errorf("check password failed matching %s to %s, got: %v, want: %v,", table.x, table.y, valid, table.n)
		}
	}
}

func TestUserPasswordHash(t *testing.T) {
	tables := []struct {
		x string
		y string
		n bool
	}{
		{"", "", true},
		{"password", "password", true},
		{"i'm a super long password with $p3c!@L characters!!!!", "i'm a super long password with $p3c!@L characters!!!!", true},
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
