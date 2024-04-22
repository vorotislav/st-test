// Package storage предоставляет хранилище для объектов.
// хранилище реализовано на базе map. В качестве примитива синхронизации используется мьютекс.
package storage

import (
	"context"
	"errors"
	"sync"

	"st-test/internal/models"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

var errNotAvailable = errors.New("store is not available")

// repo описывает методы хранилища для сохранения и получения объектов на диск.
//
//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=repo --with-expecter=true --exported
type repo interface {
	Insert(item models.Item) error
	ReadAll() ([]models.Item, error)
	DeleteAll() error
}

// Store является локальным хранилищем объектов в оперативной памяти.
type Store struct {
	log  *zap.Logger
	s    map[int]models.Item
	m    sync.Mutex
	repo repo
}

// NewStore конструктор для хранилища. Так же после успешного создания пытаемся прочитать объекты из файла.
func NewStore(log *zap.Logger, repo repo) *Store {
	s := &Store{
		log:  log.Named("store"),
		s:    make(map[int]models.Item),
		repo: repo,
	}

	s.loadItems()

	return s
}

// Stop "останавливает" работу хранилища и записывает все текущие объекты в файл.
func (s *Store) Stop() {
	s.saveItems()
}

// SaveObject сохраняет объект в хранилище.
func (s *Store) SaveObject(_ context.Context, item models.Item) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	s.log.Info("New item request", zap.Int("id", item.ID), zap.Int64("expires", int64(item.Expires)))

	_, ok := s.s[item.ID]
	if ok {
		s.log.Info("the item has already been saved")

		return 0, nil //nolint:nlreturn
	}

	s.s[item.ID] = item

	s.log.Info("the object was saved successfully")

	return item.ID, nil
}

// GetObject возвращает объект из хранилища по id.
func (s *Store) GetObject(_ context.Context, id int) (models.Item, error) {
	s.m.Lock()
	defer s.m.Unlock()

	s.log.Info("Request on get item", zap.Int("id", id))

	item, ok := s.s[id]
	if !ok {
		s.log.Info("Item not found", zap.Int("id", id))

		return models.Item{}, models.ErrNotFound
	}

	s.log.Info("Item was found", zap.Int("id", id))

	return item, nil
}

// Check возвращает состояние хранилища. Необходим для обработчика здоровья.
func (s *Store) Check() error {
	if s.s == nil {
		return errNotAvailable
	}

	return nil
}

func (s *Store) loadItems() {
	items, err := s.repo.ReadAll()
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.Info("no items in local repo")

			return
		}

		s.log.Error("cannot load items from local repo", zap.Error(err))

		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	for _, item := range items {
		s.s[item.ID] = item
	}

	s.log.Info("successful load items from local repo", zap.Int("items size", len(items)))
}

func (s *Store) saveItems() {
	if len(s.s) == 0 {
		s.log.Info("no items to save into local repo")

		return
	}

	err := s.repo.DeleteAll()
	if err != nil {
		s.log.Error("cannot remove old items from local repo", zap.Error(err))

		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	items := maps.Values(s.s)
	for _, item := range items {
		err := s.repo.Insert(item)
		if err != nil {
			s.log.Error("cannot save item in local repo", zap.Error(err))

			continue
		}

		s.log.Info("item successful insert to local repo", zap.Int("id", item.ID))
	}
}
