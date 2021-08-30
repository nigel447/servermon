package server

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"srvmon/dcrypto"
	"sync/atomic"
)

func (srv *Server) BootStrapServer(ctx context.Context) {
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
		atomic.AddUint64(&mainServerLoopCount, 1)
		conn, e := l.Accept()
		if e != nil {
			select {
			case <-ctx.Done():
				{
					fmt.Println("cancel main server loop with signal shutdown")
					handleError(ctx.Err()) // error due to ctx cancelation
				}
			default:
				handleError(err)
			}
			handleError(e)
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, srv.config)
		handleError(err)
		fmt.Println("we have a ssh connecion ", sshConn.RemoteAddr().String())
		// fmt.Println("we have a base64 encoding for client's public  key ", srv.ClienPublicKey)
		go srv.handleRequests(reqs)

		for ch := range chans {
			//ctype := ch.ChannelType()
			channel, requests, err := ch.Accept()
			if err != nil {
				fmt.Println("could not accept channel ", err)
			}
			fmt.Printf("accept channel %s\n", ch.ChannelType())
			go ssh.DiscardRequests(requests)

			go func(channel ssh.Channel) {
				for {
					buff := make([]byte, 256)
					fmt.Println("server listening for incomming channel mssgs")
					n, err := channel.Read(buff)
					if err != nil {
						break
					}
					b := buff[:n]
					fmt.Println("incomming channel mssg ", string(b))
					if string(b) == "start" {
						// start profiling
						RunProfiler()
					}

					if string(b) == "stop" {
						// stop profiling
						StopProfiler()
					}
				}
			}(channel)

			// buff := make([]byte, 256)

			// _, err := channel.Read(buff)
			// if err != nil {
			// 	break
			// }
			// b := buff[:n]
			// fmt.Printf("[%s]\n%s", ctype, string(b))
			// ret := fmt.Sprintf("send rpc result for request [%s]\n", string(b))
			// channel.Write([]byte(ret))

			go func(channel ssh.Channel) {
				fmt.Println("server sending system profile data")
				for {
					data := <-ProfileDataCh
					fmt.Printf("sending data\nschedluer run count %d server loop count %d\n", sheduerCount, mainServerLoopCount)
					channel.Write([]byte(data))
				}
			}(channel)

		}
	}
}

func (srv *Server) getServerConfig() {
	signer := dcrypto.GenOpenSSHSigner()
	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			srv.ClienPublicKey = key
			return nil, nil
		}}

	config.AddHostKey(signer)
	srv.config = config
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
