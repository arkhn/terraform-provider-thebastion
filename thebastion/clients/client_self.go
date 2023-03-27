package clients

import (
	"context"
	"encoding/json"
)

// Methods about the user
func (c *Client) GetInfo(ctx context.Context) (*ResponseBastionInfo, error) {
	command := "--osh info --json"
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionInfo ResponseBastionInfo
	err = json.Unmarshal(marshal, &responseBastionInfo)
	if err != nil {
		return nil, err
	}

	return &responseBastionInfo, nil
}
