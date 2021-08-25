package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"srvmon/server"
	"syscall"
)

func main() {
	InitConfig()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	srv := server.Server{Addr: config.Address}
	fmt.Println("calling bootstap on ", srv.Addr)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	srv.BootStrapServer(ctx)

}
