package domains

type DeviceRtcm struct {
	BaseDataEntity
	DeviceGuid string  `json:"deviceGuid" gorm:"size:50;comment:设备guid"`
	Sncode     string  `json:"sncode" gorm:"comment:设备名称"`
	Agreement  string  `json:"agreement" gorm:"comment:协议"`
	Ip         string  `json:"ip" gorm:"comment:主机"`
	Port       int     `json:"port" gorm:"comment:端口"`
	MountPoint string  `json:"mountPoint" gorm:"comment:挂载点"`
	Username   string  `json:"username" gorm:"comment:用户名"`
	Password   string  `json:"password" gorm:"comment:密码"`
	Longitude  float64 `json:"longitude" gorm:"comment:经度"`
	Latitude   float64 `json:"latitude" gorm:"comment:纬度"`
	Altitude   float64 `json:"altitude" gorm:"comment:高程"`
	Status     int     `json:"status" gorm:"comment:状态 0：禁用 1：启用"`
	DataCycle  int     `json:"dataCycle" gorm:"comment:数据周期"`
}

func (DeviceRtcm) TableName() string {
	return "nav_device_rtcm"
}

func (s DeviceRtcm) GetBaseData() BaseDataEntity {
	return s.BaseDataEntity
}
