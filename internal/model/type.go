package model

import "time"

type ExpenseInfo struct {
	Id          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Value       int32     `json:"value"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

type RawMessage struct {
	Id      int64     `json:"id"`
	UserId  int64     `json:"user_id"`
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
}
