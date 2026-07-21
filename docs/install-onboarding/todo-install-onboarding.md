# Install and onboarding — follow-up

Visitor install (ZIP, prompts, Releases, health-check then enable) is shipped. Overview: [README.md](README.md).

## Still open

- [ ] Inspect dangers of bricking on boot now that `--start` auto-enables `writerdeck` after a health check. Confirm the check is strong enough, document failure modes, and whether enable should stay gated (human glance at e-ink, stricter checks, or opt-in flag).
- [ ] Windows installer — see [TODO.md](../../TODO.md).

## Out of scope

- reMarkable 2 — only if the community wants it ([../decisions.md](../decisions.md) §33).
- Toltec or jailbreak-based distribution.
- Phone-only install with no computer.
