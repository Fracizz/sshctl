package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fracizz/sshfrac/internal/config"
	"github.com/Fracizz/sshfrac/internal/crypto"
)

func TestSearchCaseInsensitiveContains(t *testing.T) {
	f := &config.File{Servers: []config.Server{
		{Name: "Lab-Alpha", Host: "192.0.2.10", Description: "Primary LAB"},
		{Name: "other", Host: "198.51.100.1", Description: "unused"},
	}}
	hits := f.Search("lab")
	if len(hits) != 1 || hits[0].Host != "192.0.2.10" {
		t.Fatalf("search lab: got %#v", hits)
	}
	hits = f.Search("192.0.2")
	if len(hits) != 1 {
		t.Fatalf("search ip fragment: got %d", len(hits))
	}
}

func TestFindExactThenFuzzy(t *testing.T) {
	f := &config.File{Servers: []config.Server{
		{Name: "web", Host: "192.0.2.10", User: "root"},
		{Name: "db", Host: "192.0.2.20", User: "root"},
	}}
	s, err := f.Find("web")
	if err != nil || s.Host != "192.0.2.10" {
		t.Fatalf("exact: %v %#v", err, s)
	}
	s, err = f.Find("192.0.2.20")
	if err != nil || s.Name != "db" {
		t.Fatalf("host: %v %#v", err, s)
	}
	if _, err := f.Find("192.0.2"); err == nil {
		t.Fatal("expected ambiguous error")
	}
}

func TestEncryptRoundTripOnSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "servers.json")
	f := &config.File{}
	if _, err := f.Add(config.Server{Name: "t", Host: "192.0.2.10", User: "root", Password: "secret", OS: "Linux"}); err != nil {
		t.Fatal(err)
	}
	if !crypto.IsEncrypted(f.Servers[0].Password) {
		t.Fatal("expected encrypted password after Add")
	}
	if err := config.Save(path, f); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := loaded.Servers[0].PlainPassword()
	if err != nil || plain != "secret" {
		t.Fatalf("decrypt: %v %q", err, plain)
	}
}

func TestAddReplacesDuplicateHost(t *testing.T) {
	f := &config.File{}
	if updated, err := f.Add(config.Server{Name: "a", Host: "192.0.2.10", User: "root", Password: "one", OS: "Linux"}); err != nil || updated {
		t.Fatalf("first add: updated=%v err=%v", updated, err)
	}
	if updated, err := f.Add(config.Server{Name: "b", Host: "192.0.2.10", User: "admin", Password: "two", OS: "Windows"}); err != nil || !updated {
		t.Fatalf("second add: updated=%v err=%v", updated, err)
	}
	if len(f.Servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(f.Servers))
	}
	if f.Servers[0].Name != "b" || f.Servers[0].User != "admin" {
		t.Fatalf("unexpected server: %#v", f.Servers[0])
	}
	if err := f.ValidateUniqueHosts(); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRejectsDuplicateHost(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "servers.json")
	raw := `{
  "servers": [
    {"name":"a","host":"192.0.2.10","port":22,"user":"root","password":"","os":"Linux"},
    {"name":"b","host":"192.0.2.10","port":22,"user":"admin","password":"","os":"Windows"}
  ]
}
`
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err == nil {
		t.Fatal("expected duplicate host error")
	}
}

func TestDefaultConfigPathOutsideCwd(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	p := config.DefaultConfigPath()
	if filepath.Base(p) != "servers.json" {
		t.Fatalf("base: %s", p)
	}
	dir := filepath.Base(filepath.Dir(p))
	if dir != ".sshfrac" && dir != ".invossh" && dir != ".sshctl" {
		t.Fatalf("dir: %s", p)
	}
}

func TestResolvePathPriority(t *testing.T) {
	t.Setenv("SSHFRAC_CONFIG", "")
	t.Setenv("INVOSSH_CONFIG", "")
	t.Setenv("SSHCTL_CONFIG", "")
	if got := config.ResolvePath("/tmp/custom.json"); got != "/tmp/custom.json" {
		t.Fatalf("flag: %s", got)
	}
	t.Setenv("SSHFRAC_CONFIG", "/env/servers.json")
	if got := config.ResolvePath(""); got != "/env/servers.json" {
		t.Fatalf("env: %s", got)
	}
	t.Setenv("SSHFRAC_CONFIG", "")
	t.Setenv("SSHCTL_CONFIG", "/legacy/servers.json")
	if got := config.ResolvePath(""); got != "/legacy/servers.json" {
		t.Fatalf("legacy env: %s", got)
	}
}
