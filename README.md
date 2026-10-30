# typhonhill.dev

The site behind `typhonhill.dev/<project>` import paths, served by GitHub Pages. Each importable package has a page with the `go-import` and `go-source` tags that the `go` command and pkg.go.dev read, pointing at the project's repository in this organization.

To add a project, or a package someone may import or `go get -tool`, add it to `projects.json` and run:

```sh
go run gen.go
```

Commit the generated pages with the change. Because the import path names this domain rather than GitHub, the repositories can move without breaking anyone's imports; update `owner` in `gen.go` if they do.
