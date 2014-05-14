package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Address string `short:"a" long:"address" description:"Address of interface on which to bind" default:"0.0.0.0"`
	Port uint16 `short:"p" long:"port" description:"Local port on which to listen" default:"8080"`
	Root string `short:"r" long:"root" description:"Root directory of server" default:"."`
	Verbose bool `short:"v" long:"verbose" description:"Print absurd amounts of debugging information"`

	// SSH Stuff
	UseSSHTunnel bool `short:"s" long:"ssh-tunnel" description:"Use an SSH tunnel instead of listening on a local port"`
	PrivateKeyPath string `short:"k" long:"ssh-key" description:"Path to private key as produced by ssh-keygen" default:"~/.ssh/id_rsa"`
	SSHServerEndpoint string `short:"e" long:"ssh-endpoint" description:"Request remote SSH server to listen on remote_interface (default all interfaces) at remote_port: [remote_interface:]remote_host:remote_port"`
}

func parseOptions() (*Options, error) {
	options := &Options{}
	parser := flags.NewParser(options, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return options, nil
}
