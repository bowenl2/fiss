package main

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

// The private key we will authenticate using
// currently only supports one key
type ClientKeyring rsa.PrivateKey

func LoadKeyring(keyPath string) (*ClientKeyring, error) {
	// Read the key material
	privateKeyPEM, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Decode the key material
	block, _ := pem.Decode(privateKeyPEM)

	// Parse the key
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Could not parse the key from the decoded PEM
		return nil, err
	}

	// Everything turned out fine!
	keyring := ClientKeyring(*rsaKey)
	return &keyring, nil
}

func (keyring ClientKeyring) Key(i int) (ssh.PublicKey, error) {
	// Only support one key
	if i != 0 {
		return nil, nil
	}

	// Wrap the RSA public key in the SSH package's PublicKey wrapper
	publicKey, err := ssh.NewPublicKey(rsa.PrivateKey(keyring).PublicKey)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (keyring *ClientKeyring) Sign(i int,
	rand io.Reader,
	data []byte) (sig []byte, err error) {

	// Only support one key
	if i != 0 {
		return nil, nil
	}

	hashImpl := crypto.SHA1
	hashFunc := hashImpl.New()
	hashFunc.Write(data)
	digest := hashFunc.Sum(nil)
	privateKey := rsa.PrivateKey(*keyring)
	return rsa.SignPKCS1v15(rand, &privateKey, hashImpl, digest)
}

// Turns 1337:1.2.3.4[:127.0.0.1] into components:
// remote_port, remote_host, remote_listen_interface (optional)
func parseEndpoint(endpoint string) (int, *net.IPAddr, *net.IPAddr, error) {
	segs := strings.Split(endpoint, ":")
	if len(segs) < 2 || len(segs) > 3 {
		return 0, nil, nil, fmt.Errorf("endpoint '%s' has invalid number of segments", endpoint)
	}

	remotePort, err := strconv.Atoi(segs[0])
	if err != nil {
		return 0, nil, nil, fmt.Errorf("given non-numeric listening port: %s", segs[0])
	}

	remoteHost := net.ParseIP(segs[1])
	if remoteHost == nil {
		return 0, nil, nil, fmt.Errorf("given invalid IP address: %s", remoteHost)
	}

	var remoteListenInterface net.IP
	if len(segs) == 3 {
		// tell the remote host on which interface he should listen
		remoteListenInterface = net.ParseIP(segs[2])
		if remoteListenInterface == nil {
			err = fmt.Errorf("remote listen interface %s specified but invalid", segs[2])
			return 0, nil, nil, err
		}
	}

	return remotePort, &net.IPAddr{IP: remoteHost}, &net.IPAddr{IP: remoteListenInterface}, nil
}

func makeSSHTunnel(endpoint string, username string, keyPath string) (net.Listener, error) {
	remotePort, remoteAddr, remoteListenInterface, err := parseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	keyring, err := LoadKeyring(keyPath)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.ClientAuth{ssh.ClientAuthKeyring(keyring)},
	}

	conn, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%d", remoteAddr.String(), 22),
		sshConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// remoteListenEndpoint := fmt.Sprintf("%s:%d",
	// 	remoteListenInterface.String(), remotePort)
	listener, err := conn.ListenTCP(&net.TCPAddr{
		IP: remoteListenInterface.IP,
		Port: remotePort,
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}
