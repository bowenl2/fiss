package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

var (
	rootPath = "/"
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

	tmplFuncs := map[string]interface{}{
		"fmtsize": func(s int64) string {
			return ByteSize(s).String()
		},
		"abspath": func(f os.FileInfo) string {
			return filepath.Join(path, f.Name())
		},
	}

	tmplAsset, err := Asset("directory-list.html")
	if err != nil {
		internalErrorHandler(err, rw, r)
		return
	}
	tmplString := string(tmplAsset)

	tmpl := template.Must(
		template.New("directory-list.html").Funcs(tmplFuncs).Parse(tmplString))

	err = tmpl.Execute(rw, dl)
	if err != nil {
		fmt.Println("Error executing the template: %v\n", err)
	}
}

func handleFile(path string, fileInfo os.FileInfo, rw http.ResponseWriter, r *http.Request) {
	content, err := os.Open(path)
	if err != nil {

		io.WriteString(rw, err.Error())
	}

	http.ServeContent(rw, r, path, fileInfo.ModTime(), content)
}

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		p := filepath.Join(rootPath, path.Clean(req.URL.Path))
		p, err := filepath.Abs(p)
		if err != nil {
			internalErrorHandler(err, rw, req)
			return
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			internalErrorHandler(err, rw, req)
			return
		}

		fmt.Printf("req: %v %v\n", req.RemoteAddr, p)

		// Intercept directories to perform listing
		if fileInfo.IsDir() {
			if req.FormValue("r") == "" {
				handleListDir(p, fileInfo, rw, req)
			} else {
				handleListDirRecursive(p, fileInfo, rw, req)
			}
			return
		}

		handleFile(p, fileInfo, rw, req)

	})
	http.ListenAndServe(":8888", nil)
}
