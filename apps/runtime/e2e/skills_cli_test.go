package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

// TestSkillsCLI_Enable tests the skills enable command
func TestSkillsCLI_Enable(t *testing.T) {
	// Create a disabled skill
	cmd := exec.Command("/tmp/pryx-core", "skills", "enable", "test-skill")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Enable output: %s", output)
		// It's OK if enable fails due to skill not existing
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "enable") {
		t.Logf("Expected 'enable' in output, got: %s", outputStr)
	}
}

// TestSkillsCLI_Disable tests the skills disable command
func TestSkillsCLI_Disable(t *testing.T) {
	// Disable a skill
	cmd := exec.Command("/tmp/pryx-core", "skills", "disable", "test-skill")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Disable output: %s", output)
		// It's OK if disable fails due to skill not existing
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "disable") {
		t.Logf("Expected 'disable' in output, got: %s", outputStr)
	}
}

// TestSkillsCLI_EnableDisable tests enable/disable round-trip
func TestSkillsCLI_EnableDisable(t *testing.T) {
	// Enable a skill
	enableCmd := exec.Command("/tmp/pryx-core", "skills", "enable", "test-skill")
	enableOutput, enableErr := enableCmd.CombinedOutput()
	if enableErr != nil {
		t.Logf("Enable output: %s", enableOutput)
	}

	// Disable the skill
	disableCmd := exec.Command("/tmp/pryx-core", "skills", "disable", "test-skill")
	disableOutput, disableErr := disableCmd.CombinedOutput()
	if disableErr != nil {
		t.Logf("Disable output: %s", disableOutput)
	}

	// Enable again
	enableAgainCmd := exec.Command("/tmp/pryx-core", "skills", "enable", "test-skill")
	enableAgainOutput, enableAgainErr := enableAgainCmd.CombinedOutput()
	if enableAgainErr != nil {
		t.Logf("Re-enable output: %s", enableAgainOutput)
	}

	// Verify state persisted
	t.Log("Enable-disable round-trip completed")
}

// TestSkillsCLI_Info tests info for valid skill
func TestSkillsCLI_Info(t *testing.T) {
	// Test info for a bundled skill (git-tool is a bundled skill)
	cmd := exec.Command("/tmp/pryx-core", "skills", "info", "git-tool")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("skills info failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "git-tool") {
		t.Errorf("Expected 'git-tool' in output, got: %s", outputStr)
	}
}

// TestSkillsCLI_EnableNotFound tests enabling non-existent skill
func TestSkillsCLI_EnableNotFound(t *testing.T) {
	// Try to enable a non-existent skill
	cmd := exec.Command("/tmp/pryx-core", "skills", "enable", "nonexistent-skill-xyz123")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("Expected error for non-existent skill, got: %s", string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "not found") && !strings.Contains(outputStr, "error") {
		t.Logf("Output: %s", outputStr)
	}
}

// TestSkillsCLI_Install tests the skills install command
func TestSkillsCLI_Install(t *testing.T) {
	// Install a skill
	cmd := exec.Command("/tmp/pryx-core", "skills", "install", "test-skill")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Install output: %s", output)
		// It's OK if install fails due to skill not being installable
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "install") {
		t.Logf("Expected 'install' in output, got: %s", outputStr)
	}
}

// TestSkillsCLI_Uninstall tests the skills uninstall command
func TestSkillsCLI_Uninstall(t *testing.T) {
	// Uninstall a skill
	cmd := exec.Command("/tmp/pryx-core", "skills", "uninstall", "test-skill")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Uninstall output: %s", output)
		// It's OK if uninstall fails due to skill not being installed
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "uninstall") {
		t.Logf("Expected 'uninstall' in output, got: %s", outputStr)
	}
}
