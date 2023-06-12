package Model

import (
	"fmt"
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity/Pin"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type outPin struct {
	id        string
	status    bool
	periphPin gpio.PinIO
}

func NewOutPin(id string) Pin.OutPin {
	_, err := host.Init()
	if err != nil {
		panic("Periph.io Adapter error")
	}
	e := outPin{id: id, status: false, periphPin: gpioreg.ByName(id)}
	return &e
}

func (outPin *outPin) TearDown() {
	outPin.periphPin.DefaultPull()
	outPin.periphPin.Out(gpio.Low)
}

func (outPin *outPin) Up() {
	outPin.status = true
	outPin.periphPin.Out(gpio.High)
	fmt.Printf("Pin %v: %v\n", outPin.id, outPin.status)
}

func (outPin *outPin) Down() {
	outPin.status = false
	outPin.periphPin.Out(gpio.Low)
	fmt.Printf("Pin %v: %v\n", outPin.id, outPin.status)
}
