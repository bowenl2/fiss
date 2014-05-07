package main

import (
	"fmt"
	"github.com/AllenDang/w32"
)

func listRoots() []string {
	driveBitset := w32.GetLogicalDrives()
	roots := make([]string, 0, 32)
	for i := uint(0); i < 32; i++ {
		if (driveBitset & (1 << i)) != 0 {
			roots = append(roots, fmt.Sprintf(`%c:\`, 'A'+i))
		}
	}
	return roots
}
