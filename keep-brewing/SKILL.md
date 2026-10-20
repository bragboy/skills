---
name: keep-brewing
description: For the rest of the session, never leave work running in the background and end the turn. Run commands in the foreground and stay in the turn until every job and subagent finishes, so tools that keep the machine awake "until agents finish" (like Herdr's caffeine cup) stay on exactly as long as work is running. Use on /keep-brewing.
---

Applies for the rest of the session.

Some tools keep the machine awake only while an agent is working. Ending your turn marks you idle, so the machine may sleep mid-job. So: **do not end the turn while anything you started is still running.**

- Run commands in the foreground.
- If a job must outlive one command's timeout, start it, then wait on it in the foreground in chunks of a few minutes, with a one-line progress update between chunks. Do not poll in a tight loop.
- Wait for subagents to finish instead of ending the turn to be notified later.
- Ceiling: 2 hours of waiting on one job. Then stop waiting, say what is still running, and ask the user whether to keep waiting or kill it.
- When everything is done, end the turn so the machine can sleep again.
