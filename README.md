# skills

My Claude Code skills, free for anyone to use.

## Install

Copy a skill folder into `~/.claude/skills/` (all projects) or `.claude/skills/` (one project):

```sh
git clone https://github.com/bragboy/skills.git
cp -R skills/skills/assume-reviewer ~/.claude/skills/
```

Then run it in Claude Code with `/<skill-name>`.

## Skills

### assume-reviewer

For long sessions. Claude becomes the lead: it plans the work, hands the heavy parts to subagents, reviews what they send back, and merges or commits it. This keeps the main conversation small, so it can run for a long time without losing track. Small jobs, like a one-line fix or a quick check, it just does itself.

At the start it asks you how to work:

- **Landing changes:** commit each accepted change, merge without committing, stage for you to commit, or ask each time.
- **Isolation:** a separate git worktree per subagent, or one shared working tree.
- **Parallelism:** one task at a time, or several at once when they are independent.
- **Review depth:** skim the summary, read the diff, or read the diff and run the tests.

It keeps a running list of what is done, in progress and next, and ends with a summary of what landed, what was rejected and what is left.
