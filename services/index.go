package services

import (
	"errors"
	"nav-receive-go/domains"
	"nav-receive-go/global"
	"nav-receive-go/utils"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	DeviceService
	DeviceRtcmService
}

type HasBaseData interface {
	GetBaseData() domains.BaseDataEntity
}

type CrudService[T HasBaseData] struct {
}

var dbMutex sync.Mutex

func (s *CrudService[T]) Create(entity T) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	if entity.GetBaseData().Guid == "" {
		return global.NAV_DB.Create(&entity).Error
	}
	var existing T
	tx := global.NAV_DB.Where("deleted = ?", 0).Where("guid = ?", entity.GetBaseData().Guid).First(&existing)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return global.NAV_DB.Create(&entity).Error
		}
		return tx.Error
	}
	return global.NAV_DB.Model(&existing).Where("deleted = ?", 0).Where("guid = ?", entity.GetBaseData().Guid).Updates(&entity).Error
}

func (s *CrudService[T]) Updates(entity T) error {
	return global.NAV_DB.Model(&entity).Where("deleted = ?", 0).Where("guid = ?", entity.GetBaseData().Guid).Updates(&entity).Error
}

func (s *CrudService[T]) Update(entity T, field string, value interface{}) error {
	return global.NAV_DB.Model(&entity).Where("deleted = ?", 0).Where("guid = ?", entity.GetBaseData().Guid).Update(field, value).Error
}

func (s *CrudService[T]) DeleteByGuid(guid string) error {
	return global.NAV_DB.Where("deleted = ?", 0).Where("deleted = ?", 0).Where("guid = ?", guid).Delete(new(T)).Error
}

func (s *CrudService[T]) GetByGuid(guid string) (*T, error) {
	var result T
	err := global.NAV_DB.Where("deleted = ?", 0).Where("guid = ?", guid).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (s *CrudService[T]) GetById(id uint) (*T, error) {
	var result T
	err := global.NAV_DB.Where("deleted = ?", 0).Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (s *CrudService[T]) GetByField(field string, value string) (*T, error) {
	var result T
	err := global.NAV_DB.Where("deleted = ?", 0).Where(utils.CamelToSnake(field)+" = ?", value).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (s *CrudService[T]) List(pageInfo domains.PageInfo, likeFields string) (list interface{}, total int64, err error) {
	limit := pageInfo.Size
	offset := pageInfo.Size * (pageInfo.Page - 1)
	var result []T
	db := global.NAV_DB.Where("deleted = ?", 0)
	if pageInfo.Content != "" && likeFields != "" {
		likeFieldSlice := strings.Split(likeFields, ",")
		for index, field := range likeFieldSlice {
			if index == 0 {
				db = db.Where(field+" like ?", "%"+pageInfo.Content+"%")
			} else {
				db = db.Or(field+" like ?", "%"+pageInfo.Content+"%")
			}
		}
	}
	db = db.Model(new(T))
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	order := "id desc"
	if pageInfo.Asc != "" {
		order = pageInfo.Asc + " asc"
	}
	if pageInfo.Desc != "" {
		order = pageInfo.Desc + " desc"
	}
	err = db.Order(order).Limit(limit).Offset(offset).Find(&result).Error
	return result, total, err
}

func (s *CrudService[T]) ListAll(params map[string]string) ([]T, error) {
	var results []T
	db := global.NAV_DB.Where("deleted = ?", 0)
	for key, value := range params {
		db = db.Where(utils.CamelToSnake(key)+" = ?", value)
	}
	err := db.Find(&results).Error
	return results, err
}

func (s *CrudService[T]) SafeFirst(result T) (*T, error) {
	err := global.NAV_DB.Where("deleted = ?", 0).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (s *CrudService[T]) SafeLast(result T) (*T, error) {
	err := global.NAV_DB.Where("deleted = ?", 0).Last(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
