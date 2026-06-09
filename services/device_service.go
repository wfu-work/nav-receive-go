package services

import (
	"nav-receive-go/domains"
)

type DeviceService struct {
	CrudService[domains.Device]
}

var DeviceServiceApp = new(DeviceService)

func (s DeviceService) GetBySncode(sncode string) (*domains.Device, error) {
	device, err := s.GetByField("sncode", sncode)
	return device, err
}

func (s DeviceService) GetByDeviceId(deviceId string) (*domains.Device, error) {
	device, err := s.GetByField("device_id", deviceId)
	return device, err
}

func (s DeviceService) GetBySn(sn string) []domains.Device {
	params := map[string]string{
		"sn": sn,
	}
	result, err := s.ListAll(params)
	if err != nil {
		return []domains.Device{}
	}
	return result
}
