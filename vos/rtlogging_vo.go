package vos

type FileChannelVo struct {
	Sncode   string `json:"sncode"`
	LastTime int64  `json:"lastTime"`
	Size     int    `json:"size"`
	Status   int    `json:"status"`
}
