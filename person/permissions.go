package person

type Permission string

const (
	PermissionLookup      Permission = "lookup"
	PermissionBan         Permission = "ban"
	PermissionInformation Permission = "information"
	PermissionDonator     Permission = "donator"
	PermissionGiveItem    Permission = "give_item"
	PermissionTakeItem    Permission = "take_item"
	PermissionReset       Permission = "reset"

	PermissionAll Permission = "all"
)