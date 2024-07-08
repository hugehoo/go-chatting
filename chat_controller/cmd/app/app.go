package app

import (
	"chat_controller/config"
	"chat_controller/network"
	"chat_controller/repository"
	"chat_controller/service"
)

type App struct {
	config     *config.Config
	service    *service.Service
	network    *network.Server
	repository *repository.Repository
}

func NewApp(config *config.Config) *App {
	a := &App{config: config}
	var err error
	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.repository)
		a.network = network.NewNetwork(a.service, config.Info.Port)
	}
	return a
}

func (app *App) Start() {
	app.network.Start()
}
