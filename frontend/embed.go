package frontend

import (
	"embed"
	"io/fs"
)

//go:embed all:dist/*
var build embed.FS

var Content fs.FS

func init() {
	sub, err := fs.Sub(build, "dist")
	if err != nil {
		panic(err)
	}
	Content = sub
}
