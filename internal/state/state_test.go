package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(s.Reviews))
	}
}

func TestAddAndSave(t *testing.T) {
	dir := t.TempDir()
	s := &State{}
	s.Add(ReviewSession{
		Branch:       "feature-a",
		WorktreePath: "/tmp/wt/feature-a",
		StartedAt:    time.Now(),
	})

	if err := s.Save(dir); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(loaded.Reviews))
	}
	if loaded.Reviews[0].Branch != "feature-a" {
		t.Fatalf("expected branch feature-a, got %s", loaded.Reviews[0].Branch)
	}
}

func TestRemove(t *testing.T) {
	s := &State{}
	s.Add(ReviewSession{Branch: "a", WorktreePath: "/tmp/a", StartedAt: time.Now()})
	s.Add(ReviewSession{Branch: "b", WorktreePath: "/tmp/b", StartedAt: time.Now()})

	s.Remove("a")
	if len(s.Reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(s.Reviews))
	}
	if s.Reviews[0].Branch != "b" {
		t.Fatalf("expected branch b, got %s", s.Reviews[0].Branch)
	}
}

func TestFindByBranch(t *testing.T) {
	s := &State{}
	s.Add(ReviewSession{Branch: "feature-x", WorktreePath: "/tmp/x", StartedAt: time.Now()})

	found := s.FindByBranch("feature-x")
	if found == nil {
		t.Fatal("expected to find session")
	}

	notFound := s.FindByBranch("nonexistent")
	if notFound != nil {
		t.Fatal("expected nil for nonexistent branch")
	}
}

func TestFindByCwd(t *testing.T) {
	dir := t.TempDir()
	wtPath := filepath.Join(dir, "wt")
	os.MkdirAll(wtPath, 0755)

	s := &State{}
	s.Add(ReviewSession{Branch: "test", WorktreePath: wtPath, StartedAt: time.Now()})

	found := s.FindByCwd(wtPath)
	if found == nil {
		t.Fatal("expected to find session by cwd")
	}

	notFound := s.FindByCwd("/nonexistent")
	if notFound != nil {
		t.Fatal("expected nil for nonexistent cwd")
	}
}

func TestCleanStale(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing")
	os.MkdirAll(existing, 0755)

	s := &State{}
	s.Add(ReviewSession{Branch: "valid", WorktreePath: existing, StartedAt: time.Now()})
	s.Add(ReviewSession{Branch: "stale", WorktreePath: "/nonexistent/path", StartedAt: time.Now()})

	s.CleanStale()
	if len(s.Reviews) != 1 {
		t.Fatalf("expected 1 review after cleanup, got %d", len(s.Reviews))
	}
	if s.Reviews[0].Branch != "valid" {
		t.Fatalf("expected valid branch, got %s", s.Reviews[0].Branch)
	}
}
