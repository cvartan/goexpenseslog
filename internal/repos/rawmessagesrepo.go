package repos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/cvartan/goexpenseslog/internal/model"
)

type RawMessageRepository struct {
	node            *snowflake.Node
	insertQ         *sql.Stmt
	deleteQ         *sql.Stmt
	deleteBeforeIdQ *sql.Stmt
	getAllQ         *sql.Stmt
	getUserDataQ    *sql.Stmt
}

func NewRawMessgeRepository(db *sql.DB) *RawMessageRepository {
	n, err := snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("can't create node for generating snowflake id error: %v", err))
	}

	repo := &RawMessageRepository{
		node: n,
	}

	if repo.insertQ, err = db.Prepare("INSERT INTO raw_messages (id, user_id, msg_value, created) VALUES (?,?,?,?)"); err != nil {
		panic(fmt.Sprintf("error creating statement error: %v", err))
	}

	if repo.deleteQ, err = db.Prepare("DELETE FROM raw_messages WHERE id=?"); err != nil {
		panic(fmt.Sprintf("error creating statement error: %v", err))
	}

	if repo.deleteBeforeIdQ, err = db.Prepare("DELETE FROM raw_messages WHERE id<?"); err != nil {
		panic(fmt.Sprintf("error creating statement error: %v", err))
	}

	if repo.getAllQ, err = db.Prepare("SELECT id, user_id, msg_value, created FROM raw_messages ORDER BY id DESC LIMIT 20;"); err != nil {
		panic(fmt.Sprintf("error creating statement error: %v", err))
	}

	if repo.getUserDataQ, err = db.Prepare("SELECT id, user_id, msg_value, created FROM raw_messages WHERE user_id=? ORDER BY id DESC LIMIT 20;"); err != nil {
		panic(fmt.Sprintf("error creating statement error: %v", err))
	}

	return repo
}

func (r *RawMessageRepository) Add(message *model.RawMessage) error {
	id := r.node.Generate()
	created := time.Now()

	_, err := r.insertQ.Exec(id, message.UserId, message.Value, created.Format(time.StampMilli))
	if err != nil {
		return fmt.Errorf("insert raw message error: %v", err)
	}

	message.Id = id.Int64()
	message.Created = created.Round(time.Millisecond)

	return nil
}

func (r *RawMessageRepository) Delete(id int64) error {
	_, err := r.deleteQ.Exec(id)
	if err != nil {
		return fmt.Errorf("delete raw message error: %v", err)
	}

	return nil
}

func (r *RawMessageRepository) DeleteBeforeId(id int64) error {
	_, err := r.deleteBeforeIdQ.Exec(id)
	if err != nil {
		return fmt.Errorf("delete raw message error: %v", err)
	}

	return nil
}

func (r *RawMessageRepository) GetAll() (*[]model.RawMessage, error) {
	res, err := r.getAllQ.Query()
	if err != nil {
		return nil, fmt.Errorf("query raw messages cancel with error: %v", err)
	}

	defer res.Close()

	result := make([]model.RawMessage, 0, 200)
	for res.Next() {
		var (
			id      int64
			userId  int64
			value   sql.NullString
			created string
		)
		if err := res.Scan(&id, &userId, &value, &created); err != nil {
			return nil, fmt.Errorf("parse query result error: %v", err)
		}

		rm := model.RawMessage{
			Id:     id,
			UserId: userId,
		}
		if value.Valid {
			rm.Value = value.String
		}

		if c, err := time.Parse(time.StampMilli, created); err == nil {
			rm.Created = c
		}

		result = append(result, rm)

	}

	return &result, nil
}

func (r *RawMessageRepository) GetUserData(user_id int64) (*[]model.RawMessage, error) {
	res, err := r.getUserDataQ.Query(user_id)
	if err != nil {
		return nil, fmt.Errorf("query raw messages cancel with error: %v", err)
	}

	defer res.Close()

	result := make([]model.RawMessage, 0, 200)
	for res.Next() {
		var (
			id      int64
			userId  int64
			value   sql.NullString
			created string
		)
		if err := res.Scan(&id, &userId, &value, &created); err != nil {
			return nil, fmt.Errorf("parse query result error: %v", err)
		}

		rm := model.RawMessage{
			Id:     id,
			UserId: userId,
		}
		if value.Valid {
			rm.Value = value.String
		}

		if c, err := time.Parse(time.StampMilli, created); err == nil {
			rm.Created = c
		}

		result = append(result, rm)

	}

	return &result, nil
}
