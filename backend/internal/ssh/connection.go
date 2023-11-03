package ssh

import (
	"golang.org/x/crypto/ssh"
)

func CreateConfig(user string, privateKey string) (ssh.ClientConfig, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	return ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}, nil
}

func CreateConnection(user, addr string, privateKey string) (*ssh.Client, error) {
	cfg, err := CreateConfig(user, privateKey)
	if err != nil {
		return nil, err
	}
	return ssh.Dial("tcp", addr, &cfg)
}

func CheckConnection(user, addr string, privateKey string) error {
	conn, err := CreateConnection(user, addr, privateKey)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
