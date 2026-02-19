package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = "gr-state.json"

type ReviewSession struct {
	Branch       string    `json:"branch"`
	WorktreePath string    `json:"worktreePath"`
	StartedAt    time.Time `json:"startedAt"`
}

type State struct {
	Reviews []ReviewSession `json:"reviews"`
}

func statePath(commonDir string) string {
	return filepath.Join(commonDir, stateFileName)
}

func Load(commonDir string) (*State, error) {
	data, err := os.ReadFile(statePath(commonDir))
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *State) Save(commonDir string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(commonDir), data, 0644)
}

func (s *State) Add(session ReviewSession) {
	s.Reviews = append(s.Reviews, session)
}

func (s *State) Remove(branch string) {
	filtered := s.Reviews[:0]
	for _, r := range s.Reviews {
		if r.Branch != branch {
			filtered = append(filtered, r)
		}
	}
	s.Reviews = filtered
}

func (s *State) FindByBranch(branch string) *ReviewSession {
	for i := range s.Reviews {
		if s.Reviews[i].Branch == branch {
			return &s.Reviews[i]
		}
	}
	return nil
}

func (s *State) FindByCwd(cwd string) *ReviewSession {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil
	}
	for i := range s.Reviews {
		absPath, err := filepath.Abs(s.Reviews[i].WorktreePath)
		if err != nil {
			continue
		}
		if absCwd == absPath {
			return &s.Reviews[i]
		}
	}
	return nil
}

func (s *State) CleanStale() {
	filtered := s.Reviews[:0]
	for _, r := range s.Reviews {
		if _, err := os.Stat(r.WorktreePath); err == nil {
			filtered = append(filtered, r)
		}
	}
	s.Reviews = filtered
}
