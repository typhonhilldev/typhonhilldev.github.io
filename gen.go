//go:build ignore

// Gen writes the pages that make typhonhill.dev/<project> import paths work: one page per importable package,
// with the go-import and go-source tags the go command and pkg.go.dev read, and a landing page.
//
// Run it after changing projects.json:
//
//	go run gen.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

const (
	domain = "typhonhill.dev"
	owner  = "https://github.com/typhonhilldev"
)

type project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Packages    []string `json:"packages"` // importable package directories, "" for the module root
}

var pkgPage = template.Must(template.New("pkg").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="go-import" content="{{ .Module }} git {{ .Repo }}">
<meta name="go-source" content="{{ .Module }} {{ .Repo }} {{ .Repo }}/tree/main{/dir} {{ .Repo }}/blob/main{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url={{ .Target }}">
<title>{{ .Import }}</title>
</head>
<body>
<p><a href="{{ .Target }}">{{ .Import }}</a></p>
</body>
</html>
`))

var indexPage = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>` + domain + `</title>
<style>
body { font-family: system-ui, sans-serif; max-width: 40rem; margin: 3rem auto; padding: 0 1rem; line-height: 1.5; }
code { font-size: 0.95em; }
</style>
</head>
<body>
<h1>` + domain + `</h1>
<p>Open-source Go projects.</p>
<ul>
{{- range . }}
<li><a href="` + owner + `/{{ .Name }}"><code>` + domain + `/{{ .Name }}</code></a>: {{ .Description }}</li>
{{- end }}
</ul>
</body>
</html>
`))

func main() {
	data, err := os.ReadFile("projects.json")
	if err != nil {
		log.Fatal(err)
	}
	var projects []project
	if err := json.Unmarshal(data, &projects); err != nil {
		log.Fatal(err)
	}
	for _, p := range projects {
		repo := owner + "/" + p.Name
		for _, pkg := range p.Packages {
			importPath := domain + "/" + p.Name
			target := repo
			if pkg != "" {
				importPath += "/" + pkg
				target = "https://pkg.go.dev/" + importPath
			}
			write(filepath.Join(p.Name, filepath.FromSlash(pkg), "index.html"), pkgPage, map[string]string{
				"Module": domain + "/" + p.Name, "Repo": repo, "Import": importPath, "Target": target,
			})
		}
	}
	write("index.html", indexPage, projects)
	for name, content := range map[string]string{"CNAME": domain + "\n", ".nojekyll": ""} {
		if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
			log.Fatal(err)
		}
	}
}

func write(name string, tmpl *template.Template, data any) {
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(f, data); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", name)
}
