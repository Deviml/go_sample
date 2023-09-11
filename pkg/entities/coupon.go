package entities

type Coupon struct {
	CouponName string  `json:"coupon_name"`
	Discount   float32 `json:"discount"`
	IsValid    bool    `json:"is_valid"`
	UserID     string  `json:"user_id"`
	Message    string  `json:"message"`
}
