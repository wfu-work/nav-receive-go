package services

import (
	"fmt"
	"nav-receive-go/domains"
	"nav-receive-go/global"
	"strings"
	"sync"
	"time"
)

type DeviceRtcmService struct {
	CrudService[domains.DeviceRtcm]
}

var DeviceRtcmServiceApp = new(DeviceRtcmService)

type deviceMapEntry struct {
	Sncode     string
	UpdatedAt  time.Time
	Persistent bool
}

var (
	mountDeviceMap      = make(map[string]deviceMapEntry)
	missingMountLogTime = make(map[string]time.Time)
	mu                  sync.RWMutex
	missingMountLogMu   sync.Mutex
)

const (
	runtimeDeviceMapTTL  = 24 * time.Hour
	missingMountLogTTL   = 30 * time.Minute
	maxMissingMountItems = 10000
)

func RegisterDevice(mount, sncode string) {
	registerDevice(mount, sncode, false)
}

func registerDevice(mount, sncode string, persistent bool) {
	mount = strings.TrimSpace(mount)
	sncode = strings.TrimSpace(sncode)
	if mount == "" || sncode == "" {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	mountDeviceMap[mount] = deviceMapEntry{
		Sncode:     sncode,
		UpdatedAt:  time.Now(),
		Persistent: persistent,
	}
}

func RegisterDeviceRtcm(deviceRtcm domains.DeviceRtcm) {
	registerDevice(deviceRtcm.Sncode, deviceRtcm.Sncode, true)
	registerDevice(deviceRtcm.MountPoint, deviceRtcm.Sncode, true)
}

func RefreshDeviceMappings(devices []domains.Device, deviceRtcms []domains.DeviceRtcm) {
	now := time.Now()
	mu.Lock()

	for mount, entry := range mountDeviceMap {
		if entry.Persistent || now.Sub(entry.UpdatedAt) > runtimeDeviceMapTTL {
			delete(mountDeviceMap, mount)
		}
	}
	for _, device := range devices {
		setDeviceMapEntry(device.Sncode, device.Sncode, now, true)
	}
	for _, deviceRtcm := range deviceRtcms {
		setDeviceMapEntry(deviceRtcm.Sncode, deviceRtcm.Sncode, now, true)
		setDeviceMapEntry(deviceRtcm.MountPoint, deviceRtcm.Sncode, now, true)
	}
	mu.Unlock()
	pruneMissingMountLogs(now)
}

func setDeviceMapEntry(mount, sncode string, updatedAt time.Time, persistent bool) {
	mount = strings.TrimSpace(mount)
	sncode = strings.TrimSpace(sncode)
	if mount == "" || sncode == "" {
		return
	}
	mountDeviceMap[mount] = deviceMapEntry{
		Sncode:     sncode,
		UpdatedAt:  updatedAt,
		Persistent: persistent,
	}
}

func GetDeviceSncode(mount string) (string, bool) {
	mount = strings.TrimSpace(mount)
	mu.RLock()
	defer mu.RUnlock()
	entry, ok := mountDeviceMap[mount]
	return entry.Sncode, ok
}

func pruneMissingMountLogs(now time.Time) {
	missingMountLogMu.Lock()
	defer missingMountLogMu.Unlock()
	pruneMissingMountLogsLocked(now)
}

func pruneMissingMountLogsLocked(now time.Time) {
	for mount, lastTime := range missingMountLogTime {
		if now.Sub(lastTime) > missingMountLogTTL {
			delete(missingMountLogTime, mount)
		}
	}
}

func (s DeviceRtcmService) Init() {
	var result []domains.DeviceRtcm
	err := global.NAV_DB.Where("deleted = ? ", 0).Where("status = ? ", 1).Find(&result).Error
	if err != nil {
		fmt.Printf("❌init device rtcm find db error: %v\n", err)
		return
	}
	go func() {
		for _, bean := range result {
			RegisterDeviceRtcm(bean)
		}
	}()
}

func (s DeviceRtcmService) GetBySncode(sncode string) (*domains.DeviceRtcm, error) {
	deviceRtcm, err := s.GetByField("sncode", sncode)
	return deviceRtcm, err
}
