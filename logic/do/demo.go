package do

import "time"

type DemoOrder struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	BillMoney int64     `json:"bill_money"`
	OrderNo   string    `json:"order_no"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	PaidAt    time.Time `json:"paid_at"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
}
