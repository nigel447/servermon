package main

import (
	"context"
	"fmt"
	"fyne-client/client"
	"fyne-client/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"os"
	"os/signal"
	"syscall"
)

var (
	bus = make(chan [2]string)

	screenDims [2]int
)

func main() {
	InitConfig()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	client := client.Client{Addr: config.Address}
	fmt.Printf("calling bootstap on %s\n with ctrl^C shutdown\n", client.Addr)
	go func() {
		// recieve signal
		<-c
		fmt.Println("\nsend shut down via context cancel")
		cancel()

	}()

	a := app.New()
	w := a.NewWindow("Server Mon")
	go func() {
		gui.ScreenDims()
	}()

	go func(chan [2]string) {
		for {
			select {
			case guiEvent := <-bus:
				fmt.Printf("login event usr:%s pswd:%s\n", guiEvent[0], guiEvent[1])
				w.Resize(fyne.NewSize(float32(screenDims[0]), float32(screenDims[1])))
				gui.MainScreen(w, screenDims)
			case bootEvent := <-gui.Boot:
				screenDims = bootEvent
				fmt.Println(screenDims)
				gui.LoginForm(w, screenDims, bus)
			}
		}
	}(bus)

	go func() {
		client.BootStrapClient(ctx)
	}()

	w.ShowAndRun()
}
