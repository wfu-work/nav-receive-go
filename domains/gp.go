package domains

type GP struct {
	SensorModel[GP]
	GpsInitial string  `json:"gpsInitial"`
	GpsTotalX  float64 `json:"gpsTotalX"`
	GpsTotalY  float64 `json:"gpsTotalY"`
	GpsTotalZ  float64 `json:"gpsTotalZ"`
}
