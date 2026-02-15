// Package performance provides startup profiling and performance monitoring utilities.
package performance

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

// StartupPhase represents a single phase of the startup sequence
type StartupPhase struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Error     error
}

// Completed returns true if the phase has completed
func (p *StartupPhase) Completed() bool {
	return !p.EndTime.IsZero()
}

// StartupProfiler tracks timing for all startup phases
type StartupProfiler struct {
	mu        sync.RWMutex
	startTime time.Time
	endTime   time.Time
	phases    []StartupPhase
	phaseMap  map[string]*StartupPhase
	logger    *log.Logger
	enabled   bool
}

// NewStartupProfiler creates a new startup profiler
func NewStartupProfiler() *StartupProfiler {
	return &StartupProfiler{
		startTime: time.Now(),
		phases:    make([]StartupPhase, 0),
		phaseMap:  make(map[string]*StartupPhase),
		logger:    log.New(log.Writer(), "[STARTUP] ", log.LstdFlags|log.Lmicroseconds),
		enabled:   true,
	}
}

// NewStartupProfilerWithLogger creates a profiler with a custom logger
func NewStartupProfilerWithLogger(logger *log.Logger) *StartupProfiler {
	return &StartupProfiler{
		startTime: time.Now(),
		phases:    make([]StartupPhase, 0),
		phaseMap:  make(map[string]*StartupPhase),
		logger:    logger,
		enabled:   true,
	}
}

// StartPhase begins timing a new startup phase
func (sp *StartupProfiler) StartPhase(name string) *StartupPhase {
	if !sp.enabled {
		return nil
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	phase := &StartupPhase{
		Name:      name,
		StartTime: time.Now(),
	}

	sp.phaseMap[name] = phase
	sp.phases = append(sp.phases, *phase)

	sp.logger.Printf("▶ %s", name)
	return phase
}

// EndPhase ends timing for a startup phase
func (sp *StartupProfiler) EndPhase(name string, err error) {
	if !sp.enabled {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	phase, exists := sp.phaseMap[name]
	if !exists {
		sp.logger.Printf("⚠ EndPhase called for unknown phase: %s", name)
		return
	}

	phase.EndTime = time.Now()
	phase.Duration = phase.EndTime.Sub(phase.StartTime)
	phase.Error = err

	status := "✓"
	if err != nil {
		status = "✗"
	}

	sp.logger.Printf("%s %s (%s)", status, name, formatDuration(phase.Duration))
}

// EndPhaseWithResult ends a phase and logs a custom result
func (sp *StartupProfiler) EndPhaseWithResult(name string, result string, err error) {
	if !sp.enabled {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	phase, exists := sp.phaseMap[name]
	if !exists {
		sp.logger.Printf("⚠ EndPhase called for unknown phase: %s", name)
		return
	}

	phase.EndTime = time.Now()
	phase.Duration = phase.EndTime.Sub(phase.StartTime)
	phase.Error = err

	status := "✓"
	if err != nil {
		status = "✗"
	}

	sp.logger.Printf("%s %s (%s) - %s", status, name, formatDuration(phase.Duration), result)
}

// TimeFunc wraps a function call with timing
func (sp *StartupProfiler) TimeFunc(name string, fn func() error) error {
	sp.StartPhase(name)
	err := fn()
	sp.EndPhase(name, err)
	return err
}

// TimeFuncWithResult wraps a function call with timing and returns a result string
func (sp *StartupProfiler) TimeFuncWithResult(name string, fn func() (string, error)) error {
	sp.StartPhase(name)
	result, err := fn()
	sp.EndPhaseWithResult(name, result, err)
	return err
}

// MarkComplete marks the entire startup as complete
func (sp *StartupProfiler) MarkComplete() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sp.endTime = time.Now()
	totalDuration := sp.endTime.Sub(sp.startTime)

	sp.logger.Printf("═ Startup complete: %s", formatDuration(totalDuration))
}

// GetTotalDuration returns the total startup duration
func (sp *StartupProfiler) GetTotalDuration() time.Duration {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.endTime.IsZero() {
		return time.Since(sp.startTime)
	}
	return sp.endTime.Sub(sp.startTime)
}

// GetPhase returns a specific phase by name
func (sp *StartupProfiler) GetPhase(name string) (StartupPhase, bool) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	phase, exists := sp.phaseMap[name]
	if !exists {
		return StartupPhase{}, false
	}
	return *phase, true
}

// GetAllPhases returns all phases sorted by start time
func (sp *StartupProfiler) GetAllPhases() []StartupPhase {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	phases := make([]StartupPhase, len(sp.phases))
	copy(phases, sp.phases)

	// Sort by start time
	sort.Slice(phases, func(i, j int) bool {
		return phases[i].StartTime.Before(phases[j].StartTime)
	})

	return phases
}

// GenerateReport creates a formatted startup report
func (sp *StartupProfiler) GenerateReport() string {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	var sb strings.Builder
	totalDuration := sp.GetTotalDuration()

	sb.WriteString("\n")
	sb.WriteString("╔════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║           STARTUP PERFORMANCE REPORT                   ║\n")
	sb.WriteString("╠════════════════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ Total Time: %-42s ║\n", formatDuration(totalDuration)))
	sb.WriteString("╠════════════════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ %-20s %-12s %s ║\n", "Phase", "Duration", "Status"))
	sb.WriteString("╠════════════════════════════════════════════════════════╣\n")

	// Sort phases by duration (descending) to show slowest first
	sortedPhases := make([]StartupPhase, len(sp.phases))
	copy(sortedPhases, sp.phases)
	sort.Slice(sortedPhases, func(i, j int) bool {
		return sortedPhases[i].Duration > sortedPhases[j].Duration
	})

	for _, phase := range sortedPhases {
		status := "OK"
		if phase.Error != nil {
			status = "ERR"
		}
		sb.WriteString(fmt.Sprintf("║ %-20s %-12s %s ║\n",
			truncateString(phase.Name, 20),
			formatDuration(phase.Duration),
			status))
	}

	sb.WriteString("╚════════════════════════════════════════════════════════╝\n")

	return sb.String()
}

// PrintReport logs the startup report
func (sp *StartupProfiler) PrintReport() {
	if !sp.enabled {
		return
	}
	sp.logger.Print(sp.GenerateReport())
}

// Enable enables profiling
func (sp *StartupProfiler) Enable() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.enabled = true
}

// Disable disables profiling
func (sp *StartupProfiler) Disable() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.enabled = false
}

// IsEnabled returns true if profiling is enabled
func (sp *StartupProfiler) IsEnabled() bool {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.enabled
}

// Reset clears all timing data and restarts the clock
func (sp *StartupProfiler) Reset() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sp.startTime = time.Now()
	sp.endTime = time.Time{}
	sp.phases = make([]StartupPhase, 0)
	sp.phaseMap = make(map[string]*StartupPhase)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// AsyncProfiler provides async initialization profiling
type AsyncProfiler struct {
	profiler *StartupProfiler
	name     string
	done     chan struct{}
	result   string
	err      error
}

// StartAsync begins an async initialization
func (sp *StartupProfiler) StartAsync(name string) *AsyncProfiler {
	ap := &AsyncProfiler{
		profiler: sp,
		name:     name,
		done:     make(chan struct{}),
	}
	sp.StartPhase(name)
	return ap
}

// Complete marks an async initialization as complete
func (ap *AsyncProfiler) Complete(result string, err error) {
	ap.result = result
	ap.err = err
	ap.profiler.EndPhaseWithResult(ap.name, result, err)
	close(ap.done)
}

// Wait waits for async initialization to complete
func (ap *AsyncProfiler) Wait(ctx context.Context) error {
	select {
	case <-ap.done:
		return ap.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitWithTimeout waits for async initialization with a timeout
func (ap *AsyncProfiler) WaitWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ap.Wait(ctx)
}

// IsDone returns true if the async operation is complete
func (ap *AsyncProfiler) IsDone() bool {
	select {
	case <-ap.done:
		return true
	default:
		return false
	}
}
