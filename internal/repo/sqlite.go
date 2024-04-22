// Package repo содержит методы для получения, хранения и удаления объектов в файл. Данное решение является на базе sqlite.
package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"st-test/internal/models"
	"st-test/internal/settings"

	// Use CGO-free sqlite driver.
	_ "modernc.org/sqlite"
)

// Repo хранит объект для работы с БД и предоставляет методы для удобной работы.
type Repo struct {
	db *sql.DB
}

// NewRepo создаёт хранилище sqlite с единственной ключ-значение таблицей и возвращает объект Repo.
func NewRepo(set settings.LocalStorageSettings) (*Repo, error) {
	var db *sql.DB

	var errOpen error

	db, errOpen = sql.Open("sqlite", set.Path)

	if errOpen != nil {
		return nil, fmt.Errorf("opening sqlite repo: %w", errOpen)
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS storage (key INTEGER PRIMARY KEY, value BLOB NOT NULL)")
	if err != nil {
		cerr := db.Close()
		if cerr != nil {
			slog.Warn(fmt.Sprintf("closing sqlite repo: %v; ignore", cerr.Error()))
		}

		return nil, fmt.Errorf("creating 'storage' table: %w", err)
	}

	return &Repo{db: db}, nil
}

// Insert вставляет объект в таблицу.
func (r *Repo) Insert(item models.Item) error {
	_, err := r.db.Exec("INSERT INTO storage (key, value) VALUES (?, ?)", item.ID, item.Body)
	if err != nil {
		return fmt.Errorf("inserting key %d: %w", item.ID, err)
	}

	return nil
}

// Read возвращает объект по ключу.
func (r *Repo) Read(key int) (models.Item, error) {
	var body []byte

	err := r.db.QueryRow("SELECT value FROM storage WHERE key = ? LIMIT 1", key).Scan(&body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, models.ErrNotFound
		}

		return models.Item{}, fmt.Errorf("read from repo: %w", err)
	}

	return models.Item{
		ID:      key,
		Body:    body,
		Expires: 0,
	}, nil
}

// ReadAll возвращает все объекты из таблицы.
func (r *Repo) ReadAll() ([]models.Item, error) {
	rows, err := r.db.Query("SELECT key, value FROM storage")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("read from repo: %w", err)
	}

	items := make([]models.Item, 0)

	for rows.Next() {
		i := models.Item{}

		err = rows.Scan(&i.ID, &i.Body)
		if err != nil {
			return nil, fmt.Errorf("scan row from repo: %w", err)
		}

		items = append(items, i)
	}

	return items, nil
}

// Delete удаляет объект из таблицы по ключу.
func (r *Repo) Delete(key int) error {
	_, err := r.db.Exec("DELETE FROM storage WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("deleting key %d: %w", key, err)
	}

	return nil
}

// DeleteAll удаляет все объекты из таблицы.
func (r *Repo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM storage")
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

// Close закрывает sqlite-базу.
func (r *Repo) Close() {
	err := r.db.Close()
	if err != nil {
		slog.Warn(fmt.Sprintf("closing repo: %v; ignore", err.Error()))
	}
}
