package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

var (
	rootPath = "/"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		filePath := filepath.Join(path.Clean(req.URL.Path), filePath)
		filepath.Abs(filePath)

		fmt.Println(filePath)
	})
	http.ListenAndServe(":8888", nil)
}
