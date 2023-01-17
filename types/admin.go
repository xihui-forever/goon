package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelAdmin struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_admin_username,unique"`

	Username string `json:"username,omitempty" gorm:"column:username;not null;index:idx_admin_username,unique"`
	Password string `json:"password,omitempty" gorm:"column:password;not null"`
}

/*
	username: admin
	deleted_at: 0

	idx_admin_username: admin-0
*/

func (m *ModelAdmin) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelAdmin) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelAdmin) TableName() string {
	return "goon_admin"
}

type (
	AdminLoginReq struct {
		Username string `json:"username,omitempty" form:"username" binding:"required"`
		Password string `json:"password,omitempty" form:"password" binding:"required"`
	}

	AdminLoginRsp struct {
		Token  string `json:"token,omitempty"`
		Expire uint32 `json:"expire,omitempty"`
	}

	AdminSession struct {
		Id       uint64 `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	}
)
