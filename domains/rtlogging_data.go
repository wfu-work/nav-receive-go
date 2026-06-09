package domains

import "time"

// RtloggingData 泛型结构体
type RtloggingData[T any] struct {
	Data     T      `json:"data"`
	Device   string `json:"device"`
	DeviceId string `json:"deviceId"`
	Time     string `json:"time"`
	Service  string `json:"service"`
}

// NewRtloggingData 构造函数
func NewRtloggingData[T any](device, deviceId string, data T, service string) *RtloggingData[T] {
	return &RtloggingData[T]{
		Data:     data,
		Device:   device,
		DeviceId: deviceId,
		Time:     time.Now().UTC().Format("2006-01-02T15:04:05"),
		Service:  service,
	}
}

// SensorModel 泛型结构体
type SensorModel[T any] struct {
	Sncode string `json:"sncode"`
	Time   int64  `json:"time"`
}
