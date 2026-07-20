# Manual prompts for cleanup & tersifying at irregular intervals

Main goal (applies to every prompt below): Make all documentation more terse while still making it enable the project proceed forward in a prudent and high quality manner. Preserve durable operational commands and durable lessons (what worked *and* what didn't, and why). When in doubt, keep. Remove any chatty references to previous, uninteresting history. Prune narration, not valuable, forward-looking knowledge.

Also, consider your writing style. By default, you are too liberal with text formatting. E.g., too much bold. Rather, imagine that you are John Gruber of Daring Fireball. Consider what his documentation would look like, both semantically and syntactically.

* Assess `DONE.md`. Remove everything that's superfluous in terms of understanding what's done; both what worked and what didn't. Future reader might be LLM or dev.

* Assess `TODO.md`. Remove everything superfluous to understanding the plan — the architecture, the next actionable phase, the open risks, and the constraints to honor. Keep what a fresh reader needs to pick up the first unchecked work; prune settled history, dead paths (beyond a one-line gravestone), and redundant narration. Future reader might be LLM or dev.

* Assess the other non-main READMEs. Remove everything superfluous.

* Assess the rest of the documentation. Remove everything superfluous.

* I'm guessing some scripts, probes, logs etc. might not be interesting in the future, since they have been used to uncover items now documented e.g. in the todo and done files. Might that be true? if yes, prune these files. Aggressive is OK. Update any references to avoid link rot. Consider whether there are opportunities to make it more terse, as it consumes attention / tokens.

* Prune any design decisions that are uninteresting for anyone extending the project later. E.g., the Lobby screen was initially white-on-black. It was then made black-on-white. This is entirely inconsequential for other to know. Nothing to learn from it that will inform anyone extending the app in the future. So identify and prune such information. Aggressive is OK.
