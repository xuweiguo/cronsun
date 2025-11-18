// node 服务
// 用于在所需要执行 cron 任务的机器启动服务，替代 cron 执行所需的任务
package main

import (
	"flag"
	slog "log"

	"cronsun/internal/app/bootstrap"
)

var (
	level    = flag.Int("l", 0, "log level, -1:debug, 0:info, 1:warn, 2:error")
	confFile = flag.String("conf", "./conf/base.json", "config file path")
)

func main() {
	flag.Parse()

	if err := bootstrap.RunNode(*confFile, *level); err != nil {
		slog.Fatalln(err)
	}
}
