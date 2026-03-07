package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gonethernet/legacy-ghopertunnel/legacy"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player"
)

func main() {
	srv, err := legacy.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		for p := range srv.Accept() {
			go func(pl *player.Player) {
				time.Sleep(20 * time.Millisecond)
			}(p)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
