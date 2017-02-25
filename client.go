package quasar

import (
	"net/rpc"

	"github.com/pkg/errors"
)

type Client struct {
	client *rpc.Client
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewClient(c config) (*Client, error) {
	client, err := rpc.DialHTTP("tcp", c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "cannot dialing RPC client:")
	}

	return &Client{client: client}, nil
}

func (c *Client) GetEnv(daemon, envname string) (string, error) {
	args := RPCArgs{
		Name:    daemon,
		Envname: envname,
	}
	var resp string

	err := c.client.Call("Server.GetEnv", &args, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *Client) EnvClose(daemon, envname string) error {
	args := RPCArgs{
		Name:    daemon,
		Envname: envname,
	}
	resp := ""

	err := c.client.Call("Server.Close", &args, &resp)
	if err != nil {
		return err
	}

	return nil
}
