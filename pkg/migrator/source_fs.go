package migrator

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileSource reads migrations from a directory.
type FileSource struct {
	Dir string
}

func (s FileSource) Migrations() ([]Migration, error) {
	return readMigrations(os.DirFS(s.Dir), ".")
}

// EmbeddedSource reads migrations from an embedded filesystem.
func (s EmbeddedSource) Migrations() ([]Migration, error) {
	dir := s.Dir
	if dir == "" {
		dir = "."
	}
	return readMigrations(s.FS, dir)
}

func readMigrations(fsys fs.FS, root string) ([]Migration, error) {
	tmp := map[string]*Migration{}
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".surql") {
			return nil
		}
		full := filepath.ToSlash(path)
		content, readErr := fs.ReadFile(fsys, full)
		if readErr != nil {
			return readErr
		}
		base := name
		isUp := strings.HasSuffix(base, ".up.surql")
		isDown := strings.HasSuffix(base, ".down.surql")
		var id string
		switch {
		case isUp:
			id = strings.TrimSuffix(base, ".up.surql")
		case isDown:
			id = strings.TrimSuffix(base, ".down.surql")
		default:
			id = strings.TrimSuffix(base, ".surql")
		}
		m, ok := tmp[id]
		if !ok {
			m = &Migration{ID: id}
			tmp[id] = m
		}
		if isDown {
			m.DownSQL = string(content)
		} else {
			m.UpSQL = string(content)
		}
		return nil
	}
	if err := fs.WalkDir(fsys, root, walk); err != nil {
		return nil, err
	}

	out := make([]Migration, 0, len(tmp))
	for _, m := range tmp {
		m.Checksum = checksum(m.UpSQL + "\n" + m.DownSQL)
		out = append(out, *m)
	}
	// sort by ID (timestamp prefix ensures order)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func checksum(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
