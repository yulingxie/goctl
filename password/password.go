package password

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/nitro/nitro/v3/crypto"
)

func Crypt(cliContext *cli.Context) error {
	password := cliContext.String("password")
	if len(password) == 0 {
		return errors.New("未指定密钥，通过 -password 指定")
	}
	content := cliContext.String("content")
	if len(content) == 0 {
		return errors.New("未指定加密内容，通过 -content 指定")
	}
	cryptor := crypto.NewPBE512AndAES256Cryptor(password)
	crypted, err := cryptor.Encrypt([]byte(content))
	if err != nil {
		return err
	}
	fmt.Printf("crypted result: %s\n", crypted)
	return nil
}

func Decrypt(cliContext *cli.Context) error {
	password := cliContext.String("password")
	if len(password) == 0 {
		return errors.New("未指定密钥，通过 -password 指定")
	}
	content := cliContext.String("content")
	if len(content) == 0 {
		return errors.New("未指定解密内容，通过 -content 指定")
	}
	decryptor := crypto.NewPBE512AndAES256Cryptor(password)
	decrypted, err := decryptor.Decrypt([]byte(content))
	if err != nil {
		return err
	}
	fmt.Printf("decrypted result: %s\n", decrypted)
	return nil
}
