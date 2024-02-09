package person

import "github.com/ectrc/snow/aid"

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

	permissionRealAll = 1<<iota - 1
	// every permission except owner
	PermissionAll Permission = permissionRealAll ^ PermissionOwner
)

func (p Permission) StringifyMe() string {
	aid.Print(p)
	switch p {
	case PermissionLookup:
		return "Lookup"
	case PermissionBan:
		return "Ban"
	case PermissionInformation:
		return "Information"
	case PermissionItemControl:
		return "ItemControl"
	case PermissionLockerControl:
		return "LockerControl"
	case PermissionPermissionControl:
		return "PermissionControl"
	case PermissionOwner:
		return "Owner"
	case PermissionDonator:
		return "Donator"
	case PermissionAll:
		return "All"
	default:
		return "Unknown"
	}
}