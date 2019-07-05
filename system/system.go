package system

type Resource struct {
	Version string
	Uptime  float64
	CPULoad float64
}

type Health struct {
	Voltage          float64
	Current          float64
	Temperature      float64
	CPUTemperature   float64
	PowerConsumption float64
}
