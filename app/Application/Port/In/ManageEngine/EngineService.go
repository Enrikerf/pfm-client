package ManageEngine

import "github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity"

type Service struct {
	Engine Entity.Engine
}

func (service *Service) GetPosition() int16 {
	return service.Engine.GetPosition()
}

func (service *Service) Turnaround() {
	service.Engine.MakeLap()
}

func (service *Service) RpmControl(rpm float64) {
	service.Engine.RpmControl(rpm)
}

func (service *Service) StopRpmControl() {
	service.Engine.StopRmpControl()
}

func (service *Service) SetGas(gas int) {
	service.Engine.SetGas(Entity.GasLevel(gas))
}

func (service *Service) StepResponse() {
	service.Engine.StepResponse()
}

func (service *Service) UnBrake() {
	service.Engine.UnBrake()
}

func (service *Service) TearDown() {
	service.Engine.TearDown()
}

func (service *Service) Reset() {
	service.Engine.InitialState()
}


func (service *Service)GetCurrentAngularSpeed()float64{
	return service.Engine.GetCurrentAngularSpeed()
}