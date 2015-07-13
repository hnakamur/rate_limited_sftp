package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hnakamur/rate_limited_sftp"

	"golang.org/x/crypto/ssh"
)

var (
	hostAndPort            string
	user                   string
	topDir                 string
	limitKiloBitsPerSecond int64
)

func init() {
	flag.StringVar(&hostAndPort, "hostandport", "localhost:22", "sftp server host and port in host:port format")
	flag.StringVar(&topDir, "topdir", "/", "top directory")
	flag.StringVar(&user, "user", "root", "sftp user name")
	flag.Int64Var(&limitKiloBitsPerSecond, "limit", 64, "transfer rate limit in kilobits per second")
}

func run() error {
	password := os.Getenv("SFTP_PASSWORD")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", hostAndPort, config)
	if err != nil {
		return err
	}
	limitBytsPerSecond := limitKiloBitsPerSecond * 1000 / 8
	sftp, err := rate_limited_sftp.NewClient(client, limitBytsPerSecond, limitBytsPerSecond)
	if err != nil {
		return err
	}
	defer sftp.Close()

	walker := sftp.Walk(topDir)
	for walker.Step() {
		if walker.Err() != nil {
			continue
		}
		fmt.Println(walker.Path())
	}
	return nil
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		panic(err)
	}
}
