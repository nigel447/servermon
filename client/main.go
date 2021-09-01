package main

import (
	"climon/sshclient"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	InitConfig()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	client := sshclient.Client{Addr: config.Address}
	fmt.Printf("calling bootstap on %s\n with ctrl^C shutdown\n", client.Addr)
	go func() {
		// recieve signal
		<-c
		fmt.Println("\nsend shut down via context cancel")
		cancel()

	}()
	go client.BootStrapClient(ctx)

}
