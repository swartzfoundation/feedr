package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist/*
var f embed.FS
var dist fs.FS

func init() {
	sub, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}
	dist = sub
}

func HandlerFn() http.HandlerFunc {
	f := http.FS(dist)
	return http.FileServer(f).ServeHTTP
}
