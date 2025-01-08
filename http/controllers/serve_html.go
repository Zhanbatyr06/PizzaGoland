package controllers

import (
	"net/http"
)

type StaticController struct {
	staticDir string
}

func NewStaticController(staticDir string) *StaticController {
	return &StaticController{staticDir: staticDir}
}

//func (s *StaticController) ServeHTML(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, s.staticDir+"/main.html")
//}

func (s *StaticController) ServeStaticFiles() http.Handler {
	return http.FileServer(http.Dir(s.staticDir))
}
