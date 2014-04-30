package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Port uint16 `short:"p" long:"port" description:"Local port on which to listen"`
	Address string `short:"a" long:"address" description:"Address to bind to (0.0.0.0 for all)"`
	Root string `short:"r" long:"root" description:"Root directory of server"`
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
