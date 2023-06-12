package ApiGrcp

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/In/ApiGrcp/Controller"
	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/In/ApiGrcp/gen/call"
	"github.com/Enrikerf/pfm/commandExecutor/app/Application/Port/In/ManageEngine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiGrpc struct {
	serverHost          string
	serverPort          string
	grpcServer          *grpc.Server
	listener            net.Listener
	manageEngineUseCase ManageEngine.UseCase
}

func (api *ApiGrpc) Initialize(host string, port string, manageEngineUseCase ManageEngine.UseCase) {
	fmt.Println("Starting Sentence Executor...")
	api.serverHost = host
	api.serverPort = port
	api.manageEngineUseCase = manageEngineUseCase
	api.loadServer()
	api.configControllers(manageEngineUseCase)
	api.loadListener()
}

func (api *ApiGrpc) Run() {
	if os.Getenv("APP_DEBUG") == "true" {
		reflection.Register(api.grpcServer)
	}
	go func() {
		fmt.Println("Starting at: " + api.serverHost + api.serverPort)
		if err := api.grpcServer.Serve(api.listener); err != nil {
			log.Fatalf("fatal")
		}
	}()
}

func (api *ApiGrpc) Stop() {
	// Wait for control C to exit
	// channel := make(chan os.Signal, 1)
	// signal.Notify(channel, os.Interrupt)
	// Bock until a signal is received
	// <-channel
	fmt.Println("Stopping engine")
	api.manageEngineUseCase.TearDown()
	fmt.Println("Stopping apiGrpc")
	api.grpcServer.Stop()
	fmt.Println("closing the Listener")
	err := api.listener.Close()
	if err != nil {
		return
	}
	fmt.Println("End of program")
}

func (api *ApiGrpc) configControllers(manageEngineUseCase ManageEngine.UseCase) {

	var callController = Controller.CallController{
		ManageEngineUseCase: manageEngineUseCase,
	}
	call.RegisterCallServiceServer(api.grpcServer, &callController)
}

func (api *ApiGrpc) loadServer() {
	var serverOptions []grpc.ServerOption
	api.grpcServer = grpc.NewServer(serverOptions...)
	if os.Getenv("APP_DEBUG") == "true" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func (api *ApiGrpc) loadListener() {
	listener, err := net.Listen("tcp", api.serverHost+api.serverPort)
	if err != nil {
		log.Fatalf("failed to listen at: " + api.serverHost + api.serverPort)
	}
	api.listener = listener
}
