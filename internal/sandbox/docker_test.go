package sandbox

import (
	"testing"
)

// Docker sandbox tests are skipped since we have a stub implementation.
// Full implementation requires Docker SDK setup.
// See docker_impl.go.reference for full implementation.

func TestDockerSandbox_NotImplemented(t *testing.T) {
	cfg := &Config{
		Type: "docker",
	}

	_, err := NewDockerSandbox(cfg)
	if err == nil {
		t.Error("expected error for docker sandbox")
	}
}
