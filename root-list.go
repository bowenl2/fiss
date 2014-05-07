package main

import (
	"os"
)

type RootListViewModel struct {
	RootInfos []os.FileInfo
	Errors []error
}
