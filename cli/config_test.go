package cli_test

import (
	"os"
	"path"
	"testing"

	"github.com/BrunoQuaresma/openwritter/cli"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()

		path := "invalid-path"
		_, err := cli.NewConfig(path)
		require.Error(t, err)
		require.EqualError(t, err, "open invalid-path: no such file or directory")
	})

	t.Run("invalid yaml", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		configPath := path.Join(tmp, "config.yaml")
		err := os.WriteFile(configPath, []byte("invalid-yaml"), 0644)
		require.NoError(t, err)

		_, err = cli.NewConfig(configPath)
		require.Error(t, err)
		require.ErrorContains(t, err, "unmarshal errors")
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		configPath := path.Join(tmp, "config.yaml")
		yaml := `
profiles:
  - name: "Development"
    description: "Profile for development environment"
`
		err := os.WriteFile(configPath, []byte(yaml), 0644)
		require.NoError(t, err)

		cfg, err := cli.NewConfig(configPath)
		require.NoError(t, err)
		require.Equal(t, cfg.Profiles[0].Name, "Development")
		require.Equal(t, cfg.Profiles[0].Description, "Profile for development environment")
	})
}
