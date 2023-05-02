package web

import "embed"

// `go-cook/web/dist/` as embedded filesystem
//
//go:embed dist/.* dist/*
var WebDist embed.FS
