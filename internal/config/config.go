package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Fracizz/sshfrac/internal/crypto"
)

// Server is one SSH endpoint entry.
type Server struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	OS          string `json:"os"`
	KeyFile     string `json:"key_file,omitempty"`
}

// File is the on-disk JSON document.
type File struct {
	Servers []Server `json:"servers"`
}

// DefaultConfigPath returns ~/.sshfrac/servers.json (outside any repo).
// Falls back to legacy ~/.invossh or ~/.sshctl if present.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".sshfrac", "servers.json")
	}
	primary := filepath.Join(home, ".sshfrac", "servers.json")
	if _, err := os.Stat(primary); err == nil {
		return primary
	}
	for _, dir := range []string{".invossh", ".sshctl"} {
		legacy := filepath.Join(home, dir, "servers.json")
		if _, err := os.Stat(legacy); err == nil {
			return legacy
		}
	}
	return primary
}

// ResolvePath picks config path: flag > SSHFRAC_CONFIG > INVOSSH_CONFIG > SSHCTL_CONFIG > default.
func ResolvePath(flagPath string) string {
	if flagPath != "" {
		return flagPath
	}
	for _, key := range []string{"SSHFRAC_CONFIG", "INVOSSH_CONFIG", "SSHCTL_CONFIG"} {
		if env := os.Getenv(key); env != "" {
			return env
		}
	}
	return DefaultConfigPath()
}

// Load reads JSON, encrypts any plaintext passwords, and rewrites the file when needed.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	changed := false
	for i := range f.Servers {
		s := &f.Servers[i]
		if s.Port == 0 {
			s.Port = 22
		}
		if s.Password == "" || crypto.IsEncrypted(s.Password) {
			continue
		}
		enc, err := crypto.Encrypt(s.Password)
		if err != nil {
			return nil, fmt.Errorf("encrypt password for %s: %w", s.Name, err)
		}
		s.Password = enc
		changed = true
	}
	if changed {
		if err := Save(path, &f); err != nil {
			return nil, err
		}
	}
	return &f, nil
}

// Save writes JSON with indentation.
func Save(path string, f *File) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Find locates a server by name, host, or "user@host".
// Exact match first; if none, case-insensitive substring on name/host (must be unique).
func (f *File) Find(query string) (*Server, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("empty server query")
	}
	userHint := ""
	hostHint := q
	if strings.Contains(q, "@") {
		parts := strings.SplitN(q, "@", 2)
		userHint, hostHint = parts[0], parts[1]
	}
	for i := range f.Servers {
		s := &f.Servers[i]
		if s.Name == q || s.Host == q || s.Host == hostHint {
			if userHint != "" && s.User != "" && s.User != userHint {
				continue
			}
			return s, nil
		}
		if s.User+"@"+s.Host == q {
			return s, nil
		}
	}
	hits := f.Search(q)
	if userHint != "" {
		filtered := hits[:0]
		for _, s := range hits {
			if s.User == userHint {
				filtered = append(filtered, s)
			}
		}
		hits = filtered
	}
	switch len(hits) {
	case 1:
		return hits[0], nil
	case 0:
		return nil, fmt.Errorf("server not found: %s", query)
	default:
		return nil, fmt.Errorf("ambiguous server %q: %d matches (use sshfrac search -s %q)", query, len(hits), query)
	}
}

// Search returns servers whose name, host, or description contains keyword (case-insensitive).
func (f *File) Search(keyword string) []*Server {
	k := strings.ToLower(strings.TrimSpace(keyword))
	if k == "" {
		out := make([]*Server, 0, len(f.Servers))
		for i := range f.Servers {
			out = append(out, &f.Servers[i])
		}
		return out
	}
	var out []*Server
	for i := range f.Servers {
		s := &f.Servers[i]
		if strings.Contains(strings.ToLower(s.Name), k) ||
			strings.Contains(strings.ToLower(s.Host), k) ||
			strings.Contains(strings.ToLower(s.Description), k) {
			out = append(out, s)
		}
	}
	return out
}

// PlainPassword decrypts the stored password for runtime use.
func (s *Server) PlainPassword() (string, error) {
	return crypto.Decrypt(s.Password)
}

// Add appends a server, encrypting password immediately.
func (f *File) Add(s Server) error {
	if s.Port == 0 {
		s.Port = 22
	}
	if s.Password != "" && !crypto.IsEncrypted(s.Password) {
		enc, err := crypto.Encrypt(s.Password)
		if err != nil {
			return err
		}
		s.Password = enc
	}
	f.Servers = append(f.Servers, s)
	return nil
}
