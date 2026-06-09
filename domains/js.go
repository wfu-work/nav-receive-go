package domains

type JS struct {
	SensorModel[JS]
	GX float64 `json:"gX"`
	GY float64 `json:"gY"`
	GZ float64 `json:"gZ"`
}
