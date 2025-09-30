package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type DemoOrder struct {
	Id        int64                 `gorm:"column:id;primary_key" json:"id"`
	UserId    int64                 `gorm:"column:user_id" json:"user_id"`
	BillMoney int64                 `gorm:"column:bill_money" json:"bill_money"`
	OrderNo   string                `gorm:"column:order_no;type:varchar(32)" json:"order_no"`
	State     int8                  `gorm:"column:state;default:1" json:"state"`
	PaidAt    time.Time             `gorm:"column:paid_at;default:\"1970-01-01 00:00:00\"" json:"paid_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`
	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at" json:"updated_at"`
}

func (DemoOrder) TableName() string {
	return "demo_orders"
}
