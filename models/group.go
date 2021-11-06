package models

import "github.com/google/uuid"

// GroupDNSAdmin is the uuid of the Dns Administrators group
var GroupDNSAdmin = uuid.Must(uuid.Parse("0b1e0a88-37a4-48f0-8060-2814906fa9f7"))

// GroupSuperAdmin is the uuid of the Super Administrators group
var GroupSuperAdmin = uuid.Must(uuid.Parse("71df8f2b-f293-4fde-93b1-e40dbe5c97ea"))

// GroupUserAdmin is the uuid of the User Administrators group
var GroupUserAdmin = uuid.Must(uuid.Parse("fbc827a0-32db-4d71-b95e-632b414e7993"))

// GroupTitle contains the titles of the groups.
var GroupTitle = map[uuid.UUID]string{
	GroupDNSAdmin:   "DNS Admin",
	GroupSuperAdmin: "Super Admin",
	GroupUserAdmin:  "User Admin",
}

// groups of groups

// GroupsAll contains the uuids of all groups
var GroupsAll = []uuid.UUID{
	GroupDNSAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

// GroupsAllAdmins contains the uuids of all admin groups
var GroupsAllAdmins = []uuid.UUID{
	GroupDNSAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

// GroupsDNSAdmin contains the uuids of all admin who have access to DNS admin functions
var GroupsDNSAdmin = []uuid.UUID{
	GroupDNSAdmin,
	GroupSuperAdmin,
}

// GroupsUserAdmin contains the uuids of all admin who have access to user admin functions
var GroupsUserAdmin = []uuid.UUID{
	GroupSuperAdmin,
	GroupUserAdmin,
}
