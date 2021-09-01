package client

import (
	"context"
	"errors"
	"fmt"
	"fyne-client/dcrypto"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"sync/atomic"
	"time"

	// ctrl channels
	"fyne-client/gui"
)

var dataCount uint64

// StartStop = make(chan [1]int)
// DataPipe  = make(chan [1]string)

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
	channel, requests, err := conn.OpenChannel("rpc-remote", s("%s extra data", "init string"))
	// handle start stop with channel
	go func(channel ssh.Channel) {
		go ssh.DiscardRequests(requests)
		for {
			// should blocak and wait
			startDataPipe := <-gui.StartStop
			if startDataPipe == "start" {
				fmt.Println("BootStrapClient opening data channel")
				_, err := channel.Write(s("start"))
				handleError(err)
			} else {
				fmt.Println("BootStrapClient closing data channel")
				_, err = channel.Write(s("stop"))
				handleError(err)
			}
		}
	}(channel)

	// pipe data to gui with channel
	go func(channel ssh.Channel) {
		buff := make([]byte, 256)
		for {
			n, err := channel.Read(buff)
			if err != nil {
				break
			}
			b := buff[:n]
			fmt.Println("data coming in")
			fmt.Println(string(b))
			// should blok and wait for data
			gui.DataPipe <- b
			fmt.Println("recieved data count ", dataCount)
			atomic.AddUint64(&dataCount, 1)
			time.Sleep(200 * time.Millisecond)

		}
	}(channel)
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
