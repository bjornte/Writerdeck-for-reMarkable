# Manual prompts for cleanup & tersifying at irregular intervals

Guard (applies to every prompt below): preserve durable operational commands (e.g. the dev-loop shortcuts in TODO §2) and durable lessons (what worked *and* what didn't, and why). When in doubt, keep. Prune narration, not knowledge.

Consider your writing style. By default, you are too liberal with text formatting. E.g., too much bold. Rather, imagine that you are John Gruber of Daring Fireball. Consider what his documentation would look like, both semantically and syntactically.

Do the items below one by one. Halt for every work item and wait for my signal before continuing. But do the items in any order you like.

* Assess `DONE.md`. Remove anything that's superfluous in terms of understanding what's done; both what worked and what didn't. Future reader might be LLM or dev.

* Assess `TODO.md`. Remove anything superfluous to understanding the plan — the architecture, the next actionable phase, the open risks, and the constraints to honor. Keep what a fresh reader needs to pick up the first unchecked work; prune settled history, dead paths (beyond a one-line gravestone), and redundant narration. Future reader might be LLM or dev.

* Assess the main `README.md`. Remove anything superfluous to understanding what this project is, why it exists, and how to get started. Keep the purpose, the high-level architecture, and the entry points; cut stale status, deep implementation detail (that lives in TODO/DONE), and duplication. Future reader might be LLM, dev or a casual user from e.g. the writerdeck community, who will need a human friendly intro.

* Assess the other READMEs. Remove anything superfluous

* I'm guessing some scripts, probes, logs etc. might not be interesting in the future, since they have been used to uncover items now documented e.g. in the todo and done files. Might that be true? if yes, prune these files. Aggressive is OK. Update any references to avoid link rot. Consider whether there are opportunities to make it more terse, as it consumes attention / tokens.

* Prune any design decisions that are uninteresting for anyone extending the project later. E.g., the Lobby screen was initially white-on-black. It was then made black-on-white. This is entirely inconsequential for other to know. Nothing to learn from it that will inform anyone extending the app in the future. So identify and prune such information. Aggressive is OK.
