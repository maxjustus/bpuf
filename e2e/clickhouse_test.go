package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unionFindOutput struct {
	Result []struct {
		Value string `json:"value"`
		Root  string `json:"root"`
	} `json:"result"`
}

type bipartiteOutput struct {
	Result []struct {
		U     string `json:"u"`
		VRoot string `json:"v_root"`
	} `json:"result"`
}

func TestClickHouseUDFs(t *testing.T) {
	if _, err := exec.LookPath("clickhouse"); err != nil {
		t.Skip("clickhouse not found, skipping e2e tests")
	}

	udfBinary := buildUDF(t)
	userScriptsDir, configFile := setupClickHouseEnv(t, udfBinary)

	tests := []struct {
		name     string
		query    string
		validate func(t *testing.T, output []byte)
	}{
		{
			name:  "unionFind",
			query: "SELECT unionFind([('user1', 'user2'), ('user2', 'user3'), ('user4', 'user5')]) as result FORMAT JSONEachRow",
			validate: func(t *testing.T, output []byte) {
				var result unionFindOutput
				err := json.Unmarshal(output, &result)
				require.NoError(t, err)
				require.Len(t, result.Result, 5)

				roots := make(map[string]string)
				for _, r := range result.Result {
					roots[r.Value] = r.Root
				}

				assert.Equal(t, roots["user1"], roots["user2"])
				assert.Equal(t, roots["user2"], roots["user3"])
				assert.Equal(t, roots["user4"], roots["user5"])
				assert.NotEqual(t, roots["user1"], roots["user4"])
			},
		},
		{
			name: "bipartiteUnionFind",
			query: `SELECT bipartiteUnionFind([
				('entity1', 'group100'),
				('entity1', 'group101'),
				('entity2', 'group101'),
				('entity3', 'group200')
			]) as result FORMAT JSONEachRow`,
			validate: func(t *testing.T, output []byte) {
				var result bipartiteOutput
				err := json.Unmarshal(output, &result)
				require.NoError(t, err)
				require.Len(t, result.Result, 3)

				vRoots := make(map[string]string)
				for _, r := range result.Result {
					vRoots[r.U] = r.VRoot
				}

				assert.Equal(t, vRoots["entity1"], vRoots["entity2"])
				assert.NotEqual(t, vRoots["entity1"], vRoots["entity3"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := runClickHouseQuery(t, tt.query, userScriptsDir, configFile)
			tt.validate(t, output)
		})
	}
}

func buildUDF(t *testing.T) string {
	projectRoot, err := findProjectRoot()
	require.NoError(t, err)

	udfBinary := filepath.Join(projectRoot, "bin", "bpuf-clickhouse")
	cmd := exec.Command("go", "build", "-o", udfBinary, "./cmd/bpuf-clickhouse") //nolint:gosec // test code with known safe args
	cmd.Dir = projectRoot
	require.NoError(t, cmd.Run())

	return udfBinary
}

func setupClickHouseEnv(t *testing.T, udfBinary string) (userScriptsDir, configFile string) {
	userScriptsDir = t.TempDir()
	udfInScripts := filepath.Join(userScriptsDir, "bpuf-clickhouse")

	cmd := exec.Command("cp", "-f", udfBinary, udfInScripts) //nolint:gosec // test code with known safe args
	require.NoError(t, cmd.Run())
	require.NoError(t, os.Chmod(udfInScripts, 0o755)) //nolint:gosec // test needs executable permissions

	cmd = exec.Command(udfInScripts, "--udf-xml") //nolint:gosec // test code with known safe args
	xmlConfig, err := cmd.Output()
	require.NoError(t, err)

	configFile = filepath.Join(t.TempDir(), "udf_config.xml")
	require.NoError(t, os.WriteFile(configFile, xmlConfig, 0o644)) //nolint:gosec // test file permissions

	return userScriptsDir, configFile
}

func runClickHouseQuery(t *testing.T, query, userScriptsDir, configFile string) []byte {
	cmd := exec.Command("clickhouse", "local", //nolint:gosec // test code with known safe args
		"-q", query,
		"--",
		"--user_scripts_path="+userScriptsDir,
		"--user_defined_executable_functions_config="+configFile,
	)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "ClickHouse command failed: %s", string(output))
	return output
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}
