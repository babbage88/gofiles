package main

import "embed"

//go:embed all:templates
var viewtmpl embed.FS

//go:embed static/*
var staticfs embed.FS

//go:embed all:test
var testfile embed.FS

func main() {
	server := NewFromEnv(".env")
	server.Start()
}
