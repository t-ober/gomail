package mail

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Mail struct {
	ID        int
	HistoryId int
	Received  time.Time
	Sender    string
	Recipient string
	Subject   string
	Msg       string
}

const (
	DB_PATH    = "./mail_db.db"
	TABLE_NAME = "Mail"
)

type dbRepo struct {
	db   *sql.DB
	lock sync.RWMutex
}

func NewDBRepo(dbPath string) (*dbRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Could not establish connection to sqlite db: %v", err)
	}
	_, err = createTableIfNotExists(db)
	if err != nil {
		return nil, fmt.Errorf("Error creating db table: %v", err)
	}
	return &dbRepo{db: db}, nil
}

func createTableIfNotExists(db *sql.DB) (sql.Result, error) {
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
    id TEXT,
    history_id INTEGER,
    received TEXT,
    sender TEXT,
    recipient TEXT,
    subject TEXT,
    msg TEXT	
	);`, TABLE_NAME)
	return db.Exec(createTableSQL)
}

func (repo *dbRepo) ReadRecent(limit int) ([]*Mail, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	stmt := fmt.Sprintf("SELECT id, history_id, received, sender, recipient, subject, msg from %s ORDER BY received DESC LIMIT %d", TABLE_NAME, limit)
	rows, err := repo.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving data from database: %v", err)
	}
	defer rows.Close()
	return fromRows(rows)
}

func fromRows(rows *sql.Rows) ([]*Mail, error) {
	var mails []*Mail
	for rows.Next() {
		fmt.Println("Scanning row")
		var m Mail
		var recStr string
		err := rows.Scan(&m.ID, &m.HistoryId, &recStr, &m.Sender, &m.Recipient, &m.Subject, &m.Msg)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning row: %v", err)
		}
		m.Received, err = time.Parse(time.RFC3339, recStr)
		if err != nil {
			return nil, fmt.Errorf("time parsing failed: %v", err)
		}
		mails = append(mails, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during row iteration: %v", err)
	}
	return mails, nil
}
