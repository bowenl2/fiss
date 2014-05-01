package main

import (
	"fmt"
	"io/ioutil"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}


// directory_list_go_html reads file data from disk. It returns an error on failure.
func directory_list_go_html() ([]byte, error) {
	return bindata_read(
		"/home/liam/go/src/github.com/bowenl2/fiss/templates/directory-list.go.html",
		"directory-list.go.html",
	)
}

// error_go_html reads file data from disk. It returns an error on failure.
func error_go_html() ([]byte, error) {
	return bindata_read(
		"/home/liam/go/src/github.com/bowenl2/fiss/templates/error.go.html",
		"error.go.html",
	)
}

// layout_go_html reads file data from disk. It returns an error on failure.
func layout_go_html() ([]byte, error) {
	return bindata_read(
		"/home/liam/go/src/github.com/bowenl2/fiss/templates/layout.go.html",
		"layout.go.html",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	if f, ok := _bindata[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string] func() ([]byte, error) {
	"directory-list.go.html": directory_list_go_html,
	"error.go.html": error_go_html,
	"layout.go.html": layout_go_html,

}
