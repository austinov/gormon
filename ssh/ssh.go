package ssh

import (
	"bytes"
	"strings"

	c "github.com/austinov/gormon/config"
	"golang.org/x/crypto/ssh"
)

type Client interface {
	Host() string
	Connect() error
	Disconnect() error
	Run() (string, error)
}

type sshClient struct {
	config    c.HostConfig
	sshConfig *ssh.ClientConfig
	sshClient *ssh.Client
}

func NewClient(cfg c.HostConfig) Client {
	auths := make([]ssh.AuthMethod, 0)
	auths = addKeyAuth(auths, cfg.Keypath)

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: auths,
	}
	return &sshClient{
		config:    cfg,
		sshConfig: sshConfig,
	}
}

func (c *sshClient) Host() string {
	return strings.Split(c.config.Addr, ":")[0]
}

func (c *sshClient) Connect() error {
	var err error
	c.sshClient, err = ssh.Dial("tcp", c.config.Addr, c.sshConfig)
	return err
}

func (c *sshClient) Disconnect() error {
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

func (c *sshClient) Run() (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(c.config.Command + " && echo '# PS' && ps aux | awk '{print $2, $3, $11}' | grep redis-server")
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
