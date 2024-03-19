package sshx

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewClient(opts ...Option) (*ssh.Client, error) {
	o := NewOptions(opts...)
	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         o.TimeOut, //ssh 连接time out
		User:            o.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
	}
	if o.AuthType == AuthTypePass {
		config.Auth = []ssh.AuthMethod{ssh.Password(o.Auth)}
	} else if o.AuthType == AuthTypePublicKey {
		// signer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(o.Auth), []byte("jenkins"))
		signer, err := ssh.ParsePrivateKey([]byte(o.Auth))
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	addr := fmt.Sprintf("%s:%d", o.Host, o.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		err = errors.Wrapf(err, "")
		return nil, err
	}
	return client, err
}
