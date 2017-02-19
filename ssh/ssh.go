package ssh

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

type Client interface {
	Connect() error
	Run() (string, error)
}

type sshClient struct {
	config    Config
	sshConfig *ssh.ClientConfig
	sshClient *ssh.Client
}

func New(cfg Config) Client {
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

func (c *sshClient) Connect() error {
	var err error
	c.sshClient, err = ssh.Dial("tcp", c.config.Addr, c.sshConfig)
	return err
}

func (c *sshClient) Run() (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(c.config.Command)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
