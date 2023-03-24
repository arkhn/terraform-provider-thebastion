package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"terraform-provider-thebastion/utils"
)

// Ingress keys methods

func (c *Client) GetListIngressKeys(ctx context.Context, nameAccount string) (*ResponseBastionListIngressKeys, error) {
	command := fmt.Sprintf("--osh accountListIngressKeys --account %s --json", nameAccount)
	responseBastion, err := c.SendCommandBastion(ctx, command)
	if err != nil {
		return nil, err
	}

	// map to struct
	marshal, err := json.Marshal(responseBastion)
	if err != nil {
		return nil, err
	}

	responseBastionListIngressKeys := ResponseBastionListIngressKeys{}
	err = json.Unmarshal(marshal, &responseBastionListIngressKeys)
	if err != nil {
		return nil, err
	}

	return &responseBastionListIngressKeys, nil
}

// Return the ResponseBastion of selfAddIngressKey command on account
func (c *Client) AddListIngressKeys(ctx context.Context, nameAccount string, ingressKeys []string) error {
	for _, key := range ingressKeys {
		command := fmt.Sprintf("--osh adminSudo -- --sudo-as %s --sudo-cmd selfAddIngressKey -- --public-key '%s' --json", nameAccount, key)
		_, err := c.SendCommandBastion(ctx, command)
		if err != nil {
			return err
		}
	}
	return nil
}

// Return the ResponseBastion of selfDelIngressKey command on account
func (c *Client) DelListIngressKeys(ctx context.Context, nameAccount string, ingressKeys []string) error {
	responseBastion, err := c.GetListIngressKeys(ctx, nameAccount)
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
		_, err := c.SendCommandBastion(ctx, command)
		if err != nil {
			// If the key is already created on the account
			return err
		}
	}
	return nil
}

func (c *Client) UpdateListIngressKeys(ctx context.Context, name string, oldIngressKeys []string, newIngressKeys []string) error {
	if !reflect.DeepEqual(oldIngressKeys, newIngressKeys) {
		leftOnly, rightOnly := utils.CompareLists(oldIngressKeys, newIngressKeys)

		// An account cannot have only one key
		// to avoid that we add before we remove any keys
		err := c.AddListIngressKeys(ctx, name, rightOnly)
		if err != nil {
			return err
		}

		err = c.DelListIngressKeys(ctx, name, leftOnly)
		if err != nil {
			return err
		}
	}
	return nil
}
