package deploy

import (
	"fmt"
	"io/ioutil"
	"time"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/sshx"
	"golang.org/x/crypto/ssh"
)

type Deployer struct {
	cfg        *config
	gen        *ScriptGenerator
	client     *ssh.Client
	session    *ssh.Session
	sftpClient *sshx.SftpClient
}

var (
	testAndDevHostOpts = []sshx.Option{
		sshx.Host("172.16.1.155"),
		sshx.Port(65022),
		sshx.TimeOut(3 * time.Second),
	}
	qaHostOpts = []sshx.Option{
		sshx.Host("172.16.1.174"),
		sshx.Port(65022),
		sshx.TimeOut(3 * time.Second),
	}
)

func getHostOpts(env string) []sshx.Option {
	switch env {
	case "dev", "test":
		return testAndDevHostOpts
	case "qa":
		return qaHostOpts
	default:
		panic(fmt.Sprintf("环境配置错误%v", env))
	}
}

func NewDeployer(cfg *config, gcfg *generatorConfig) (*Deployer, error) {
	opts := getHostOpts(cfg.Env)
	if len(cfg.Password) != 0 {
		opts = append(opts, sshx.User(cfg.User), sshx.Auth(cfg.Password), sshx.AuthType(sshx.AuthTypePass))
	} else if len(cfg.IdentityFile) != 0 {
		if key, err := ioutil.ReadFile(cfg.IdentityFile); err != nil {
			return nil, err
		} else {
			opts = append(opts, sshx.User(cfg.User), sshx.Auth(string(key)), sshx.AuthType(sshx.AuthTypePublicKey))
		}
	}
	if client, err := sshx.NewClient(opts...); err != nil {
		return nil, err
	} else {
		if session, err := client.NewSession(); err != nil {
			client.Close()
			return nil, err
		} else {
			d := &Deployer{
				cfg:     cfg,
				gen:     &ScriptGenerator{gcfg},
				client:  client,
				session: session,
			}
			if sftpCli, err := sshx.NewSftpClient(client); err != nil {
				session.Close()
				return nil, err
			} else {
				d.sftpClient = sftpCli
			}
			return d, nil
		}
	}
}

func (d *Deployer) Deploy() error {
	defer func() {
		d.session.Close()
		d.sftpClient.Close()
	}()
	script, err := d.gen.Generate()
	if err != nil {
		return err
	}

	if len(d.gen.cfg.RemoteFile) != 0 && len(d.gen.cfg.LocalFile) != 0 {
		fmt.Printf("文件上传中 %v\n", d.gen.cfg.LocalFile)
		if err := d.sftpClient.UploadFile(d.gen.cfg.RemoteFile, d.gen.cfg.LocalFile); err != nil {
			return err
		}
		fmt.Println("上传完成")
	}
	fmt.Printf("部署中, 请勿关闭进程 ...\n")
	if out, err := d.session.Output(script); err != nil {
		return err
	} else {
		fmt.Println("部署完成")
		fmt.Printf("输出日志: \n%v", string(out))
		return nil
	}
}
