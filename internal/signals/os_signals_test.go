package signals

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOSSignals(t *testing.T) {
	oss := NewOSSignals(context.Background())

	gotSignal := make(chan os.Signal, 1)

	oss.Subscribe(func(signal os.Signal) {
		gotSignal <- signal
	})

	defer oss.Stop()

	proc, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)

	assert.NoError(t, proc.Signal(syscall.SIGUSR2)) // send the signal

	time.Sleep(time.Millisecond * 5)
}

func TestNewOSSignalCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	oss := NewOSSignals(ctx)

	gotSignal := make(chan os.Signal, 1)

	oss.Subscribe(func(signal os.Signal) {
		gotSignal <- signal
	})

	defer oss.Stop()

	proc, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)

	cancel()

	assert.NoError(t, proc.Signal(syscall.SIGUSR2)) // send the signal

	assert.Empty(t, gotSignal)
}
