package app

import (
	"chat_controller/config"
	"chat_controller/network"
	"chat_controller/service"
)

type App struct {
	config  *config.Config
	service *service.Service
	network *network.Server
}

func NewApp(config *config.Config) *App {
	a := &App{config: config}
	// netWork 에 대해 세팅해줘야함.

	a.network = network.NewNetwork(a.service, config.Info.Port)
	return a
}

func (app *App) Start() {
	app.network.Start()
}
