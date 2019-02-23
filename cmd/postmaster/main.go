package main

import (
	"flag"

	"github.com/openware/postmaster/pkg/consumer"
)

func main() {
	config := flag.String(
		"config",
		"config/postmaster.yml",
		"Path to postmaster config file",
	)
	flag.Parse()

	consumer.Run(*config)
}
