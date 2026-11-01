---
name: compact-prep
description: Prepare a long session for context compaction. Save the current working state to handoff Markdown files in the repo, save only lasting critical facts to memory, ask a few questions if anything is unclear, then put a continuation prompt on the clipboard and tell the user to compact and paste it. Use only when the user invokes /compact-prep. Can be run many times in one session.
---

The user is about to compact this session. Your job is to make sure nothing important is lost. Do not track context size or suggest running this yourself. It runs only when the user asks for it.

## 1. Take stock

From the conversation so far, collect:
- **Goal**: what the session is trying to achieve, in one or two lines.
- **State**: what is done, in progress, and next (concrete steps, not vague themes).
- **Decisions**: what was decided and why, including approaches that were tried and rejected.
- **Gotchas**: things that cost time to find out: odd commands, flaky tests, env quirks, paths, IDs, URLs.
- **User corrections and preferences** given during this session.
- **Loose ends**: uncommitted changes (`git status`), background jobs or subagents still running, open questions.

## 2. Ask, only if needed

If anything is unclear, ask at most 3-4 short questions in one go (use your question tool if you have one). Ask only about things that change what you save, for example:
- an open decision that the handoff would have to guess at
- uncommitted changes or running jobs: should they be committed, stashed, waited on, or just noted?
- whether something from the session is still relevant or was abandoned

If nothing is unclear, skip this step. Never commit, stash or kill anything on your own. Note it in the handoff instead.

## 3. Write the handoff files

Work in `.claude/handoff/` at the repo root (or the working directory if this is not a git repo). Create it if needed. In a git repo, keep it out of git without touching tracked files: add `.claude/handoff/` to `.git/info/exclude` if it is not already ignored (check with `git check-ignore`).

**`HANDOFF.md`**: rewrite it completely each run, so it always shows the current state. Keep it dense, written for an AI picking up the work cold:

```markdown
# Handoff (updated <date time>)

## Critical
<lasting facts from step 4 if your agent has no memory system; else omit>

## Goal
## State
- Done:
- In progress:
- Next: (ordered, concrete)
## Key Files and Commands
## Gotchas
## Loose Ends
<uncommitted changes, running jobs, open questions, with what the user said to do>
```

**`decisions.md`**: append only, never rewrite. Add each new decision since the last run as a dated bullet with a one-line reason. Read the file first so you do not repeat entries.

Reference files and code by path. Do not paste large code blocks or logs.

## 4. Save to memory, sparingly

Memory loads into every future session, so save only lasting facts: user preferences, corrections the user gave you, project constraints that are not written down anywhere else. Task progress, next steps and session details go in the handoff files only, never in memory.

Before saving, check whether an existing memory already covers it and update that instead of adding a duplicate. Most runs save nothing new to memory, and that is fine.

If your agent has no memory system, put these facts in the "Critical" section of `HANDOFF.md` instead.

## 5. Clipboard and hand back

Write the continuation prompt the user will paste after compaction. It must still work if the handoff file is only skimmed:

```
Continuing after compaction. First read .claude/handoff/HANDOFF.md (and decisions.md if needed).
Goal: <one line>
Next: <the immediate next step>
Watch out: <2-3 most important gotchas or loose ends>
```

Copy it to the clipboard: `pbcopy` on macOS, `wl-copy` or `xclip -selection clipboard` on Linux, `clip.exe` on Windows/WSL. If none works, say so.

Then end with a short message:
1. What was saved: the handoff files, and any memory entries added or updated (one line each).
2. Suggested command, with a one-line focus hint taken from the goal, for example: `/compact focus on the auth refactor and its next steps`.
3. "After compacting, paste the clipboard to continue." Also show the continuation prompt in the message, in case the clipboard copy failed.

Do not start new work after this. Stop and wait for the user.

## After compaction

When the user pastes the continuation prompt, read `HANDOFF.md`, check that its loose ends still hold (for example, run `git status`), and continue with the next step.
