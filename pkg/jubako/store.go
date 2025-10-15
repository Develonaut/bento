package jubako

import (
	"os"
	"path/filepath"
	"strings"

	"bento/pkg/neta"
)

// Store manages bento file storage.
type Store struct {
	workDir string
	parser  *Parser
}

// NewStore creates a new store.
func NewStore(workDir string) (*Store, error) {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}

	return &Store{
		workDir: workDir,
		parser:  NewParser(),
	}, nil
}

// Load reads a bento by name.
func (s *Store) Load(name string) (neta.Definition, error) {
	path := s.pathFor(name)
	return s.parser.Parse(path)
}

// Save writes a bento to disk.
func (s *Store) Save(name string, def neta.Definition) error {
	path := s.pathFor(name)

	data, err := s.parser.Format(def)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Delete removes a bento.
func (s *Store) Delete(name string) error {
	path := s.pathFor(name)
	return os.Remove(path)
}

// List returns all bentos in the store.
func (s *Store) List() ([]BentoInfo, error) {
	pattern := filepath.Join(s.workDir, "*.bento.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	infos := make([]BentoInfo, 0, len(matches))
	for _, path := range matches {
		info, err := s.getInfo(path)
		if err != nil {
			continue // Skip invalid files
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// pathFor returns the file path for a bento name.
func (s *Store) pathFor(name string) string {
	if !strings.HasSuffix(name, ".bento.yaml") {
		name += ".bento.yaml"
	}
	return filepath.Join(s.workDir, name)
}

// getInfo extracts bento info from a file.
func (s *Store) getInfo(path string) (BentoInfo, error) {
	def, err := s.parser.Parse(path)
	if err != nil {
		return BentoInfo{}, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		return BentoInfo{}, err
	}

	return BentoInfo{
		Name:     filepath.Base(path),
		Path:     path,
		Type:     def.Type,
		Modified: stat.ModTime(),
	}, nil
}
