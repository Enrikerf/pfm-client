package Pin

type Duty int32
type Frequency int64

type PWMPin interface {
	GetMaxDuty()Duty
	GetMinDuty()Duty
	GetMaxFrequency() Frequency
	SetPWM(duty Duty, frequency Frequency)
	StopPWM()
	TearDown()
}
