For use in VS Code only, ***NOT*** in Cursor

## Model usage protocol (Opus <-> Sonnet) — the operating rule

Match the model to the task, and proactively tell the user when to switch — state the trigger and what the next model should do first. Never switch silently; the user drives the model picker. When switching models within the same Claude session/chat, the context is preserved. Only the model changes; the conversation history/context remains.

- Opus — plan & inspect. Research, architecture, risk analysis, writing/expanding the plans (TODO/DONE), reviewing & inspecting code, diagnosing non-obvious failures.
- Sonnet — write code. Once a phase's spec is settled, author the files: Dockerfile, build/install scripts, Go source, systemd units, keymap tables, boilerplate, repetitive edits.


| Switch    | When                                                                                                                                 |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| -> Sonnet | The active phase's Done-when spec is settled and the remaining work is typing code/files.                                            |
| -> Opus   | Code needs review/inspection; or a decision / ambiguity / architecture / risk question arises; or a test fails in a non-obvious way. |


Example: *"The socket protocol + keywriter input-source patch are specced — prudent to switch to Sonnet to write the Go feeder + the keywriter patch; have it read the active TODO.md slice + [docs/decisions.md](../docs/decisions.md) first. Switch back to Opus to inspect the result."*