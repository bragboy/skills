---
name: assume-reviewer
description: Run a long session as lead reviewer and context keeper. The main thread plans, delegates heavy work to subagents (doing small things itself), reviews what comes back, and integrates it per the user's chosen mode. Each subagent keeps a small status file the lead can read to answer "what is it doing?" and to spot a stuck agent. Use on /assume-reviewer or when the user asks you to take over a long task and delegate.
---

You are the lead. Keep the main thread small: plan, delegate, review, integrate.

Use judgment on what to delegate. Do it yourself when it is small and cheap in context (a one-line fix, a quick check, a known file, a commit). Delegate when it would pull lots of code or output into the thread, or can run in parallel.

## Start

Ask the user these in one go (use your question tool if you have one):
1. **Integration**: auto-commit each accepted change / merge as-is without commits / stage only, I commit / ask me each time
2. **Isolation**: a git worktree per subagent / shared working tree
3. **Parallelism**: one at a time / parallel when independent
4. **Review depth**: skim the summary / read the diff / diff + run tests

Then confirm the goal in one line and break it into small, independent tasks.

Keep status files out of git: in a git repo, add `.agent-status/` to `.git/info/exclude` (shared by every worktree of the repo) if it is not ignored already.

## Loop

1. Delegate each task to a subagent at **medium** effort (pass `effort: "medium"`), unless the user asks for another level for that task. Prompt: goal, relevant paths, constraints, done criteria, and the **status file rules** below (paste them in). Ask for a short report back: files changed, what was verified, open issues. No file dumps.
2. Review per the chosen depth. Reject or send back with concrete fixes if wrong.
3. Integrate per the chosen mode. Commit messages: one line of intent, no attribution.
4. Keep a running ledger in the thread: done / in progress / next / decisions. Update it after every task, so it survives context compaction. Note each running agent's status file path in the ledger.

## Status file (paste into every subagent prompt)

> Keep a status file at `.agent-status/<task-name>.md` at the root of your working tree (your worktree, or the shared tree), creating the folder if needed. It is for the lead, not for you: keep it short and overwrite it, never append.
> - Start it with your own plan as numbered steps. You choose the steps and may change them as you learn more (add, split, merge, drop); the total is whatever it is now.
> - Update it when you start, whenever a step finishes or the plan changes, and before and after any command that may take more than a minute.
> - Format:
>   ```
>   updated: <ISO time>   started: <ISO time>
>   step: 7 of 12 - <what step 7 is>, running since <ISO time>
>   steps:
>   - [x] 1. <step> (4 min)
>   - [x] 2. <step> (11 min)
>   - [ ] 7. <step> (running)
>   - [ ] 8. <step>
>   eta: about <N> min left (<one line on why, e.g. "the test suite takes 8 min">)
>   waiting on: <command or process, with PID and log path if it runs in the background, else nothing>
>   blocked: <no, or what blocks you>
>   ```
>   Take every time from the clock (`date -u +%FT%TZ`, or `TZ=Asia/Kolkata date` for local time) when you write it; never estimate or invent a timestamp. A step's duration is the difference between two clock readings. Each finished step carries how long it took. The estimate is your honest guess from the steps so far and what is left; revise it at every update.
> - Never leave a background process running unattended: give long commands a timeout, and record background PIDs in the file. Before you finish, stop every process you started.
> - Delete the file just before you send your final report.

## When the user asks what an agent is doing

Read its status file and answer from it in one or two lines: which step of how many it is on, what that step is and how long it has been running, and the estimate left. For example: "Step 7 of 12, porting the scroll animations, running 6 min; about 25 min left." Do not message the agent for a progress report, and do not read its transcript. If the file is missing while the agent is still running, say so and fall back to cheap signals: recent commits and recently modified files in its worktree.

## Stuck agents

You cannot see inside a subagent, so never assume "still running" means "still working". Check whenever the user asks about an agent, and at least whenever you are otherwise idle with agents running:

- **Fresh:** `updated:` within about 15 minutes, or within the time a named long command needs. Leave it alone.
- **Stale:** older than that with no new commits or file changes in its worktree. Check what it says it is waiting on: is the PID alive, is its log still growing? A dead process or a log that stopped growing means it is stuck.
- **Stuck:** send the agent one short message naming what you see (for example "your status says waiting on PID 4242, which exited 20 minutes ago"). If nothing changes within a few minutes, stop the agent, kill any processes it left behind, keep its committed work, and relaunch the remaining part with a fresh prompt that says what was already done.
- Tell the user plainly when you stopped or relaunched an agent and why. Never report "waiting for the agent" for an agent whose status is stale.

## Rules

- Delegate broad searches and keep only conclusions. A single known lookup you can do yourself.
- If you cannot spawn subagents, do the work yourself in small steps and keep the ledger anyway.
- Stop and ask the user only for real decisions or blockers, not progress check-ins.
- At the end: summary of what landed, what was rejected, and what is left. Remove any status files and processes left behind by finished or stopped agents.
