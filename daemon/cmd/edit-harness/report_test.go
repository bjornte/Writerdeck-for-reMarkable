package main

import (
	"strings"
	"testing"
	"time"
)

func TestLinkContamination(t *testing.T) {
	results := []scenarioResult{
		{Name: "a", Kind: outcomeFail},
		{Name: "b", Kind: outcomePrepareFail},
		{Name: "c", Kind: outcomePass},
	}
	linkContamination(results)
	if results[1].ContaminatedBy != "a" {
		t.Fatalf("ContaminatedBy = %q want a", results[1].ContaminatedBy)
	}
	if !results[0].PossiblePoisoner {
		t.Fatal("expected a marked possible poisoner")
	}
	if results[2].ContaminatedBy != "" {
		t.Fatal("pass after prepare fail should not be linked")
	}
}
func TestLinkContaminationRecovery(t *testing.T) {
	results := []scenarioResult{
		{Name: "a", Kind: outcomeFail},
		{Name: "b", Kind: outcomePass, PrepareRecovered: true},
	}
	linkContamination(results)
	if results[1].ContaminatedBy != "a" {
		t.Fatalf("ContaminatedBy = %q want a", results[1].ContaminatedBy)
	}
	if !results[0].PossiblePoisoner {
		t.Fatal("expected a marked possible poisoner after recovery")
	}
}

func TestFormatContaminationReportEmpty(t *testing.T) {
	if got := formatContaminationReport([]scenarioResult{{Name: "ok", Kind: outcomePass}}); got != "" {
		t.Fatalf("unexpected report: %q", got)
	}
}

func TestFormatMarkdownReport(t *testing.T) {
	meta := runMeta{
		StartedAt: time.Date(2026, 7, 14, 19, 0, 0, 0, time.UTC),
		Target:    "tablet:8000",
		Mode:      "soft-reset (single launch)",
	}
	results := []scenarioResult{
		{Name: "ok", Kind: outcomePass, Duration: 1500 * time.Millisecond},
		{Name: "bad", Kind: outcomeFail, Duration: 2 * time.Second, Err: "cursor want 3 got 5", PossiblePoisoner: true},
		{Name: "next", Kind: outcomePrepareFail, Duration: 3 * time.Second, Err: "textLen mismatch", ContaminatedBy: "bad", PrepareRecovered: true},
	}
	got := formatMarkdownReport(meta, results)
	for _, want := range []string{
		"| Name | Result | Time (s) | Recovery | Cascade | Comments |",
		"| ok | pass | 1.5 | no | — | — |",
		"| bad | fail | 2.0 | no | suspect |",
		"| next | prepare fail | 3.0 | yes | bad |",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("report missing %q:\n%s", want, got)
		}
	}
}

func TestFormatContaminationReport(t *testing.T) {
	results := []scenarioResult{
		{Name: "bad", Kind: outcomeFail},
		{Name: "next", Kind: outcomePrepareFail, PrepareRecovered: true},
	}
	got := formatContaminationReport(results)
	for _, want := range []string{"POISON_SUSPECT bad", "CASCADE_SUSPECT next"} {
		if !strings.Contains(got, want) {
			t.Fatalf("report missing %q:\n%s", want, got)
		}
	}
}

func TestEscapeMDCell(t *testing.T) {
	if got := escapeMDCell("a|b"); got != "a\\|b" {
		t.Fatalf("got %q", got)
	}
}
