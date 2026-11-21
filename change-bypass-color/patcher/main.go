// Command patcher recolors Claude Code's "bypass permissions on" footer label.
//
// Claude Code ships as a Bun single-file executable. Each bundled module has
// its JavaScript source and precompiled JSC bytecode; Bun runs the bytecode
// when present. The mode table lives in source as
//
//	bypassPermissions:{...,color:"error",external:"bypassPermissions"}
//
// so the patch (1) rewrites that color in the source, keeping the module's
// byte length, (2) clears the module's bytecode length in Bun's module table
// so Bun compiles the patched source instead, and (3) blanks the module's
// "@bytecode" pragma. The result is re-signed ad hoc by the caller.
//
// Layout (Bun 1.4, macOS): the __BUN,__bun section holds a u64 payload length
// then the payload. The payload ends with a 32-byte Offsets struct and the
// trailer "\n---- Bun! ----\n". Offsets.modules is {u32 offset, u32 length}
// into the payload; each module record is 13 u32 words: name, contents,
// sourcemap, bytecode (words 6,7 = offset, length), module_info,
// bytecode_origin_path, then packed flags.
package main

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

const (
	trailer        = "\n---- Bun! ----\n"
	offsetsSize    = 32
	recordWords    = 13
	recordSize     = recordWords * 4
	bytecodePragma = "// @bun @bytecode\n"
	plainPragma    = "// @bun          \n"
)

var colorRe = regexp.MustCompile(`color:"([^"]*)"( *),external:"bypassPermissions"`)

func main() {
	color := flag.String("color", "#888", `new color: a theme key ("text", "inactive") or "#rgb"/"#rrggbb"`)
	dryRun := flag.Bool("dry-run", false, "report what would change without writing")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: patcher [-color #888] [-dry-run] <claude binary>")
		os.Exit(2)
	}
	if err := run(flag.Arg(0), *color, *dryRun); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(path, color string, dryRun bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sectOff, err := bunSection(path)
	if err != nil {
		return err
	}
	payloadLen := binary.LittleEndian.Uint64(data[sectOff:])
	payload := sectOff + 8
	end := payload + payloadLen
	if end > uint64(len(data)) || !bytes.Equal(data[end-uint64(len(trailer)):end], []byte(trailer)) {
		return errors.New("Bun trailer not found: unknown executable layout")
	}
	offs := data[end-uint64(len(trailer))-offsetsSize:]
	modOff := uint64(binary.LittleEndian.Uint32(offs[8:]))
	modLen := uint64(binary.LittleEndian.Uint32(offs[12:]))
	if modLen%recordSize != 0 || payload+modOff+modLen > end {
		return fmt.Errorf("module table size %d is not a multiple of %d: unknown Bun version", modLen, recordSize)
	}

	matches := 0
	for i := uint64(0); i < modLen/recordSize; i++ {
		rec := payload + modOff + i*recordSize
		word := func(n uint64) uint64 { return uint64(binary.LittleEndian.Uint32(data[rec+n*4:])) }
		srcStart, srcLen := payload+word(2), word(3)
		if srcStart+srcLen > end {
			continue
		}
		src := data[srcStart : srcStart+srcLen]
		loc := colorRe.FindSubmatchIndex(src)
		if loc == nil {
			continue
		}
		matches++
		if len(colorRe.FindAllIndex(src, -1)) != 1 {
			return errors.New("bypass color appears more than once in one module")
		}
		nameStart, nameLen := payload+word(0), word(1)
		name := string(data[nameStart : nameStart+nameLen])
		oldColor := string(src[loc[2]:loc[3]])
		// The value, then its closing quote, then padding spaces: the value may
		// grow into the padding, and the quote moves to stay right after it.
		room := (loc[3] - loc[2]) + (loc[5] - loc[4])
		if len(color) > room {
			return fmt.Errorf("color %q is %d chars; only %d fit (old value %q)", color, len(color), room, oldColor)
		}
		bcLen := word(7)
		hasPragma := bytes.HasPrefix(src, []byte(bytecodePragma))
		fmt.Printf("module %s: color %q -> %q, bytecode length %d -> 0, pragma blanked: %v\n",
			name, oldColor, color, bcLen, hasPragma)
		if oldColor == color && bcLen == 0 && !hasPragma {
			fmt.Println("already patched")
			continue
		}
		if dryRun {
			continue
		}
		repl := []byte(color + `"` + string(bytes.Repeat([]byte(" "), room-len(color))))
		copy(src[loc[2]:loc[5]], repl)
		if hasPragma {
			copy(src, plainPragma)
		}
		binary.LittleEndian.PutUint32(data[rec+7*4:], 0)
	}
	if matches == 0 {
		return errors.New("bypass mode color not found: Claude Code's source changed, patch by hand")
	}
	if matches > 1 {
		return fmt.Errorf("bypass mode color found in %d modules: refusing to guess", matches)
	}
	if dryRun {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, info.Mode().Perm())
}

// bunSection returns the file offset of the __BUN,__bun section.
func bunSection(path string) (uint64, error) {
	f, err := macho.Open(path)
	if err != nil {
		return 0, fmt.Errorf("not a Mach-O executable (macOS only): %w", err)
	}
	defer f.Close()
	s := f.Section("__bun")
	if s == nil || s.Seg != "__BUN" {
		return 0, errors.New("no __BUN,__bun section: not a Bun single-file executable")
	}
	return uint64(s.Offset), nil
}
