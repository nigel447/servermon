package main

import (
	"context"
	"fmt"
	sshCli "fyne-client/client"
	"fyne-client/gui"
	"fyne-client/icon"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"os"
	"os/signal"
	"syscall"
)

var (
	bus        = make(chan bool)
	screenDims [2]int
)

func main() {
	InitConfig()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	client := sshCli.Client{Addr: config.Address}
	fmt.Printf("calling bootstap on %s\n with ctrl^C shutdown\n", client.Addr)
	go func() {
		// recieve signal
		<-c
		fmt.Println("\nsend shut down via context cancel")
		cancel()

	}()

	a := app.New() // icon
	a.SetIcon(icon.Metericon)

	w := a.NewWindow("Server Mon")

	go func(chan bool) {
		for {
			select {
			case boot := <-bus:
				if boot {
					w.Resize(fyne.NewSize(float32(screenDims[0]), float32(screenDims[1])))
					gui.MainScreen(w, screenDims)
				} else {
					cancel()
				}
			case bootEvent := <-gui.Boot:
				fmt.Println("bootEvent  ", bootEvent)
				screenDims = bootEvent
				gui.BootDialog(w, screenDims, bus, true)
			case <-sshCli.Start:
				fmt.Println("case sshCli.Start")
				go func() { gui.ScreenDims() }()

			}
		}
	}(bus)

	go func() {
		err := client.BootStrapClient(ctx)
		if err != nil {
			fmt.Println("BootStrapClient err != nil ", err)
			var a [2]int
			a[0] = 800
			a[1] = 600
			gui.BootDialog(w, a, bus, false)
			// run login screen as dialog message serveroffline
			//gui.BootDialog(w, screenDims, bus, false)
		}
	}()

	w.ShowAndRun()
}
