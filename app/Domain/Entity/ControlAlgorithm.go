package Entity

type ControlAlgorithm interface {
	SetGoal(goal float64)
	SetP(p float64)
	SetI(i float64)
	SetD(d float64)
	SetSampleTime(d float64)
	GetIntegralTerm() float64
	Calculate(currentValue float64) float64
}

type controlAlgorithm struct {
	goal         float64
	P            float64
	I            float64
	D            float64
	integralTerm float64
	sampleTime   float64
	currentValue float64
	currentError float64
	pastError    float64
}

func NewControlAlgorithm() ControlAlgorithm {
	return &controlAlgorithm{}
}

func (ca *controlAlgorithm) SetGoal(goal float64) {
	ca.goal = goal
}

func (ca *controlAlgorithm) SetP(p float64) {
	ca.P = p
}

func (ca *controlAlgorithm) SetI(i float64) {
	ca.I = i
}

func (ca *controlAlgorithm) SetD(d float64) {
	ca.D = d
}

func (ca *controlAlgorithm) GetIntegralTerm() float64 {
	return ca.integralTerm
}

func (ca *controlAlgorithm) SetSampleTime(st float64) {
	ca.sampleTime = st
}

func (ca *controlAlgorithm) Calculate(currentValue float64) float64 {
	ca.currentValue = currentValue
	ca.currentError = ca.goal - ca.currentValue
	proportionalTerm := ca.P * ca.currentError
	ca.integralTerm = ca.integralTerm + ca.currentError*ca.sampleTime
	derivativeTerm := ca.D * (ca.currentError - ca.pastError) / ca.sampleTime
	ca.pastError = ca.currentError
	return proportionalTerm + ca.I*ca.integralTerm + derivativeTerm

}
