package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/santonov10/otus_hw_project/internal/config"
	"github.com/santonov10/otus_hw_project/internal/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/default.json", "Path to configuration file")
}

func main() {
	config.SetFilePath(configFile)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	serv := server.NewServer()
	serv.Start(ctx)

	<-ctx.Done()

	fmt.Println("закрываем приложение")
}
