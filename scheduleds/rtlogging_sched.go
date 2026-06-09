package scheduleds

import (
	"fmt"
	"nav-receive-go/global"
	"nav-receive-go/services"
)

// RtloggingSched
// @description: 设备数据状态状态
// @return: error
func RtloggingSched() error {
	if global.NAV_DB == nil {
		return fmt.Errorf("db is not init")
	}
	params := make(map[string]string)
	deviceList, err := services.DeviceServiceApp.ListAll(params)
	if err != nil {
		fmt.Printf("❌device service find db error: %v\n", err)
		return err
	}
	deviceRtcmList, err := services.DeviceRtcmServiceApp.ListAll(map[string]string{
		"status": "1",
	})
	if err != nil {
		fmt.Printf("❌device rtcm service find db error: %v\n", err)
		return err
	}
	services.RefreshDeviceMappings(deviceList, deviceRtcmList)
	return nil
}
