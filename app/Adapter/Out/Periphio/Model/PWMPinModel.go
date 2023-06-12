package Model

import (
	"fmt"
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity/Pin"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
)

const (
	minDuty = 1258292 // PWM don't move the EMG30 below this duty
)

type pwmPin struct {
	periphPin gpio.PinIO
	id        string
	status    bool
}

func NewPwmPin(id string) Pin.PWMPin {
	_, err := host.Init()
	if err != nil {
		panic("Periph.io Adapter error")
	}
	e := pwmPin{id: id, status: false, periphPin: gpioreg.ByName(id)}
	return &e
}

func (outPin *pwmPin) TearDown() {
	outPin.periphPin.Halt()
	outPin.periphPin.DefaultPull()
	outPin.periphPin.Out(gpio.Low)
}

func (outPin *pwmPin) GetMaxDuty() Pin.Duty {
	return Pin.Duty(gpio.DutyMax)
}

func (outPin *pwmPin) GetMinDuty() Pin.Duty {
	return Pin.Duty(minDuty)
}

func (outPin *pwmPin) GetMaxFrequency() Pin.Frequency {
	return Pin.Frequency(10 * physic.KiloHertz)
}

func (pwmPin *pwmPin) SetPWM(duty Pin.Duty, frequency Pin.Frequency) {
	if err := pwmPin.periphPin.PWM(gpio.Duty(duty), physic.Frequency(frequency)); err != nil {
		fmt.Println(err)
	}
}

func (pwmPin *pwmPin) StopPWM() {
	if err := pwmPin.periphPin.PWM(0, 100*physic.KiloHertz); err != nil {
		panic("%d: pwm panic")
	}
}
