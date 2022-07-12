package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var defaultHandleTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Choose your own adventure</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
      {{range .Options}}
      <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
      {{end}}
    </ul>
  </body>
</html>`

// Use design pattern: Functional Options
type HandlerOption func(h *Handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *Handler) {
		h.template = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *Handler) {
		h.pathFunc = fn
	}
}

// func WithDatabase(username, password string) HandlerOption {}
// WithOther... -> positive of use Functional option pattern

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	defaultTpml := template.Must(template.New("").Parse(defaultHandleTemplate))
	h := Handler{
		s,
		defaultTpml,
		defaultPathFunc,
	}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type Handler struct {
	story    Story
	template *template.Template
	pathFunc func(r *http.Request) string
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := handler.pathFunc(r)
	if chapter, ok := handler.story[path]; ok {
		err := handler.template.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found...", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)
	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
