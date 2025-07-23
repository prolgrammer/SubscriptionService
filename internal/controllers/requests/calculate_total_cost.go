package requests

type CalculateTotalCost struct {
	StartPeriod string `json:"start_period" binding:"required"`
	EndPeriod   string `json:"end_period" binding:"required"`
	UserID      string `json:"user_id,omitempty" binding:"omitempty,uuid"`
	ServiceName string `json:"service_name,omitempty"`
}
