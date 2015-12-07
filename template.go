package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"

	"github.com/kr/pretty"
)

func render(viewName string, viewModel interface{}, w io.Writer) error {
	layoutAsset, err := Asset("layout.go.html")
	if err != nil {
		return err
	}
	layoutString := string(layoutAsset)

	contentAsset, err := Asset(viewName)
	if err != nil {
		return err
	}
	contentString := string(contentAsset)

	tmplFuncs := map[string]interface{}{
		"sizefmt": func(s int64) string {
			return ByteSize(s).String()
		},
		"pathfmt": func(f os.FileInfo) template.URL {
			// Path relative to the base of root of the server
			path := path.Join(
				viewModel.(DirectoryList).Path,
				f.Name())
			if f.IsDir() {
				path = fmt.Sprintf("%s/", path)
			}
			return template.URL(path)
		},
		"archivepath": func() template.URL {
			return template.URL(
				path.Join("/", "thumbs",
					viewModel.(DirectoryList).Path))

		},
		"prettyfmt": pretty.Formatter,
	}

	tmpl := template.New("template").Funcs(tmplFuncs)
	tmpl = template.Must(tmpl.Parse(layoutString))
	tmpl = template.Must(tmpl.Parse(contentString))

	return tmpl.Execute(w, viewModel)
}
