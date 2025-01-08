package main

import (
	"flag"
	"log"
	"strings"
)

type config struct {
	debug  bool
	hosts  []string
	listen string
}

var c config

// Read our command line args and fire of the workers
func main() {
	var hosts string
	flag.BoolVar(&c.debug, "debug", false, "enable debug mode")
	flag.StringVar(&hosts, "hosts", "", "comma separated list of hosts to connect to")
	flag.StringVar(&c.listen, "listen", "localhost:8888", "address:port to listen on")
	flag.Parse()
	if hosts == "" {
		log.Fatal("Hosts list must not be empty. Use -h for help.")
	}
	c.hosts = strings.Split(hosts, ",")
	if len(c.hosts) == 0 {
		log.Fatal("Hosts list must not be empty. Use -h for help.")
	}
	cchan := startCarChan()
	for x := range c.hosts {
		go startStream(c.hosts[x], cchan)
	}
	httpserver()
}
