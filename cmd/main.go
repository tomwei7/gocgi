package main

import (
	"context"
	"os"
	"log"
	"flag"

	"github.com/tomwei7/gocgi"
)

func ParseCommandLine() *gocgi.GoCGIOptions {
	opts := gocgi.NewGoCGIOptions()
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	return opts
}

func main() {
	opts := ParseCommandLine()
	logger := log.New(os.Stderr, "gocgi ", log.LstdFlags)

	server, err := gocgi.New(logger, opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer server.Shutdown(context.Background())

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
