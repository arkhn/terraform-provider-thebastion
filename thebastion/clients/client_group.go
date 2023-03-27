package clients

import (
	"context"
	"encoding/json"
	"fmt"
)

// Methods about the user
func (c *Client) GetListGroup(ctx context.Context) (*ResponseBastionGroupList, error) {
	command := "--osh groupList --all --json"
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionGroupList ResponseBastionGroupList
	err = json.Unmarshal(marshal, &responseBastionGroupList)
	if err != nil {
		return nil, err
	}

	return &responseBastionGroupList, nil
}

// Cannot create a group encrypted yet
func (c *Client) CreateGroup(ctx context.Context, groupName, owner, algo string, size int) (*ResponseBastionCreateGroup, error) {
	command := fmt.Sprintf("--osh groupCreate --group %s --owner %s --algo %s --size %s --json", groupName, owner, algo, fmt.Sprint(size))
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionCreateGroup ResponseBastionCreateGroup
	err = json.Unmarshal(marshal, &responseBastionCreateGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionCreateGroup, nil
}
