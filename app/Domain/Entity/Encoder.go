package Entity


type Encoder interface {
	Watchdog()
	ResetPosition() 
	GetPosition() int16
	TearDown()
}

