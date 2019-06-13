package main

import (
	"flag"

	"github.com/openware/postmaster/pkg/consumer"
)

const (
	path = "config/postmaster.yml"
)

var (
	config = flag.String("config", path, "Path to postmaster config file")
)

func main() {
	flag.Parse()

	consumer.Run(*config)
}
