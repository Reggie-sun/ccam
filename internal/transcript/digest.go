package transcript

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"unicode"
)

// Digest is what a brief (ccam next --continue --brief) tells of a
// conversation, read from its transcript without a model. Rows of kinds or
// shapes it does not know are skipped, so a field stays empty when Claude
// Code writes it differently.
type Digest struct {
	Title      string   // as Title: the /rename name, else Claude Code's ai-title
	Prompts    []string // the prompts a person sent, oldest first (entry.prompt)
	Answer     string   // the text of Claude's last reply (its text blocks, a blank line between)
	Unanswered bool     // a prompt came after that reply
	Summary    string   // what Claude Code wrote at its last compaction, without the lines it puts around it
	Files      []string // files Claude changed (Edit, Write, MultiEdit, NotebookEdit), least recently changed first, each once
	Todos      []Todo   // Claude's todo list as it last stood
}

// Todo is an item of Claude's todo list: a task (TaskCreate, TaskUpdate) or
// an item of TodoWrite.
type Todo struct {
	Text   string
	Status string // pending, in_progress or completed
}

// What Claude Code (2.1.304-2.1.310) puts around the summary it writes when
// it compacts a conversation: a first paragraph, and a last one that tells
// the model to go on without asking.
const (
	compactLead  = "This session is being continued from a previous conversation that ran out of context."
	compactTrail = "Continue the conversation from where it left off without asking the user any further questions."
)

// ReadDigest reads the conversation at path once, from its first row to its
// last, decoding only rows that can hold what it keeps (useful). Subagents'
// rows (isSidechain) are skipped. The error is one of reading the file.
func ReadDigest(path string) (Digest, error) {
	f, err := os.Open(path)
	if err != nil {
		return Digest{}, err
	}
	defer f.Close()
	var d digest
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64<<10), 16<<20)
	for sc.Scan() {
		l := sc.Bytes()
		if !useful(l, len(d.waiting) > 0) {
			continue
		}
		var e struct {
			entry
			APIError   bool `json:"isApiErrorMessage"`
			Attachment struct {
				Type   string          `json:"type"`
				Prompt json.RawMessage `json:"prompt"`
				Meta   bool            `json:"isMeta"`
				Origin *origin         `json:"origin"`
			} `json:"attachment"`
			Custom string `json:"customTitle"`
			AI     string `json:"aiTitle"`
		}
		if json.Unmarshal(l, &e) != nil || e.Sidechain {
			continue
		}
		switch e.Type {
		case "assistant":
			if !e.APIError && e.Message.Model != synthetic {
				d.reply(e.Message.ID, e.Message.Content)
			}
		case "user":
			if e.Compact {
				text, _ := contentText(e.Message.Content)
				d.Summary = cleanSummary(text)
			} else if text, ok := e.prompt(); ok {
				d.Prompts, d.Unanswered = append(d.Prompts, text), true
			} else if len(d.waiting) > 0 {
				d.results(e.Message.Content)
			}
		case "attachment":
			// A prompt the person typed while Claude worked, delivered
			// within the turn.
			a := e.Attachment
			if a.Type == "queued_command" && !a.Meta && a.Origin != nil && a.Origin.Kind == "human" {
				if text, tool := contentText(a.Prompt); !tool {
					if text, ok := personText(text); ok {
						d.Prompts, d.Unanswered = append(d.Prompts, text), true
					}
				}
			}
		case "custom-title":
			d.custom = e.Custom
		case "ai-title":
			d.ai = e.AI
		}
	}
	if err := sc.Err(); err != nil {
		return Digest{}, err
	}
	d.Title = cmp.Or(d.custom, d.ai)
	for _, t := range slices.Concat(d.tasks, d.waiting) {
		d.Todos = append(d.Todos, t.Todo)
	}
	if len(d.Todos) == 0 {
		d.Todos = d.listed
	}
	return d.Digest, nil
}

// useful reports whether the line l may hold what a digest keeps: Claude's
// text and the tool calls it reads, user rows (one of tool results only
// while a TaskCreate waits for its result), queued prompts and titles.
// Claude Code writes compact JSON, and quotes inside a string are escaped,
// so no text can hold these keys and values; the row's type is decoded
// again all the same. A line is looked at as the first kind it names: no
// user row holds an assistant's type (2.1.303-2.1.310), and most bytes are
// in rows of tool calls and results, which this keeps to one or two scans.
func useful(l []byte, waiting bool) bool {
	switch {
	case bytes.Contains(l, []byte(`"type":"assistant"`)):
		if bytes.Contains(l, []byte(`"type":"text"`)) {
			return true
		}
		for _, call := range []string{`"name":"Edit"`, `"name":"Write"`, `"name":"MultiEdit"`, `"name":"NotebookEdit"`, `"name":"TodoWrite"`, `"name":"TaskCreate"`, `"name":"TaskUpdate"`} {
			if bytes.Contains(l, []byte(call)) {
				return true
			}
		}
		return false
	case bytes.Contains(l, []byte(`"type":"user"`)):
		return !bytes.Contains(l, []byte(`"type":"tool_result"`)) ||
			waiting && (bytes.Contains(l, []byte("created successfully")) || bytes.Contains(l, []byte(`"is_error":true`)))
	}
	return bytes.Contains(l, []byte(`"queued_command"`)) || bytes.Contains(l, []byte(`-title"`))
}

// digest is ReadDigest's state while it reads.
type digest struct {
	Digest
	custom, ai string
	last       string // the message id of the last row with text of Claude's
	tasks      []task // the tasks of TaskCreate and TaskUpdate, in the order they were made
	waiting    []task // TaskCreate calls whose result has not named their task (id: the call's)
	listed     []Todo // the list of the last TodoWrite
}

type task struct {
	id string
	Todo
}

// reply reads the content of a row of Claude's reply with message id: its
// text, and the tool calls that change files or the todo list. The rows of
// one message share its id, one block each.
func (d *digest) reply(id string, content json.RawMessage) {
	var blocks []struct {
		Type  string          `json:"type"`
		Text  string          `json:"text"`
		ID    string          `json:"id"`
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}
	if json.Unmarshal(content, &blocks) != nil {
		return
	}
	var texts []string
	for _, b := range blocks {
		switch b.Type {
		case "text":
			if t := strings.TrimSpace(b.Text); t != "" {
				texts = append(texts, t)
			}
		case "tool_use":
			d.call(b.ID, b.Name, b.Input)
		}
	}
	if len(texts) == 0 {
		return
	}
	text := strings.Join(texts, "\n\n")
	if id != "" && id == d.last {
		d.Answer += "\n\n" + text
	} else {
		d.Answer = text
	}
	d.last, d.Unanswered = id, false
}

// call reads Claude's call (id) of tool name with input: the file it
// changes, or the todo list, as Claude Code 2.1.310 replays it when it
// resumes a conversation. A call it would refuse is skipped.
func (d *digest) call(id, name string, input json.RawMessage) {
	switch name {
	case "Edit", "Write", "MultiEdit", "NotebookEdit":
		var in struct {
			File     string `json:"file_path"`
			Notebook string `json:"notebook_path"`
		}
		if json.Unmarshal(input, &in) != nil {
			return
		}
		file := in.File
		if name == "NotebookEdit" {
			file = in.Notebook
		}
		if file != "" {
			d.Files = append(slices.DeleteFunc(d.Files, func(f string) bool { return f == file }), file)
		}
	case "TodoWrite":
		var in struct {
			Todos []struct {
				Content string `json:"content"`
				Status  string `json:"status"`
			} `json:"todos"`
		}
		d.listed = nil
		if json.Unmarshal(input, &in) == nil {
			for _, t := range in.Todos {
				if t.Content != "" {
					d.listed = append(d.listed, Todo{t.Content, t.Status})
				}
			}
		}
	case "TaskCreate":
		// The task's id is in the result.
		var in struct {
			Subject string `json:"subject"`
		}
		if json.Unmarshal(input, &in) == nil && in.Subject != "" {
			d.waiting = append(d.waiting, task{id, Todo{in.Subject, "pending"}})
		}
	case "TaskUpdate":
		var in struct {
			ID      string `json:"taskId"`
			Alias   string `json:"id"`
			Alias2  string `json:"task_id"`
			Status  string `json:"status"`
			Subject string `json:"subject"`
		}
		if json.Unmarshal(input, &in) != nil {
			return
		}
		key := cmp.Or(in.ID, in.Alias, in.Alias2)
		switch {
		case key == "":
		case in.Status == "deleted":
			d.tasks = slices.DeleteFunc(d.tasks, func(t task) bool { return t.id == key })
		case in.Status == "" || in.Status == "pending" || in.Status == "in_progress" || in.Status == "completed":
			i := slices.IndexFunc(d.tasks, func(t task) bool { return t.id == key })
			if i < 0 {
				// Claude Code shows a task it does not know by its id.
				d.tasks, i = append(d.tasks, task{key, Todo{key, "pending"}}), len(d.tasks)
			}
			if in.Subject != "" {
				d.tasks[i].Text = in.Subject
			}
			if in.Status != "" {
				d.tasks[i].Status = in.Status
			}
		}
	}
}

// results reads the tool results in the content of a user row for the
// TaskCreate calls that wait for theirs: one that failed is dropped, and
// "Task #<id> created successfully" names the task, unless a task of that
// id exists. Any other result leaves the call waiting, as in Claude Code.
func (d *digest) results(content json.RawMessage) {
	var blocks []struct {
		Type    string          `json:"type"`
		UseID   string          `json:"tool_use_id"`
		IsError bool            `json:"is_error"`
		Content json.RawMessage `json:"content"`
	}
	if json.Unmarshal(content, &blocks) != nil {
		return
	}
	for _, b := range blocks {
		i := slices.IndexFunc(d.waiting, func(t task) bool { return t.id == b.UseID })
		if b.Type != "tool_result" || b.UseID == "" || i < 0 {
			continue
		}
		if !b.IsError {
			text, _ := contentText(b.Content)
			id := createdID(text)
			if id == "" {
				continue
			}
			if !slices.ContainsFunc(d.tasks, func(t task) bool { return t.id == id }) {
				d.tasks = append(d.tasks, task{id, d.waiting[i].Todo})
			}
		}
		d.waiting = slices.Delete(d.waiting, i, i+1)
	}
}

// createdID is the id of the task a TaskCreate result names: "Task #<id>
// created successfully: <subject>"; "" when it names none.
func createdID(result string) string {
	rest, ok := strings.CutPrefix(result, "Task #")
	if !ok {
		return ""
	}
	i := strings.IndexFunc(rest, unicode.IsSpace)
	if i <= 0 || !strings.HasPrefix(rest[i:], " created successfully") {
		return ""
	}
	return rest[:i]
}

// cleanSummary is the summary Claude Code wrote at a compaction without
// what it puts around it: compactLead's paragraph, and compactTrail with
// what follows it. A note Claude Code puts before the summary (about
// Artifact content others wrote) stays.
func cleanSummary(s string) string {
	if i := strings.Index(s, compactLead); i >= 0 {
		end := len(s)
		if j := strings.Index(s[i:], "\n\n"); j >= 0 {
			end = i + j + 2
		}
		s = s[:i] + s[end:]
	}
	if i := strings.Index(s, compactTrail); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
