package message

type AddBalanceReq struct {
	TargetID string  `json:"t_id" binding:"required, uuid"`   //充值目标ID
	Value    float32 `json:"value" binding:"required, min=0"` // 充值金额
}

type BuyReq struct {
	ProductID string `json:"p_id" binding:"required, uuid"`    //购买商品ID
	Amount    int    `json:"amount" binding:"required, min=0"` // 购买数量
}
