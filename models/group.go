package models

import "github.com/google/uuid"

var GroupDnsAdmin = uuid.Must(uuid.Parse("0b1e0a88-37a4-48f0-8060-2814906fa9f7"))
var GroupSuperAdmin = uuid.Must(uuid.Parse("71df8f2b-f293-4fde-93b1-e40dbe5c97ea"))
var GroupUserAdmin = uuid.Must(uuid.Parse("fbc827a0-32db-4d71-b95e-632b414e7993"))

var GroupTitle = map[uuid.UUID]string{
	GroupDnsAdmin:   "DNS Admin",
	GroupSuperAdmin: "Super Admin",
	GroupUserAdmin:  "User Admin",
}

// groups of groups

var GroupsAll = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

var GroupsAllAdmins = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

var GroupsDnsAdmin = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
}

var GroupsUserAdmin = []uuid.UUID{
	GroupSuperAdmin,
	GroupUserAdmin,
}
