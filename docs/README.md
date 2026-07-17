# docs/

Reference for the project. Test script output goes in `recon/` — gitignored text files, safe to delete; lessons belong in DONE, TODO, and lessons.

`architecture.md` — how the system works, device facts, deploy loop. Document integrity is the product contract.

`decisions.md` — why each choice was made.

`lessons.md` — operational gotchas.

`browser-vs-tablet.md` — phone vs tablet; upload vs download vs paste-at-cursor.

`integrity-audit.md` — integrity status and open risks.

`improvements.md` — design notes for encryption, locales, future UX.

`user-should-test.md` — checklist of what the owner should try by hand (scripts cannot cover these).

`todo-handoff-physical-home-input.md` — Physical Home exclusive gpio grab (**done**; [decisions.md](decisions.md) §28).

`todo-install-onboarding.md` — visitor install friction: README, releases, preflight/install scripts.

`editor-migration-1-to-QML/` — keywriter fork migration (**done**; [todo-handoff-keywriter-fork.md](editor-migration-1-to-QML/todo-handoff-keywriter-fork.md)); ownership in [decisions.md](decisions.md) §3.

`editor-migration-2-to-cpp/` — QML edit helpers → C++ `EditHelper` (**Phase A+B done**; Phase C optional in [todo-handoff-edit-helper-cpp.md](editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md)).

`editor-testing/` — scenario catalog ([scenario-catalog.md](editor-testing/scenario-catalog.md)), harness handoff ([todo.md](editor-testing/todo.md)), scoreboard ([milestone-runs.md](editor-testing/milestone-runs.md)), run history ([recon/harness-runs.md](recon/harness-runs.md)).

`server-sync-implementation.md` — shipped server-side GitHub sync reference (token restore, `needtoken`, log verification).

`manual-cleanup-tersify-prompts-at-intervals.md` — periodic doc cleanup prompts.

External: [keywriter](https://github.com/dps/remarkable-keywriter) (Qt 5 / C++ / QML editor engine; we ship it as Writerdeck) · [input docs](https://remarkable.guide/devel/device/input.html) · [awesome-reMarkable](https://github.com/reHackable/awesome-reMarkable)
