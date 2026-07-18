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
		// Shift-only arrows extend selection; End-prime jumps to EOF first (breaks shift-right-from-home).
		if k.isModifiedNav() && (k.Alt || k.Ctrl) {
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
// Escape toggles edit/preview on key-up in Writerdeck, so send an explicit release.
func (k Key) needsExplicitRelease() bool {
	if k.Name == "Escape" {
		return true
	}
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

// criticalScenarios are keyboard behaviors that must work for basic editing.
// Sourced from Microsoft/Obsidian/macOS text-editing conventions: plain and
// shift+arrow navigation, backspace/delete, enter, select-all, typing over
// selection, word/line delete, doc home/end, undo/redo, and copy/cut/paste.
var criticalScenarios = map[string]bool{
	"load-cursor-at-start": true,
	"home-clears-selection": true,
	"shift-right-from-home": true,
	"shift-left-from-end": true,
	"shift-right-after-home-no-stale-anchor": true,
	"shift-down-after-arrow-down":            true,
	"shift-up-after-arrow-down":              true,
	"ctrl-shift-left-select-line":            true,
	"down-one-logical-line":                  true,
	"shift-left-repeat-from-end":             true,
	"alt-backspace-deletes-word":             true,
	"ctrl-backspace-deletes-line":            true,
	"shift-left-repeat-mid-doc":              true,
	"cm-line-down-basic":                     true,
	"cm-line-down-last-line":                 true,
	"combo-alt-left":                         true,
	"combo-alt-right":                        true,
	"combo-ctrl-home":                        true,
	"combo-ctrl-end":                         true,
	"bs-plain":                               true,
	"wrap-down-one-visual-line":              true,
	"wrap-up-from-visual-line-2":             true,
	"wrap-ctrl-left":                         true,
	"wrap-ctrl-right":                        true,
	"wrap-ctrl-right-then-left":              true,
	"wrap-end-then-up":                       true,
	"wrap-combo-ctrl-bs-line":                true,
	"wrap-down-goal-column":                  true,
	"wrap-shift-ctrl-left":                   true,
	"wrap-shift-ctrl-right":                  true,
	"cm-line-down-goal-col":                  true,
	"combo-alt-up":                           true,
	"combo-alt-down":                         true,
	"combo-alt-up-double-blank":              true,
	"combo-alt-down-double-blank":            true,
	"combo-alt-up-prose-double-blank":        true,
	"undo-redo-len":                          true,
	"gap-up-at-doc-start":                    true,
	"gap-plain-left-moves-caret":              true,
	"gap-plain-right-moves-caret":             true,
	"gap-collapse-selection-left":            true,
	"gap-collapse-selection-right":           true,
	"gap-delete-forward":                     true,
	"gap-delete-with-selection":              true,
	"gap-select-all":                         true,
	"gap-copy-paste":                         true,
	"gap-cut-paste":                          true,
	"gap-enter-new-line":                     true,
	"gap-type-replaces-selection":            true,
	"gap-redo-shift-ctrl-z":                  true,
	"gap-undo-chain":                         true,
	"gap-empty-doc-backspace":                true,
	// Mid-sentence Shift+vertical on wrapping paragraphs (not short \n lines).
	"gap-shift-down-mid-wrapping-paras": true,
	"gap-shift-up-mid-wrapping-paras":   true,
}

func isCriticalScenario(name string) bool {
	return criticalScenarios[name]
}

func inferScenarioTags(name string) []string {
	var tags []string
	if isCriticalScenario(name) {
		tags = append(tags, "critical")
	}
	switch {
	case strings.HasPrefix(name, "wrap-"):
		tags = append(tags, "wrap")
	case strings.HasPrefix(name, "combo-"):
		tags = append(tags, "combo")
	case strings.HasPrefix(name, "cm-"):
		tags = append(tags, "cm")
	case strings.HasPrefix(name, "gap-"):
		tags = append(tags, "gap")
	case strings.HasPrefix(name, "hw-"):
		tags = append(tags, "hw")
	case strings.HasPrefix(name, "read-"):
		tags = append(tags, "read")
	case strings.HasPrefix(name, "touch-"):
		tags = append(tags, "touch")
	case strings.HasPrefix(name, "shift-left-then-right"), strings.HasPrefix(name, "shift-right-then-left"),
		strings.HasPrefix(name, "shift-up-then-down"):
		tags = append(tags, "selection")
	case strings.HasPrefix(name, "undo-"), strings.HasPrefix(name, "redo-"):
		tags = append(tags, "undo")
	case strings.HasPrefix(name, "bs-"), strings.HasPrefix(name, "delete-"):
		tags = append(tags, "backspace")
	case name == "load-cursor-at-start" || strings.HasPrefix(name, "shift-") ||
		strings.HasPrefix(name, "home-") || strings.HasPrefix(name, "ctrl-shift-"):
		tags = append(tags, "core")
	default:
		switch name {
		case "down-one-logical-line", "up-one-logical-line", "shift-down-then-up-shrinks", "shift-left-repeat-from-end",
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
		} else if isCriticalScenario(sc.Name) {
			sc.Tags = appendCriticalTag(sc.Tags)
		}
		out[i] = sc
	}
	return out
}

func appendCriticalTag(tags []string) []string {
	for _, t := range tags {
		if t == "critical" {
			return tags
		}
	}
	return append([]string{"critical"}, tags...)
}
