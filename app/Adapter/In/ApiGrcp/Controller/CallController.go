package Controller

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/In/ApiGrcp/gen/call"
	"github.com/Enrikerf/pfm/commandExecutor/app/Application/Port/In/ManageEngine"
)

type CallController struct {
	ManageEngineUseCase ManageEngine.UseCase
	call.UnimplementedCallServiceServer
	currentCommand string
}

// CallUnary TODO: don't break on error never. return responses to allows keep going the system
func (s *CallController) CallUnary(ctx context.Context, request *call.CallRequest) (*call.CallResponse, error) {
	fmt.Println("executing command: " + request.Command)
	resultContent := execCommand(request.GetCommand())
	return &call.CallResponse{Result: resultContent}, nil
}

func (s *CallController) CallServerStream(request *call.CallRequest, server call.CallService_CallServerStreamServer) error {
	fmt.Printf("Server stream %v\n", request)
	//TODO: sense?
	for i := 0; i < 3; i++ {
		resultContent := execCommand(request.GetCommand())
		err := server.Send(&call.CallResponse{Result: resultContent})
		if err != nil {
			log.Fatalf("error")
			return nil
		}
	}
	return nil
}

func (s *CallController) CallClientStream(server call.CallService_CallClientStreamServer) error {
	fmt.Printf("average invoked")
	var result string
	for {
		request, err := server.Recv()
		if err == io.EOF {
			return server.SendAndClose(&call.CallResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error")
		}
		result += execCommand(request.GetCommand())
	}
}

func (s *CallController) CallBidirectional(server call.CallService_CallBidirectionalServer) error {
	fmt.Println("Bidirectional")
	s.ManageEngineUseCase.Reset()
	go s.bidiRecv(server)
	for {
		var valueToSend float64 = 0
		if s.currentCommand == "makeLap" {
			valueToSend = float64(s.ManageEngineUseCase.GetPosition())
		}
		if s.currentCommand == "rpmControl" {
			valueToSend = s.ManageEngineUseCase.GetCurrentAngularSpeed()
		}
		sendError := server.Send(&call.CallResponse{
			Result: strconv.Itoa(int(valueToSend)),
		})
		if sendError != nil {
			fmt.Println("finished Bidirectional")
			s.ManageEngineUseCase.StopRpmControl()
			s.ManageEngineUseCase.Reset()
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (s *CallController) bidiRecv(server call.CallService_CallBidirectionalServer) {
	request, err := server.Recv()
	if err == io.EOF {
		return
	}
	if err != nil {
		return
	}
	log.Println(request.GetCommand())
	switch request.GetCommand() {
	case "makeLap":
		s.currentCommand = request.GetCommand()
		s.ManageEngineUseCase.Turnaround()
	case "rpmControl":
		s.currentCommand = request.GetCommand()
		s.ManageEngineUseCase.RpmControl(120)
	case "Stop":
	default:
	}
}

func execCommand(command string) string {
	var resultContent string
	parts := strings.Fields(command)
	head := parts[0]
	// Head at this point is "sudo"
	parts = parts[1:]
	cmd := exec.Command(head, parts...)
	stdout, err := cmd.Output()
	if err != nil {
		resultContent = err.Error()
		fmt.Println(err.Error())
	} else {
		resultContent = string(stdout)
		//TODO: print only in debug mode
		fmt.Println("result: " + string(stdout))
	}
	return resultContent
}
