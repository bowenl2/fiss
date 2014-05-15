package main

import (
	"os"
)

type RootListViewModel struct {
	Machine   string
	RootInfos []os.FileInfo
	Errors    []error
}
