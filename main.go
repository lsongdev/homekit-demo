package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

func main() {
	// dnssdlog.Debug.Enable()

	a := accessory.NewSwitch(accessory.Info{
		Name:         "Light",
		Manufacturer: "Brutella",
		Model:        "Switch",
		SerialNumber: "123456789",
		Firmware:     "1.0.0",
	})

	s, err := hap.NewServer(hap.NewFsStore("./db"), a.A)
	if err != nil {
		log.Fatal(err)
	}

	// Log to console when client (e.g. iOS app) changes the value of the on
	// characteristic.
	a.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			log.Println("Client changed switch to on")
		} else {
			log.Println("Client changed switch to off")
		}
	})

	// Periodically toggle the switch's on characteristic
	// go func() {
	// 	for {
	// 		on := !a.Switch.On.Value()
	// 		if on == true {
	// 			log.Println("Switch is on")
	// 		} else {
	// 			log.Println("Switch is off")
	// 		}
	// 		a.Switch.On.SetValue(on)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	if err := s.ListenAndServe(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
