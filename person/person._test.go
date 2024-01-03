package person

import (
	"testing"
)

func TestPerson(t *testing.T) {
	person := NewPersonWithCustomID("test")

	person.AddPermission("test")
	if !person.HasPermission("test") {
		t.Error("person should have permission")
	}

	person.RemovePermission("test")
	if person.HasPermission("test") {
		t.Error("person should not have permission")
	}

	person.AddFriend("test")
	if len(person.Friends) != 1 {
		t.Error("person should have 1 friend")
	}

	person.RemoveFriend("test")
	if len(person.Friends) != 0 {
		t.Error("person should have no friends")
	}

	if person.ID != "test" {
		t.Error("person should have id of test")
	}

	profilesToTest := []string{ "common_core", "athena", "common_public", "profile0", "collections", "creative" }
	for _, profile := range profilesToTest {
		if person.GetProfileFromType(profile) == nil {
			t.Error("person should have profile")
		}

		if person.GetProfileFromType(profile).Type != profile {
			t.Error("person should have profile with id of " + profile)
		}

		if person.GetProfileFromType(profile).PersonID != person.ID {
			t.Error("person should have profile with person id of " + person.ID)
		}
	}

	item := NewItem("Test:Test", 1)

	person.AthenaProfile.Items.AddItem(item)
	if person.AthenaProfile.Items.GetItemByTemplateID("Test:Test") == nil {
		t.Error("person should have item")
	}

	if person.AthenaProfile.Items.GetItem(item.ID) == nil {
		t.Error("person should have item")
	}
}