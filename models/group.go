package models

import "github.com/google/uuid"

var GroupSuperadmin = uuid.Must(uuid.Parse("71df8f2b-f293-4fde-93b1-e40dbe5c97ea"))

var GroupTitle = map[uuid.UUID]string{
	GroupSuperadmin: "SuperAdmin",
}

// groups of groups
var GroupsAllAdmins = []uuid.UUID{
	GroupSuperadmin,
}

var GroupsUserAdmin = []uuid.UUID{
	GroupSuperadmin,
}
