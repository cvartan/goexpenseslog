package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cvartan/goexpenseslog/internal/model"
)

type RestService struct {
	rawMessageRepo model.RawMessageRepositoryHandler
	expensesRepo   model.ExpensesRepositoryHandler
}

func NewRestService(rawMessageRepo model.RawMessageRepositoryHandler, expensesRepo model.ExpensesRepositoryHandler) *RestService {
	return &RestService{
		rawMessageRepo: rawMessageRepo,
		expensesRepo:   expensesRepo,
	}
}

func (s *RestService) GetAllRawMessages(resp http.ResponseWriter, req *http.Request) {
	result, err := s.rawMessageRepo.GetAll()
	if err != nil {
		log.Printf("can't get raw messages from database with error: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("can't marshal raw messages with error: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Write(bytes)
	resp.WriteHeader(http.StatusOK)
}

func (s *RestService) GetAllExpenses(resp http.ResponseWriter, req *http.Request) {
	result, err := s.expensesRepo.GetAll()
	if err != nil {
		log.Printf("can't get expenses from database with error: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("can't marshal expenses with error: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Write(bytes)
	resp.WriteHeader(http.StatusOK)
}
