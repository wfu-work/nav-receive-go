package services

import (
	"log"
	"nav-receive-go/domains"
	"strings"
	"sync"
	"time"
)

var (
	deviceSncodeCache      = make(map[string]struct{})
	deviceSncodeCacheMu    sync.RWMutex
	unknownDeviceLogTime   = make(map[string]time.Time)
	unknownDeviceLogTimeMu sync.Mutex
)

const (
	unknownDeviceLogInterval = time.Minute
	unknownDeviceLogTTL      = 30 * time.Minute
	maxUnknownDeviceLogItems = 10000
)

func RefreshDeviceCache(devices []domains.Device) {
	nextCache := make(map[string]struct{}, len(devices))
	for _, device := range devices {
		sncode := strings.TrimSpace(device.Sncode)
		if sncode == "" {
			continue
		}
		nextCache[sncode] = struct{}{}
	}

	deviceSncodeCacheMu.Lock()
	deviceSncodeCache = nextCache
	deviceSncodeCacheMu.Unlock()

	pruneUnknownDeviceLogs(time.Now())
}

func IsDeviceSncodeCached(sncode string) bool {
	sncode = strings.TrimSpace(sncode)
	if sncode == "" {
		return false
	}
	deviceSncodeCacheMu.RLock()
	defer deviceSncodeCacheMu.RUnlock()
	_, ok := deviceSncodeCache[sncode]
	return ok
}

func LogUnknownDevice(sncode string) {
	if shouldLogUnknownDevice(sncode) {
		log.Printf("rtlogging设备不存在，跳过写入: sncode=%s", strings.TrimSpace(sncode))
	}
}

func shouldLogUnknownDevice(sncode string) bool {
	sncode = strings.TrimSpace(sncode)
	if sncode == "" {
		return false
	}
	now := time.Now()
	unknownDeviceLogTimeMu.Lock()
	defer unknownDeviceLogTimeMu.Unlock()
	if len(unknownDeviceLogTime) > maxUnknownDeviceLogItems {
		pruneUnknownDeviceLogsLocked(now)
	}
	lastTime, ok := unknownDeviceLogTime[sncode]
	if ok && now.Sub(lastTime) < unknownDeviceLogInterval {
		return false
	}
	unknownDeviceLogTime[sncode] = now
	return true
}

func pruneUnknownDeviceLogs(now time.Time) {
	unknownDeviceLogTimeMu.Lock()
	defer unknownDeviceLogTimeMu.Unlock()
	pruneUnknownDeviceLogsLocked(now)
}

func pruneUnknownDeviceLogsLocked(now time.Time) {
	for sncode, lastTime := range unknownDeviceLogTime {
		if now.Sub(lastTime) > unknownDeviceLogTTL {
			delete(unknownDeviceLogTime, sncode)
		}
	}
}
