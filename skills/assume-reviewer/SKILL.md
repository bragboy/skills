---
name: assume-reviewer
description: Run a long session as lead reviewer and context keeper. The main thread plans, delegates heavy work to subagents (doing small things itself), reviews what comes back, and integrates it per the user's chosen mode. Use on /assume-reviewer or when the user asks you to take over a long task and delegate.
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

## Loop

1. Delegate each task to a subagent. Prompt: goal, relevant paths, constraints, done criteria. Ask for a short report back: files changed, what was verified, open issues. No file dumps.
2. Review per the chosen depth. Reject or send back with concrete fixes if wrong.
3. Integrate per the chosen mode. Commit messages: one line of intent, no attribution.
4. Keep a running ledger in the thread: done / in progress / next / decisions. Update it after every task, so it survives context compaction.

## Rules

- Delegate broad searches and keep only conclusions. A single known lookup you can do yourself.
- If you cannot spawn subagents, do the work yourself in small steps and keep the ledger anyway.
- Stop and ask the user only for real decisions or blockers, not progress check-ins.
- At the end: summary of what landed, what was rejected, and what is left.
