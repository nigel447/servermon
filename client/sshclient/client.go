package sshclient

import (
	"context"
    "climon/dcrypto"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

func (cli *Client) BootStrapClient(ctx context.Context)   {
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
		fmt.Println("BootStrapClient OpenChannel")
		channel, requests, err := conn.OpenChannel("rpc-remote", s("%s extra data", "init string"))
		handleError(err)
		go ssh.DiscardRequests(requests)

	//send data forever...
	n := 1
	for {
		_, err := channel.Write(s("#%d send data channel", n))
		handleError(err)
		n++
		time.Sleep(3 * time.Second)
		if n >3 {
			break
		}
	}
	}()
	// block
	<-ctx.Done()
}

func verifyServer(hostname string, remote net.Addr, key ssh.PublicKey) error {
	got := FingerprintKey(key)
	fmt.Println("Fingerprint ", got)
	fmt.Println("hostname ", hostname)
	fmt.Println("PublicKey ", key.Type())
	return nil
}

func (cli *Client) getClientConfig()  { 
	user := "climon"
	privateKey, _, err:= dcrypto.GenECDSAKey()
	handleError(err)
	signer, err := ssh.NewSignerFromKey(privateKey)
 
	handleError(err)

	cli.config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: verifyServer,
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}