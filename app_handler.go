package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// AppHandlerFunc defines a function which acts as a context-aware HTTP handler
// In case of error, it returns the error which is handled separately
type AppHandlerFunc func(http.ResponseWriter, *http.Request, Context) error

// AppHandler is the application's http.Handler
type AppHandler struct {
	// Request paths are considered relative to RootPath
	RootPath string
}

func (h *AppHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := Context{
		App:       h,
		Recursive: r.URL.Query().Get("r") != "",
		Format:    parseFmt(r.URL.Query().Get("fmt")),
	}

	err := fissBaseHandlerFunc(rw, r, c)
	if err != nil {
		internalErrorHandlerFunc(rw, r, c, err)
	}
}

func fissBaseHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {
	// Fill in filesystem details
	p := filepath.Join(
		c.App.RootPath,
		path.Clean(r.URL.Path))
	p, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(p)
	if err != nil {
		return err
	}

	c.FSPath = p
	c.FSInfo = fileInfo

	// Choose next handler to invoke
	return routeFSHandlerFunc(rw, r, c)
}

func parseFmt(f string) ResponseFormat {
	switch f {
	case "json":
		return FmtJSON
	case "csv":
		return FmtCSV
	case "html":
		return FmtHTML
	case "dl":
		return FmtForceDownload
	}
	return FmtAuto
}

func routeFSHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {
	fmt.Printf("req: %v %v\n", r.RemoteAddr, c.FSPath)

	if c.FSInfo.IsDir() {
		return directoryHandlerFunc(rw, r, c)
	}

	return fileHandlerFunc(rw, r, c)
}

func directoryHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {
	switch c.Format {
	case FmtForceDownload:
		return archiveHandlerFunc(rw, r, c)
	case FmtCSV:
		return recursiveDirectoryHandlerFunc(rw, r, c)
	}
	// FIXME: Implement JSON
	// FIXME: Separate CSV from Recursive
	// -- there should be non-recursive CSV and archives
	// -- there should be recursive HTML and JSON
	return directoryListHandlerFunc(rw, r, c)
}
