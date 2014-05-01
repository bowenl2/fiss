package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Address string `short:"a" long:"address" description:"Address of interface on which to bind" default:"0.0.0.0"`
	Port uint16 `short:"p" long:"port" description:"Local port on which to listen" default:"8080"`
	Root string `short:"r" long:"root" description:"Root directory of server" default:"."`
	Verbose bool `short:"v" long:"verbose" description:"Print absurd amounts of debugging information"`
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
