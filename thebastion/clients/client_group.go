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

func (c *Client) DestroyGroup(ctx context.Context, groupName string) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh groupDestroy --group %s --no-confirm --json", groupName)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionDeleteGroup ResponseBastion
	err = json.Unmarshal(marshal, &responseBastionDeleteGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionDeleteGroup, nil
}

// Add a server to a group
func (c *Client) AddServerToGroup(ctx context.Context, groupName string, host string, user string, port int, comment string) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh groupAddServer --group %s --host %s --user %s --port %s --comment %s --json", groupName, host, user, fmt.Sprint(port), comment)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionAddServerToGroup ResponseBastion
	err = json.Unmarshal(marshal, &responseBastionAddServerToGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionAddServerToGroup, nil
}

// Delete a server from a group
func (c *Client) DeleteServerFromGroup(ctx context.Context, groupName string, host string, user string, port int) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh groupDelServer --group %s --host %s --user %s --port %s --json", groupName, host, user, fmt.Sprint(port))
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionDeleteServerFromGroup ResponseBastion
	err = json.Unmarshal(marshal, &responseBastionDeleteServerFromGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionDeleteServerFromGroup, nil
}

// Add owner to a group
func (c *Client) AddOwnerToGroup(ctx context.Context, groupName string, owner string) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh groupAddOwner --group %s --account %s --json", groupName, owner)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionAddOwnerToGroup ResponseBastion
	err = json.Unmarshal(marshal, &responseBastionAddOwnerToGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionAddOwnerToGroup, nil
}

// Delete owner from a group
func (c *Client) DeleteOwnerFromGroup(ctx context.Context, groupName string, owner string) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh groupDelOwner --group %s --account %s --json", groupName, owner)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	var responseBastionDeleteOwnerFromGroup ResponseBastion
	err = json.Unmarshal(marshal, &responseBastionDeleteOwnerFromGroup)
	if err != nil {
		return nil, err
	}

	return &responseBastionDeleteOwnerFromGroup, nil
}
