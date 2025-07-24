package requests

type SubRequest struct {
	ServiceName string `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int    `json:"price" binding:"required" example:"400"`
	UserID      string `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}
