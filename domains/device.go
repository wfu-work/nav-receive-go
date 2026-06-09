package domains

type Device struct {
	BaseDataEntity
	CompanyGuid  string  `json:"companyGuid" gorm:"size:50;comment:企业guid"`
	ProjectGuid  string  `json:"projectGuid" gorm:"size:50;comment:项目guid"`
	ProjectName  string  `json:"projectName" gorm:"size:50;comment:项目名"`
	Sncode       string  `json:"sncode" gorm:"comment:设备名称"`
	Alias        string  `json:"alias" gorm:"comment:别名"`
	Sn           string  `json:"sn" gorm:"comment:序列号"`
	DeviceType   string  `json:"deviceType" gorm:"comment:设备类型"`
	DeviceId     string  `json:"deviceId" gorm:"comment:设备id"`
	DeviceKey    string  `json:"deviceKey" gorm:"comment:设备key"`
	Manufacturer string  `json:"manufacturer" gorm:"comment:制造商"`
	Model        string  `json:"model" gorm:"comment:型号"`
	Agreement    string  `json:"agreement" gorm:"comment:接入协议　0:MQTT 1-HTTP 2-CoAP"`
	Tags         string  `json:"tags" gorm:"comment:标签"`
	Longitude    float64 `json:"longitude" gorm:"comment:经度"`
	Latitude     float64 `json:"latitude" gorm:"comment:纬度"`
	Altitude     float64 `json:"altitude" gorm:"comment:高程"`
	Azi          float64 `json:"azi" gorm:"comment:方位角"`
	Status       int     `json:"status" gorm:"comment:状态 0：禁用 1：启用"`
	Base         int     `json:"base" gorm:"comment:是否基准站"`
	Warn         int     `json:"warn" gorm:"comment:是否告警"`
	WarnSwitch   int     `json:"warnSwitch" gorm:"comment:告警开关"`
	RawSwitch    int     `json:"rawSwitch" gorm:"comment:离线开关"`
	Remark       string  `json:"remark" gorm:"comment:备注"`
}

func (Device) TableName() string {
	return "nav_device"
}

func (s Device) GetBaseData() BaseDataEntity {
	return s.BaseDataEntity
}
