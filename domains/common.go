package domains

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	Id uint `gorm:"primarykey;autoIncrement" json:"-"`
}

type BaseDataEntity struct {
	BaseEntity
	Guid       string    `gorm:"column:guid;size:50;unique" json:"guid"`
	Creater    string    `gorm:"column:creater;size:50" json:"-"`
	Updater    string    `gorm:"column:updater;size:50" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
	Deleted    int       `gorm:"index" json:"-"`
}

func (s *BaseDataEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Guid == "" {
		s.Guid = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	if s.CreateTime.IsZero() {
		s.CreateTime = time.Now()
	}
	s.UpdateTime = time.Now()
	return nil
}

func (s BaseDataEntity) GetBaseData() BaseDataEntity {
	return s
}

type PageResult struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type PageInfo struct {
	Page    int    `json:"page" form:"page"`       // 页码
	Size    int    `json:"size" form:"size"`       // 每页大小
	Desc    string `json:"desc" form:"desc"`       // 倒序
	Asc     string `json:"asc" form:"asc"`         // 正序
	Content string `json:"content" form:"content"` // 查询内容
}

type Empty struct{}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.Size > 100:
			r.Size = 100
		case r.Size <= 0:
			r.Size = 10
		}
		offset := (r.Page - 1) * r.Size
		return db.Offset(offset).Limit(r.Size)
	}
}
