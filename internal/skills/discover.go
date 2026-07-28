package skills

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Skill is one agent skill directory that contains SKILL.md.
type Skill struct {
	Agent       string // claude / cursor / codex / local
	Name        string
	Description string
	Path        string
	HasBin      bool
}

// Discover finds skills under known agent skill roots and beside the current executable.
func Discover(exePath string) ([]Skill, error) {
	roots := defaultRoots()
	if beside := skillRootBesideExe(exePath); beside != "" {
		// Parent may be an agent skills root (e.g. ~/.claude/skills) — scan siblings as local too.
		roots = append(roots, rootRef{agent: "local", path: filepath.Dir(beside)})
	}

	seen := make(map[string]struct{})
	out := make([]Skill, 0)
	for _, r := range roots {
		entries, err := os.ReadDir(r.path)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			dir := filepath.Join(r.path, e.Name())
			sk, ok := loadSkill(r.agent, dir)
			if !ok {
				continue
			}
			key := strings.ToLower(sk.Path)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, sk)
		}
	}
	return out, nil
}

// Filter keeps skills whose name, path base, or description contains q (case-insensitive).
func Filter(skills []Skill, q string) []Skill {
	q = strings.TrimSpace(q)
	if q == "" {
		return skills
	}
	ql := strings.ToLower(q)
	out := make([]Skill, 0, len(skills))
	for _, s := range skills {
		if strings.Contains(strings.ToLower(s.Name), ql) ||
			strings.Contains(strings.ToLower(filepath.Base(s.Path)), ql) ||
			strings.Contains(strings.ToLower(s.Description), ql) {
			out = append(out, s)
		}
	}
	return out
}

type rootRef struct {
	agent string
	path  string
}

func defaultRoots() []rootRef {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return nil
	}
	return []rootRef{
		{agent: "claude", path: filepath.Join(home, ".claude", "skills")},
		{agent: "cursor", path: filepath.Join(home, ".cursor", "skills")},
		{agent: "codex", path: filepath.Join(home, ".codex", "skills")},
	}
}

// skillRootBesideExe returns the skill root if exe is .../<skill>/bin/sshctl[.exe].
func skillRootBesideExe(exePath string) string {
	if exePath == "" {
		return ""
	}
	abs, err := filepath.Abs(exePath)
	if err != nil {
		abs = exePath
	}
	binDir := filepath.Dir(abs)
	if !strings.EqualFold(filepath.Base(binDir), "bin") {
		return ""
	}
	root := filepath.Dir(binDir)
	if _, err := os.Stat(filepath.Join(root, "SKILL.md")); err != nil {
		return ""
	}
	return root
}

func loadSkill(agent, dir string) (Skill, bool) {
	md := filepath.Join(dir, "SKILL.md")
	if _, err := os.Stat(md); err != nil {
		return Skill{}, false
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	name, desc := parseFrontmatter(md)
	if name == "" {
		name = filepath.Base(dir)
	}
	return Skill{
		Agent:       agent,
		Name:        name,
		Description: desc,
		Path:        abs,
		HasBin:      hasBin(dir),
	}, true
}

func hasBin(skillRoot string) bool {
	binDir := filepath.Join(skillRoot, "bin")
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() {
			return true
		}
	}
	return false
}

func parseFrontmatter(path string) (name, description string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	if !sc.Scan() || strings.TrimSpace(sc.Text()) != "---" {
		return "", ""
	}

	var (
		inDesc    bool
		descLines []string
	)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)
		if trim == "---" {
			break
		}
		if inDesc {
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || trim == "" {
				descLines = append(descLines, strings.TrimSpace(line))
				continue
			}
			inDesc = false
		}
		if strings.HasPrefix(trim, "name:") {
			name = strings.TrimSpace(strings.TrimPrefix(trim, "name:"))
			name = strings.Trim(name, `"'`)
			continue
		}
		if strings.HasPrefix(trim, "description:") {
			rest := strings.TrimSpace(strings.TrimPrefix(trim, "description:"))
			if rest == "|" || rest == ">" || rest == "" {
				inDesc = true
				continue
			}
			description = strings.Trim(rest, `"'`)
			continue
		}
	}
	if len(descLines) > 0 {
		description = strings.TrimSpace(strings.Join(descLines, " "))
	}
	return name, description
}
