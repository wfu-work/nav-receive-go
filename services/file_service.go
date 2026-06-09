package services

import (
	"fmt"
	"log"
	"nav-receive-go/global"
	"nav-receive-go/utils"
	"nav-receive-go/vos"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SetRtloggingSize(sncode string, length int) {
	if strings.TrimSpace(sncode) == "" {
		return
	}
	key := utils.REDIS_RTLOGGING_PREFIX + sncode
	var fileChannelVo vos.FileChannelVo
	err := utils.RedisUtilApp.GetObj(key, &fileChannelVo)
	if err != nil {
		fileChannelVo = vos.FileChannelVo{
			Sncode:   sncode,
			LastTime: time.Now().UnixMilli(),
			Size:     length,
			Status:   1,
		}
		err = utils.RedisUtilApp.SetObj(key, fileChannelVo, utils.RedisNotExpire)
		if err != nil {
			log.Printf("rtlogging设置缓存异常：%s - %v\n", sncode, err)
		}
		return
	}
	fileChannelVo.LastTime = time.Now().UnixMilli()
	fileChannelVo.Size += length
	fileChannelVo.Sncode = sncode
	fileChannelVo.Status = 1
	err = utils.RedisUtilApp.SetObj(key, fileChannelVo, utils.RedisNotExpire)
	if err != nil {
		log.Printf("rtlogging设置缓存异常：%s - %v\n", sncode, err)
	}
}

func WriteRtcmFile(sncode string, utcTime time.Time, bytes []byte) {
	sncode = strings.TrimSpace(sncode)
	if sncode == "" {
		return
	}
	if !IsDeviceSncodeCached(sncode) {
		LogUnknownDevice(sncode)
		return
	}
	SetRtloggingSize(sncode, len(bytes))
	fileName := getRtcmFileName(sncode, utcTime)
	path := filepath.Join(
		global.NAV_CONFIG.Iot.RawPath,
		utils.LocalTimeUtilGetUtcTimeYearFileName(utcTime),
		utils.LocalTimeUtilGetDayOfYear(utcTime),
		utils.LocalTimeUtilGetHourOfDay(utcTime),
		fileName,
	)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Printf("创建目录失败: %v", err)
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开文件失败 %s: %v\n", path, err)
		return
	}
	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		fmt.Printf("写入数据失败 %s-%d: %v", sncode, len(bytes), err)
	}
}

func getRtcmFileName(sncode string, utcTime time.Time) string {
	utc := utcTime.UTC()
	year := utc.Year()
	dayOfYear := utc.YearDay()
	return fmt.Sprintf("%s.%d%03dbinRTCM3", sncode, year, dayOfYear)
}
