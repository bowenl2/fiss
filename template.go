package main

import (
	"github.com/kr/pretty"
	"html/template"
	"os"
	"io"
	"path/filepath"
)

func render(viewName string, viewModel interface{}, w io.Writer) (error) {
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
		"fmtsize": func(s int64) string {
			return ByteSize(s).String()
		},
		"relpath": func(f os.FileInfo) string {
			// Path relative to the base of root of the server
			return filepath.Join(viewModel.(DirectoryList).Path, f.Name())
		},
		"prettyfmt": pretty.Formatter,
	}

	tmpl := template.New("template").Funcs(tmplFuncs)
	tmpl = template.Must(tmpl.Parse(layoutString))
	tmpl = template.Must(tmpl.Parse(contentString))

	return tmpl.Execute(w, viewModel)
}
