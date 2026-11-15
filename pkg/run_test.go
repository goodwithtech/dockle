package pkg

import (
	"context"
	"os"
	"testing"

	"github.com/urfave/cli/v3"
)

func setOptionalEnv(t *testing.T, key, value string) {
	t.Helper()
	original, exists := os.LookupEnv(key)
	if value == "" {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
	} else if err := os.Setenv(key, value); err != nil {
		t.Fatalf("set %s: %v", key, err)
	}
	t.Cleanup(func() {
		if exists {
			_ = os.Setenv(key, original)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}

func TestGetDockerSockPath(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		dockleHost string
		dockerHost string
		xdgRuntime string
		expected   string
	}{
		{
			name:       "explicit host takes precedence",
			args:       []string{"--host", "tcp://explicit:2375"},
			dockerHost: "tcp://docker:2375",
			expected:   "tcp://explicit:2375",
		},
		{
			name:       "dockle host environment takes precedence",
			dockleHost: "tcp://dockle:2375",
			dockerHost: "tcp://docker:2375",
			expected:   "tcp://dockle:2375",
		},
		{
			name:       "standard docker host environment",
			dockerHost: "tcp://docker:2375",
			expected:   "tcp://docker:2375",
		},
		{
			name:       "xdg socket",
			args:       []string{"--use-xdg"},
			xdgRuntime: "/run/user/1000",
			expected:   "unix:///run/user/1000/docker.sock",
		},
		{
			name:     "default socket",
			expected: "unix:///var/run/docker.sock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setOptionalEnv(t, "DOCKLE_HOST", tt.dockleHost)
			setOptionalEnv(t, "DOCKER_HOST", tt.dockerHost)
			setOptionalEnv(t, "XDG_RUNTIME_DIR", tt.xdgRuntime)

			var actual string
			cmd := NewApp()
			cmd.Action = func(_ context.Context, cmd *cli.Command) error {
				actual = getDockerSockPath(cmd)
				return nil
			}
			if err := cmd.Run(context.Background(), append([]string{"dockle"}, tt.args...)); err != nil {
				t.Fatalf("run command: %v", err)
			}
			if actual != tt.expected {
				t.Fatalf("getDockerSockPath() = %q, want %q", actual, tt.expected)
			}
		})
	}
}
