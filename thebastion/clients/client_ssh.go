package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		out, read_err := io.ReadAll(stdout)
		if read_err != nil {
			return string(out), read_err
		}
		return string(out), err
	}

	out, err := io.ReadAll(stdout)
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

// Convert string from bastion into struct ResponseBastion
func unmarshalBastionResponse(outputSsh string) (*ResponseBastion, error) {
	r, _ := regexp.Compile("(?s)JSON_START(.*?)JSON_END.*")
	jsonOutput := r.FindStringSubmatch(outputSsh)[1]

	// Convert string json into responseBastion struct
	results := ResponseBastion{}
	err := json.Unmarshal([]byte(jsonOutput), &results)
	if err != nil {
		return &results, err
	}

	return &results, nil
}

// Wrapper for connect and unmarshall function
func (c *Client) SendCommandBastion(ctx context.Context, command string) (*ResponseBastion, error) {
	tflog.Info(ctx, fmt.Sprintf("Request bastion: %s", command))
	outputSsh, err := c.connectSsh(command)
	if err != nil {
		// Check if the command run exits with a non zero exit status.
		if _, ok := err.(*ssh.ExitError); ok {
			res, err := unmarshalBastionResponse(outputSsh)
			if err != nil {
				return res, err
			}
			return res, fmt.Errorf("thebastion error code: %s / msg: %s", res.ErrorCode, res.ErrorMessage)
		}
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("Response bastion: %s", outputSsh))
	res, err := unmarshalBastionResponse(outputSsh)
	if err != nil {
		return res, err
	}

	tflog.Info(ctx, fmt.Sprintf("Struct from response bastion: %s", res))
	return res, nil
}
