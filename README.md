# invossh

> Renamed from `sshctl` to avoid widespread name collisions (PyPI/AUR/GitHub).

**English** | [中文](#invossh-中文)

High-performance cross-platform SSH/SCP CLI with an encrypted local server inventory.

> **Built primarily for AI agents** (Cursor, Claude Code, Codex, and similar).
> Stable, scriptable subcommands and JSON inventory make remote exec / SCP easy for tools to call without interactive prompts.

Humans can use it too — the UX prioritizes machine-friendly flags, exit codes, and search (`-s`) over a rich TUI.

## Features

- Single static binaries: Linux / Windows / macOS (amd64 + arm64)
- Remote command execution and interactive shells
- Fast file transfer via SFTP (`scp` subcommand, directories supported)
- AES-256-GCM encrypted passwords at rest (`enc:v1:`)
- Case-insensitive contains search on name / IP / description
- Host key verification via `~/.ssh/known_hosts` (optional `--insecure` for labs)
- Default config outside the repo: `~/.invossh/servers.json`

## Why AI-first?

| Need | How invossh helps |
|------|------------------|
| Non-interactive | `exec` / `scp` / `list` / `search` — no TTY prompts for common flows |
| Discover hosts | `search -s <keyword>` on name, IP, description |
| Safe-ish secrets | Passwords encrypted on disk; agents pass `--password` only when adding |
| Predictable I/O | Tabular list output + non-zero exit on remote failure |

## Install

```bash
go install github.com/Fracizz/invossh@latest
```

Or from source:

```bash
git clone https://github.com/Fracizz/invossh.git
cd invossh
go build -o bin/invossh .
```

## Quick start

```bash
# Create ~/.invossh/servers.json (placeholder only — no real secrets)
invossh init

# Add a host (password is encrypted on save)
invossh add --name lab --host 192.0.2.10 --user root --password 'secret' --desc "lab box"

# Or key-based auth
invossh add --name prod --host 192.0.2.11 --user root --key ~/.ssh/id_ed25519 --desc "prod"

invossh list
invossh search -s 192.0.2
invossh exec lab -- hostname
invossh shell lab
invossh scp ./app.tar.gz lab:/tmp/app.tar.gz
```

Config path priority: `--config` > `$INVOSSH_CONFIG` > `~/.invossh/servers.json`

### Master password (enc:v2)

```bash
export INVOSSH_MASTER_PASSWORD='your-strong-secret'
# optional: also bind ciphertext to this machine
export INVOSSH_BIND_MACHINE=1

invossh add --name lab --host 192.0.2.10 --user root --password 'host-pass' --desc "lab"
invossh exec lab -- hostname
```

Or pass `--master-password` / `--bind-machine` on each invocation. See [SECURITY.md](SECURITY.md).

### Exit codes & completion

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | local error |
| 2 | usage (documented; reserved) |
| N | remote process status (`exec`) |

```bash
invossh version
invossh completion powershell > invossh.ps1
invossh completion bash > /etc/bash_completion.d/invossh
```

First SSH to a host should populate OpenSSH `known_hosts` (or pass `--insecure` only in trusted labs).

## JSON format

```json
{
  "servers": [
    {
      "name": "example-host",
      "description": "demo entry",
      "host": "192.0.2.10",
      "port": 22,
      "user": "root",
      "password": "enc:v1:...",
      "os": "Linux"
    }
  ]
}
```

See `configs/servers.example.json` for placeholders. Never commit real passwords.

## Security

Read [SECURITY.md](SECURITY.md). Ciphertext is bound to the local machine user material; moving the file to another PC requires re-entering passwords.

## License

MIT — see [LICENSE](LICENSE).

---

# invossh (中文)

高性能跨平台 SSH/SCP 命令行工具，带本地加密主机清单。

> **主要面向 AI 编程助手使用**（Cursor、Claude Code、Codex 等）。
> 子命令稳定、可脚本化，配合 JSON 主机清单，方便 Agent 非交互地执行远程命令与文件传输。

人工使用同样可以；交互式体验让位于：机器友好的参数、退出码，以及 `-s` 模糊搜索。

## 特性

- 单文件二进制：Linux / Windows / macOS（amd64 + arm64）
- 远程命令执行、交互 Shell
- 基于 SFTP 的高速传文件（`scp` 子命令，支持目录）
- 密码落盘 AES-256-GCM 加密（`enc:v1:`）
- 对 name / IP / description 不区分大小写的包含搜索
- 默认校验 `~/.ssh/known_hosts`（实验环境可用 `--insecure`）
- 配置默认在仓库外：`~/.invossh/servers.json`

## 为何适合 AI？

| 需求 | invossh 的做法 |
|------|----------------|
| 非交互 | `exec` / `scp` / `list` / `search`，常见流程不弹密码提示 |
| 找机器 | `search -s <关键词>` 匹配名称、IP、描述 |
| 密钥/密码 | 密码加密存储；Agent 仅在 `add` 时传入一次 `--password` |
| 输出稳定 | 表格列表；远端失败返回非 0 退出码 |

## 安装

```bash
go install github.com/Fracizz/invossh@latest
```

或源码编译：

```bash
git clone https://github.com/Fracizz/invossh.git
cd invossh
go build -o bin/invossh .
```

## 快速开始

```bash
invossh init
invossh add --name lab --host 192.0.2.10 --user root --password 'secret' --desc "实验机"
invossh add --name prod --host 192.0.2.11 --user root --key ~/.ssh/id_ed25519 --desc "生产"

invossh list
invossh search -s 192.0.2
invossh exec lab -- hostname
invossh shell lab
invossh scp ./app.tar.gz lab:/tmp/app.tar.gz
```

配置优先级：`--config` > `$INVOSSH_CONFIG` > `~/.invossh/servers.json`

### 主密码（enc:v2）

```bash
export INVOSSH_MASTER_PASSWORD='强密码'
export INVOSSH_BIND_MACHINE=1   # 可选：绑定本机

invossh add --name lab --host 192.0.2.10 --user root --password 'host-pass' --desc "实验机"
```

详见 [SECURITY.md](SECURITY.md)。

### 退出码与补全

| 码 | 含义 |
|----|------|
| 0 | 成功 |
| 1 | 本地错误 |
| N | 远端进程状态（`exec`） |

```bash
invossh version
invossh completion powershell > invossh.ps1
```

首次连接请先写入 OpenSSH `known_hosts`（或仅在可信实验环境使用 `--insecure`）。

## 安全说明

详见 [SECURITY.md](SECURITY.md)。密文与本机用户环境绑定，换机器需重新录入密码。

## 许可证

MIT — 见 [LICENSE](LICENSE)。

