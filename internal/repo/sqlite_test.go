package repo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"st-test/internal/models"
	"st-test/internal/settings"
)

const (
	storagePath = "storage.db"
)

func removeStorage(t *testing.T) {
	t.Helper()

	_ = os.Remove(storagePath)
}

// testRepo returns repo for testing purposes.
func testRepo(t *testing.T) *Repo {
	t.Helper()

	sets := settings.LocalStorageSettings{Path: storagePath}
	repo, err := NewRepo(sets)
	require.NoError(t, err)
	require.NotNil(t, repo)

	return repo
}

func TestNewRepo(t *testing.T) {
	sets := settings.LocalStorageSettings{Path: storagePath}
	repo, err := NewRepo(sets)
	require.NoError(t, err)
	require.NotNil(t, repo)

	removeStorage(t)
}

func TestRepo_Insert(t *testing.T) {
	repo := testRepo(t)
	defer removeStorage(t)

	err := repo.Insert(models.Item{
		ID:   1,
		Body: []byte(`{"some":"body"}`),
	})
	require.NoError(t, err)

	_ = repo.Delete(1)

	repo.Close()
}

func TestRepo_ReadAll(t *testing.T) {
	repo := testRepo(t)
	defer removeStorage(t)

	err := repo.Insert(models.Item{
		ID:   1,
		Body: []byte(`{"some":"body"}`),
	})
	require.NoError(t, err)

	err = repo.Insert(models.Item{
		ID:   2,
		Body: []byte(`{"some2":"body2"}`),
	})
	require.NoError(t, err)

	gotItems, err := repo.ReadAll()
	require.NoError(t, err)
	require.Len(t, gotItems, 2)

	_ = repo.Delete(1)
	_ = repo.Delete(2)

	repo.Close()
}

func TestRepo_Delete(t *testing.T) {
	repo := testRepo(t)
	defer removeStorage(t)

	err := repo.Insert(models.Item{
		ID:   1,
		Body: []byte(`{"some":"body"}`),
	})
	require.NoError(t, err)

	err = repo.Delete(1)
	require.NoError(t, err)

	_, err = repo.Read(1)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrNotFound)

	repo.Close()
}

func TestRepo_DeleteAll(t *testing.T) {
	repo := testRepo(t)
	defer removeStorage(t)

	err := repo.Insert(models.Item{
		ID:   1,
		Body: []byte(`{"some":"body"}`),
	})
	require.NoError(t, err)

	err = repo.DeleteAll()
	require.NoError(t, err)

	_, err = repo.Read(1)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrNotFound)

	repo.Close()
}
