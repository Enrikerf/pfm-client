package Pin

type EncoderPin interface {
	Read()bool
	TearDown()
	WaitForEdge()
}
