package responses

type CalculateTotalCost struct {
	Total int `json:"total" binding:"required"`
}
