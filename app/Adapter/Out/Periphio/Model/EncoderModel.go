package Model

import (
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity"
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity/Pin"
	"sync"
)

type state int8

const (
	stateOO state = 0
	stateOI state = 1
	stateIO state = 2
	stateII state = 3
)

type risingEdgeDetected struct {
	pinAStatus bool
	pinBStatus bool
}

type EncoderModel struct {
	encoderPinA Pin.EncoderPin
	encoderPinB Pin.EncoderPin
	lastState   state
	position    int16
	forward     bool
	risingEdges []*risingEdgeDetected
	lock        *sync.Mutex
}

func NewEncoder(
	encoderPinA Pin.EncoderPin,
	encoderPinB Pin.EncoderPin,
) Entity.Encoder {
	v := []*risingEdgeDetected{}
	e := EncoderModel{
		encoderPinA: encoderPinA,
		encoderPinB: encoderPinB,
		lastState:   stateOO,
		position:    0,
		risingEdges: v,
		lock:        &sync.Mutex{},
	}
	return &e
}

func (encoder *EncoderModel) GetPosition() int16 {
	return encoder.position
}

func (encoder *EncoderModel) ResetPosition() {
	encoder.lastState = stateOO
	encoder.position = 0
	encoder.forward = true
}

func (encoder *EncoderModel) TearDown() {
	encoder.encoderPinA.TearDown()
	encoder.encoderPinB.TearDown()
}

func (encoder *EncoderModel) Watchdog() {
	encoder.ResetPosition()
	go encoder.waitA()
	go encoder.waitB()
	for {
		for len(encoder.risingEdges) > 0 {
			encoder.lock.Lock()
			risingEdge := encoder.risingEdges[0]
			encoder.risingEdges = encoder.risingEdges[1:]
			encoder.calculatePos(*risingEdge)
			encoder.lock.Unlock()
			//TODO: logg in dev environment
			//fmt.Printf("%v,%v,%v\n", encoder.position, risingEdge.pinAStatus, risingEdge.pinBStatus)
		}
	}
}

func (encoder *EncoderModel) calculatePos(t risingEdgeDetected) {
	if !t.pinAStatus && !t.pinBStatus {
		if encoder.lastState == stateIO {
			encoder.position--
			encoder.forward = false
		} else {
			encoder.position++
			encoder.forward = true
		}
		encoder.lastState = stateOO
	} else if !t.pinAStatus && t.pinBStatus {
		if encoder.lastState == stateOO {
			encoder.position--
			encoder.forward = false
		} else {
			encoder.position++
			encoder.forward = true
		}
		encoder.lastState = stateOI
	} else if t.pinAStatus && !t.pinBStatus {
		if encoder.lastState == stateII {
			encoder.position--
			encoder.forward = false
		} else {
			encoder.position++
			encoder.forward = true
		}
		encoder.lastState = stateIO
	} else if t.pinAStatus && t.pinBStatus {
		if encoder.lastState == stateOI {
			encoder.position--
			encoder.forward = false
		} else {
			encoder.position++
			encoder.forward = true
		}
		encoder.lastState = stateII
	} else {
		panic("encoder at impossible position")
	}
}

func (encoder *EncoderModel) waitA() {
	for {
		encoder.encoderPinA.WaitForEdge()
		encoder.lock.Lock()
		encoder.risingEdges = append([]*risingEdgeDetected{{
			encoder.encoderPinA.Read(),
			encoder.encoderPinB.Read(),
		}}, encoder.risingEdges...)
		encoder.lock.Unlock()
	}
}

func (encoder *EncoderModel) waitB() {
	for {
		encoder.encoderPinB.WaitForEdge()
		encoder.lock.Lock()
		encoder.risingEdges = append([]*risingEdgeDetected{{
			encoder.encoderPinA.Read(),
			encoder.encoderPinB.Read(),
		}}, encoder.risingEdges...)
		encoder.lock.Unlock()
	}
}
