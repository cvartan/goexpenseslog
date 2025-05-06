package model

import "time"

type ExpenseInfo struct {
	Id          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Value       int32     `json:"value"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

type ExpensesRepositoryHandler interface {
	Add(*ExpenseInfo) error
	Delete(int64) error
	DeleteBeforeId(int64) error
	GetAll() (*[]ExpenseInfo, error)
	GetMonthSummary() (int32, error)
	GetPrevMonthSummary() (int32, error)
}
