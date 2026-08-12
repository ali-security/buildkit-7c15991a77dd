package containerdexecutor

import (
	"context"
	"testing"

	"github.com/containerd/containerd"
	"github.com/moby/buildkit/executor"
	gatewayapi "github.com/moby/buildkit/frontend/gateway/pb"
	"github.com/stretchr/testify/require"
)

func TestContainerdUnknownExitStatus(t *testing.T) {
	// There are assumptions in the containerd executor that the UnknownExitStatus
	// used in errdefs.ExitError matches the variable in the containerd package.
	if containerd.UnknownExitStatus != gatewayapi.UnknownExitStatus {
		t.Fatalf("containerd.UnknownExitStatus != errdefs.UnknownExitStatus")
	}
}

// TestRunInvalidContainerID checks that a container id that could be used to
// reference paths outside of the executor root is rejected before it reaches
// any of the container handling.
func TestRunInvalidContainerID(t *testing.T) {
	w := &containerdExecutor{
		running: make(map[string]*containerState),
	}

	for _, id := range []string{"../../../tmp/pwned", "..", "a/b", "a-b"} {
		_, err := w.Run(context.TODO(), id, executor.Mount{}, nil, executor.ProcessInfo{}, nil)
		require.Error(t, err, "id %q", id)
		require.Contains(t, err.Error(), "invalid container id", "id %q", id)
	}
}
