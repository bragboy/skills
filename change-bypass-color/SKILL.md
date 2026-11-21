---
name: change-bypass-color
description: Re-apply the local patch that draws Claude Code's red "bypass permissions on" footer label in a neutral gray (or another color). Use when the user invokes /change-bypass-color, typically after Claude Code auto-updated and the label went red again, or asks to change that label's color.
---

# Change the bypass permissions label color

Claude Code draws "bypass permissions on" in the theme's `error` red. No
setting or mod can change just that label, so this skill patches the user's
own Claude Code binary. The user asked for this and accepts that every update
undoes it.

## Run it

```sh
sh ~/.claude/skills/change-bypass-color/apply.sh          # gray #888
sh ~/.claude/skills/change-bypass-color/apply.sh '#999'   # another color
```

Run it in the foreground from the session's own directory (it is trusted, so the
test run skips the folder trust prompt); it takes about 20 seconds. It:

1. resolves the active binary from `~/.local/bin/claude`,
2. builds `patcher/` (Go) and exits early if that binary already has the color,
3. patches a copy beside it, re-signs it ad hoc, and starts it in a pseudo-terminal in bypass mode to check the footer is no longer red and nothing failed to parse,
4. only then keeps the unpatched binary as `<version>.orig` and moves the copy over the original.

When it finishes, tell the user to run `/restart` in each Claude session. The
running process keeps the old binary until then.

`CLAUDE_BIN=/path/to/binary` patches another binary instead, for testing on a copy.

## Colors

A color must fit in the bytes of the old value: up to 5 characters over the
original `"error"`. That allows `#rgb` hex (`#888`, `#999`, `#aaa`) or a theme
key of up to 5 letters (`text`). `#rrggbb` does not fit, and the patcher refuses it.

## When it fails

The script never touches the installed binary unless the patched copy was
verified. Report its output, then:

- `bypass mode color not found`: Claude Code's mode table changed. Search the
  binary for `external:"bypassPermissions"` (`LC_ALL=C grep -a -b -o -F`) and
  update `colorRe` in `patcher/main.go`.
- `Bun trailer not found` or `module table size ... not a multiple of 52`: Bun's
  executable layout changed (it was Bun 1.4.4). Recheck the trailer's 32-byte
  Offsets struct and the 13-word module record documented in `patcher/main.go`.
- `footer is still red`: Bun ran bytecode anyway. Check how this Bun version
  decides to use a module's bytecode.
- `could not see the bypass footer`: the copy showed another screen first
  (trust prompt, onboarding, update notice). Read the capture and adjust the test.

## Undo

```sh
V="$(readlink -f ~/.local/bin/claude)"; mv "$V.orig" "$V"
```

Then `/restart`.
