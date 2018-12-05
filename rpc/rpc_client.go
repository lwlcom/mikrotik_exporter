package rpc

import (
	"fmt"

	"log"

	"github.com/lwlcom/mikrotik_exporter/connector"
)

// Client sends commands to mikrotik
type Client struct {
	conn  *connector.SSHConnection
	Debug bool
}

// NewClient creates a new client to connect to
func NewClient(ssh *connector.SSHConnection, debug bool) *Client {
	rpc := &Client{conn: ssh, Debug: debug}

	return rpc
}

// RunCommand runs a command on a mikrotik device
func (c *Client) RunCommand(cmd string) (string, error) {
	if c.Debug {
		log.Printf("Running command on %s: %s\n", c.conn.Host, cmd)
	}
	output, err := c.conn.RunCommand(fmt.Sprintf("%s without-paging", cmd))
	if err != nil {
		println(err.Error())
		return "", err
	}

	return output, nil
}
