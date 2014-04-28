package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var (
	rootPath = "/"
)

// Directory List
type FileSort []os.FileInfo

func (l FileSort) Len() int {
	return len(l)
}

func (l FileSort) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l FileSort) Less(i, j int) bool {
	if l[i].IsDir() && !l[j].IsDir() {
		return true
	}
	if !l[i].IsDir() && l[j].IsDir() {
		return false
	}
	return l[i].Name() < l[j].Name()
}

type DirectoryList struct {
	Machine  string
	Path     string
	BaseInfo os.FileInfo
	Entries  []os.FileInfo
}

func recursiveDirectoryList(fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

}

func handleDir(fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Template
	tmpl := template.New("DirectoryList")
	tmp, _ = tmpl.ParseFiles("templates/directory-list.html")

func handleDir(path string, fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
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

	tmplFuncs := map[string]interface{}{
		"fmtsize": func(s int64) string {
			return ByteSize(s).String()
		},
	}

	tmpl := template.Must(
		template.New("directory-list.html").Funcs(tmplFuncs).ParseFiles("templates/directory-list.html"))

	err = tmpl.Execute(rw, dl)
	if err != nil {
		fmt.Println("Error executing the template: %v\n", err)
	}
}

func handleFile(fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Length", string(fileInfo.Size()))
	handle, err := os.Open(fileInfo.Name())
	if err != nil {
		rw.WriteHeader(500)
		io.WriteString(rw, err.Error())
	}
	io.Copy(rw, handle)
}

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		p := filepath.Join(rootPath, path.Clean(req.URL.Path))
		filepath.Abs(p)
		fileInfo, err := os.Stat(p)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			io.WriteString(rw, err.Error())
			return
		}

		fmt.Println(fileInfo)

		if fileInfo.Mode().IsRegular() {
			handleFile(fileInfo, rw, req)
			return
		}

		if fileInfo.IsDir() {
			handleDir(fileInfo, rw, req)
			return
		}

	})
	http.ListenAndServe(":8888", nil)
}
