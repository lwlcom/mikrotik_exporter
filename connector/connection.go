package connector

import (
	"bytes"
	"os"
	"strings"

	"sync"

	"time"

	"golang.org/x/crypto/ssh"
)

const timeoutInSeconds = 5

var (
	cachedConfig *ssh.ClientConfig
	lock         = &sync.Mutex{}
	identity ssh.AuthMethod
)

func ReadIdentity(filename string) error {
	key, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	identity = ssh.PublicKeys(signer)
	return nil
}

// NewSSSHConnection connects to device
func NewSSSHConnection(host, user, pass string) (*SSHConnection, error) {
	if !strings.Contains(host, ":") {
		host = host + ":22"
	}

	c := &SSHConnection{Host: host}
	err := c.Connect(user, pass)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// SSHConnection encapsulates the connection to the device
type SSHConnection struct {
	conn *ssh.Client
	Host string
}

// Connect connects to the device
func (c *SSHConnection) Connect(user, pass string) error {
	config := &ssh.ClientConfig{
		User:            user,
		Timeout:         timeoutInSeconds * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.KeyExchanges = []string{
		"diffie-hellman-group-exchange-sha256"}

	if identity != nil {
		config.Auth = append(config.Auth, identity)
	}
	if pass != "" {
		config.Auth = append(config.Auth, ssh.Password(pass))
	}

	var err error
	c.conn, err = ssh.Dial("tcp", c.Host, config)
	return err
}

// RunCommand runs a command against the device
func (c *SSHConnection) RunCommand(cmd string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b = &bytes.Buffer{}
	session.Stdout = b

	err = session.Run(cmd)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// Close closes connection
func (c *SSHConnection) Close() {
	if c.conn == nil {
		return
	}

	c.conn.Close()
	c.conn = nil
}
