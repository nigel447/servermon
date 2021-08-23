package server

import (
	"context"
	"fmt"
	//"net"
	"os"
	"golang.org/x/crypto/ssh"
    "srvmon/dcrypto"
)

func (srv *Server) BootStrapServer(ctx context.Context)   {
	srv.getServerConfig()

	l, err := lc.Listen(ctx, "tcp", srv.Addr)
	handleError(err)
	defer l.Close()
	fmt.Println("BootStrapServer Listen ", srv.Addr)

	go func() {
		<-ctx.Done()
		fmt.Println("... shutting service down")
		l.Close()
		os.Exit(1)
	}()

	for {
		conn, e := l.Accept()
		if e != nil {
			select {
			case <-ctx.Done(): {
				fmt.Println("cancel main server loop with signal shutdown" )
				handleError( ctx.Err()) // error due to ctx cancelation
			}
			default:
				handleError(err)
			}
			handleError(e)
		}

		fmt.Println("BootStrapServer we have a TCP connecion ", conn.RemoteAddr().String())

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, srv.config)
		if err != nil {
			fmt.Println("Failed to listen on 2222")

		}
		fmt.Println("we have a ssh connecion ", sshConn.RemoteAddr().String())
		fmt.Println("we have a client  public key ", srv.ClienPublicKey)
		go srv.handleRequests(reqs)
		for ch := range chans {
			ctype := ch.ChannelType()
 
			channel, requests, err := ch.Accept()
			if err != nil {
				fmt.Println("could not accept channel ", err)
			}
			go ssh.DiscardRequests(requests)

			buff := make([]byte, 256)
		 
			for {
				n, err := channel.Read(buff)
				if err != nil {
					break
				}
				b := buff[:n]
				fmt.Printf("[%s]\n%s", ctype, string(b))
			}
 
		}
	}

}

func (srv *Server) getServerConfig()  { 
	privateKey, _, err:= dcrypto.GenECDSAKey()
	handleError(err)
	signer, err := ssh.NewSignerFromKey(privateKey)
	handleError(err)

	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			fmt.Println("we have a client presents public key ", key)
			srv.ClienPublicKey = key
			return nil, nil
		}}
 
	config.AddHostKey(signer)
	srv.config =config

}

func (srv *Server) handleRequests(in <-chan *ssh.Request) {
	for req := range in {
		fmt.Println("handleRequests ", req.Type)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}