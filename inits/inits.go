package inits

import (
	"nav-receive-go/dbs"
	"nav-receive-go/kafkas"
	"nav-receive-go/middlewares"
	"nav-receive-go/scheduleds"
	"nav-receive-go/utils"
)

func Init() {
	middlewares.Viper()
	dbs.GormMysql()
	utils.InitRedis()
	scheduleds.Init()
	_ = scheduleds.RtloggingSched()
	kafkas.StartRtloggingConsumer()
}
