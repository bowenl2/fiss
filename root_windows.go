package main

import (
	"github.com/AllenDang/w32"
	"fmt"
)

func listRoots() []string {
	driveBitset := w32.GetLogicalDrives()
	roots = make([]string, 0, 32)
	for i=0; i<32; i++ {
		if driveBitset & (1 << i) {
			append(roots, fmt.Sprintf(`%c:\`, 'A'+i))
		}
	}
	return roots
}
