package server

import (
	"io/fs"
	"net/http"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/web"
)

func (h *Web) Static() http.Handler {
	fsys, _ := fs.Sub(web.StaticFS, "static/app")
	fileServer := http.FileServer(http.FS(fsys))
	if h.Dev {
		fileServer = http.FileServer(http.Dir("web/static/app"))
	}
	return fileServer
}
