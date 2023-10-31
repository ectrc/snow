package person

import "github.com/ectrc/snow/aid"

type ProfileChange struct {
	ChangeType string `json:"changeType"`
}

type FullProfileUpdate struct {
	ChangeType string `json:"changeType"`
	Profile aid.JSON `json:"profile"`
}

type StatModified struct {
	ChangeType string `json:"changeType"`
	Name string `json:"name"`
	Value interface{} `json:"value"`
}

type ItemAdded struct {
	ChangeType string `json:"changeType"`
	ItemId string `json:"itemId"`
	Item aid.JSON `json:"item"`
}

type ItemRemoved struct {
	ChangeType string `json:"changeType"`
	ItemId string `json:"itemId"`
}

type ItemAttributeChanged struct {
	ChangeType string `json:"changeType"`
	ItemId string `json:"itemId"`
	AttributeName string `json:"attributeName"`
	AttributeValue interface{} `json:"attributeValue"`
}

type ItemQuantityChanged struct {
	ChangeType string `json:"changeType"`
	ItemId string `json:"itemId"`
	Quantity int `json:"quantity"`
}