package performance

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewStartupProfiler(t *testing.T) {
	profiler := NewStartupProfiler()

	if profiler == nil {
		t.Fatal("NewStartupProfiler returned nil")
	}

	if !profiler.IsEnabled() {
		t.Error("New profiler should be enabled by default")
	}

	if len(profiler.GetAllPhases()) != 0 {
		t.Error("New profiler should have no phases")
	}
}

func TestStartupProfiler_StartPhase(t *testing.T) {
	profiler := NewStartupProfiler()

	phase := profiler.StartPhase("test-phase")
	if phase == nil {
		t.Fatal("StartPhase returned nil")
	}

	if phase.Name != "test-phase" {
		t.Errorf("Expected phase name 'test-phase', got '%s'", phase.Name)
	}

	if phase.StartTime.IsZero() {
		t.Error("Phase should have non-zero start time")
	}

	// Verify phase was recorded
	phases := profiler.GetAllPhases()
	if len(phases) != 1 {
		t.Errorf("Expected 1 phase, got %d", len(phases))
	}
}

func TestStartupProfiler_EndPhase(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.StartPhase("test-phase")
	time.Sleep(10 * time.Millisecond)
	profiler.EndPhase("test-phase", nil)

	phase, exists := profiler.GetPhase("test-phase")
	if !exists {
		t.Fatal("Phase should exist")
	}

	if !phase.Completed() {
		t.Error("Phase should be marked as completed")
	}

	if phase.Duration < 10*time.Millisecond {
		t.Errorf("Duration should be at least 10ms, got %v", phase.Duration)
	}
}

func TestStartupProfiler_EndPhase_Error(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.StartPhase("error-phase")
	testErr := errors.New("test error")
	profiler.EndPhase("error-phase", testErr)

	phase, _ := profiler.GetPhase("error-phase")
	if phase.Error == nil {
		t.Error("Phase should have error recorded")
	}

	if phase.Error.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%v'", phase.Error)
	}
}

func TestStartupProfiler_EndPhase_Unknown(t *testing.T) {
	profiler := NewStartupProfiler()

	// This should not panic, just log a warning
	profiler.EndPhase("non-existent-phase", nil)
}

func TestStartupProfiler_TimeFunc(t *testing.T) {
	profiler := NewStartupProfiler()

	executed := false
	err := profiler.TimeFunc("timed-func", func() error {
		executed = true
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !executed {
		t.Error("Function should have been executed")
	}

	phase, exists := profiler.GetPhase("timed-func")
	if !exists {
		t.Fatal("Phase should exist")
	}

	if phase.Duration < 5*time.Millisecond {
		t.Errorf("Duration should be at least 5ms, got %v", phase.Duration)
	}
}

func TestStartupProfiler_TimeFuncWithResult(t *testing.T) {
	profiler := NewStartupProfiler()

	err := profiler.TimeFuncWithResult("timed-func-with-result", func() (string, error) {
		time.Sleep(5 * time.Millisecond)
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	phase, exists := profiler.GetPhase("timed-func-with-result")
	if !exists {
		t.Fatal("Phase should exist")
	}

	if phase.Duration < 5*time.Millisecond {
		t.Errorf("Duration should be at least 5ms, got %v", phase.Duration)
	}
}

func TestStartupProfiler_MarkComplete(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.MarkComplete()

	duration := profiler.GetTotalDuration()
	if duration <= 0 {
		t.Errorf("Total duration should be positive, got %v", duration)
	}
}

func TestStartupProfiler_GetAllPhases(t *testing.T) {
	profiler := NewStartupProfiler()

	// Add phases in specific order
	profiler.StartPhase("phase-1")
	time.Sleep(5 * time.Millisecond)
	profiler.EndPhase("phase-1", nil)

	profiler.StartPhase("phase-2")
	time.Sleep(5 * time.Millisecond)
	profiler.EndPhase("phase-2", nil)

	phases := profiler.GetAllPhases()
	if len(phases) != 2 {
		t.Fatalf("Expected 2 phases, got %d", len(phases))
	}

	// Should be sorted by start time
	if phases[0].Name != "phase-1" || phases[1].Name != "phase-2" {
		t.Error("Phases should be sorted by start time")
	}
}

func TestStartupProfiler_GenerateReport(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.StartPhase("fast-phase")
	profiler.EndPhase("fast-phase", nil)

	profiler.StartPhase("slow-phase")
	time.Sleep(20 * time.Millisecond)
	profiler.EndPhase("slow-phase", nil)

	profiler.MarkComplete()

	report := profiler.GenerateReport()

	if report == "" {
		t.Error("Report should not be empty")
	}

	// Should contain phase names
	if !contains(report, "fast-phase") || !contains(report, "slow-phase") {
		t.Error("Report should contain phase names")
	}

	// Should contain total time
	if !contains(report, "Total Time:") {
		t.Error("Report should contain total time")
	}
}

func TestStartupProfiler_EnableDisable(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.Disable()
	if profiler.IsEnabled() {
		t.Error("Profiler should be disabled")
	}

	// Should not record when disabled
	phase := profiler.StartPhase("disabled-phase")
	if phase != nil {
		t.Error("StartPhase should return nil when disabled")
	}

	profiler.Enable()
	if !profiler.IsEnabled() {
		t.Error("Profiler should be enabled")
	}
}

func TestStartupProfiler_Reset(t *testing.T) {
	profiler := NewStartupProfiler()

	profiler.StartPhase("phase-1")
	profiler.EndPhase("phase-1", nil)

	profiler.Reset()

	if len(profiler.GetAllPhases()) != 0 {
		t.Error("Should have no phases after reset")
	}

	if profiler.GetTotalDuration() < 0 {
		t.Error("Total duration should be positive after reset")
	}
}

func TestAsyncProfiler(t *testing.T) {
	profiler := NewStartupProfiler()

	async := profiler.StartAsync("async-task")

	// Complete in background
	go func() {
		time.Sleep(10 * time.Millisecond)
		async.Complete("done", nil)
	}()

	// Wait with timeout
	err := async.WaitWithTimeout(100 * time.Millisecond)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !async.IsDone() {
		t.Error("Async task should be done")
	}

	phase, exists := profiler.GetPhase("async-task")
	if !exists {
		t.Fatal("Phase should exist")
	}

	if phase.Duration < 10*time.Millisecond {
		t.Errorf("Duration should be at least 10ms, got %v", phase.Duration)
	}
}

func TestAsyncProfiler_Timeout(t *testing.T) {
	profiler := NewStartupProfiler()

	async := profiler.StartAsync("slow-async-task")

	// Never complete
	err := async.WaitWithTimeout(10 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{500 * time.Microsecond, "500µs"},
		{5 * time.Millisecond, "5ms"},
		{1500 * time.Millisecond, "1.50s"},
		{2 * time.Second, "2.00s"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is..."},
		{"exactly", 7, "exactly"},
	}

	for _, test := range tests {
		result := truncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateString(%s, %d) = %s, expected %s", test.input, test.maxLen, result, test.expected)
		}
	}
}

func TestNewStartupProfilerWithLogger(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	profiler := NewStartupProfilerWithLogger(logger)

	if profiler == nil {
		t.Fatal("NewStartupProfilerWithLogger returned nil")
	}

	if !profiler.IsEnabled() {
		t.Error("Profiler should be enabled by default")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
