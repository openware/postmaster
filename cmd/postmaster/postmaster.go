package main

import (
	"flag"

	"github.com/openware/postmaster/pkg/consumer"
)

const (
	DefaultPath = "config/postmaster.yml"
	DefaultTag  = "postmaster"
)

var (
	config = flag.String("config", DefaultPath, "Path to postmaster config file")
	tag    = flag.String("tag", DefaultTag, "Tag for RabbitMQ consumer")
)

func main() {
	flag.Parse()

	consumer.Run(*config, *tag)
}
