---
name: keep-brewing
description: For the rest of the session, never leave shell commands running in the background and end the turn. Running subagents already keep the machine awake (Herdr's caffeine cup stays on while any agent works), so they need no foreground watching; only background shell jobs (builds, uploads, test runs, servers you wait on) must be run or waited on in the foreground. Use on /keep-brewing.
---

Applies for the rest of the session.

Tools like Herdr's caffeine cup keep the machine awake while any agent is working: the main session in a turn, or any subagent still running. A background shell command is not an agent: once the turn ends, nothing counts it as work, so the machine may sleep in the middle of it.

**Subagents:** a running subagent keeps the machine awake by itself. Launch it (in the background is fine), end the turn, and let its completion notification bring you back. Never watch or poll a subagent from the foreground.

**Shell jobs** (builds, uploads, test suites, deploy watchers, anything you would run with `run_in_background`):

- Run them in the foreground.
- If a job must outlive one command's timeout, start it, then wait on it in the foreground in chunks of a few minutes, with a one-line progress update between chunks. Do not poll in a tight loop.
- Ceiling: 2 hours of waiting on one job. Then stop waiting, say what is still running, and ask the user whether to keep waiting or kill it.

When no shell job of yours is running, end the turn, even if subagents are still working.
