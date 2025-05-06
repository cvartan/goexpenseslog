package repos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cvartan/goexpenseslog/internal/model"
)

type ExpensesRepository struct {
	insertQ             *sql.Stmt
	deleteQ             *sql.Stmt
	deleteBeforeIdQ     *sql.Stmt
	getAllQ             *sql.Stmt
	getSummaryByPeriodQ *sql.Stmt
}

func NewExpensesRepository(db *sql.DB) *ExpensesRepository {
	var err error

	repo := &ExpensesRepository{}

	if repo.insertQ, err = db.Prepare("INSERT INTO expenses (id, exp_date, exp_value, exp_description, created) VALUES (?,?,?,?,?)"); err != nil {
		panic(fmt.Sprintf("error creating statement with error: %v", err))
	}

	if repo.deleteQ, err = db.Prepare("DELETE FROM expenses WHERE id=?"); err != nil {
		panic(fmt.Sprintf("error creating statement with error: %v", err))
	}

	if repo.deleteBeforeIdQ, err = db.Prepare("DELETE FROM expenses WHERE id<?"); err != nil {
		panic(fmt.Sprintf("error creating statement with error: %v", err))
	}

	if repo.getAllQ, err = db.Prepare("SELECT id, exp_date, exp_value, exp_description, created FROM expenses"); err != nil {
		panic(fmt.Sprintf("error creating statement with error: %v", err))
	}

	if repo.getSummaryByPeriodQ, err = db.Prepare("SELECT SUM(exp_value) summary FROM expenses WHERE exp_date BETWEEN ? AND ?"); err != nil {
		panic(fmt.Sprintf("error creating statement with error: %v", err))
	}

	return repo
}

func (r *ExpensesRepository) Add(expense *model.ExpenseInfo) error {
	created := time.Now()

	_, err := r.insertQ.Exec(expense.Id, expense.Date.Format(time.DateOnly), expense.Value, expense.Description, created.Format(time.StampMilli))
	if err != nil {
		return fmt.Errorf("insert expense error: %v", err)
	}

	return nil
}

func (r *ExpensesRepository) Delete(id int64) error {
	_, err := r.deleteQ.Exec(id)
	if err != nil {
		return fmt.Errorf("delete expense error: %v", err)
	}
	return nil
}

func (r *ExpensesRepository) DeleteBeforeId(id int64) error {
	_, err := r.deleteBeforeIdQ.Exec(id)
	if err != nil {
		return fmt.Errorf("delete expenses error: %v", err)
	}

	return nil
}

func (r *ExpensesRepository) GetAll() (*[]model.ExpenseInfo, error) {
	res, err := r.getAllQ.Query()
	if err != nil {
		return nil, fmt.Errorf("get expenses query error: %v", err)
	}
	defer res.Close()

	result := make([]model.ExpenseInfo, 0, 200)

	for res.Next() {
		var (
			id          int64
			date        string
			value       int32
			description sql.NullString
			created     string
		)
		err := res.Scan(&id, &date, &value, &description, &created)
		if err != nil {
			return nil, fmt.Errorf("parse expenses query result error: %v", err)
		}

		exp := model.ExpenseInfo{
			Id:    id,
			Value: value,
		}

		if description.Valid {
			exp.Description = description.String
		}

		if dt, err := time.Parse(time.DateOnly, date); err == nil {
			exp.Date = dt
		}

		if ct, err := time.Parse(time.StampMilli, created); err == nil {
			exp.Created = ct
		}

		result = append(result, exp)
	}

	return &result, nil
}

func (r *ExpensesRepository) GetMonthSummary() (int32, error) {
	year, month, _ := time.Now().Date()

	startMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
	endMonth := startMonth.AddDate(0, 1, -1)

	res, err := r.getSummaryByPeriodQ.Query(startMonth, endMonth)
	if err != nil {
		return 0, fmt.Errorf("get month summary query error: %v", err)
	}
	defer res.Close()

	var sum int32
	for res.Next() {
		err := res.Scan(&sum)
		if err != nil {
			return 0, fmt.Errorf("parse expenses query result error: %v", err)
		}
	}

	return sum, nil
}
func (r *ExpensesRepository) GetPrevMonthSummary() (int32, error) {
	year, month, _ := time.Now().AddDate(0, -1, 0).Date()

	startMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
	endMonth := startMonth.AddDate(0, 1, -1)

	res, err := r.getSummaryByPeriodQ.Query(startMonth, endMonth)
	if err != nil {
		return 0, fmt.Errorf("get month summary query error: %v", err)
	}
	defer res.Close()

	var sum sql.NullInt32
	for res.Next() {
		err := res.Scan(&sum)
		if err != nil {
			return 0, fmt.Errorf("parse expenses query result error: %v", err)
		}
	}

	var result int32
	if sum.Valid {
		result = sum.Int32
	}

	return result, nil
}
