package web

import "embed"

//go:embed ui/dist/*
var webUi embed.FS
