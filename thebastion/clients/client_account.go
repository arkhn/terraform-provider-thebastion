package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// Account methods
func (c *Client) GetListAccount(ctx context.Context) (*ResponseBastionAccountList, error) {
	command := "--osh accountList --json"
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	responseBastionAccountList := ResponseBastionAccountList{}
	err = json.Unmarshal(marshal, &responseBastionAccountList)
	if err != nil {
		return nil, err
	}

	return &responseBastionAccountList, nil
}

func (c *Client) GetAccount(ctx context.Context, name string) (*Account, error) {
	responseBastionAccountList, err := c.GetListAccount(ctx)
	if err != nil {
		return nil, err
	}

	account, ok := responseBastionAccountList.Value[name]
	if !ok {
		return nil, nil
	}

	return &account, nil
}

func (c *Client) CreateAccount(ctx context.Context, name string, uid int64, ingress_keys []string) (*ResponseBastion, error) {
	key, ingress_keys_remaining := ingress_keys[0], ingress_keys[1:]
	command := fmt.Sprintf("--osh accountCreate --account %s --uid %s --public-key \"%s\" --json", name, fmt.Sprint(uid), key)

	// Command to create the account
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		// If the key is already created on the account
		if (!reflect.DeepEqual(responseBastion, ResponseBastion{}) && responseBastion.ErrorCode == "KO_DUPLICATE_KEY") {
			return nil, fmt.Errorf("bastion account already created: %s", name)
		}
		return nil, err
	}

	if len(ingress_keys_remaining) > 0 {
		err = c.AddListIngressKeys(ctx, name, ingress_keys_remaining)
		if err != nil {
			return nil, err
		}
	}

	return responseBastion, nil
}

func (c *Client) DeleteAccount(ctx context.Context, name string) (*ResponseBastion, error) {
	command := fmt.Sprintf("--osh accountDelete --account %s --no-confirm --json", name)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	return responseBastion, nil
}
