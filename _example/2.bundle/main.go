package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/_example/2.bundle/httpserver"
)

// StartServer
func StartServer(server *http.Server) error {
	fmt.Println("Addr: ", server.Addr)
	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		fmt.Println(err)
		done <- struct{}{}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
		if err := server.Close(); err != nil {
			return err
		}

		<-done
	case <-done:
	}

	return nil
}

func main() {
	slice.Run(
		slice.SetName("bundle-example"),
		slice.SetDispatcher(StartServer),
		slice.RegisterBundles(
			&httpserver.Bundle{},
		),
	)
}
