package sources

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vshulcz/deja-vu/internal/model"
)

// A Gemini tool call whose result reports an error changed nothing (#4974):
// the generic reader took every edit call it found, so a `replace` Gemini
// refused showed as code the session wrote, beside the edit that succeeded.
func TestGenericChangesSkipsFailedEditCalls(t *testing.T) {
	root := t.TempDir()
	t.Setenv("DEJA_GEMINI_ROOT", root)
	_, chats := geminiTree(t)
	lines := `{"sessionId":"sess-failed-edit-1","projectHash":"abc","startTime":"2026-07-20T21:07:00.000Z","lastUpdated":"2026-07-20T21:08:00.000Z","kind":"main"}
{"id":"u1","timestamp":"2026-07-20T21:07:01.000Z","type":"user","content":[{"text":"edit pool.go"}]}
{"id":"g1","timestamp":"2026-07-20T21:07:02.000Z","type":"gemini","content":"","model":"gemini-3","toolCalls":[{"id":"replace__1","name":"replace","args":{"file_path":"/Users/x/coding/gateway/pool.go","old_string":"MaxConns: 10","new_string":"MaxConns: 40"},"result":[{"functionResponse":{"id":"replace__1","name":"replace","response":{"error":"old_string not found"}}}],"status":"error"},{"id":"replace__2","name":"replace","args":{"file_path":"/Users/x/coding/gateway/pool.go","old_string":"Timeout: 1s","new_string":"Timeout: 5s"},"result":[{"functionResponse":{"id":"replace__2","name":"replace","response":{"output":"ok"}}}],"status":"success"}]}
{"id":"u2","timestamp":"2026-07-20T21:07:03.000Z","type":"user","content":[{"functionResponse":{"id":"replace__1","name":"replace","response":{"error":"old_string not found"}}},{"functionResponse":{"id":"replace__2","name":"replace","response":{"output":"ok"}}}]}
{"id":"g2","timestamp":"2026-07-20T21:07:04.000Z","type":"gemini","content":"Done.","model":"gemini-3"}
`
	p := filepath.Join(chats, "session-2026-07-20T21-07-sess-failed-edit-1.jsonl")
	if err := os.WriteFile(p, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range Changes(model.Session{Harness: "gemini", ID: "sess-failed-edit-1", Path: p}) {
		got = append(got, c.Old+"|"+c.New)
	}
	if want := "Timeout: 1s|Timeout: 5s"; strings.Join(got, " ") != want {
		t.Fatalf("got  %q\nwant %q", strings.Join(got, " "), want)
	}
}
