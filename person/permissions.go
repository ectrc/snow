package person

type Permission int64

// DO NOT MOVE THE ORDER OF THESE PERMISSIONS AS THEY ARE USED IN THE DATABASE
const (
	PermissionLookup Permission = 1 << iota
	PermissionBan
	PermissionInformation
	PermissionItemControl
	PermissionLockerControl
	PermissionPermissionControl
	// user roles, not really permissions but implemented as such
	PermissionOwner
	PermissionDonator

	PermissionAll Permission = 1<<iota - 1
)