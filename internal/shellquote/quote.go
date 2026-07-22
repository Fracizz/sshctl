// Package shellquote 为远程经 shell 执行的命令做 POSIX 风格转义。
package shellquote

import (
	"strings"
	"unicode"
)

// Quote 将单个参数转为可安全嵌入远程 shell 命令行的形式。
func Quote(s string) string {
	if s == "" {
		return "''"
	}
	if isSafeUnquoted(s) {
		return s
	}
	// POSIX：单引号内除 ' 外原样保留；' 写成 '\'' 
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// JoinRemoteCommand 拼接远程命令。
//
// 单参数原样返回，保留 "cd /tmp && pwd" 这类需经 shell 解释的写法；
// 多参数逐项 Quote 后空格连接，避免 bash -lc 等场景被拆坏。
func JoinRemoteCommand(args []string) string {
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	default:
		parts := make([]string, len(args))
		for i, arg := range args {
			parts[i] = Quote(arg)
		}
		return strings.Join(parts, " ")
	}
}

func isSafeUnquoted(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		switch r {
		case '@', '%', '+', '=', ':', ',', '.', '/', '-', '_':
			continue
		default:
			return false
		}
	}
	return true
}
