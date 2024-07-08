package main

import (
	"chat_controller/cmd/app"
	"chat_controller/config"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "configuration setting")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)
	a := app.NewApp(c)
	a.Start()
}
