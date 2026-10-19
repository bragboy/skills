---
name: meeting-video-to-doc
description: Turn a recorded meeting (MP4 + VTT transcript) into a dense, AI-readable Markdown document, using targeted video frames for the visual context the transcript is missing. Use when the user has a call recording and wants notes, a runbook, a decision log, or context to feed to another AI - especially screen-share demos where the transcript says "you can see here" and "type this" without saying what. Invoke on /meeting-video-to-doc or when the user supplies a .vtt and .mp4 of the same session.
---

# Meeting video → document

A meeting VTT is nearly useless on its own for screen-share demos: the speaker says "you can see here", "type this", "this one", and the referent is on screen, not in the words. The MP4 has the referent but is far too expensive to watch end to end.

The method: **read the transcript to find out where to look, then look only there.** The transcript is the index; the video is random-access storage. Never process the video linearly.

## Non-negotiables

- **Never dump the whole video into context.** No full-video reads, no frame-per-second extraction, no "let me sample every 30 seconds across 60 minutes."
- **Never invent on-screen content.** If a frame is unreadable, extract a neighbouring frame or say the detail is unverified. A plausible-looking IP address or command that you did not actually read off the screen is the single worst failure mode of this skill.
- **Ask before producing.** The scoping conversation in Phase 2 is the point of the skill, not overhead.
- **Transcribe, don't screenshot, anything that is text.** For a downstream AI reader, a fenced code block beats an image every time. Images earn their place only when they carry something text cannot: a UI layout, a topology, a table with meaningful structure.

## Phase 1 - Cheap reconnaissance

Do all of this before talking to the user.

1. **Check tooling.** `which ffmpeg ffprobe`. If missing, `brew install ffmpeg` (macOS) in the background and continue with the transcript work while it installs. ffmpeg is the right tool here; there is no better option for frame extraction.
2. **Probe the video:** `ffprobe -v error -show_entries format=duration -show_entries stream=codec_type,width,height,r_frame_rate -of default=noprint_wrappers=1 <file>`
3. **Verify timeline alignment.** The VTT's last cue end time should be within a few seconds of the video duration. If it is not, the files are from different recordings or the video has a leading offset - establish the offset with one probe frame at a moment the transcript describes unambiguously, and apply it to every later extraction.
4. **Read the VTT** and build a section map: contiguous blocks of one topic, with start/end timestamps and a one-line summary each. For a Teams-style VTT the cue IDs (`<uuid>/14-0`, `/14-1`) group fragments of one utterance - merging on that group number gives clean turns.
5. **Extract one probe frame** mid-recording to learn the layout: is it a full-screen share, a share plus participant thumbnails, picture-in-picture? Note the pixel box of the shared region so every later crop can exclude webcam tiles and browser chrome.

## Phase 2 - Scope with the user

Present the section map (timestamps + one line each) and ask, in one batched round of questions:

- **Which sections** to convert. Rarely all of them.
- **What the document is for** - runbook, decision log, design record, onboarding doc, context to paste into another AI. This decides the structure more than anything else.
- **How to treat on-screen content** - transcribe to code blocks with a few supporting screenshots (the usual right answer), screenshots with prose, or text only.
- **Screenshot budget** - offer concrete numbers (3-4 / 6-10 / 15-20). "As few as possible" means one per distinct idea, not one per step.

If the user already answered some of these in their request, don't re-ask.

## Phase 3 - Targeted extraction

1. **Pick candidate moments from the transcript, not from the clock.** Deictic language is the signal: "you can see here", "type this", "this one", "if we go to", "the important thing is". Also every command the speaker reads aloud, and every claim about what the screen shows.
2. **Aim slightly late.** Land 2-8 seconds *after* the sentence starts - the speaker describes what they are about to do, then does it. For command output, aim after the command has had time to render.
3. **Extract a candidate batch in one pass**, cropped to the shared-screen region learned in Phase 1:
   ```bash
   ffmpeg -hide_banner -loglevel error -ss 00:35:28 -i in.mp4 -frames:v 1 \
     -vf "crop=W:H:X:Y" -q:v 3 frames/f003528.jpg
   ```
   Put `-ss` before `-i` for fast seek. Name frames by timestamp so every one is traceable back to the video.
4. **Read the frames and transcribe what they show.** Budget roughly 3k tokens per frame read; 20-30 frames is normally affordable and plenty. Retype commands and output into code blocks as you go. Terminal text at 1080p is readable when cropped to the window.
5. **Re-aim when a frame misses** - mid-scroll, dialog in the way, output already scrolled off. One retry a few seconds either side usually lands it.
6. **Cut ruthlessly.** Most extracted frames exist to inform your writing, not to appear in the document. Keep only the ones inside the agreed budget.

### Crop, then verify

Final images go in an `assets/` folder beside the document. Crop tightly to the content - the terminal window, the table, the panel - not the whole desktop. Estimate crop boxes from a frame you have already read, then **re-read one cropped output to confirm the framing** before writing the doc around it. PNG for text-heavy screens.

## Phase 4 - Write

Structure for a downstream AI reader, not for a human skimming:

- **Provenance header:** source files, the time range covered, and how current the information is.
- **A short primer** so the document stands alone, even when the user scoped you to one section.
- **Commands and output in fenced code blocks**, transcribed exactly - including the failed attempts and typos visible on screen. A speaker's mistyped command that produced `Invalid token [ro]` is genuinely useful; it documents the CLI's real grammar.
- **Tables** for anything enumerable: inventory, field meanings, pass criteria, known issues.
- **Screenshots captioned with their video timestamp**, e.g. `*video 51:50*`, so any claim can be re-verified at the source.
- **Explicit pass/fail criteria** wherever the source was a demo. "What proves this works" and "what it looks like when it's broken" are the highest-value sentences in the whole document, and speakers usually state both in passing.
- **A caveats section.** Anything the speaker said was untested, unexplained, or known-broken. This is what stops a downstream reader trusting the doc too far.

Fix systematic ASR errors consistently and **list the corrections** in your reply so the user can catch a wrong guess. Flag any correction you are unsure about explicitly rather than burying it. Keep genuinely ambiguous audio verbatim.

## Failure modes to avoid

| Symptom | Fix |
|---|---|
| Screenshot of a terminal, no transcription | Retype the content; the image is the citation, not the payload |
| "The speaker demonstrated several verification methods" | Name them, give the commands, give the expected output |
| Frame shows a scrollbar mid-drag or a stale screen | Re-extract a few seconds later |
| Filler and backchannels reproduced verbatim | This is a document, not a transcript - the cleaned VTT already exists for that |
| Screenshots that duplicate the code blocks | Cut them; the budget is for what text cannot carry |
| A detail asserted that no frame or cue actually supports | Delete it or mark it unverified |
