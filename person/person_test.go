package person

import (
	"testing"
)

func TestPerson(t *testing.T) {
	person := NewPersonWithCustomID("test")
	if person.ID != "test" {
		t.Error("person should have id of test. has " + person.ID)
	}

	profilesToTest := []string{ "common_core", "athena", "common_public", "profile0", "collections", "creative" }
	for _, profile := range profilesToTest {
		if person.GetProfileFromType(profile) == nil {
			t.Error("person should have profile")
		}

		if person.GetProfileFromType(profile).Type != profile {
			t.Error("person should have profile with id of " + profile)
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