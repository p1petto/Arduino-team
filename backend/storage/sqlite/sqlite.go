package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"arduinoteam/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	*sql.DB
}

func NewStorage() *Storage {
	db, err := sql.Open("sqlite3", "database/store.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(username TEXT PRIMARY KEY,score INTEGER DEFAULT 0, room TEXT, admin INTEGER, apikey TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rooms(id TEXT PRIMARY KEY, name TEXT )`)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db}
}

func (s *Storage) SaveRoom(name string, id string) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.Prepare("INSERT INTO rooms(name, id) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(name, id)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrRoomExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row_id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return row_id, nil
}
func (s *Storage) SaveKey(name string, apikey string) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.Prepare("INSERT INTO rooms(name, apikey) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(name, apikey)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrRoomExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row_id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return row_id, nil
}
