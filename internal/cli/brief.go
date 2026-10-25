package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// The first prompt of a brief conversation: one line (a batch claude
// refuses line breaks), starting with neither / nor !. %s names the
// summary (briefPrompt).
const (
	briefAsk  = `Read %s, the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here.`
	briefWait = " Then reply with one or two sentences on where the work stands, and wait for my next message."
	briefThen = " Then: "
)

// What a brief holds, in characters. Together they keep a brief under
// briefMax (TestWriteBriefStaysUnderMax): a few thousand tokens where
// resuming reads the whole conversation.
const (
	briefMax                 = 16000
	briefPrompts             = 4 // the last prompts it quotes
	titleMax                 = 100
	firstHead, firstTail     = 400, 200
	summaryHead, summaryTail = 2000, 4000
	promptHead, promptTail   = 250, 250
	answerHead, answerTail   = 300, 1200
	briefTodos, todoMax      = 10, 100
	briefFiles, fileMax      = 10, 100
	gitMax                   = 800
)

// briefIntro opens a brief; the transcript's path follows it.
const briefIntro = "ccam wrote this from the saved transcript of that conversation, without a model, so details are missing. To look one up, search the transcript (a JSONL file) rather than reading it whole: "

// gitTimeout is how long a brief waits for git status.
const gitTimeout = 3 * time.Second

// briefKeep is how long a summary stays in the state folder, for a person
// to look at what was sent: the new conversation keeps its text, which
// Claude Code attaches when the first prompt is sent or, when the prompt
// names the file (briefPrompt), Claude reads with its Read tool.
const briefKeep = 7 * 24 * time.Hour

// briefArgs are the claude arguments that start a new conversation on the
// account, which first reads ccam's summary of conv (--brief), and the lines
// ccam prints about it. The summary is saved in the state folder or, in a
// dry run, printed instead; a prompt among args follows the summary.
func briefArgs(e *env, conv *conversation, args []string) ([]string, string, error) {
	d, err := transcript.ReadDigest(conv.path)
	if err != nil {
		return nil, "", fmt.Errorf("--brief could not read the conversation (%v); `ccam next --continue` without --brief resumes it whole", err)
	}
	if len(d.Prompts) == 0 && d.Answer == "" && d.Summary == "" {
		return nil, "", errors.New("--brief found nothing to summarize in the conversation (a newer Claude Code may save it differently); `ccam next --continue` without --brief resumes it whole")
	}
	// The folder as Claude Code has it: links resolved.
	dir := conv.cwd
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		dir = real
	}
	text := writeBrief(d, conv.path, dir, gitStatus(conv.cwd))
	folder := filepath.Join(e.paths.State, "brief")
	file := filepath.Join(folder, transcript.ID(conv.path)+".md")
	if e.dryRun {
		fmt.Fprintln(e.stdout, text)
	} else {
		pruneBriefs(folder, clock())
		if err := fsutil.AtomicWrite(file, []byte(text), 0o600); err != nil {
			return nil, "", fmt.Errorf("--brief could not save the summary (%v)", err)
		}
	}
	what := "the last conversation"
	if d.Title != "" {
		what = fmt.Sprintf("%q", d.Title)
	}
	note := "starting a new conversation with ccam's summary of " + what + " instead of resuming it"
	if n := transcript.Size(conv.path); n > 0 {
		note += fmt.Sprintf(" (about %s tokens)", tokens(n))
	}
	note += fmt.Sprintf("\n  summary (%s characters): %s", tokens(utf8.RuneCountInString(text)), file)
	prompt, rest := launcher.TakePrompt(args)
	return append([]string{briefPrompt(file, dir, prompt)}, rest...), note, nil
}

// briefPrompt is the first prompt for the summary at file, ending in the
// person's own prompt when they gave one. It @-mentions file, which Claude
// Code then attaches without asking: by its path or, when that has a "
// (where the mention would end) or a # (which Claude Code reads as a line
// range or an anchor), by its path from dir, the folder the session starts
// in with its links resolved, which Claude Code resolves it from. When that
// has one too, or no path leads from dir (another drive), it names the file
// for Claude to read with its Read tool, which asks first outside the
// working folders (a -p run in Manual mode denies it).
func briefPrompt(file, dir, prompt string) string {
	name := `@"` + file + `"`
	if strings.ContainsAny(file, `"#`) {
		name = "the file " + file
		if rel, err := filepath.Rel(dir, file); err == nil && !strings.ContainsAny(rel, `"#`) {
			name = `@"` + rel + `"`
		}
	}
	ask := fmt.Sprintf(briefAsk, name)
	if prompt == "" {
		return ask + briefWait
	}
	return ask + briefThen + prompt
}

// writeBrief is the summary of conversation d, whose transcript is at path
// and whose folder is dir, with git (gitStatus there): Markdown, its parts
// cut so that it stays under briefMax characters. A part with nothing in
// it is left out.
func writeBrief(d transcript.Digest, path, dir, git string) string {
	head := "# Summary of an earlier Claude Code conversation"
	if d.Title != "" {
		head += `: "` + cut(oneLine(d.Title), titleMax) + `"`
	}
	parts := []string{head, briefIntro + path}
	recent := d.Prompts[max(0, len(d.Prompts)-briefPrompts):]
	earlier := len(d.Prompts) - len(recent)
	if earlier > 0 {
		parts = append(parts, "## First request\n\n"+clip(d.Prompts[0], firstHead, firstTail))
		earlier--
	}
	if d.Summary != "" {
		parts = append(parts, "## Summary Claude Code wrote when it compacted the conversation\n\n"+clip(d.Summary, summaryHead, summaryTail))
	}
	if len(recent) > 0 {
		s := "## Last requests, oldest first"
		if earlier > 0 {
			s += "\n\n" + leftOut(earlier, "earlier request is", "earlier requests are")
		}
		for i, p := range recent {
			s += fmt.Sprintf("\n\n%d. %s", i+1, clip(p, promptHead, promptTail))
		}
		parts = append(parts, s)
	}
	if d.Answer != "" {
		s := "## Claude's last reply\n\n" + clip(d.Answer, answerHead, answerTail)
		if d.Unanswered {
			s += "\n\n(This came before the last request above; Claude wrote no text after that request.)"
		}
		parts = append(parts, s)
	}
	if todos := d.Todos; len(todos) > 0 {
		s := "## Todo list\n\n"
		if n := len(todos) - briefTodos; n > 0 {
			var done int
			todos, done = keepTodos(todos, n)
			switch done {
			case n:
				s += leftOut(n, "finished item is", "finished items are")
			case 0:
				s += leftOut(n, "unfinished item is", "unfinished items are")
			default:
				s += fmt.Sprintf("(%d items are left out: %d finished, %d unfinished.)", n, done, n-done)
			}
			s += "\n"
		}
		lines := make([]string, len(todos))
		for i, t := range todos {
			switch text := cut(oneLine(t.Text), todoMax); t.Status {
			case "completed":
				lines[i] = "- [x] " + text
			case "in_progress":
				lines[i] = "- [ ] " + text + " (in progress)"
			default:
				lines[i] = "- [ ] " + text
			}
		}
		parts = append(parts, s+strings.Join(lines, "\n"))
	}
	if files := d.Files; len(files) > 0 {
		s := "## Files Claude changed\n\n"
		if n := len(files) - briefFiles; n > 0 {
			s += leftOut(n, "other file is", "other files are") + "\n"
			files = files[n:]
		}
		lines := make([]string, len(files))
		for i, f := range files {
			// Relative to the folder when inside it.
			if dir != "" {
				if rel, err := filepath.Rel(dir, f); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					f = rel
				}
			}
			lines[i] = "- " + cutLeft(f, fileMax)
		}
		parts = append(parts, s+strings.Join(lines, "\n"))
	}
	if git != "" {
		lines := strings.Split(git, "\n")
		var kept []string
		size := 0
		for _, l := range lines {
			l = strings.TrimSuffix(l, "\r")
			n := utf8.RuneCountInString(l)
			if len(kept) > 0 {
				n++ // the line break before it
			}
			if size+n > gitMax {
				break
			}
			size += n
			kept = append(kept, l)
		}
		if len(kept) == 0 {
			kept = []string{cut(strings.TrimSuffix(lines[0], "\r"), gitMax)}
		}
		s := "## git status\n\n```\n" + strings.Join(kept, "\n") + "\n```"
		if n := len(lines) - len(kept); n > 0 {
			s += "\n\n" + leftOut(n, "more line is", "more lines are")
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, "\n\n") + "\n"
}

// keepTodos leaves n items out of todos so that the rest show where the
// work stands: finished items first, the oldest first, then pending ones
// from the end of the list (the work furthest ahead), and items in
// progress last. The rest keep their order; done is how many finished
// items it left out.
func keepTodos(todos []transcript.Todo, n int) (kept []transcript.Todo, done int) {
	out := make([]bool, len(todos))
	for i := 0; i < len(todos) && n > 0; i++ {
		if todos[i].Status == "completed" {
			out[i] = true
			n--
			done++
		}
	}
	for i := len(todos) - 1; i >= 0 && n > 0; i-- {
		if s := todos[i].Status; s != "completed" && s != "in_progress" {
			out[i] = true
			n--
		}
	}
	for i := 0; i < len(todos) && n > 0; i++ {
		if todos[i].Status == "in_progress" {
			out[i] = true
			n--
		}
	}
	for i, t := range todos {
		if !out[i] {
			kept = append(kept, t)
		}
	}
	return kept, done
}

// clip shortens s, a text of lines, to its first head and last tail
// characters with a line between them that says how many it left out.
func clip(s string, head, tail int) string {
	r := []rune(s)
	if len(r) <= head+tail {
		return s
	}
	return string(r[:head]) + fmt.Sprintf("\n[... %d characters left out ...]\n", len(r)-head-tail) + string(r[len(r)-tail:])
}

// cut shortens s, one line, to n characters ending in "...".
func cut(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-3]) + "..."
	}
	return s
}

// cutLeft shortens s, a path, to n characters starting with "...".
func cutLeft(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return "..." + string(r[len(r)-n+3:])
	}
	return s
}

// oneLine is s with every run of spaces and line breaks one space.
func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// leftOut says that n things were left out: one names one of them, many
// several.
func leftOut(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("(1 %s left out.)", one)
	}
	return fmt.Sprintf("(%d %s left out.)", n, many)
}

// gitStatus is `git status --short --branch` in dir, or "" outside a git
// work tree, without git, or after gitTimeout. It takes no optional locks,
// so it does not rewrite the index, and shows paths as they are
// (core.quotepath off). A variable so tests can replace it.
var gitStatus = func(dir string) string {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "--no-optional-locks", "-c", "core.quotepath=off", "status", "--short", "--branch")
	cmd.Dir = dir
	cmd.WaitDelay = time.Second
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(out), "\r\n")
}

// pruneBriefs removes the summaries in dir older than briefKeep at now.
func pruneBriefs(dir string, now time.Time) {
	entries, _ := os.ReadDir(dir)
	for _, de := range entries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".md") {
			continue
		}
		if info, err := de.Info(); err == nil && now.Sub(info.ModTime()) > briefKeep {
			os.Remove(filepath.Join(dir, de.Name()))
		}
	}
}
