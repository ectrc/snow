package storage

import (
	"embed"
	"strings"
)

var (
	//go:embed assets/*
	Assets embed.FS
)

func Asset(file string) (*[]byte) {
	data, err := Assets.ReadFile("assets/" + strings.ToLower(file))
	if err != nil {
		return nil
	}

	return &data
}