package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strconv"
)

func makeSSHTunnel(
	username string,
	sshServer string,
	sshOutboundPort int,
	sshListenInterface string,
	sshInboundPort int,
	privateKeyPath string) (net.Listener, error) {

	pemBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	conn, err := ssh.Dial("tcp",
		net.JoinHostPort(sshServer, strconv.Itoa(sshOutboundPort)),
		sshConfig)
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}

	remoteListenEndpoint := net.JoinHostPort(
		sshListenInterface,
		strconv.Itoa(sshInboundPort))
	listener, err := conn.Listen("tcp", remoteListenEndpoint)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
