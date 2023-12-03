package person

type PersonSnapshot struct {
	ID                  string
	DisplayName         string
	AthenaProfile       ProfileSnapshot
	CommonCoreProfile   ProfileSnapshot
	CommonPublicProfile ProfileSnapshot
	Profile0Profile     ProfileSnapshot
	CollectionsProfile  ProfileSnapshot
}

type ProfileSnapshot struct {
	ID         string
	Items      map[string]ItemSnapshot
	Gifts      map[string]GiftSnapshot
	Quests     map[string]Quest
	Attributes map[string]Attribute
	Loadouts   map[string]Loadout
	Revision   int
	Type       string
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