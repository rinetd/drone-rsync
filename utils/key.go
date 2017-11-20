package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return "/root"
}

// WriteKey writes the private key.
// chmod 700 ~/.ssh
// chmod 600 ~/.ssh/id_rsa
// chmod 644 ~/.ssh/id_rsa.pub
// chmod 400 ~/.ssh/authorized_keys
func WriteKey(privateKey string) error {
	if privateKey == "" {
		log.Println("private key is none")
		return nil
	}

	home := homeDir()
	// log.Println("【home】", home)
	sshpath := filepath.Join(
		home,
		".ssh")

	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return err
	}

	confpath := filepath.Join(
		sshpath,
		"config")

	privpath := filepath.Join(
		sshpath,
		"id_rsa")

	ioutil.WriteFile(
		confpath,
		[]byte("StrictHostKeyChecking no\n"),
		0700)

	return ioutil.WriteFile(
		privpath,
		[]byte(privateKey),
		0600)
}

const netrcFile = `
machine %s
login %s
password %s
`

// WriteNetrc writes the netrc file.
func WriteNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}

	netrcContent := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	home := homeDir()

	netpath := filepath.Join(
		home,
		".netrc")

	return ioutil.WriteFile(
		netpath,
		[]byte(netrcContent),
		0600)
}
