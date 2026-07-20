package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNoteTitleKey(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Doc.md", "doc"},
		{"doc.md", "doc"},
		{"DOC.md.enc", "doc"},
		{"Doc", "doc"},
		{"a B.md", "a b"},
	}
	for _, c := range cases {
		if got := noteTitleKey(c.in); got != c.want {
			t.Errorf("noteTitleKey(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestNoteNameConflictCaseAndEnc(t *testing.T) {
	dir := t.TempDir()
	old := notesDirPath
	notesDirPath = dir
	t.Cleanup(func() { notesDirPath = old })

	if err := os.WriteFile(filepath.Join(dir, "Doc.md"), []byte("# Doc\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if !noteNameConflict("doc.md", "") {
		t.Fatal("doc.md should conflict with Doc.md")
	}
	if !noteNameConflict("DOC.md.enc", "") {
		t.Fatal("DOC.md.enc should conflict with Doc.md")
	}
	if noteNameConflict("Other.md", "") {
		t.Fatal("Other.md should be free")
	}
	if noteNameConflict("doc.md", filepath.Join(dir, "Doc.md")) {
		t.Fatal("case-only rename of Doc.md to doc.md should not conflict")
	}
	if noteNameConflict("Doc.md.enc", filepath.Join(dir, "Doc.md")) {
		t.Fatal("encrypting Doc.md to Doc.md.enc should not conflict with itself")
	}
}

func TestSortNotesByTitle(t *testing.T) {
	notes := []noteInfo{
		{Name: "zebra.md"},
		{Name: "Apple.md"},
		{Name: "banana.md.enc"},
	}
	sortNotesByTitle(notes)
	if notes[0].Name != "Apple.md" || notes[1].Name != "banana.md.enc" || notes[2].Name != "zebra.md" {
		t.Fatalf("order = %v, %v, %v", notes[0].Name, notes[1].Name, notes[2].Name)
	}
}
