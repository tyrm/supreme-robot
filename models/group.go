package models

import "github.com/google/uuid"

// GroupDnsAdmin is the uuid of the Dns Administrators group
var GroupDnsAdmin = uuid.Must(uuid.Parse("0b1e0a88-37a4-48f0-8060-2814906fa9f7"))

// GroupSuperAdmin is the uuid of the Super Administrators group
var GroupSuperAdmin = uuid.Must(uuid.Parse("71df8f2b-f293-4fde-93b1-e40dbe5c97ea"))

// GroupUserAdmin is the uuid of the User Administrators group
var GroupUserAdmin = uuid.Must(uuid.Parse("fbc827a0-32db-4d71-b95e-632b414e7993"))

// GroupTitle contains the titles of the groups.
var GroupTitle = map[uuid.UUID]string{
	GroupDnsAdmin:   "DNS Admin",
	GroupSuperAdmin: "Super Admin",
	GroupUserAdmin:  "User Admin",
}

// groups of groups

// GroupsAll contains the uuids of all groups
var GroupsAll = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

// GroupsAllAdmins contains the uuids of all admin groups
var GroupsAllAdmins = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
	GroupUserAdmin,
}

// GroupsDnsAdmin contains the uuids of all admin who have access to DNS admin functions
var GroupsDnsAdmin = []uuid.UUID{
	GroupDnsAdmin,
	GroupSuperAdmin,
}

// GroupsUserAdmin contains the uuids of all admin who have access to user admin functions
var GroupsUserAdmin = []uuid.UUID{
	GroupSuperAdmin,
	GroupUserAdmin,
}
