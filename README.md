# My Agent Skills

Skills I use with my coding agents, free for anyone to use. They follow the open `SKILL.md` format and are written to work with any agent that supports it. So far I use them with Claude Code and Codex.

## Install

Clone the repo and copy a skill folder into your agent's skills directory:

```sh
git clone https://github.com/bragboy/skills.git
cp -R skills/skills/assume-reviewer ~/.claude/skills/   # Claude Code
cp -R skills/skills/assume-reviewer ~/.codex/skills/    # Codex
```

Then invoke it by name, for example `/assume-reviewer` in Claude Code.

## Skills

### Assume Reviewer

For long sessions. The agent becomes the lead: it plans the work, hands the heavy parts to subagents, reviews what they send back, and merges or commits it. This keeps the main conversation small, so it can run for a long time without losing track. Small jobs, like a one-line fix or a quick check, it just does itself. If the agent cannot spawn subagents, it does the work itself in small steps.

At the start it asks you how to work:

- **Landing Changes:** commit each accepted change, merge without committing, stage for you to commit, or ask each time.
- **Isolation:** a separate git worktree per subagent, or one shared working tree.
- **Parallelism:** one task at a time, or several at once when they are independent.
- **Review Depth:** skim the summary, read the diff, or read the diff and run the tests.

It keeps a running list of what is done, in progress and next, and ends with a summary of what landed, what was rejected and what is left.
