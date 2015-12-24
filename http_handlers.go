package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func recursiveDirectoryHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {

	rw.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w := csv.NewWriter(rw)
	w.Write([]string{"Path", "Modified", "Size", "Mode"})
	return filepath.Walk(c.FSPath,
		func(path string, info os.FileInfo, err error) error {
			w.Write([]string{
				filepath.Join(c.FSPath, path),
				info.ModTime().Format("2006-01-02 15:04:05 -0700 MST"),
				strconv.Itoa(int(info.Size())),
				info.Mode().String(),
			})
			return nil // Never stop the function!
		})
}

// Handle errors encountered while processing requests
// Not an AppHandlerFunc
func internalErrorHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context, err error) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(500)

	err = render("error.go.html", map[string]interface{}{
		"err": err,
		"req": r,
	}, rw)
	if err != nil {
		io.WriteString(rw, "Internal server error.  Additionally, an error was encountered while loading the error page: "+err.Error())
	}
}

func handleListRoots(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	hostname, _ := os.Hostname()
	rootInfos, errors := listRootInfos()

	viewModel := RootListViewModel{
		Machine:   hostname,
		RootInfos: rootInfos,
		Errors:    errors,
	}

	err := render("root-list.go.html", viewModel, rw)
	if err != nil {
		fmt.Printf("template rendering error: %v\n", err)
	}
}

func directoryListHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	dir, err := os.Open(c.FSPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	// FIXME: Use gen to simplify this stuff
	sort.Sort(FileSort(entries))

	// The view should see the path as relative to the root
	// (it should not care where the root is)
	relPath, _ := filepath.Rel(c.App.RootPath, c.FSPath)
	relPath = filepath.Clean(
		filepath.Join(string(filepath.Separator), relPath))
	hostname, _ := os.Hostname()

	// ViewModel
	dl := DirectoryList{
		Machine:     hostname,
		Path:        relPath,
		BaseInfo:    c.FSInfo,
		Entries:     entries,
		BreadCrumbs: makeBreadCrumbs(relPath),
	}

	return render("directory-list.go.html", dl, rw)
}

func fileHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {
	content, err := os.Open(c.FSPath)
	if err != nil {
		return err
	}
	if c.Format == FmtForceDownload {
		rw.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(rw, r, c.FSPath, c.FSInfo.ModTime(), content)
	return nil
}

func archiveHandlerFunc(
	rw http.ResponseWriter, r *http.Request, c Context) error {

	p, err := MakeArchive(c.FSPath)
	if err != nil {
		return err
	}
	defer os.Remove(p) // once served, don't hang around.

	archiveFile, err := os.Open(p)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	http.ServeContent(rw, r,
		fmt.Sprintf("%s.zip", filepath.Base(p)),
		time.Now(),
		archiveFile)
	return nil
}
