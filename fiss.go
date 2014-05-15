package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func makeTCPListener(string localInterface, port int) (*net.Listener, error) {
	addr := net.TCPAddr{
		IP:   net.ParseIP(localInterface),
		Port: port,
	}
	listener, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return nil, err
	}

	return listener.(*net.Listener), err
}

func main() {
	options, err := parseOptions()
	if err != nil {
		return
	}
	absRoot, _ := filepath.Abs(options.Root)

	fmt.Printf(
		"%s: %s:%d %s\n",
		os.Args[0],
		options.Address,
		options.Port,
		absRoot)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		p := filepath.Join(absRoot, path.Clean(req.URL.Path))
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
				handleListDir(absRoot, p, fileInfo, rw, req)
			} else {
				handleListDirRecursive(p, fileInfo, rw, req)
			}
			return
		}

		handleFile(p, fileInfo, rw, req)

	})

	// Determine where to listen for connections
	var listener *net.Listener
	switch options.UseSSHTunnel {
	case false:
		listener, err = makeTCPListener(options.Address, options.Port)
	case true:
		listener, err = makeSSHTunnel(options.SSHServerEndpoint, options.Username, options.PrivateKeyPath)
	}
	if err != nil {
		fmt.Printf(err)
		return
	}

	http.Serve(listener, nil)
}
