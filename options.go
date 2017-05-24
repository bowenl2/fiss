package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Address string `short:"a" long:"address" description:"Address of interface on which to bind" default:"0.0.0.0"`
	Port    int    `short:"p" long:"port" description:"Local port on which to listen" default:"8080"`
	Root    string `short:"r" long:"root" description:"Root directory of server" default:"."`
	Verbose bool   `short:"v" long:"verbose" description:"Print absurd amounts of debugging information"`

	// SSH Stuff
	UseSSHTunnel       bool   `short:"t" long:"ssh-tunnel" description:"Use an SSH tunnel instead of listening on a local port"`
	Username           string `short:"u" long:"ssh-username" description:"Username used to authenticate to SSH server"`
	PrivateKeyPath     string `short:"k" long:"ssh-key" description:"Path to private key as produced by ssh-keygen" default:"~/.ssh/id_rsa"`
	SSHServer          string `short:"s" long:"ssh-server" description:"Remote SSH server to request reverse port forwarding"`
	SSHOutboundPort    int    `long:"ssh-outbound-port" description:"Port on which to connect to SSH server" default:"22"`
	SSHInboundPort     int    `short:"l" long:"ssh-inbound-port" description:"Port on which the SSH server should listen for incoming requests"`
	SSHListenInterface string `short:"i" long:"ssh-listen-interface" description:"Interface on which the SSH server should listen" default:"0.0.0.0"`
	HTTPPassword       string `long:"password" description:"Password required to authenticate (otherwise, anonymous access is allowed)"`
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
