//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func waitForFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for file: %s", path)
}

func TestCLI_SkillsListJSON_IncludesBundledSkills(t *testing.T) {
	home := t.TempDir()
	bundled := filepath.Join(home, "bundled-skills")
	for _, id := range []string{"docker-manager", "git-tool", "cloud-deploy"} {
		writeSkill(t, bundled, id, true)
	}

	out, code := runPryxCoreWithEnv(t, home, map[string]string{
		"PRYX_BUNDLED_SKILLS_DIR": bundled,
		"PRYX_MANAGED_SKILLS_DIR": filepath.Join(home, "managed-skills"),
		"PRYX_WORKSPACE_ROOT":     filepath.Join(home, "workspace"),
	}, "skills", "list", "--json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}

	var skills []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(out), &skills); err != nil {
		t.Fatalf("expected json output, got:\n%s\nerror: %v", out, err)
	}

	got := map[string]bool{}
	for _, s := range skills {
		got[s.ID] = true
	}

	for _, id := range []string{"docker-manager", "git-tool", "cloud-deploy"} {
		if !got[id] {
			t.Fatalf("expected bundled skill %q in list, got ids: %v", id, keys(got))
		}
	}
}

func TestCLI_SkillsInfo_WorksForBundledSkill(t *testing.T) {
	home := t.TempDir()
	bundled := filepath.Join(home, "bundled-skills")
	writeSkill(t, bundled, "docker-manager", true)

	out, code := runPryxCoreWithEnv(t, home, map[string]string{
		"PRYX_BUNDLED_SKILLS_DIR": bundled,
		"PRYX_MANAGED_SKILLS_DIR": filepath.Join(home, "managed-skills"),
		"PRYX_WORKSPACE_ROOT":     filepath.Join(home, "workspace"),
	}, "skills", "info", "docker-manager")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}
	if !strings.Contains(out, "Skill: docker-manager") {
		t.Fatalf("expected output to contain skill header, got:\n%s", out)
	}
}

func TestCLI_MCPConfig_RoundTrip(t *testing.T) {
	home := t.TempDir()

	out, code := runPryxCoreWithEnv(t, home, nil, "mcp", "list", "--json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}

	out, code = runPryxCoreWithEnv(t, home, nil, "mcp", "add", "test-server", "--url", "https://example.com")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}

	out, code = runPryxCoreWithEnv(t, home, nil, "mcp", "list", "--json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}

	var servers map[string]any
	if strings.TrimSpace(out) != "" && strings.TrimSpace(out) != "null" {
		if err := json.Unmarshal([]byte(out), &servers); err != nil {
			t.Fatalf("expected json output, got:\n%s\nerror: %v", out, err)
		}
	}
	if servers == nil || servers["test-server"] == nil {
		t.Fatalf("expected test-server in config, got:\n%s", out)
	}

	out, code = runPryxCoreWithEnv(t, home, nil, "mcp", "remove", "test-server")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}
}

func TestCLI_Config_SetThenGet(t *testing.T) {
	home := t.TempDir()

	out, code := runPryxCoreWithEnv(t, home, nil, "config", "set", "listen_addr", ":12345")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}

	out, code = runPryxCoreWithEnv(t, home, nil, "config", "get", "listen_addr")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out)
	}
	if strings.TrimSpace(out) != ":12345" {
		t.Fatalf("expected listen_addr to be updated, got: %q", strings.TrimSpace(out))
	}
}

func TestRuntime_HealthAndWebsocket(t *testing.T) {
	home := t.TempDir()
	bin := buildPryxCore(t)
	bundled := filepath.Join(home, "bundled-skills")
	writeSkill(t, bundled, "alpha-skill", true)

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(),
		"HOME="+home,
		"PRYX_DB_PATH="+filepath.Join(home, "pryx.db"),
		"PRYX_WORKSPACE_ROOT="+filepath.Join(home, "workspace"),
		"PRYX_BUNDLED_SKILLS_DIR="+bundled,
		"PRYX_KEYCHAIN_FILE="+filepath.Join(home, ".pryx", "keychain.json"),
		"PRYX_TELEMETRY_DISABLED=true",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("start runtime: %v", err)
	}
	waitCh := make(chan error, 1)
	go func() { waitCh <- cmd.Wait() }()
	t.Cleanup(func() {
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-time.After(2 * time.Second):
			_ = cmd.Process.Kill()
			<-waitCh
		case <-waitCh:
		}
	})

	portFile := filepath.Join(home, ".pryx", "runtime.port")
	if err := waitForFile(portFile, 5*time.Second); err != nil {
		t.Fatalf("%v", err)
	}

	portBytes, err := os.ReadFile(portFile)
	if err != nil {
		t.Fatalf("read port file: %v", err)
	}
	port := strings.TrimSpace(string(portBytes))
	if port == "" {
		t.Fatalf("empty port file")
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected /health response: %d %q", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("expected json /health output, got:\n%s\nerror: %v", strings.TrimSpace(string(body)), err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("expected status ok, got: %v (body: %s)", payload["status"], strings.TrimSpace(string(body)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := fmt.Sprintf("ws://localhost:%s/ws?surface=e2e&event=trace.event", port)
	c, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial ws: %v", err)
	}
	_ = c.Close(websocket.StatusNormalClosure, "")
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func runtimeRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	for i := 0; i < 8; i++ {
		if filepath.Base(dir) == "runtime" && filepath.Base(filepath.Dir(dir)) == "apps" {
			return dir
		}
		next := filepath.Dir(dir)
		if next == dir {
			break
		}
		dir = next
	}

	t.Fatalf("could not locate apps/runtime from cwd")
	return ""
}

func repoRoot(t *testing.T) string {
	t.Helper()

	runtime := runtimeRoot(t)
	apps := filepath.Dir(runtime)
	return filepath.Dir(apps)
}

func startPryxCore(t *testing.T, bin string, home string) (port string, cancel context.CancelFunc) {
	t.Helper()

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(),
		"HOME="+home,
		"PRYX_DB_PATH="+filepath.Join(home, "pryx.db"),
		"PRYX_LISTEN_ADDR=:0",
		"PRYX_WORKSPACE_ROOT="+repoRoot(t),
		"PRYX_BUNDLED_SKILLS_DIR="+filepath.Join(runtimeRoot(t), "internal", "skills", "bundled"),
		"PRYX_KEYCHAIN_FILE="+filepath.Join(home, ".pryx", "keychain.json"),
		"PRYX_TELEMETRY_DISABLED=true",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("start runtime: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-time.After(2 * time.Second):
			_ = cmd.Process.Kill()
		case <-func() chan struct{} {
			ch := make(chan struct{})
			go func() {
				_ = cmd.Wait()
				close(ch)
			}()
			return ch
		}():
		}
	}()

	portFile := filepath.Join(home, ".pryx", "runtime.port")
	if err := waitForFile(portFile, 5*time.Second); err != nil {
		cancel()
		t.Fatalf("%v", err)
	}

	portBytes, err := os.ReadFile(portFile)
	if err != nil {
		cancel()
		t.Fatalf("read port file: %v", err)
	}
	port = strings.TrimSpace(string(portBytes))
	if port == "" {
		cancel()
		t.Fatalf("empty port file")
	}

	return port, cancel
}

func waitForServer(t *testing.T, port string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://localhost:" + port + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for server to be ready")
}
