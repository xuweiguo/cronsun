package main

import (
	"flag"
	slog "log"

	"cronsun/internal/app/bootstrap"
)

var (
	level    = flag.Int("l", 0, "log level, -1:debug, 0:info, 1:warn, 2:error")
	confFile = flag.String("conf", "./conf/base.json", "config file path")
	network  = flag.String("network", "0.0.0.0", "network protocol of listen address: ipv4/ipv6, or empty use both")
)

func main() {
	flag.Parse()

	if err := bootstrap.RunWeb(*confFile, *network, *level); err != nil {
		slog.Fatalln(err)
	}
}
