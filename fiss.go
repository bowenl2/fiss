package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// ResponseFormat specifies format in which to render response
type ResponseFormat uint32

const (
	// FmtAuto automatically determines a reasonable choice
	FmtAuto ResponseFormat = iota
	// FmtHTML forces rendering in HTML (directory only)
	FmtHTML
	// FmtJSON forces rendering in JSON (directory only)
	FmtJSON
	// FmtCSV forces rendering in CSV (directory only)
	FmtCSV
	// FmtForceDownload forces the response to download a file,
	// which involves creating a ZIP archive for a directory
	FmtForceDownload
)

// Context of environment relevant to handlers
type Context struct {
	// AppHandler describes the application context
	App *AppHandler
	// Recursive specifies that the response associated with the
	// request path should recursively include data from subdirectories
	// FIXME: Not implemented
	Recursive bool
	// Response Format
	Format ResponseFormat
	// Absolute filesystem path to file
	FSPath string
	// FSPath Stat
	FSInfo os.FileInfo
}

func makeTCPListener(localInterface string, port int) (net.Listener, error) {
	addr := &net.TCPAddr{
		IP:   net.ParseIP(localInterface),
		Port: port,
	}
	listener, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return nil, err
	}

	return listener, err
}

func main() {
	options, err := parseOptions()
	if err != nil {
		fmt.Printf("fatal: %v", err)
		os.Exit(1)
	}

	absRoot, _ := filepath.Abs(options.Root)
	app := &AppHandler{
		RootPath: absRoot,
	}
	http.Handle("/", app)

	// Determine where to listen for connections
	var listener net.Listener
	if options.UseSSHTunnel {
		listener, err = makeSSHTunnel(
			options.Username,
			options.SSHServer,
			options.SSHOutboundPort,
			options.SSHListenInterface,
			options.SSHInboundPort,
			options.PrivateKeyPath)
		fmt.Printf("ssh: %s@%s:%d (listen on %s:%d) using key: %s\n",
			options.Username,
			options.SSHServer,
			options.SSHOutboundPort,
			options.SSHListenInterface,
			options.SSHInboundPort,
			options.PrivateKeyPath)
	} else {
		listener, err = makeTCPListener(options.Address, options.Port)
		fmt.Printf(
			"%s: %s:%d %s\n",
			os.Args[0],
			options.Address,
			options.Port,
			absRoot)
	}

	if err != nil {
		fmt.Printf("fatal: %v\n", err)
		return
	}

	http.Serve(listener, nil)
}
