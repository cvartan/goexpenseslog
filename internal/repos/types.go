package repos

import "github.com/cvartan/goexpenseslog/internal/model"

type ExpensesRepository interface {
	Add(*model.ExpenseInfo) error
	Delete(int64) error
	DeleteAllBeforeId(int64) error
	GetAll() (*[]model.ExpenseInfo, error)
	GetMonthSummary() (int32, error)
	GetPrevMonthSummary() (int32, error)
}

type RawMessageRepository interface {
	Add(*model.RawMessage) error
	Delete(int64) error
	DeleteAllBeforeId(int64) error
	GetAll() (*[]model.RawMessage, error)
	GetUserData(int64) (*[]model.RawMessage, error)
}
