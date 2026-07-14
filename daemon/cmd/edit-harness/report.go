package main

import (
	"fmt"
	"strings"
	"time"
)

type outcomeKind int

const (
	outcomePass outcomeKind = iota
	outcomeFail
	outcomePrepareFail
)

type scenarioResult struct {
	Name             string
	Kind             outcomeKind
	Err              string
	Duration         time.Duration
	ResetDuration    time.Duration // editor quit/relaunch attributed to this scenario
	PrepareRecovered bool
	ContaminatedBy   string
	PossiblePoisoner bool
	HealthNotes      []string
}

type runMeta struct {
	StartedAt    time.Time
	Target       string
	Mode         string
	Fast         bool
	ScenarioCount int
	SetupDuration time.Duration // one cold start before first scenario (soft suite)
}

func (r scenarioResult) label() string {
	switch r.Kind {
	case outcomePass:
		return "PASS"
	case outcomePrepareFail:
		return "PREPARE_FAIL"
	default:
		return "FAIL"
	}
}

func (r scenarioResult) resultCell() string {
	switch r.Kind {
	case outcomePass:
		return "pass"
	case outcomePrepareFail:
		return "prepare fail"
	default:
		return "fail"
	}
}

func (r scenarioResult) comments() string {
	var parts []string
	if r.Err != "" {
		parts = append(parts, r.Err)
	}
	if r.PrepareRecovered {
		parts = append(parts, "auto hard-reset before run")
	}
	if r.ResetDuration > 0 {
		parts = append(parts, fmt.Sprintf("editor reset %.1fs", r.ResetDuration.Seconds()))
	}
	if r.ContaminatedBy != "" {
		parts = append(parts, "cascade suspect after "+r.ContaminatedBy)
	}
	if r.PossiblePoisoner {
		parts = append(parts, "may have poisoned next scenario")
	}
	parts = append(parts, r.HealthNotes...)
	return strings.Join(parts, "; ")
}

func escapeMDCell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// linkContamination marks prepare failures and recovery-after-failure patterns.
func linkContamination(results []scenarioResult) {
	for i := 1; i < len(results); i++ {
		cur := &results[i]
		prev := results[i-1]
		if cur.Kind == outcomePrepareFail {
			if prev.Kind == outcomeFail || prev.Kind == outcomePrepareFail {
				cur.ContaminatedBy = prev.Name
				if prev.Kind == outcomeFail {
					results[i-1].PossiblePoisoner = true
				}
			}
			continue
		}
		if cur.PrepareRecovered && prev.Kind == outcomeFail {
			cur.ContaminatedBy = prev.Name
			results[i-1].PossiblePoisoner = true
		}
	}
}

func formatContaminationReport(results []scenarioResult) string {
	linkContamination(results)
	var lines []string
	for _, r := range results {
		if r.PossiblePoisoner {
			lines = append(lines, fmt.Sprintf("POISON_SUSPECT %s (next scenario could not prepare cleanly)", r.Name))
		}
		if r.ContaminatedBy != "" {
			lines = append(lines, fmt.Sprintf("CASCADE_SUSPECT %s (prepare failed after %s failed)", r.Name, r.ContaminatedBy))
		}
		if r.Kind == outcomePrepareFail && r.PrepareRecovered {
			lines = append(lines, fmt.Sprintf("PREPARE_RETRIES %s (prepare failed after retries)", r.Name))
		}
		if r.PrepareRecovered && r.Kind != outcomePrepareFail {
			lines = append(lines, fmt.Sprintf("PREPARE_RETRY %s (needed %d prepare retries; result=%s)", r.Name, 1, r.label()))
		}
	}
	if len(lines) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("=== contamination report ===\n")
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString("Re-check POISON_SUSPECT / CASCADE_SUSPECT with: bash scripts/test-keyboard-harness.sh -s NAME --fast\n")
	return b.String()
}

func formatMarkdownReport(meta runMeta, results []scenarioResult) string {
	linkContamination(results)
	var b strings.Builder
	b.WriteString("# Keyboard harness results\n\n")
	b.WriteString(fmt.Sprintf("Run: %s\n\n", meta.StartedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("Target: `%s`\n\n", meta.Target))
	b.WriteString(fmt.Sprintf("Mode: %s\n\n", meta.Mode))
	if meta.Fast {
		b.WriteString("Timing: fast pauses\n\n")
	}
	if meta.SetupDuration > 0 {
		b.WriteString(fmt.Sprintf("Suite setup: %.1fs (one cold start, included in first scenario time)\n\n", meta.SetupDuration.Seconds()))
	}
	pass, fail, prep := 0, 0, 0
	var total time.Duration
	for _, r := range results {
		total += r.Duration
		switch r.Kind {
		case outcomePass:
			pass++
		case outcomePrepareFail:
			prep++
		default:
			fail++
		}
	}
	b.WriteString(fmt.Sprintf("Summary: %d pass, %d fail, %d prepare fail; total %.1fs\n\n", pass, fail, prep, total.Seconds()))
	b.WriteString("| Name | Result | Time (s) | Recovery | Cascade | Comments |\n")
	b.WriteString("|------|--------|----------|----------|---------|----------|\n")
	for _, r := range results {
		recovery := "no"
		if r.PrepareRecovered || r.ResetDuration > 0 {
			recovery = "yes"
		}
		cascade := "—"
		if r.ContaminatedBy != "" {
			cascade = r.ContaminatedBy
		} else if r.PossiblePoisoner {
			cascade = "suspect"
		}
		comments := r.comments()
		if comments == "" {
			comments = "—"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %.1f | %s | %s | %s |\n",
			escapeMDCell(r.Name),
			r.resultCell(),
			r.Duration.Seconds(),
			recovery,
			escapeMDCell(cascade),
			escapeMDCell(comments),
		))
	}
	return b.String()
}
