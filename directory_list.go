package main

import (
	"fmt"
	"os"
	"strings"
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
	Machine     string
	Path        string
	BaseInfo    os.FileInfo
	Entries     []os.FileInfo
	BreadCrumbs []breadCrumb
}

type breadCrumb struct {
	Link     string
	Basename string
}

func joinBreadCrumb(segs ...string) string {
	return fmt.Sprintf("/%s",
		strings.TrimSpace(
			strings.Join(segs, "/")))
}

func makeBreadCrumbs(path string) []breadCrumb {
	segs := strings.Split(path, "/")
	breadCrumbs := make([]breadCrumb, 0, len(segs))
	for i, seg := range segs {
		breadCrumbs = append(breadCrumbs,
			breadCrumb{
				Basename: seg,
				Link:     joinBreadCrumb(segs[:i+1]...),
			})
	}
	return breadCrumbs
}
