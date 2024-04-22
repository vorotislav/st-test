package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"st-test/internal/storage/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"st-test/internal/models"
)

func TestNewStore(t *testing.T) {
	t.Parallel()
	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	repo := mocks.NewRepo(t)
	require.NotNil(t, repo)
	repo.EXPECT().ReadAll().Once().Return(nil, models.ErrNotFound)

	s := NewStore(log, repo)
	require.NotNil(t, s)
	require.NotNil(t, s.s)
}

func TestStore_SaveObject(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	repo := mocks.NewRepo(t)
	require.NotNil(t, repo)
	repo.EXPECT().ReadAll().Once().Return(nil, models.ErrNotFound)

	s := NewStore(log, repo)
	require.NotNil(t, s)

	gotID, err := s.SaveObject(context.Background(), models.Item{
		ID:      1,
		Body:    []byte(`{"some":"body"}`),
		Expires: 0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, gotID)
	require.Len(t, s.s, 1)

	updateGotID, err := s.SaveObject(context.Background(), models.Item{
		ID:      1,
		Body:    []byte(`{"some":"body"}`),
		Expires: 0,
	})
	require.NoError(t, err)
	require.Equal(t, 0, updateGotID)
	require.Len(t, s.s, 1)
}

func TestStore_GetObject(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	repo := mocks.NewRepo(t)
	require.NotNil(t, repo)
	repo.EXPECT().ReadAll().Once().Return(nil, models.ErrNotFound)

	s := NewStore(log, repo)
	require.NotNil(t, s)

	gotID, err := s.SaveObject(context.Background(), models.Item{
		ID:      1,
		Body:    []byte(`{"some":"body"}`),
		Expires: 0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, gotID)
	require.Len(t, s.s, 1)

	gotItem, err := s.GetObject(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, 1, gotItem.ID)
	require.Equal(t, []byte(`{"some":"body"}`), gotItem.Body)

	notFoundItem, err := s.GetObject(context.Background(), 2)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrNotFound)
	require.Nil(t, notFoundItem.Body)
}

func TestStore_Stop(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	cases := []struct {
		name         string
		prepareRepo  func(rep *mocks.Repo)
		prepareStore func(s *Store)
	}{
		{
			name: "without objects",
			prepareRepo: func(rep *mocks.Repo) {
				rep.EXPECT().ReadAll().Once().Return(nil, models.ErrNotFound)
			},
		},
		{
			name: "with objects",
			prepareRepo: func(rep *mocks.Repo) {
				rep.EXPECT().ReadAll().Once().Return(nil, models.ErrNotFound)
				rep.EXPECT().DeleteAll().Once().Return(nil)
				rep.EXPECT().Insert(mock.AnythingOfType("models.Item")).Once().Return(nil)
			},
			prepareStore: func(s *Store) {
				s.s[1] = models.Item{}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepo(t)
			require.NotNil(t, repo)
			if tc.prepareRepo != nil {
				tc.prepareRepo(repo)
			}

			s := NewStore(log, repo)
			require.NotNil(t, s)

			if tc.prepareStore != nil {
				tc.prepareStore(s)
			}

			s.Stop()
		})
	}
}
