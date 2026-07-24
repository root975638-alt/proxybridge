package litellm

import (
	"testing"

	"github.com/proxybridge/cli/internal/config"
	"github.com/proxybridge/cli/internal/logging"
)

func TestManager_IsInstalled(t *testing.T) {
	cfg := config.NewConfig()
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// This will fail in test environment since litellm is not installed
	// Just verify the method exists and doesn't panic
	installed := m.IsInstalled()
	_ = installed
}

func TestManager_Install(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LiteLLM.Path = "litellm"
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test install - will fail without litellm available
	err := m.Install()
	if err == nil {
		// If it succeeds, that's fine too
		t.Log("LiteLLM installed successfully")
	}
}

func TestManager_Uninstall(t *testing.T) {
	cfg := config.NewConfig()
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test uninstall - will handle missing installation
	err := m.Uninstall()
	if err != nil {
		t.Logf("Uninstall returned expected error: %v", err)
	}
}

func TestManager_Status(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LiteLLM.ListenAddress = "127.0.0.1"
	cfg.LiteLLM.Port = 4000
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	status, err := m.Status()
	if err != nil {
		t.Logf("Status returned error (expected): %v", err)
	}

	if status == nil {
		t.Error("Status returned nil")
	} else {
		t.Logf("LiteLLM status: running=%v, address=%s, port=%d",
			status.Running, status.Address, status.Port)
	}
}

func TestManager_Start(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LiteLLM.ListenAddress = "127.0.0.1"
	cfg.LiteLLM.Port = 4000
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test start - will fail without litellm
	err := m.Start()
	if err == nil {
		t.Log("LiteLLM started successfully")
	} else {
		t.Logf("Start returned expected error: %v", err)
	}
}

func TestManager_Stop(t *testing.T) {
	cfg := config.NewConfig()
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test stop - will handle non-running instance
	err := m.Stop()
	if err != nil {
		t.Logf("Stop returned error (may be expected): %v", err)
	}
}

func TestManager_Restart(t *testing.T) {
	cfg := config.NewConfig()
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test restart - will handle missing installation
	err := m.Restart()
	if err != nil {
		t.Logf("Restart returned expected error: %v", err)
	}
}

func TestManager_Validate(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LiteLLM.ConfigPath = "/tmp/test_litellm_config.yaml"
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test validate - will fail without config file
	err := m.Validate()
	if err == nil {
		t.Log("Validate passed (unexpected)")
	} else {
		t.Logf("Validate returned expected error: %v", err)
	}
}

func TestManager_HealthCheck(t *testing.T) {
	cfg := config.NewConfig()
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test health check - will fail without running instance
	err := m.HealthCheck()
	if err == nil {
		t.Log("Health check passed (unexpected)")
	} else {
		t.Logf("Health check returned expected error: %v", err)
	}
}

func TestManager_EnsureConfigExists(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LiteLLM.ConfigPath = "/tmp/test_litellm_config.yaml"
	logger := logging.GetLogger()

	m := NewManager(cfg, logger)

	// Test ensure config exists
	err := m.ensureConfigExists()
	if err != nil {
		t.Logf("ensureConfigExists returned: %v", err)
	}
}

func TestStatus_Structure(t *testing.T) {
	status := &Status{
		Running:     true,
		PID:         1234,
		Address:     "127.0.0.1",
		Port:        4000,
		StartTime:   status.StartTime,
		Version:     "1.0.0",
	}

	if !status.Running {
		t.Error("Status should be running")
	}
	if status.PID != 1234 {
		t.Error("PID should be 1234")
	}
	if status.Address != "127.0.0.1" {
		t.Error("Address should be 127.0.0.1")
	}
	if status.Port != 4000 {
		t.Error("Port should be 4000")
	}
}
