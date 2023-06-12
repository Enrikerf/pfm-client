package ManageEngine

type UseCase interface {
	Reset()
	GetPosition() int16
	Turnaround()
	RpmControl(rpm float64)
	StopRpmControl()
	UnBrake()
	SetGas(gas int)
	StepResponse()
	TearDown()
	GetCurrentAngularSpeed()float64
}
