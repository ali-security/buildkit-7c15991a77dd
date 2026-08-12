//go:build linux
// +build linux

package runcexecutor

import (
	"context"
	"testing"

	"github.com/moby/buildkit/executor"
	"github.com/stretchr/testify/require"
)

// TestRunInvalidContainerID checks that a container id that would let the
// bundle directory escape the executor root is rejected before any directory
// is created for it.
func TestRunInvalidContainerID(t *testing.T) {
	w := &runcExecutor{
		root:    t.TempDir(),
		running: make(map[string]chan error),
	}

	for _, id := range []string{"../../../tmp/pwned", "..", "a/b", "a-b"} {
		_, err := w.Run(context.TODO(), id, executor.Mount{}, nil, executor.ProcessInfo{}, nil)
		require.Error(t, err, "id %q", id)
		require.Contains(t, err.Error(), "invalid container id", "id %q", id)
	}
}
