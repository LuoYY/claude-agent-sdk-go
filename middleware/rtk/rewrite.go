package rtk

import (
	"path/filepath"
	"strings"
)

// rewrite walks a shell command line and prefixes each pipeline segment's
// leading program with the configured rtk binary when that program is in
// the allow-list (and not blocked).
//
// It is intentionally a lexical, best-effort rewriter: it understands
// simple separators (&&, ||, ;, |) but does not try to emulate a full
// shell parser. Commands inside subshells, backticks, heredocs, or process
// substitutions are left alone. This mirrors how rtk's own PreToolUse
// hook behaves for Claude Code and keeps the middleware predictable.
//
// rewrite returns the new command string and a bool indicating whether
// anything changed.
func (c *config) rewrite(cmd string) (string, bool) {
	if strings.TrimSpace(cmd) == "" {
		return cmd, false
	}

	segments, seps := splitSegments(cmd)
	changed := false

	for i, seg := range segments {
		newSeg, ok := c.rewriteSegment(seg)
		if ok {
			segments[i] = newSeg
			changed = true
		}
	}

	if !changed {
		return cmd, false
	}
	return joinSegments(segments, seps), true
}

// rewriteSegment rewrites a single pipeline segment (no top-level
// separators inside). Environment-variable assignments and common shell
// builtins that simply pass control along (cd, env, sudo, time, nice,
// ionice, stdbuf, command, exec) are skipped so we operate on the real
// program.
func (c *config) rewriteSegment(seg string) (string, bool) {
	trimmed, leading, trailing := trimSurroundingSpace(seg)
	if trimmed == "" {
		return seg, false
	}

	// Already prefixed: do nothing.
	tokens := tokenize(trimmed)
	if len(tokens) == 0 {
		return seg, false
	}
	if basename(tokens[0]) == basename(c.binary) {
		return seg, false
	}

	// Walk past leading env-var assignments and pass-through wrappers
	// (sudo / env / time / nice / ionice / stdbuf / command / exec) until
	// we land on the real program. env-assignments can appear both before
	// any wrapper ("FOO=1 git status") and after one ("env FOO=1 git status").
	idx := 0
	for idx < len(tokens) {
		tok := tokens[idx]
		if isEnvAssignment(tok) {
			idx++
			continue
		}
		switch basename(tok) {
		case "sudo", "env", "time", "nice", "ionice", "stdbuf", "command", "exec":
			idx++
			for idx < len(tokens) && strings.HasPrefix(tokens[idx], "-") {
				idx++
			}
			continue
		}
		break
	}
	if idx >= len(tokens) {
		return seg, false
	}

	prog := basename(tokens[idx])
	if prog == "" {
		return seg, false
	}
	if _, blocked := c.blocked[prog]; blocked {
		return seg, false
	}
	if _, ok := c.commands[prog]; !ok {
		return seg, false
	}

	// Build the prefix, preserving leading whitespace and any skipped
	// wrapper tokens verbatim from the original segment.
	before := strings.Join(tokens[:idx], " ")
	rest := strings.Join(tokens[idx:], " ")

	prefix := c.binary
	if c.ultraCompact {
		prefix += " -u"
	}

	var rebuilt string
	switch {
	case before == "":
		rebuilt = prefix + " " + rest
	default:
		rebuilt = before + " " + prefix + " " + rest
	}
	return leading + rebuilt + trailing, true
}

// splitSegments splits a command line on top-level shell operators
// (&&, ||, ;, |) and returns the segments together with the separators
// that appeared between them (in order). Quoted regions are preserved:
// operators inside single or double quotes do not split.
func splitSegments(cmd string) (segments []string, seps []string) {
	var cur strings.Builder
	inSingle, inDouble, escape := false, false, false

	flush := func() {
		segments = append(segments, cur.String())
		cur.Reset()
	}

	for i := 0; i < len(cmd); i++ {
		ch := cmd[i]
		if escape {
			cur.WriteByte(ch)
			escape = false
			continue
		}
		switch {
		case ch == '\\' && !inSingle:
			cur.WriteByte(ch)
			escape = true
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			cur.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			cur.WriteByte(ch)
		case !inSingle && !inDouble && (ch == '&' || ch == '|') && i+1 < len(cmd) && cmd[i+1] == ch:
			flush()
			seps = append(seps, string([]byte{ch, ch}))
			i++
		case !inSingle && !inDouble && ch == ';':
			flush()
			seps = append(seps, ";")
		case !inSingle && !inDouble && ch == '|':
			flush()
			seps = append(seps, "|")
		default:
			cur.WriteByte(ch)
		}
	}
	flush()
	return segments, seps
}

func joinSegments(segments, seps []string) string {
	var b strings.Builder
	for i, s := range segments {
		b.WriteString(s)
		if i < len(seps) {
			b.WriteString(seps[i])
		}
	}
	return b.String()
}

// tokenize produces a naive whitespace tokenization that honours single
// and double quotes. Quote characters are kept in the returned tokens so
// the segment can be round-tripped faithfully.
func tokenize(s string) []string {
	var tokens []string
	var cur strings.Builder
	inSingle, inDouble, escape := false, false, false

	flush := func() {
		if cur.Len() > 0 {
			tokens = append(tokens, cur.String())
			cur.Reset()
		}
	}

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escape {
			cur.WriteByte(ch)
			escape = false
			continue
		}
		switch {
		case ch == '\\' && !inSingle:
			cur.WriteByte(ch)
			escape = true
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			cur.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			cur.WriteByte(ch)
		case (ch == ' ' || ch == '\t') && !inSingle && !inDouble:
			flush()
		default:
			cur.WriteByte(ch)
		}
	}
	flush()
	return tokens
}

// trimSurroundingSpace splits s into (middle, leading, trailing) where
// leading+middle+trailing == s, middle has no outer spaces/tabs, and the
// spans are byte-exact so the segment can be reassembled without loss.
func trimSurroundingSpace(s string) (middle, leading, trailing string) {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end], s[:start], s[end:]
}

func basename(p string) string {
	if p == "" {
		return ""
	}
	// Strip any surrounding quotes left by tokenize.
	p = strings.Trim(p, "'\"")
	return filepath.Base(p)
}

func isEnvAssignment(tok string) bool {
	if tok == "" {
		return false
	}
	eq := strings.IndexByte(tok, '=')
	if eq <= 0 {
		return false
	}
	// Name must be a valid identifier prefix.
	for i := 0; i < eq; i++ {
		ch := tok[i]
		switch {
		case ch >= 'A' && ch <= 'Z':
		case ch >= 'a' && ch <= 'z':
		case ch >= '0' && ch <= '9' && i > 0:
		case ch == '_':
		default:
			return false
		}
	}
	return true
}
