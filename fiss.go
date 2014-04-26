package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	rootPath = "/"
)

type ByFilename []os.FileInfo

func (l ByFilename) Len() int {
	return len(l)
}

func (l ByFilename) Swap(i, j int) {
	l[i], l[j] = a[j], a[i]
}

func (l ByFilename) Less(i, j) bool {
	l[i].Name() < l[j].Name()
}

func handleDir(fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	io.WriteString(rw, "Hello directory<br/>")
	dir, err := os.Open(fileInfo.Name())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fileEntries := make([]os.FileInfo, 0, len(files))
	dirEntries := make([]os.FileInfo, 0, len(files))
	otherEntries := make([]os.FileInfo, 0, len(files))

	// TODO: Sort Files
	for _, file := range files {
		if file.Mode().IsRegular() {
			append(fileEntries, file)
		} else if file.Mode().IsDir() {
			append(dirEntries, file)
		} else {
			append(otherEntries, file)
		}
	}

	io.WriteString(rw, fmt.Sprintf("%v<br/>", fileEntries))
	io.WriteString(rw, fmt.Sprintf("%v<br/>", dirEntries))
	io.WriteString(rw, fmt.Sprintf("%v<br/>", otherEntries))
}

func handleFile(fileInfo os.FileInfo, rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Length", size(fileInfo.Size()))
	handle, err := os.Open(fileInfo.Name())
	if err != nil {
		rw.WriteHeader(500)
		io.WriteString(rw, error.Error())
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
			handleDir(rw, req)
			return
		}

	})
	http.ListenAndServe(":8888", nil)
}
