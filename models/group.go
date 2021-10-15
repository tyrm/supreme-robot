package models

const GroupSuperadmin = "71df8f2b-f293-4fde-93b1-e40dbe5c97ea"

var GroupTitle = map[string]string{
	GroupSuperadmin:        "SuperAdmin",
}

// groups of groups
var GroupsAllAdmins = []string {
	GroupSuperadmin,
}

var GroupsUserAdmin = []string {
	GroupSuperadmin,
}