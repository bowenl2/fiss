package main

import (
	"os"
)

// Directory List
// Natural sort (by directory then name)
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
