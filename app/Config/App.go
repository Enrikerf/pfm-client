package Config

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/In/ApiGrcp"
	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/Out/Periphio/Model"
	"github.com/Enrikerf/pfm/commandExecutor/app/Application/Port/In/ManageEngine"
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity"
	"github.com/joho/godotenv"
)

type App struct {
	apiGrpc ApiGrcp.ApiGrpc
}

func (server *App) Run() {
	server.loadDotEnv()
	server.loadApiGrpc()

	// Wait for control C to exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	// Bock until a signal is received
	<-channel
	server.apiGrpc.Stop()
}

func (server *App) loadDotEnv() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
}

func (server *App) loadApiGrpc() {
	pwmPin := Model.NewPwmPin("18")
	brakePin := Model.NewOutPin("12")
	dirPin := Model.NewOutPin("7")
	encoderPinA := Model.NewEncoderPin("25")
	encoderPinB := Model.NewEncoderPin("8")
	encoder := Model.NewEncoder(encoderPinA, encoderPinB)
	control := Entity.NewControlAlgorithm()
	engine := Entity.NewEngine(
		360,
		control,
		brakePin,
		dirPin,
		pwmPin,
		encoder,
	)
	engineService := ManageEngine.Service{Engine: engine}
	server.apiGrpc = ApiGrcp.ApiGrpc{}
	server.apiGrpc.Initialize(
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_PORT"),
		&engineService,
	)
	go server.apiGrpc.Run()

}
