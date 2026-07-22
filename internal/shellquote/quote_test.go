package shellquote

import "testing"

func TestQuote(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "''"},
		{"hostname", "hostname"},
		{"-lc", "-lc"},
		{"/tmp/a", "/tmp/a"},
		{"cd /tmp && pwd", "'cd /tmp && pwd'"},
		{"it's", `'it'\''s'`},
		{"a b", "'a b'"},
	}
	for _, tc := range cases {
		if got := Quote(tc.in); got != tc.want {
			t.Fatalf("Quote(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestJoinRemoteCommand(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "empty",
			args: nil,
			want: "",
		},
		{
			name: "single shell script kept raw",
			args: []string{"cd /home/worksapce/nem-panel && pwd && git rev-parse --show-toplevel"},
			want: "cd /home/worksapce/nem-panel && pwd && git rev-parse --show-toplevel",
		},
		{
			name: "bash -lc with script is quoted",
			args: []string{"bash", "-lc", "cd /home/worksapce/nem-panel && pwd"},
			want: "bash -lc 'cd /home/worksapce/nem-panel && pwd'",
		},
		{
			name: "spaces in path argument",
			args: []string{"ls", "/tmp/my dir"},
			want: "ls '/tmp/my dir'",
		},
		{
			name: "embedded single quote",
			args: []string{"echo", "it's"},
			want: `echo 'it'\''s'`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := JoinRemoteCommand(tc.args); got != tc.want {
				t.Fatalf("JoinRemoteCommand(%q)=%q want %q", tc.args, got, tc.want)
			}
		})
	}
}
