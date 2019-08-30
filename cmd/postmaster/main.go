package main

import (
	"flag"

	"github.com/openware/postmaster/pkg/consumer"
)

var (
	config = flag.String("config", "config/postmaster.yml", "Path to postmaster config file")
	tag    = flag.String("tag", "postmaster", "RabbitMQ consumer unique tag")
)

func main() {
	flag.Parse()

	consumer.Run(*config, *tag)
}
