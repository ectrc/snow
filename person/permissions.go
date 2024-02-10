package person

type Permission int64

// DO NOT MOVE THE ORDER OF THESE PERMISSIONS AS THEY ARE USED IN THE DATABASE
const (
	// random utility permissions
	PermissionLookup Permission = 1 << iota
	PermissionInformation

	// control permissions
	PermissionBansControl
	PermissionItemControl
	PermissionLockerControl
	PermissionPermissionControl

	// user roles, not really permissions but implemented as such
	PermissionOwner
	PermissionDonator

	// special permissions
	PermissionAll          = PermissionLookup | PermissionBansControl | PermissionInformation | PermissionItemControl | PermissionLockerControl | PermissionPermissionControl
	PermissionAllWithRoles = PermissionAll | PermissionOwner | PermissionDonator
)

func (p Permission) GetName() string {
	if p == 0 {
		return "None"
	}

	if p == PermissionAll {
		return "All"
	}

	if p == PermissionAllWithRoles {
		return "AllWithRoles"
	}

	if p&PermissionLookup != 0 {
		return "Lookup"
	}

	if p&PermissionBansControl != 0 {
		return "Ban"
	}

	if p&PermissionInformation != 0 {
		return "Information"
	}

	if p&PermissionItemControl != 0 {
		return "ItemControl"
	}

	if p&PermissionLockerControl != 0 {
		return "LockerControl"
	}

	if p&PermissionPermissionControl != 0 {
		return "PermissionControl"
	}

	if p&PermissionOwner != 0 {
		return "Owner"
	}

	if p&PermissionDonator != 0 {
		return "Donator"
	}

	return "Unknown"
}

func IntToPermission(i int64) Permission {
	return Permission(i)
}