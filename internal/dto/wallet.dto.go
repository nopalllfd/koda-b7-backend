package dto

type DashboardUser struct {
	Balance float64 `json:"balance"`
	Expense float64 `json:"expense"`
	Income  float64 `json:"income"`
}

type DashboardSwaggerResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    DashboardUser `json:"data"`
}
