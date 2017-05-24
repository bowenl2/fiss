# fiss

*File System Server*

fiss is a fast, lightweight HTTP file server written in Go.

## Usage:

```none
$ fiss [OPTIONS]

Usage:
  fiss [OPTIONS]

Application Options:
  -a, --address=              Address of interface on which to bind (0.0.0.0)
  -p, --port=                 Local port on which to listen (8080)
  -r, --root=                 Root directory of server (.)
  -v, --verbose               Print absurd amounts of debugging information
  -t, --ssh-tunnel            Use an SSH tunnel instead of listening on a local port
  -u, --ssh-username=         Username used to authenticate to SSH server
  -k, --ssh-key=              Path to private key as produced by ssh-keygen (~/.ssh/id_rsa)
  -s, --ssh-server=           Remote SSH server to request reverse port forwarding
  --ssh-outbound-port=        Port on which to connect to SSH server (22)
  -l, --ssh-inbound-port=     Port on which the SSH server should listen for incoming requests
  -i, --ssh-listen-interface= Interface on which the SSH server should listen (0.0.0.0)
  --password=                 Password required to authenticate (otherwise,
                              anonymous access is allowed)

Help Options:
  -h, --help                  Show this help message
```

### SSH Tunneling
SSH tunneling (`-t`) allows someone running an SSH server to access their files from behind a firewall or NAT with no additional configuration or external port forwarding.

For example, let's say that I'm at a university with a NAT such that I have a 10.x.x.x IP address with no way to accept incoming connections from outside, and I wish to serve some files. I have a VPS account at example.com, which I configured to accept my key using `~/.ssh/authorized_keys`. Then I can execute:
```
./fiss -t -u liam -k ~/.ssh/id_rsa -s example.com -l 1337
```
and my *local* directory will be available (by magic!) at `http://example.com:1337`
