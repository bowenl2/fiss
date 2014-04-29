package main

import (
	"fmt"
	
	
	"net/http"
	"os"
	"path"
	"path/filepath"
	
)

var (
	rootPath = "/"
)


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
