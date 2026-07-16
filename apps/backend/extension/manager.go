package extension

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ExtensionManifest struct {
	Name      string   `json:"name"`
	Command   []string `json:"command"`
	Enabled   bool     `json:"enabled"`
	TimeoutMs int      `json:"timeout_ms"`
}

type ExtensionInstance struct {
	Name    string
	Port    int
	Command []string
	Process *os.Process
}

type ExtensionManager struct {
	mu         sync.RWMutex
	extensions map[string]*ExtensionInstance
	httpClient *http.Client
}

func NewExtensionManager() *ExtensionManager {
	return &ExtensionManager{
		extensions: make(map[string]*ExtensionInstance),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *ExtensionManager) Load(manifest ExtensionManifest) error {
	if !manifest.Enabled {
		slog.Info("Extension is disabled, skipping", "name", manifest.Name)
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.extensions[manifest.Name]; exists {
		return fmt.Errorf("extension %s is already loaded", manifest.Name)
	}

	if len(manifest.Command) == 0 {
		return fmt.Errorf("extension %s has no command", manifest.Name)
	}
	if manifest.TimeoutMs <= 0 {
		manifest.TimeoutMs = 10000
	}

	cmd := exec.Command(manifest.Command[0], manifest.Command[1:]...)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("extension %s failed to create stdout pipe: %w", manifest.Name, err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("extension %s failed to start: %w", manifest.Name, err)
	}

	port, err := readPortFromStdout(stdout, manifest.Name, time.Duration(manifest.TimeoutMs)*time.Millisecond)
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("extension %s: %w", manifest.Name, err)
	}

	healthURL := fmt.Sprintf("http://localhost:%d/health", port)
	if err := m.healthCheck(healthURL, time.Duration(manifest.TimeoutMs)*time.Millisecond); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("extension %s: %w", manifest.Name, err)
	}

	instance := &ExtensionInstance{
		Name:    manifest.Name,
		Port:    port,
		Command: manifest.Command,
		Process: cmd.Process,
	}

	m.extensions[manifest.Name] = instance

	slog.Info("Extension loaded",
		"name", instance.Name,
		"port", instance.Port,
		"pid", instance.Process.Pid,
	)

	return nil
}

func (m *ExtensionManager) Unload(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension %s is not loaded", name)
	}

	if err := instance.Process.Signal(os.Interrupt); err == nil {
		done := make(chan struct{})
		go func() {
			instance.Process.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			instance.Process.Kill()
		}
	} else {
		instance.Process.Kill()
	}

	delete(m.extensions, name)

	slog.Info("Extension unloaded", "name", name)
	return nil
}

func (m *ExtensionManager) Call(name string, action string, params map[string]interface{}) (*ExecuteResponse, error) {
	m.mu.RLock()
	instance, ok := m.extensions[name]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("extension %s is not loaded", name)
	}

	reqBody := ExecuteRequest{
		Action: action,
		Params: params,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%d/execute", instance.Port)
	resp, err := m.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("extension %s call failed: %w", name, err)
	}
	defer resp.Body.Close()

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("extension %s returned invalid response: %w", name, err)
	}

	return &result, nil
}

func (m *ExtensionManager) List() []map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]map[string]interface{}, 0, len(m.extensions))
	for _, instance := range m.extensions {
		result = append(result, map[string]interface{}{
			"name": instance.Name,
			"port": instance.Port,
			"pid":  instance.Process.Pid,
		})
	}
	return result
}

func (m *ExtensionManager) GetPort(name string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	instance, ok := m.extensions[name]
	if !ok {
		return 0, false
	}
	return instance.Port, true
}

func (m *ExtensionManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, instance := range m.extensions {
		slog.Info("Shutting down extension", "name", name)
		if err := instance.Process.Signal(os.Interrupt); err == nil {
			done := make(chan struct{})
			go func() {
				instance.Process.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				instance.Process.Kill()
			}
		} else {
			instance.Process.Kill()
		}
	}

	m.extensions = make(map[string]*ExtensionInstance)
	slog.Info("All extensions shut down")
}

func (m *ExtensionManager) healthCheck(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}

	deadline := time.After(timeout)
	for {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			return fmt.Errorf("health check returned status %d", resp.StatusCode)
		}

		select {
		case <-deadline:
			return fmt.Errorf("health check timed out after %s", timeout)
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func readPortFromStdout(stdout io.ReadCloser, name string, timeout time.Duration) (int, error) {
	portCh := make(chan int, 1)
	errCh := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, PortAnnouncePrefix) {
				portStr := strings.TrimPrefix(line, PortAnnouncePrefix)
				port, err := strconv.Atoi(strings.TrimSpace(portStr))
				if err != nil {
					errCh <- fmt.Errorf("invalid port from extension %s: %q", name, portStr)
					return
				}
				portCh <- port
				return
			}
		}
		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("stdout read error from extension %s: %w", name, err)
			return
		}
		errCh <- fmt.Errorf("extension %s exited without announcing port", name)
	}()

	select {
	case port := <-portCh:
		return port, nil
	case err := <-errCh:
		return 0, err
	case <-time.After(timeout):
		return 0, fmt.Errorf("extension %s port announcement timed out after %s", name, timeout)
	}
}

var globalManager *ExtensionManager

func InitGlobal(extensionsDir string) error {
	globalManager = NewExtensionManager()

	manifests, err := LoadManifests(extensionsDir)
	if err != nil {
		return fmt.Errorf("failed to scan extensions directory: %w", err)
	}

	for _, mf := range manifests {
		if err := globalManager.Load(mf); err != nil {
			slog.Error("Failed to load extension, continuing", "name", mf.Name, "error", err)
		}
	}

	slog.Info("Extension manager initialized", "loaded", len(globalManager.List()))
	return nil
}

func GetGlobal() *ExtensionManager {
	return globalManager
}

func ShutdownGlobal() {
	if globalManager != nil {
		globalManager.Shutdown()
	}
}

func LoadManifests(dir string) ([]ExtensionManifest, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read extensions directory %s: %w", dir, err)
	}

	var manifests []ExtensionManifest
	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := fmt.Sprintf("%s/%s/extension.json", dir, entry.Name())
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				continue
			}

			var manifest ExtensionManifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				slog.Warn("Invalid extension manifest, skipping", "path", manifestPath, "error", err)
				continue
			}

			manifests = append(manifests, manifest)
		}
	}

	return manifests, nil
}
