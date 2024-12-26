package main

import "embed"

//go:embed all:templates
var viewtmpl embed.FS

//go:embed static/*
var staticfs embed.FS

func main() {
	server := NewFromEnv(".env")
	server.Start()
}
