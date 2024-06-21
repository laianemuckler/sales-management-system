package models

import "time"

type Order struct {
	ID         int
	Time       time.Time
	TotalValue float64
	EmployeeID int
}
