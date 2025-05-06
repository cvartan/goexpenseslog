package model

import "time"

type RawMessage struct {
	Id      int64     `json:"id"`
	UserId  int64     `json:"user_id"`
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
}

type RawMessageRepositoryHandler interface {
	Add(*RawMessage) error
	Delete(int64) error
	DeleteBeforeId(int64) error
	GetAll() (*[]RawMessage, error)
	GetUserData(int64) (*[]RawMessage, error)
}
