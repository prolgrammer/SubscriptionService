package responses

type SubResponse struct {
	ID          string `json:"id" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}
