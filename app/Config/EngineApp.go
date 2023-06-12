package Config

import (
	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/In/Console"
	"github.com/Enrikerf/pfm/commandExecutor/app/Adapter/Out/Periphio/Model"
	"github.com/Enrikerf/pfm/commandExecutor/app/Application/Port/In/ManageEngine"
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity"
)

type EngineApp interface {
	Run()
}
type engineApp struct {
	console Console.Console
}

func NewEngineApp() EngineApp {
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
	service := ManageEngine.Service{Engine: engine}
	console := Console.NewConsole(&service)
	return &engineApp{console: console}
}

func (app *engineApp) Run() {
	app.console.Run()
}
