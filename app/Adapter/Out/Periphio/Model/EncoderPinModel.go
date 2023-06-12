package Model

import (
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity/Pin"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type EncoderPinModel struct {
	periphInPin gpio.PinIO
	id          string
}

func NewEncoderPin(id string) Pin.EncoderPin {
	periphInPin := gpioreg.ByName(id)
	periphInPin.In(gpio.PullUp, gpio.BothEdges)
	return &EncoderPinModel{
		periphInPin: periphInPin,
		id:          id,
	}
}

func (encoderPin *EncoderPinModel) Read() bool {
	return bool(encoderPin.periphInPin.Read())
}

func (encoderPin *EncoderPinModel) TearDown() {
	encoderPin.periphInPin.DefaultPull()
	encoderPin.periphInPin.Out(gpio.Low)
}

func (encoderPin *EncoderPinModel) WaitForEdge() {
	encoderPin.periphInPin.WaitForEdge(-1)
}
