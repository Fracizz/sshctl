package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	dir := t.TempDir()
	md := filepath.Join(dir, "SKILL.md")
	content := "---\nname: sshctl\ndescription: |\n  line one\n  line two\n---\n\n# Body\n"
	if err := os.WriteFile(md, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	name, desc := parseFrontmatter(md)
	if name != "sshctl" {
		t.Fatalf("name=%q", name)
	}
	if desc != "line one line two" {
		t.Fatalf("desc=%q", desc)
	}
}

func TestDiscoverAndFilter(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home) // Windows
	t.Setenv("HOME", home)

	claude := filepath.Join(home, ".claude", "skills", "sshctl")
	if err := os.MkdirAll(filepath.Join(claude, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: sshctl\ndescription: remote CLI\n---\n"
	if err := os.WriteFile(filepath.Join(claude, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claude, "bin", "sshctl.exe"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Force defaultRoots to use our home by monkeying via Discover only on home-based paths.
	// Discover uses os.UserHomeDir — set HOME/USERPROFILE above; on Windows UserHomeDir uses USERPROFILE.
	got, err := Discover("")
	if err != nil {
		t.Fatal(err)
	}
	found := Filter(got, "sshctl")
	if len(found) == 0 {
		t.Fatalf("expected sshctl skill, got %#v", got)
	}
	var hit *Skill
	for i := range found {
		if found[i].Name == "sshctl" && found[i].HasBin {
			hit = &found[i]
			break
		}
	}
	if hit == nil {
		t.Fatalf("no sshctl with bin: %#v", found)
	}
}
