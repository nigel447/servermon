package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"srvmon/server"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c,  os.Interrupt, syscall.SIGTERM)
	srv := server.Server{ Addr: ":2222" }
	fmt.Println("calling bootstap on", srv.Addr)
	ctx, cancel := context.WithCancel(context.Background())

	go func(){
		<-c
		cancel()
	}()

	srv.BootStrapServer(ctx)

}
