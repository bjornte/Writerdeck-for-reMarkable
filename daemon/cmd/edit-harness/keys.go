package main

import (
	"strings"
)

// harnessNamedKeys are browser key names the daemon forwards to Writerdeck.
var harnessNamedKeys = map[string]string{
	"Enter":      "Return",
	"Backspace":  "Backspace",
	"Delete":     "Delete",
	"Tab":        "Tab",
	"Escape":     "Escape",
	"Home":       "Home",
	"End":        "End",
	"ArrowUp":    "ArrowUp",
	"ArrowDown":  "ArrowDown",
	"ArrowLeft":  "ArrowLeft",
	"ArrowRight": "ArrowRight",
}

// harnessSocketNamedKeys are k values accepted by socket-inject (plus A-Z actions).
var harnessSocketNamedKeys = map[string]bool{
	"Escape": true, "Return": true, "Backspace": true, "Delete": true, "Tab": true,
	"Home": true, "End": true,
	"ArrowUp": true, "ArrowDown": true, "ArrowLeft": true, "ArrowRight": true,
}

func (k Key) hasModifier() bool {
	return k.Shift || k.Ctrl || k.Alt || k.Meta
}

func isNavKeyName(name string) bool {
	switch name {
	case "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown", "Home", "End", "Backspace", "Delete":
		return true
	default:
		return false
	}
}

func (k Key) isModifiedNav() bool {
	if !k.hasModifier() || !isNavKeyName(k.Name) {
		return false
	}
	return true
}

func stepHasModifiedNav(step Step) bool {
	for _, k := range step.Keys {
		if k.isModifiedNav() {
			return true
		}
	}
	return false
}

func stepNeedsModifiedPrime(step Step) bool {
	for _, k := range step.Keys {
		// End-prime before Ctrl+Home poisons positioning steps that start at 0.
		if k.Ctrl && k.Name == "Home" && !k.Shift && !k.Alt {
			continue
		}
		if k.isModifiedNav() {
			return true
		}
	}
	return false
}

// validateAllScenarioKeys returns the first invalid key in the suite.
func validateAllScenarioKeys() string {
	for _, sc := range AllScenarios() {
		for _, step := range sc.Steps {
			for _, k := range step.Keys {
				if msg := validateScenarioKey(k); msg != "" {
					return sc.Name + ": " + msg
				}
			}
		}
	}
	return ""
}

func stepExpectedCursor(step Step) *int {
	if step.Expect == nil {
		return nil
	}
	return step.Expect.Cursor
}

// needsExplicitRelease: Ctrl+Shift combos block auto release in socket-inject.
func (k Key) needsExplicitRelease() bool {
	return k.Shift && k.Ctrl && isNavKeyName(k.Name)
}

func validateScenarioKey(k Key) string {
	if k.Name == "" {
		return "empty key name"
	}
	if len([]rune(k.Name)) == 1 {
		return ""
	}
	if _, ok := harnessNamedKeys[k.Name]; ok {
		return ""
	}
	return "unknown key name " + k.Name
}

func inferScenarioTags(name string) []string {
	var tags []string
	switch {
	case strings.HasPrefix(name, "wrap-"):
		tags = append(tags, "wrap")
	case strings.HasPrefix(name, "combo-"):
		tags = append(tags, "combo")
	case strings.HasPrefix(name, "cm-"):
		tags = append(tags, "cm")
	case strings.HasPrefix(name, "gap-"):
		tags = append(tags, "gap")
	case strings.HasPrefix(name, "undo-"), strings.HasPrefix(name, "redo-"):
		tags = append(tags, "undo")
	case strings.HasPrefix(name, "bs-"):
		tags = append(tags, "backspace")
	case name == "load-cursor-at-start" || strings.HasPrefix(name, "shift-") ||
		strings.HasPrefix(name, "home-") || strings.HasPrefix(name, "ctrl-shift-"):
		tags = append(tags, "core")
	default:
		switch name {
		case "down-one-logical-line", "shift-down-then-up-shrinks", "shift-left-repeat-from-end",
			"alt-backspace-deletes-word", "ctrl-backspace-deletes-line", "shift-left-repeat-mid-doc":
			tags = append(tags, "regression")
		}
	}
	return tags
}

func findScenariosByTag(tag string) ([]Scenario, bool) {
	var out []Scenario
	for _, sc := range AllScenarios() {
		for _, t := range sc.Tags {
			if t == tag {
				out = append(out, sc)
				break
			}
		}
	}
	return out, len(out) > 0
}

func attachScenarioTags(scenarios []Scenario) []Scenario {
	out := make([]Scenario, len(scenarios))
	for i, sc := range scenarios {
		if len(sc.Tags) == 0 {
			sc.Tags = inferScenarioTags(sc.Name)
		}
		out[i] = sc
	}
	return out
}
