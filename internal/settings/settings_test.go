package settings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testSettingsPath = "../test/test-config.yaml"

func TestNewSettings(t *testing.T) {
	t.Parallel()

	sets, err := NewSettings(testSettingsPath)
	require.NoError(t, err)
	require.NotNil(t, sets)

	expected := Settings{}
	expected.API.Address = "127.0.0.1"
	expected.API.Port = 8080

	expected.Storage.Path = "st-test.db"

	expected.Log.Level = "debug"
	expected.Log.Verbose = true
	expected.Log.Format = "text"

	require.Equal(t, expected, *sets)
}
