// ------------------------------------------------------------------------------------------------
// Here are the enums used for authentication
// ------------------------------------------------------------------------------------------------
package goald

// ------------------------------------------------------------------------------------------------
// the type of a permission
// ------------------------------------------------------------------------------------------------

type PermissionType int

const (
	PermissionTypeUNKNOWN PermissionType = 0
	PermissionTypeREAD    PermissionType = 1
	PermissionTypeWRITE   PermissionType = 2
	PermissionTypeDELETE  PermissionType = 3
)

var permissionTypes = map[int]string{
	int(PermissionTypeUNKNOWN): "unknown",
	int(PermissionTypeREAD):    "read",
	int(PermissionTypeWRITE):   "write",
	int(PermissionTypeDELETE):  "delete",
}

func (thisProperty PermissionType) String() string {
	return permissionTypes[int(thisProperty)]
}

// Val helps implement the IEnum interface
func (thisProperty PermissionType) Val() int {
	return int(thisProperty)
}

// Values helps implement the IEnum interface
func (thisProperty PermissionType) Values() map[int]string {
	return permissionTypes
}

// ------------------------------------------------------------------------------------------------
// the status of a user or a user profile
// ------------------------------------------------------------------------------------------------

type UserStatus int

const (
	UserStatusINACTIVE    UserStatus = 0
	UserStatusACTIVE      UserStatus = 2
	UserStatusTOxCONFIRM  UserStatus = 1
	UserStatusTOxVALIDATE UserStatus = -1
	UserStatusSUSPENDED   UserStatus = -2
)

var userStatuses = map[int]string{
	int(UserStatusINACTIVE):    "inactive",
	int(UserStatusACTIVE):      "active",
	int(UserStatusTOxCONFIRM):  "to confirm",
	int(UserStatusTOxVALIDATE): "to validate",
	int(UserStatusSUSPENDED):   "suspended",
}

func (thisProperty UserStatus) String() string {
	return userStatuses[int(thisProperty)]
}

// Val helps implement the IEnum interface
func (thisProperty UserStatus) Val() int {
	return int(thisProperty)
}

// Values helps implement the IEnum interface
func (thisProperty UserStatus) Values() map[int]string {
	return userStatuses
}
