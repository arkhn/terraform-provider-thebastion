package clients

import (
	"context"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Client struct {
	Host            string
	SshClientConfig ssh.ClientConfig
}

func NewClient(host, username, pathPrivateKey, pathKnownHost string) (*Client, error) {
	// Setup the ssh config
	privateKey, err := os.ReadFile(pathPrivateKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	hostkeyCallback, err := knownhosts.New(pathKnownHost)
	if err != nil {
		return nil, err
	}

	conf := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: hostkeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// All possible ssh keys values
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
		},
	}

	c := Client{
		Host:            host,
		SshClientConfig: *conf,
	}

	return &c, nil
}

// Run simple command help to check client connection to TheBastion
func (c *Client) CheckClientConnection(ctx context.Context) (*ResponseBastion, error) {
	command := "--osh help --json"
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return responseBastion, err
	}
	return responseBastion, nil
}
