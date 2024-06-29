package main

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Queue struct {
	db *sql.DB
}

func InitQueue() (*Queue, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY,
			topic TEXT,
			message TEXT,
			read INTEGER DEFAULT 0
		);
	`)

	if err != nil {
		return nil, err
	}

	queue := Queue{
		db: db,
	}

	return &queue, err
}

func (q *Queue) Emit(topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = q.db.Exec("INSERT INTO messages (topic, message) VALUES (?, ?)", topic, data)
	return err
}

func (q *Queue) Peek(topic string) (int64, string, error) {
	now := time.Now().UnixMilli()
	query, err := q.db.Query("SELECT id, message FROM messages WHERE topic = ? AND read < ? ORDER BY id LIMIT 1", topic, now-60000)
	if err != nil {
		return 0, "", err
	}
	defer query.Close()

	var id int64
	var message string
	if query.Next() {
		query.Scan(&id)
		query.Scan(&message)
	} else {
		return 0, "", nil
	}

	query.Close()

	_, err = q.db.Exec("UPDATE messages SET read = ? WHERE id = ?", now, id)
	if err != nil {
		return 0, "", err
	}

	return id, message, nil
}

func (q *Queue) Delete(messageId int64) error {
	_, err := q.db.Exec("DELETE FROM messages WHERE messageId = ?", messageId)
	return err
}
