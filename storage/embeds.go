package storage

import (
	"embed"
	"io"
	"net/http"
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

func HttpAsset(file string) (*[]byte) {
	client := http.Client{}
	
	resp, err := client.Get("https://raw.githubusercontent.com/ectrc/ectrc/main/" + file)
	if err != nil {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return &data
}