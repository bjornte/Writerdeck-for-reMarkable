# Typing-test methodology: why the attempt keeps failing

Standing status: [decisions.md](../decisions.md) **Typing-test strategy is failing**. How we run checks: §13 there. Catalog: [scenario-catalog.md](scenario-catalog.md).

This note is a working theory of **methodological** failure modes for our automated typing tests — not a list of editor bugs. Update it whenever a human find (or a deliberate kill-test) shows a new way the suite can go green while basic editing is still wrong. Prefer revising this file over inventing a second essay.

---

## What this kind of testing is

We are not unit-testing pure functions. We are trying to prove that a remote, layout-dependent text editor behaves like a familiar Mac/Linux typewriter for a person: locate, select, delete, undo, wrap, blank lines, Home/End, and so on. The oracle is human judgment of motion and appearance. The harness sees a thin slice of editor state (character indexes, selection, text length or body, mode, scroll Y, and — since 18 Jul 2026 — soft-wrap `assoc` and painted `caretY`). That slice is still narrower than every writer claim; new fields only help where scenarios assert them.

That mismatch — human claim vs thin observation — is the root risk. Everything below is a way we have papered over it badly.

## Failure modes we keep repeating

### 1. Scoreboard as strategy

We treat scenario count and “critical green” as the product of the work. Growing the list and polishing dialect docs feels like progress. It does not ask whether a deliberately broken editor would fail. So the suite becomes a regression corpus for past patches, not a proof that basic editing works.

### 2. Claims without kill-tests

A claim (“Cmd+Right goes to end of this visual row”) gets a scenario that encodes our current *model* of the claim. We rarely first invent the dumb wrong behavior a writer would notice and demand a red check. Without that kill-test, we optimize the expect until it matches whatever the device already does — including broken behavior that still satisfies the numbers.

### 3. Fixtures that cannot show the bug

Short hard-broken lines make “visual line” and “paragraph” the same place. Ranges that accept almost any cursor near the end hide wrap failures. Content without blank lines cannot catch “skip empty paragraphs.” If the interesting case never appears in the note, the scenario cannot fail for that case. We have fixed some of these after the fact; the method still allows them in.

### 4. Proxies that good and bad both satisfy

We assert a bookkeeping value that two different user outcomes share (same cursor index for “end of this row” and “start of the next”; “not past paragraph end” when the bug is which row the caret prefers; length checks that ignore where the caret sits). The suite then proves a weaker statement than the product claim. Dialect accuracy does not rescue a weak proxy.

**Partial mitigation (18 Jul 2026):** soft-wrap End/Cmd+Right now publishes `assoc` / painted `caretY` and critical round-trips (`wrap-ctrl-right-then-left`, `wrap-end-then-up`). That kills the shared-index false green for that claim. Other claims can still use weak proxies — do not generalize from one fix.

### 5. Observation thinner than the claim

When the product claim is about how motion *looks* or which layout unit moved, and the harness only records indexes (or other collapsed state), right and wrong can be observationally identical. Adding more scenarios inside that observation model only multiplies false confidence. Extending what we record is a methodology step, not a one-off fix for a single chord.

### 6. Reactive scenario patches

A person finds a basic bug → we add or tweak scenarios → green again → another basic bug. That loop never forces a systematic pass over “every basic claim has a discriminating check.” Agents especially prefer the small patch that turns this failure green over questioning whether the whole class of checks can fail.

### 7. Definitions mistaken for detection

Acid-testing vocabulary against Apple/CodeMirror is necessary. It is not sufficient. Clear terms with weak expects, wrong fixtures, or missing observations still go green. We spent effort locking dialect while the detection surface stayed too coarse for the claims we wrote in that dialect.

### 8. Critical tag without fault power

Tagging something `critical` marks importance to us. It does not mean the check can catch the important failures. Inflating the basic set without kill-tests makes “44/44” look like a stronger guarantee than it is.

### 9. Agent incentives

LLM sessions are rewarded for deploy + green harness + tidy docs. Questioning whether the scenario models the user’s failure mode slows that path. Without an explicit kill-test / “would wrong fail?” gate, agents will keep shipping false completion. The handoff history in this folder already recorded that pattern; we have not yet changed the method enough to break it.

## What would count as a better method (sketch)

Not fully implemented yet — progress lives in [basic-claims.md](basic-claims.md):

1. Maintain a short list of **basic writer claims** (outcomes), separate from scenario filenames — **started** in that inventory.
2. For each claim: a **kill-test** — name a broken behavior; require today’s suite to go red for it; if not, mark Unguarded/Partial and fix fixture / assert / observation before trusting green.
3. Ban for critical claims: inert fixtures, expects the bad outcome also satisfies, proxies that do not separate good from bad.
4. When a human finds a green-suite basic bug: first encode that miss here and in a red check; then fix the editor; then scan neighbor claims for the same hole.
5. Treat green only as “these checks passed,” until misses stop recurring under that discipline.

## How to update this note

When we learn a new way the suite lied, or retire a failure mode above with proof, edit this file in the same change. Update [basic-claims.md](basic-claims.md) status rows in the same change. Point from [decisions.md](../decisions.md) **Typing-test strategy is failing**, [todo.md](todo.md), and [lessons.md](../lessons.md). Do not replace this theory with a victory lap until those standing banners are ready to come down.
