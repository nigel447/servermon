package sshclient

import (
	"climon/dcrypto"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

func (cli *Client) BootStrapClient(ctx context.Context) {
	cli.getClientConfig()
	go func() {
		// recieve cancel mssg
		<-ctx.Done()
		fmt.Println("recieve cancel mssg shutting service down")
		os.Exit(1)
	}()
	conn, err := ssh.Dial("tcp", cli.Addr, cli.config)
	handleError(err)
	defer conn.Close()
	go func() {
		fmt.Println("BootStrapClient opening a channel")
		channel, requests, err := conn.OpenChannel("rpc-remote", s("%s extra data", "init string"))
		handleError(err)
		go ssh.DiscardRequests(requests)
		//test send data
		n := 1
		for {
			_, err := channel.Write(s("#%d send data channel", n))
			handleError(err)
			n++
			time.Sleep(3 * time.Second)
			if n > 3 {
				break
			}
		}
	}()
	// block
	<-ctx.Done()
}

func verifyServer(hostname string, remote net.Addr, key ssh.PublicKey) error {
	fp := ssh.FingerprintSHA256(key)
	if dcrypto.FINGER_PRINT == fp {
		return nil
	}
	// return error to reject the remote server
	return errors.New("invalid server fingerprint")
}

func (cli *Client) getClientConfig() {
	user := "climon"
	signer := dcrypto.ParsePemOpenSSHKeyToSigner(dcrypto.PriKeyPem)
	//hostKey, err :=ssh.ParseRawPrivateKey([]byte(dcrypto.PriKeyPem) )
	//FixedHostKey requries hostKey via ssh.ParsePrivateKey(key)
	cli.config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: verifyServer,
		//HostKeyCallback: ssh.FixedHostKey(hostKey.(ssh.PublicKey)),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
