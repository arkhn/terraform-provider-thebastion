package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"golang.org/x/crypto/ssh"
)

// Connection to the bastion
func (c *Client) connectSsh(command string) (string, error) {
	// Start connection to the bastion
	conn, err := ssh.Dial("tcp", c.Host, &c.SshClientConfig)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	stdout, _ := session.StdoutPipe()
	if err := session.Run(command); err != nil {
		panic(err)
	}
	out, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	res := string(out)

	return res, nil
}

// Convert string from bastion into struct ResponseBastion
func unmarshalBastionResponse(outputSsh string) (ResponseBastion, error) {
	r, _ := regexp.Compile("(?s)JSON_START(.*?)JSON_END.*")
	jsonOutput := r.FindStringSubmatch(outputSsh)[1]

	// Convert string json into responseBastion struct
	var results ResponseBastion
	err := json.Unmarshal([]byte(jsonOutput), &results)
	if err != nil {
		return ResponseBastion{}, err
	}

	if results.ErrorCode != "OK" {
		return results, fmt.Errorf(
			"errorCode: %s, msg: %s", results.ErrorCode, results.ErrorMessage,
		)
	}

	return results, nil

}

// Wrapper for connect and unmarshall function
func (c *Client) SendCommandBastion(command string) (ResponseBastion, error) {
	outputSsh, err := c.connectSsh(command)
	if err != nil {
		return ResponseBastion{}, err
	}

	res, err := unmarshalBastionResponse(outputSsh)
	if err != nil {
		return res, err
	}

	return res, err
}
