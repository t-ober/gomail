package mail

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func TestDbConnection(t *testing.T) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal("Failed creating temporary file.")
	}
	defer file.Close()
	defer os.Remove(file.Name())
	repo, err := NewDBRepo(file.Name())
	if err != nil {
		t.Fatalf("Error getting db connection %v", err)
	}
	defer repo.db.Close()
	exists, err := tableExists(repo.db, TABLE_NAME)
	if err != nil {
		t.Fatalf("Error checking for table name %v", err)
	}
	if !exists {
		t.Fatalf("Table %s does not exist", TABLE_NAME)
	}
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name=?;
	`
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func TestReadRecent(t *testing.T) {
	fmt.Println("Test read recent print")
	repo, cleanup := dbSetup(t)
	defer cleanup()
	mails, err := repo.ReadRecent(2)
	if err != nil {
		t.Fatalf("Error reading recent mails: %v", err)
	}
	if !(len(mails) == 2) {
		t.Fatalf("Expected 2 mails, got %d", len(mails))
	}
}

func dbSetup(t *testing.T) (*dbRepo, func()) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed creating temporary file: %v", err)
	}
	repo, err := NewDBRepo(file.Name())
	if err != nil {
		t.Fatalf("Error getting db connection %v", err)
	}
	exists, err := tableExists(repo.db, TABLE_NAME)
	if err != nil {
		t.Fatalf("Error checking for table name %v", err)
	}
	if !exists {
		t.Fatalf("Table %s does not exist", TABLE_NAME)
	}
	repo.lock.Lock()
	defer repo.lock.Unlock()
	err = populateDB(repo.db)
	if err != nil {
		t.Fatalf("Error populating db: %v", err)
	}
	rowCount, err := getTableRowCount(repo.db, TABLE_NAME)
	if err != nil {
		t.Fatalf("Error getting row count: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("Poulating db failed, expected 5 rows, got %d", rowCount)
	}
	return repo, func() {
		repo.db.Close()
		file.Close()
		os.Remove(file.Name())
	}
}

func populateDB(db *sql.DB) error {
	query := fmt.Sprintf(`
		INSERT INTO
		%s (
			id,
			history_id,
			received,
			sender,
			recipient,
			subject,
			msg
		)
	VALUES
		(
			'1',
			1,
			'2023-04-01T10:00:00Z',
			'sender1@example.com',
			'recipient1@example.com',
			'Hello',
			'This is a test message'
		),
		(
			'2',
			2,
			'2023-04-01T11:30:00Z',
			'sender2@example.com',
			'recipient2@example.com',
			'Meeting Reminder',
			'Don''t forget our meeting at 2 PM'
		),
		(
			'3',
			3,
			'2023-04-02T09:15:00Z',
			'sender3@example.com',
			'recipient3@example.com',
			'Project Update',
			'The project is progressing well'
		),
		(
			'4',
			4,
			'2023-04-03T14:45:00Z',
			'sender4@example.com',
			'recipient4@example.com',
			'Question about API',
			'Can you clarify how to use the new API?'
		),
		(
			'5',
			5,
			'2023-04-04T16:20:00Z',
			'sender5@example.com',
			'recipient5@example.com',
			'Weekend Plans',
			'Are you free this weekend for a hike?'
		);
	`, TABLE_NAME)
	_, err := db.Exec(query)
	return err
}

func getTableRowCount(db *sql.DB, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
