package clients

import (
	"encoding/json"
	"fmt"
	"os"
	"terraform-provider-thebastion/utils"

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
	}

	c := Client{
		Host:            host,
		SshClientConfig: *conf,
	}

	return &c, nil
}

func (c *Client) GetListAccount() (ResponseBastionAccountList, error) {
	command := "--osh accountList --json"
	responseBastion, err := c.SendCommandBastion(command)
	if err != nil {
		return ResponseBastionAccountList{}, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return ResponseBastionAccountList{}, err
	}

	var responseBastionAccountList ResponseBastionAccountList
	if err := json.Unmarshal(marshal, &responseBastionAccountList); err != nil {
		panic(err)
	}
	return responseBastionAccountList, nil
}

func (c *Client) GetListIngressKeys(nameAccount string) (ResponseBastionListIngressKeys, error) {
	command := fmt.Sprintf("--osh accountListIngressKeys --account %s --json", nameAccount)
	responseBastion, err := c.SendCommandBastion(command)
	if err != nil {
		return ResponseBastionListIngressKeys{}, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return ResponseBastionListIngressKeys{}, err
	}

	var responseBastionListIngressKeys ResponseBastionListIngressKeys
	if err := json.Unmarshal(marshal, &responseBastionListIngressKeys); err != nil {
		panic(err)
	}
	return responseBastionListIngressKeys, err
}

// Return the ResponseBastion of selfAddIngressKey command on account
func (c *Client) AddListIngressKeys(nameAccount string, ingressKeys []string) error {
	for _, key := range ingressKeys {
		command := fmt.Sprintf("--osh adminSudo -- --sudo-as %s --sudo-cmd selfAddIngressKey -- --public-key '%s' --json", nameAccount, key)
		_, err := c.SendCommandBastion(command)
		if err != nil {
			return err
		}
	}
	return nil
}

// Return the ResponseBastion of selfDelIngressKey command on account
func (c *Client) DelListIngressKeys(nameAccount string, ingressKeys []string) error {
	responseBastion, err := c.GetListIngressKeys(nameAccount)
	if err != nil {
		return err
	}

	listIngressKeysId := make([]int64, 0)

	for _, key := range responseBastion.Value.Keys {
		if id := utils.FindStringIndex(ingressKeys, key.Line); id > -1 {
			listIngressKeysId = append(listIngressKeysId, key.Id)
		}
	}

	for _, keyId := range listIngressKeysId {
		command := fmt.Sprintf("--osh adminSudo -- --sudo-as %s --sudo-cmd selfDelIngressKey -- -I '%s' --json", nameAccount, fmt.Sprint(keyId))
		_, err := c.SendCommandBastion(command)
		if err != nil {
			// If the key is already created on the account
			return err
		}
	}
	return nil
}
