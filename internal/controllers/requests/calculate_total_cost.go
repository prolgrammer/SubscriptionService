package requests

type CalculateTotalCost struct {
	StartPeriod string `json:"start_period" binding:"required" example:"07-2025"`
	EndPeriod   string `json:"end_period" binding:"required" example:"12-2025"`
	UserID      string `json:"user_id,omitempty" binding:"omitempty,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name,omitempty" example:"Yandex Plus"`
}
