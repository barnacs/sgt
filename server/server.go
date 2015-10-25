package server

import (
	"mime"
	"strings"
	"html/template"
	"net/http"
	"github.com/barnacs/sgt/git"
)

var versionsTemplate = template.Must(template.New("versions").Parse(`
{{range .}}
<a href="{{.Id}}/index.html">{{.Message}}</a> -- {{.Author}}, {{.Time}} -- <a href="diff/{{.Id}}">diff</a><br/>
{{end}}
`))

type Server struct {
	repo *git.Repo
}

func New(repo *git.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.repo.Fetch()
		s.writeIndex(w)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/diff/") {
		s.writeDiff(w, r)
		return
	}
	s.writeFileVersion(w, r)
}

func (s *Server) writeIndex(w http.ResponseWriter) {
	versionsTemplate.Execute(w, s.repo.Versions())
}

func (s *Server) writeDiff(w http.ResponseWriter, r *http.Request) {
	version := strings.TrimPrefix(r.URL.Path, "/diff/")
	w.Write(s.repo.DiffToPrevious(version))
}

func (s *Server) writeFileVersion(w http.ResponseWriter, r *http.Request) {
	version, path := splitVersionPath(r.URL.Path)
	w.Header().Add("Content-Type", mimeType(path))
	w.Write(s.repo.FileVersion(path, version))
}

func splitVersionPath(urlPath string) (string, string) {
	path := strings.TrimLeft(urlPath, "/")
	firstSlash := strings.Index(path, "/")
	return path[:firstSlash], path[firstSlash+1:]
}

func mimeType(path string) string {
	lastDot := strings.LastIndex(path, ".")
	if lastDot < 0 {
		return ""
	}
	ext := path[lastDot:]
	return mime.TypeByExtension(ext)
}