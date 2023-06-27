package main

import (
	"github.com/brinaviv/ss-go-web-server/pkg/api"
	"github.com/brinaviv/ss-go-web-server/pkg/dal/inmemory"
	"github.com/brinaviv/ss-go-web-server/pkg/model"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ConfigureLogger()
	serverConfig := &api.ServerConfig{Host: "localhost", Port: "8080", Controllers: initControllers()} // TODO: read config from file / env

	panicIfError(api.StartWebServer(serverConfig))

	waitForTermination()
}

func initControllers() []api.Controller {
	dao := inmemory.NewDAO(func() model.UserID {
		return model.UserID(uuid.New())
	})
	controllers := []api.Controller{&api.UsersController{UserDAO: dao.Users}}
	return controllers
}

func waitForTermination() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sig
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
