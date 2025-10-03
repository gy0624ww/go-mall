package request

type DemoOrderCreate struct {
	UserId       int64 `json:"user_id"`
	BillMoney    int64 `json:"bill_money" binding:"required"`
	OrderGoodsId int64 `json:"order_goods_id" binding:"required"`
}
