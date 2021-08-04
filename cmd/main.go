package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/lixmal/url-shortener/pkg/backend"
)

const (
	VERSION = "0.1.0"
)

const (
	GracefulTimeout = 20 * time.Second
	ListenAddr      = ":8080"
)

var flagVersion *bool

func init() {
	// cmdline arg "-version"
	flagVersion = flag.Bool("version", false, "print program version")
	flag.Parse()
}

func main() {
	if *flagVersion {
		version()
		os.Exit(0)
	}

	// TODO: Make the address configurable
	srv := backend.New(ListenAddr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	<-sig
	log.Println("Server shutdown ...")

	// TODO: add small delay if running on k8s behind a reverse proxy
	// Graceful server shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), GracefulTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
	defer cancel()
	log.Println("Finished")
}

func version() {
	fmt.Printf("%s version: %s\ngo version: %s %s/%s\n", os.Args[0], VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
