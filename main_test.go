package main

import (
	_ "embed"
	"testing"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
)

func TestPersonPermissions(t *testing.T) {
	person := fortnite.NewFortnitePerson("helloWorld", false)
	if person.DisplayName != "helloWorld" {
		t.Errorf("Expected 'helloWorld', got %s", person.DisplayName)
	}

	aid.Print("intial permissions", person.Permissions)
	if person.Permissions != 0 {
		t.Errorf("Expected 0, got %d", person.Permissions)
	}

	person.AddPermission(p.PermissionAll)
	aid.Print("permissions after adding all", person.Permissions)

	if !person.HasPermission(p.PermissionAll) {
		t.Errorf("Expected to have all permissions, got false")
	}
	aid.Print("has all", person.HasPermission(p.PermissionAll))

	if !person.HasPermission(p.PermissionDonator) {
		t.Errorf("Expected donator permission as a result of all permissions, got false")
	}
	aid.Print("has donator", person.HasPermission(p.PermissionDonator))

	person.RemovePermission(p.PermissionAll)
	aid.Print("permissions after removing all", person.Permissions)

	if person.HasPermission(p.PermissionAll) {
		t.Errorf("Expected no permissions, got true")
	}
	aid.Print("does not have all", !person.HasPermission(p.PermissionAll))

	person.AddPermission(p.PermissionDonator)
	aid.Print("permissions after adding donator", person.Permissions)

	if !person.HasPermission(p.PermissionDonator) {
		t.Errorf("Expected to have donator permission, got false")
	}
	aid.Print("has donator", person.HasPermission(p.PermissionDonator))

	if person.HasPermission(p.PermissionAll) {
		t.Errorf("Expected not to have all permission, got true")
	}
	aid.Print("does not have all", !person.HasPermission(p.PermissionAll))

	if person.HasPermission(p.PermissionItemControl) {
		t.Errorf("Expected not to have give permission, got true")
	}

	person.RemovePermission(p.PermissionAll)
	aid.Print("permissions after removing all", person.Permissions)

	if person.HasPermission(p.PermissionAll) {
		t.Errorf("Expected not to have all permission, got true")
	}
	aid.Print("does not have all", !person.HasPermission(p.PermissionAll))

	if person.HasPermission(p.PermissionDonator) {
		t.Errorf("Expected not to have donator permission, got true")
	}
	aid.Print("does not have donator", !person.HasPermission(p.PermissionDonator))

	person.AddPermission(p.PermissionDonator | p.PermissionItemControl)
	aid.Print("permissions after adding donator and item control", person.Permissions)

	if !person.HasPermission(p.PermissionDonator) {
		t.Errorf("Expected to have donator permission, got false")
	}
	aid.Print("has donator", person.HasPermission(p.PermissionDonator))

	if !person.HasPermission(p.PermissionItemControl) {
		t.Errorf("Expected to have item control permission, got false")
	}
	aid.Print("has item control", person.HasPermission(p.PermissionItemControl))


	person.Delete()
}