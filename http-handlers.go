package main

import (
	"os"
	"html/template"
	"sort"
	"strconv"
	"io"
	"path/filepath"
	"fmt"
	"encoding/csv"
	"net/http"
)

func handleListDirRecursive(root string, fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w := csv.NewWriter(rw)
	w.Write([]string{"Path", "Modified", "Size", "Mode"})
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		w.Write([]string{
			filepath.Join(root, path),
			info.ModTime().Format("2006-01-02 15:04:05 -0700 MST"),
			strconv.Itoa(int(info.Size())),
			info.Mode().String(),
		})
		return nil // Never stop the function!
	})
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("\nERROR: %v", err))
	}
}

func internalErrorHandler(err error, rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(500)

	tmplAsset, tmplErr := Asset("error.html")
	if tmplErr != nil {
		io.WriteString(rw, "Internal server error.  Additionally, an error was encountered while loading the error page")
		return
	}
	tmplString := string(tmplAsset)

	tmpl := template.Must(
		template.New("error.html").Parse(tmplString))
	tmpl.Execute(rw, map[string]interface{}{
		"err": err,
		"req": r,
	})
}

func handleListDir(path string, fileInfo os.FileInfo, rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	dir, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sort.Sort(FileSort(entries))

	hostname, _ := os.Hostname()
	// ViewModel
	dl := DirectoryList{
		Machine:  hostname,
		Path:     path,
		BaseInfo: fileInfo,
		Entries:  entries,
	}

	err = render("directory-list.html", dl, rw)
	if err != nil {
		fmt.Printf("template error: %v\n", err)
	}
}

func handleFile(path string, fileInfo os.FileInfo, rw http.ResponseWriter, r *http.Request) {
	content, err := os.Open(path)
	if err != nil {

		io.WriteString(rw, err.Error())
	}

	http.ServeContent(rw, r, path, fileInfo.ModTime(), content)
}
