package main

import (
	"context"
	"climon/sshclient"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c,  os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	client := sshclient.Client{ Addr: "localhost:2222" }
	fmt.Println("calling bootstap on", client.Addr)
	go func(){
		// recieve signal
		<-c
		fmt.Println("send shut down via context cancel")
		cancel()

	}()
	client.BootStrapClient(ctx)

	done := make(chan bool, 1)

    go func() {
		time.AfterFunc(40* time.Second, func() {
			fmt.Println("exit after more than a 20s.")
			done <- true
		})
    }()

	fmt.Println("block on  channel like latch ")
	<-done
	fmt.Println("exiting")

}