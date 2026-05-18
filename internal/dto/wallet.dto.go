package dto

type DashboardUser struct {
	Balance float64 `db:"balance"`
	Expense float64 `db:"expense"`
	Income  float64 `db:"income"`
}
