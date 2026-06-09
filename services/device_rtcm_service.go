package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"nav-receive-go/domains"
	"nav-receive-go/global"
	"strconv"
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

func ShouldLogMissingMount(mount string) bool {
	mount = strings.TrimSpace(mount)
	if mount == "" {
		return false
	}
	now := time.Now()
	missingMountLogMu.Lock()
	defer missingMountLogMu.Unlock()
	if len(missingMountLogTime) > maxMissingMountItems {
		pruneMissingMountLogsLocked(now)
	}
	lastTime, ok := missingMountLogTime[mount]
	if ok && now.Sub(lastTime) < time.Minute {
		return false
	}
	missingMountLogTime[mount] = now
	return true
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

func AuthCasterServer(mount string, password string) bool {
	var result domains.DeviceRtcm
	err := global.NAV_DB.Where("deleted = ? ", 0).Where("agreement = ? and mount_point = ? and password = ?", "NTRIP Server", mount, password).First(&result).Error
	if err != nil {
		return false
	}
	RegisterDeviceRtcm(result)
	return true
}

func AuthCasterClient(mount string, username string, password string) bool {
	var result domains.DeviceRtcm
	err := global.NAV_DB.Where("deleted = ? ", 0).Where("agreement = ? and mount_point = ? and username = ? and password = ?", "NTRIP Client", mount, username, password).First(&result).Error
	if err != nil {
		return false
	}
	RegisterDeviceRtcm(result)
	return true
}

func AuthZx(mountPoint, password string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("展讯认证异常: %v", r)
		}
	}()

	splits := strings.Split(password, ":")
	if len(splits) != 2 {
		return false
	}

	sncodeBytes, err := base64.StdEncoding.DecodeString(splits[0])
	if err != nil {
		log.Printf("展讯认证异常: %v", err)
		return false
	}
	sncode := string(sncodeBytes)

	timeBytes, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		log.Printf("展讯认证异常: %v", err)
		return false
	}
	timeStr := string(timeBytes)

	if len(timeStr) <= 6 {
		log.Printf("展讯认证异常: 时间格式不正确 %s", timeStr)
		return false
	}
	substring := timeStr[:len(timeStr)-6]
	ts, err := strconv.ParseInt(substring, 10, 64)
	if err != nil {
		log.Printf("展讯认证异常(解析时间戳失败): %v", err)
		return false
	}
	if ts < 1e12 {
		ts = ts * 1000
	} else if ts > 1e15 {
		ts = ts / 1000
	}
	offTime := time.Now().UnixMilli() - ts
	if mountPoint == sncode && offTime < 60*60*1000 {
		RegisterDevice(mountPoint, sncode)
		return true
	}
	log.Printf("是否展讯认证失败: mountPoint=%s, sncode=%s, timeStr=%s", mountPoint, sncode, timeStr)
	return false
}
