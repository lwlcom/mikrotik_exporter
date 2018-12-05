package interfaces

type Interface struct {
	Name     string
	Comment  string
	RxByte   float64
	RxPacket float64
	RxDrop   float64
	RxError  float64

	TxByte   float64
	TxPacket float64
	TxDrop   float64
	TxError  float64
}
