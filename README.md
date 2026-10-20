# My Agent Skills

Skills I use with my coding agents, free for anyone to use. They follow the open `SKILL.md` format and are written to work with any agent that supports it. So far I use them with Claude Code and Codex.

## Install

Clone the repo and copy a skill folder into your agent's skills directory:

```sh
git clone https://github.com/bragboy/skills.git
cp -R skills/<skill-name> ~/.claude/skills/   # Claude Code
cp -R skills/<skill-name> ~/.codex/skills/    # Codex
```

Then invoke it by name, for example `/keep-brewing` in Claude Code.

## Skills

### Assume Reviewer

For long sessions. The agent becomes the lead: it plans the work, hands the heavy parts to subagents, reviews what they send back, and merges or commits it. This keeps the main conversation small, so it can run for a long time without losing track. Small jobs, like a one-line fix or a quick check, it just does itself. If the agent cannot spawn subagents, it does the work itself in small steps.

At the start it asks you how to work:

- **Landing Changes:** commit each accepted change, merge without committing, stage for you to commit, or ask each time.
- **Isolation:** a separate git worktree per subagent, or one shared working tree.
- **Parallelism:** one task at a time, or several at once when they are independent.
- **Review Depth:** skim the summary, read the diff, or read the diff and run the tests.

It keeps a running list of what is done, in progress and next, and ends with a summary of what landed, what was rejected and what is left.

### Meeting Video to Doc

Turns a recorded meeting (an MP4 video plus its VTT transcript) into a dense Markdown document that another person or AI can work from. It is built for screen-share demos, where the speaker says "you can see here" or "type this" and the actual content is only on screen.

It reads the transcript first to find the moments that matter, then pulls just those frames from the video with `ffmpeg` instead of watching the whole thing. Commands and output seen on screen are retyped into code blocks, and only a few screenshots are kept, each labeled with its video timestamp.

Before writing, it shows you a map of the meeting's sections and asks:

- **Sections:** which parts of the meeting to cover.
- **Purpose:** runbook, decision log, design record, onboarding doc, or context for another AI.
- **On-Screen Content:** retype it as text, use screenshots, or both.
- **Screenshot Budget:** roughly how many images to include.

It never makes up what was on screen. Anything it could not read clearly is marked as unverified, and the document ends with a caveats section. Needs `ffmpeg` installed.

### Keep Brewing

Keeps the agent working in the foreground for the rest of the session. It never starts a job in the background and ends its turn. It runs commands directly, or waits on long jobs in chunks of a few minutes with a short progress update in between, until everything (including subagents) has finished.

Why: some tools keep your computer awake only while an agent is working, such as the "until agents finish" mode of the caffeine cup in [Herdr](https://github.com/penso/herdr-gpui). An agent that ends its turn with a job still running in the background looks idle, so your computer can go to sleep in the middle of the job. With this skill on, the computer stays awake exactly as long as the work runs, and no longer.

To stop a stuck job from keeping the computer awake forever, it waits at most 2 hours on any one job, then asks you whether to keep waiting or kill it.
