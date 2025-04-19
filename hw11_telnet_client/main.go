package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	ErrRequiredHostAndPort = errors.New("host and port is required")
	ErrPortNotInt          = errors.New("port must be a integer")
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout")
	flag.Parse()

	err := run(flag.Args(), timeout)
	if err != nil {
		printError("Error:", err)
		os.Exit(1)
	}
}

func printError(prefix string, err error) {
	fmt.Fprintf(os.Stderr, "%v %v\n", prefix, err)
}

func run(args []string, timeout time.Duration) error {
	if len(args) != 2 {
		return ErrRequiredHostAndPort
	}
	port, err := strconv.Atoi(args[1])
	if err != nil || port <= 0 || port > 65535 {
		return ErrPortNotInt
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Println("...Connected to ", client.GetAddress())

	done := make(chan struct{})

	go func() {
		defer close(done)
		err := client.Send()
		if err != nil && !errors.Is(err, io.EOF) {
			printError("Send error:", err)
		}
	}()

	go func() {
		defer stop()
		err := client.Receive()
		if err != nil && !errors.Is(err, io.EOF) {
			printError("Receive error:", err)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("...Connection was closed by peer")
	case <-done:
		fmt.Println("...EOF")
	}

	return nil
}
