package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for telnet connection")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("Usage: go-telnet --timeout=10s host port or go-telnet host port")
	}

	host, port := flag.Arg(0), flag.Arg(1)

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := client.Connect(); err != nil {
		log.Println(err)
	}

	go func() {
		if err := client.Send(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()
}
