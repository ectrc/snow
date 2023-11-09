package storage

import (
	"embed"
	"strings"
)

var (
	//go:embed mem/*
	Assets embed.FS
)

func Asset(file string) (*[]byte) {
	data, err := Assets.ReadFile("mem/" + strings.ToLower(file))
	if err != nil {
		return nil
	}

	return &data
}