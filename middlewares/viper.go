package middlewares

import (
	"flag"
	"fmt"
	"nav-rtlogging-go/global"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	ConfigEnv         = "NAV_CONFIG"
	ConfigCommonFile  = "config-common.yaml"
	ConfigDefaultFile = "config.yaml"
	ConfigTestFile    = "config.test.yaml"
	ConfigDebugFile   = "config.debug.yaml"
	ConfigReleaseFile = "config.release.yaml"
)

var initFlags sync.Once

// Viper 配置
func Viper() *viper.Viper {
	config := GetConfigPath()
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	for _, key := range viper.AllKeys() {
		if !v.IsSet(key) {
			v.Set(key, viper.Get(key))
		}
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file changed: %v\n", e)
		if err = v.Unmarshal(&global.NAV_CONFIG); err != nil {
			fmt.Println(err)
			return
		}
	})
	if err = v.Unmarshal(&global.NAV_CONFIG); err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}
	return v
}

// GetConfigPath 获取配置文件路径, 优先级: 命令行 > 环境变量 > 默认值
func GetConfigPath() (config string) {
	initFlags.Do(func() {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
	})
	if config != "" {
		fmt.Printf("您正在使用命令行的 '-c' 参数传递的值, config 的路径为 %s\n", config)
		return
	}
	if env := os.Getenv(ConfigEnv); env != "" {
		config = env
		fmt.Printf("您正在使用 %s 环境变量, config 的路径为 %s\n", ConfigEnv, config)
		return
	}

	_, err := os.Stat(config)
	if err != nil || os.IsNotExist(err) {
		config = ConfigDefaultFile
		fmt.Printf("配置文件路径不存在, 使用默认配置文件路径: %s\n", config)
	}

	return
}
