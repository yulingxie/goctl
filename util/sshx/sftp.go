package sshx

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpClient struct {
	*sftp.Client
}

func NewSftpClient(sshClient *ssh.Client) (*SftpClient, error) {
	if client, err := sftp.NewClient(sshClient); err != nil {
		return nil, err
	} else {
		return &SftpClient{client}, nil
	}
}

// UploadFile 先删除远端remote文件, 再拷贝local到远端remote
func (c *SftpClient) UploadFile(remote, local string) error {
	srcFile, err := os.Open(local)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer srcFile.Close()
	if err != nil {
		return err
	}
	destFile, err := c.OpenFile(remote, os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}
	defer destFile.Close()
	buffer := make([]byte, 1024*1024*1024)
	for {
		n, err := srcFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("读取文件出错", err)
				panic(err)
			}
		}
		destFile.Write(buffer[:n])
		//注意，由于文件大小不定，不可直接使用buffer，否则会在文件末尾重复写入，以填充1024的整数倍
	}
	return nil
}
