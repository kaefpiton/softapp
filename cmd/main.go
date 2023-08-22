package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"softapi/configs"
	server "softapi/internal/server"
	"syscall"
)

const configPath = "configs/config.json"

func main() {
	cnf, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	server := server.NewEchoHTTPServer(cnf.HttpServer.Port)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("Terminating the app")
		cancel()

		fmt.Println("Stop Server")
		server.Stop(ctx)
	}()

	server.Start()
}
