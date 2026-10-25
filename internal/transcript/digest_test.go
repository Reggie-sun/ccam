package transcript

import (
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// The brief fixture is synthetic: rows of the shapes Claude Code
// 2.1.303-2.1.310 writes (key names as observed in real conversations on
// this machine, the compaction summary with the parts Claude Code puts
// around it), and Task tool rows modelled on 2.1.310's own replay of the
// todo list, as none were observed. Ids, paths and texts are placeholders.
// Its last line is cut, as a row being written is.
func TestReadDigest(t *testing.T) {
	d, err := ReadDigest(filepath.Join("testdata", "brief-2.1.310.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	want := Digest{
		Title:      "(custom title)",
		Prompts:    []string{"(first request)", "(queued request)", "(second request) see @src/x.go", "(third request)\n[image]"},
		Answer:     "(last reply, part one)\n\n(part two)",
		Unanswered: true,
		Summary: "<tools-used-before-summary>Edit, TaskCreate</tools-used-before-summary>\n(summary)\n\n" +
			"If you need specific details from before compaction (like exact code snippets, error messages, or content you generated), " +
			"read the full transcript at: /home/u/.claude/projects/-home-u-project/22222222-2222-4222-8222-222222222222.jsonl",
		Files: []string{"/home/u/project/README.md", "/home/u/project/nb.ipynb", "/home/u/project/parser.go"},
		Todos: []Todo{{"(task one)", "completed"}, {"(task two)", "in_progress"}},
	}
	if !reflect.DeepEqual(d, want) {
		t.Fatalf("got\n%#v\nwant\n%#v", d, want)
	}
	if _, err := ReadDigest(filepath.Join(t.TempDir(), "none.jsonl")); err == nil {
		t.Fatal("a missing file")
	}
	// A line longer than any Claude Code writes is an error, not a digest
	// of what came before it.
	long := filepath.Join(t.TempDir(), "long.jsonl")
	write(t, long, `{"type":"user","message":{"role":"user","content":"(first)"}}`+"\n"+strings.Repeat("x", 17<<20)+"\n", time.Now())
	if _, err := ReadDigest(long); err == nil {
		t.Fatal("a line over 16 MiB")
	}
}

// digestOf is the digest of a conversation of rows.
func digestOf(t *testing.T, rows ...string) Digest {
	t.Helper()
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, strings.Join(rows, "\n")+"\n", time.Now())
	d, err := ReadDigest(path)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// A prompt is what a person sent, as transcript.prompt counts it: not
// what Claude Code sends or shows itself. Made-up rows of the shapes Claude
// Code 2.1.303-2.1.310 writes.
func TestReadDigestPrompts(t *testing.T) {
	user := func(content string, keys string) string {
		return `{"parentUuid":"p","isSidechain":false,"type":"user","message":{"role":"user","content":` + content + `}` + keys + `,"uuid":"u"}`
	}
	lead := "This session is being continued from a previous conversation that ran out of context. The summary below covers the earlier portion of the conversation."
	trail := "Continue the conversation from where it left off without asking the user any further questions. Resume directly."
	for _, c := range []struct {
		name    string
		row     string
		prompts []string
		summary string
	}{
		{"older Claude Code, no origin", user(`"(typed)"`, ""), []string{"(typed)"}, ""},
		{"typed", user(`"  (typed)\n"`, `,"origin":{"kind":"human"},"promptSource":"typed"`), []string{"(typed)"}, ""},
		{"text and image", user(`[{"type":"text","text":"(look)"},{"type":"image","source":{"type":"base64","data":"AAAA"}}]`, `,"origin":{"kind":"human"}`), []string{"(look)\n[image]"}, ""},
		{"a slash command", user(`"<command-message>review</command-message>\n<command-name>/review</command-name>\n<command-args>(the parser)</command-args>"`, `,"origin":{"kind":"human"}`), []string{"/review (the parser)"}, ""},
		{"a slash command without arguments", user(`"<command-message>review</command-message>\n<command-name>/review</command-name>"`, `,"origin":{"kind":"human"}`), []string{"/review"}, ""},
		{"a peer", user(`"(hello)"`, `,"isMeta":true,"origin":{"kind":"peer"}`), nil, ""},
		{"a peer without isMeta", user(`"(hello)"`, `,"origin":{"kind":"peer"}`), nil, ""},
		{"a finished task", user(`"<task-notification>(done)</task-notification>"`, `,"origin":{"kind":"task-notification"}`), nil, ""},
		{"a finished task without origin", user(`"<task-notification>(done)</task-notification>"`, ""), nil, ""},
		{"the expanded slash command", user(`[{"type":"text","text":"(the command's prompt)"}]`, `,"isMeta":true`), nil, ""},
		{"a subagent", user(`"(task)"`, `,"isSidechain":true`), nil, ""},
		{"a tool's result", user(`[{"tool_use_id":"t1","type":"tool_result","content":"(ok)"}]`, ""), nil, ""},
		{"a local command", user(`"<command-name>/mcp</command-name>\n<command-message>mcp</command-message>\n<command-args></command-args>"`, ""), nil, ""},
		{"its output", user(`"<local-command-stdout>(output)</local-command-stdout>"`, ""), nil, ""},
		{"a ! command", user(`"<bash-input>ls</bash-input>"`, ""), []string{"! ls"}, ""},
		{"its output", user(`"<bash-stdout>(files)</bash-stdout><bash-stderr></bash-stderr>"`, ""), nil, ""},
		{"Esc", user(`[{"type":"text","text":"[Request interrupted by user]"}]`, ""), nil, ""},
		{"Esc during a tool", user(`[{"type":"text","text":"[Request interrupted by user for tool use]"}]`, ""), nil, ""},
		{"only spaces", user(`" \n "`, `,"origin":{"kind":"human"}`), nil, ""},
		{"shown in the transcript only", user(`"(note)"`, `,"isVisibleInTranscriptOnly":true`), nil, ""},
		{"queued while Claude worked", `{"isSidechain":false,"attachment":{"type":"queued_command","prompt":"(also this)","commandMode":"prompt","origin":{"kind":"human"},"humanTurn":true},"type":"attachment"}`, []string{"(also this)"}, ""},
		{"a queued task notification", `{"isSidechain":false,"attachment":{"type":"queued_command","prompt":"<task-notification>(done)</task-notification>","commandMode":"task-notification","origin":{"kind":"task-notification"}},"type":"attachment"}`, nil, ""},
		{"a queued peer message", `{"isSidechain":false,"attachment":{"type":"queued_command","prompt":"(hello)","commandMode":"prompt","isMeta":true,"origin":{"kind":"peer"}},"type":"attachment"}`, nil, ""},
		{"a queued prompt Claude Code marks as its own", `{"isSidechain":false,"attachment":{"type":"queued_command","prompt":"(hidden)","commandMode":"prompt","isMeta":true,"origin":{"kind":"human"}},"type":"attachment"}`, nil, ""},
		{"a queued prompt of a subagent", `{"isSidechain":true,"attachment":{"type":"queued_command","prompt":"(task)","origin":{"kind":"human"}},"type":"attachment"}`, nil, ""},
		// The summary of a compaction is no prompt; what Claude Code puts
		// around it goes, its note about Artifact content stays.
		{"a compaction summary", user(`"`+lead+`\n\n(summary)\n\n`+trail+`"`, `,"isVisibleInTranscriptOnly":true,"isCompactSummary":true`), nil, "(summary)"},
		{"a summary as Claude Code did not write it", user(`"(summary)\n\n(more)"`, `,"isVisibleInTranscriptOnly":true,"isCompactSummary":true`), nil, "(summary)\n\n(more)"},
		{"a summary with a note", user(`"<artifact-content-authored-by-others/>\n(note)\n`+lead+`\n\n(summary)"`, `,"isVisibleInTranscriptOnly":true,"isCompactSummary":true`), nil, "<artifact-content-authored-by-others/>\n(note)\n(summary)"},
		{"a summary in blocks", user(`[{"type":"text","text":"`+lead+`\n\n(summary)"}]`, `,"isCompactSummary":true`), nil, "(summary)"},
	} {
		d := digestOf(t, c.row)
		if !slices.Equal(d.Prompts, c.prompts) || d.Summary != c.summary {
			t.Errorf("%s: prompts %q, summary %q", c.name, d.Prompts, d.Summary)
		}
		if d.Unanswered != (len(c.prompts) > 0) {
			t.Errorf("%s: unanswered %v", c.name, d.Unanswered)
		}
	}
}

// Claude's reply is the text of its last message, from rows of the main
// thread that the model wrote.
func TestReadDigestReply(t *testing.T) {
	reply := func(id, model, text, keys string) string {
		return `{"isSidechain":false,"type":"assistant","message":{"id":"` + id + `","model":"` + model + `","role":"assistant","content":[{"type":"text","text":"` + text + `"}]}` + keys + `}`
	}
	prompt := `{"isSidechain":false,"type":"user","message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`
	for _, c := range []struct {
		name       string
		rows       []string
		answer     string
		unanswered bool
	}{
		{"one message over two rows", []string{reply("m1", "claude-opus-5-5", "(one)", ""), reply("m1", "claude-opus-5-5", "(two)", "")}, "(one)\n\n(two)", false},
		{"the last message", []string{reply("m1", "claude-opus-5-5", "(one)", ""), reply("m2", "claude-opus-5-5", "(two)", "")}, "(two)", false},
		{"no id", []string{reply("", "claude-opus-5-5", "(one)", ""), reply("", "claude-opus-5-5", "(two)", "")}, "(two)", false},
		{"a prompt after it", []string{reply("m1", "claude-opus-5-5", "(one)", ""), prompt}, "(one)", true},
		{"a reply to that prompt", []string{reply("m1", "claude-opus-5-5", "(one)", ""), prompt, reply("m2", "claude-opus-5-5", "(two)", "")}, "(two)", false},
		{"a refusal at a limit", []string{reply("m1", "claude-opus-5-5", "(one)", ""), reply("m2", "<synthetic>", "(limit)", `,"isApiErrorMessage":true`)}, "(one)", false},
		{"an API error", []string{reply("m1", "claude-opus-5-5", "(one)", ""), reply("m2", "claude-opus-5-5", "(error)", `,"isApiErrorMessage":true`)}, "(one)", false},
		{"a subagent's", []string{reply("m1", "claude-opus-5-5", "(one)", ""), strings.Replace(reply("m2", "claude-opus-5-5", "(two)", ""), `"isSidechain":false`, `"isSidechain":true`, 1)}, "(one)", false},
		{"blank text", []string{reply("m1", "claude-opus-5-5", "(one)", ""), reply("m2", "claude-opus-5-5", " \\n ", "")}, "(one)", false},
	} {
		d := digestOf(t, c.rows...)
		if d.Answer != c.answer || d.Unanswered != c.unanswered {
			t.Errorf("%s: %q %v", c.name, d.Answer, d.Unanswered)
		}
	}
}

// use is an assistant row that calls tool name with input (JSON).
func use(id, name, input string) string {
	return `{"isSidechain":false,"type":"assistant","message":{"id":"m` + id + `","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"` + id + `","name":"` + name + `","input":` + input + `}]}}`
}

// result is the user row with the result of tool call id.
func result(id, content string) string {
	return `{"isSidechain":false,"type":"user","message":{"role":"user","content":[{"tool_use_id":"` + id + `","type":"tool_result","content":` + content + `}]}}`
}

// The files Claude changed, least recently changed first, each once.
func TestReadDigestFiles(t *testing.T) {
	d := digestOf(t,
		use("t1", "Write", `{"file_path":"/p/a.go","content":"(a)"}`),
		use("t2", "Edit", `{"file_path":"/p/b.go","old_string":"(b)","new_string":"(c)"}`),
		use("t3", "NotebookEdit", `{"notebook_path":"/p/n.ipynb","new_source":"(cell)"}`),
		use("t4", "MultiEdit", `{"file_path":"/p/a.go","edits":[]}`),
		use("t5", "Read", `{"file_path":"/p/r.go"}`),
		use("t6", "Edit", `{"file_path":""}`),
		use("t7", "Edit", `{"file_path":7}`),
		strings.Replace(use("t8", "Write", `{"file_path":"/p/sub.go"}`), `"isSidechain":false`, `"isSidechain":true`, 1),
	)
	if want := []string{"/p/b.go", "/p/n.ipynb", "/p/a.go"}; !slices.Equal(d.Files, want) {
		t.Fatalf("%q", d.Files)
	}
}

// The todo list is replayed from the tool calls as Claude Code 2.1.310
// does when it resumes a conversation: tasks (TaskCreate, TaskUpdate) win
// over the last TodoWrite.
func TestReadDigestTodos(t *testing.T) {
	create := func(id, subject string) string {
		return use(id, "TaskCreate", `{"subject":"`+subject+`","description":"(what)"}`)
	}
	created := func(id, task string) string { return result(id, `"Task #`+task+` created successfully: (subject)"`) }
	update := func(id, input string) string { return use(id, "TaskUpdate", input) }
	todo := func(id string, items ...string) string {
		return use(id, "TodoWrite", `{"todos":[`+strings.Join(items, ",")+`]}`)
	}
	for _, c := range []struct {
		name string
		rows []string
		want []Todo
	}{
		{"the last TodoWrite", []string{
			todo("w1", `{"content":"(old)","status":"pending","activeForm":"(a)"}`),
			todo("w2", `{"content":"(new)","status":"in_progress","activeForm":"(a)"}`, `{"content":"","status":"pending","activeForm":"(a)"}`, `{"content":"(done)","status":"completed","activeForm":"(a)"}`),
		}, []Todo{{"(new)", "in_progress"}, {"(done)", "completed"}}},
		{"tasks win", []string{
			todo("w1", `{"content":"(todo)","status":"pending","activeForm":"(a)"}`),
			create("c1", "(one)"), created("c1", "1"),
		}, []Todo{{"(one)", "pending"}}},
		{"a result that names no task leaves it waiting, after the tasks", []string{
			create("c1", "(one)"), result("c1", `"(something else)"`),
			create("c2", "(two)"), result("c2", `[{"type":"text","text":"Task #2 created successfully: (two)"}]`),
		}, []Todo{{"(two)", "pending"}, {"(one)", "pending"}}},
		{"no result yet", []string{create("c1", "(one)")}, []Todo{{"(one)", "pending"}}},
		{"a failed create", []string{create("c1", "(one)"), result("c1", `"(refused)","is_error":true`)}, nil},
		{"no subject", []string{create("c1", ""), created("c1", "1")}, nil},
		{"updates by every name of the id", []string{
			create("c1", "(one)"), created("c1", "1"),
			create("c2", "(two)"), created("c2", "2"),
			create("c3", "(three)"), created("c3", "3"),
			update("u1", `{"taskId":"1","status":"completed"}`),
			update("u2", `{"id":"2","status":"in_progress","subject":"(two, renamed)"}`),
			update("u3", `{"task_id":"3","status":"deleted"}`),
		}, []Todo{{"(one)", "completed"}, {"(two, renamed)", "in_progress"}}},
		{"an unknown task is listed by its id", []string{update("u1", `{"taskId":"9","status":"in_progress"}`)}, []Todo{{"9", "in_progress"}}},
		{"a status Claude Code refuses", []string{create("c1", "(one)"), created("c1", "1"), update("u1", `{"taskId":"1","status":"bogus"}`)}, []Todo{{"(one)", "pending"}}},
		{"no id", []string{create("c1", "(one)"), created("c1", "1"), update("u1", `{"status":"completed"}`)}, []Todo{{"(one)", "pending"}}},
		{"an id that is no string", []string{create("c1", "(one)"), created("c1", "1"), update("u1", `{"taskId":1,"status":"completed"}`)}, []Todo{{"(one)", "pending"}}},
		{"the same id twice", []string{create("c1", "(one)"), created("c1", "1"), create("c2", "(again)"), created("c2", "1")}, []Todo{{"(one)", "pending"}}},
		{"a deleted task made again", []string{
			create("c1", "(one)"), created("c1", "1"), create("c2", "(two)"), created("c2", "2"),
			update("u1", `{"taskId":"1","status":"deleted"}`), update("u2", `{"taskId":"1","subject":"(back)"}`),
		}, []Todo{{"(two)", "pending"}, {"(back)", "pending"}}},
	} {
		if got := digestOf(t, c.rows...).Todos; !reflect.DeepEqual(got, c.want) && (len(got) > 0 || len(c.want) > 0) {
			t.Errorf("%s: %q", c.name, got)
		}
	}
}

// Rows the digest does not know, or that do not parse, are skipped.
func TestReadDigestSkipsOtherRows(t *testing.T) {
	d := digestOf(t,
		`{"type":"user","message":{"role":"user","content":"(first)"}}`,
		`{"type":"future-kind","value":{"type":"user","message":{"role":"user","content":"(not a prompt)"}}}`,
		`{"type":"assistant","message":{"content":"(not blocks)"}}`,
		`{"type":"assistant","message":{"model":"claude-opus-5-5","content":[{"type":"text","text":7}]}}`,
		`{"type":"ai-title","aiTitle":"(ai)"}`,
		`{"type":"user","message":{"role":"user","content":"(cut`,
	)
	if !slices.Equal(d.Prompts, []string{"(first)"}) || d.Answer != "" || d.Title != "(ai)" {
		t.Fatalf("%#v", d)
	}
	if d := digestOf(t, ""); d.Prompts != nil || d.Title != "" {
		t.Fatalf("an empty file: %#v", d)
	}
}
