package person

type PersonSnapshot struct {
	ID                string
	DisplayName       string
	AthenaProfile     ProfileSnapshot
	CommonCoreProfile ProfileSnapshot
	Loadout           Loadout
}

type ProfileSnapshot struct {
	ID         string
	Items      map[string]ItemSnapshot
	Gifts      map[string]GiftSnapshot
	Quests     map[string]Quest
	Attributes map[string]string
}

type ItemSnapshot struct {
	ID          string
	TemplateID  string
	Quantity    int
	Favorite    bool
	HasSeen     bool
	Variants    []VariantChannel
	ProfileType string
}

type GiftSnapshot struct {
	ID         string
	TemplateID string
	Quantity   int
	FromID     string
	GiftedAt   int64
	Message    string
	Loot       []Item
}

// Snapshot returns a snapshot of the person. No pointers as it has to compare the value not the address.
func CreateSnapshot(person *Person) *PersonSnapshot {
	return &PersonSnapshot{
		ID:                person.ID,
		DisplayName:       person.DisplayName,
		Loadout:           *person.Loadout,
		AthenaProfile:     *person.AthenaProfile.Snapshot(),
		CommonCoreProfile: *person.CommonCoreProfile.Snapshot(),
	}
}