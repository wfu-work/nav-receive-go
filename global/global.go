package global

import (
	"nav-receive-go/configs"
	"sync"

	"gorm.io/gorm"

	"github.com/spf13/viper"
)

var (
	NAV_DB     *gorm.DB
	NAV_CONFIG configs.Server
	NAV_VIPER  *viper.Viper
	Lock       sync.RWMutex
)
