package main

import (
	"fmt"
	"github.com/AllenDang/w32"
	"os"
)

func listRootPaths() []string {
	driveBitset := w32.GetLogicalDrives()
	roots := make([]string, 0, 32)
	for i := uint(0); i < 32; i++ {
		if (driveBitset & (1 << i)) != 0 {
			roots = append(roots, fmt.Sprintf(`%c:\`, 'A'+i))
		}
	}
	return roots
}

func listRootInfos() ([]os.FileInfo, []error) {
	paths := listRootPaths()
	infos := make([]os.FileInfo, 0, len(paths))
	errors := make([]error, 0, len(paths))
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		infos = append(infos, info)
	}
	return infos, errors
}
