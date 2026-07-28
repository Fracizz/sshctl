---
name: sshctl
description: |
  sshctl 远程主机 CLI（清单 search/exec/shell/scp/add/migrate/skills）。优先 sshctl，尽量不用原生 ssh/scp。
  二进制与技能同目录（SKILL.md 所在文件夹下的 bin/sshctl.exe），不安装到系统 PATH。
  触发词：sshctl、search -s、exec、shell、scp、skills、servers.json、SSHCTL、.sshctl。
---

# sshctl · 远程主机 CLI（Windows）

**默认用 sshctl 做所有远程操作**；不要手写 `ssh` / `scp` / `sshpass`（配免密见 `ssh-key-auth-setup` 技能）。

**Agent 必须**通过技能目录下的二进制调用，**不要**假设 `sshctl` 在系统 PATH 中。

```powershell
# skillRoot = 本 SKILL.md 所在目录（从附带技能文件路径解析，禁止硬编码仓库路径）
$skillRoot = '...'  # 例：skills/sshctl/ 或 %USERPROFILE%\.claude\skills\sshctl\
$sshctl = Join-Path $skillRoot 'bin\sshctl.exe'
& $sshctl version
```

| 操作 | 命令 |
|------|------|
| 搜主机 | `& $sshctl search -s <关键词>` |
| 远程执行 | `& $sshctl exec <host> -- <cmd>` |
| 交互 shell | `& $sshctl shell <host>` |
| 传文件 | `& $sshctl scp <src> <dst>` |
| 写入清单 | `& $sshctl add --host ... --user ... --password '...'` |
| 清单迁移 | `& $sshctl migrate` |
| 列 skills | `& $sshctl skills` / `& $sshctl skills -s sshctl` |

---

## 二进制（与 skill 同目录）

| 项 | 路径 |
|----|------|
| skillRoot | 本 `SKILL.md` 所在文件夹 |
| 二进制 | `$skillRoot\bin\sshctl.exe` |

### 构建 / 更新

在仓库根目录交叉编译 6 平台可执行文件，同步到仓库 skill 与已存在的 `.claude` / `.cursor` / `.codex` skill `bin\`：

```powershell
$env:VERSION = '0.2.6'
.\scripts\build.ps1
```

`bin/` 下二进制 **不入库**。也可从 [Releases](https://github.com/Fracizz/sshctl/releases) 获取：

| 资源 | 说明 |
|------|------|
| `sshctl-skill.zip` | **AI skills 整包**：解压到 `~/.claude/skills/` / `~/.cursor/skills/` / `~/.codex/skills/` |
| `sshctl-windows-amd64.exe` | Windows x64（Agent 默认可作 `bin/sshctl.exe`） |
| `sshctl-windows-arm64.exe` | Windows ARM64 |
| `sshctl-linux-amd64` / `sshctl-linux-arm64` | Linux x64 / ARM64 |
| `sshctl-darwin-amd64` / `sshctl-darwin-arm64` | macOS Intel / Apple Silicon |

多平台发布：推送 `v*` 标签，由 GitHub Actions 构建裸二进制 + `sshctl-skill.zip`。

### 验证

```powershell
& $sshctl version    # 0.2.6+
& $sshctl skills -s sshctl
& $sshctl list
```

---

## 配置

| 项 | 路径 / 变量 |
|----|-------------|
| 清单 | `%USERPROFILE%\.sshctl\servers.json` |
| 迁移 | `& $sshctl migrate`（`~/.sshfrac` → `~/.sshctl`，旧文件改 `.bak`） |
| 覆盖 | `$SSHCTL_CONFIG` |
| Legacy | `$SSHFRAC_CONFIG`（显式指定时） |

**规则：** 每个 IP 仅一条；`add` 同 IP 覆盖；密码特殊字符须完整引号包裹。

---

## 常用命令

```powershell
& $sshctl migrate
& $sshctl list
& $sshctl search -s 192.168
& $sshctl add --host 192.168.x.x --user administrator --password '...' --os Windows --desc "说明"
& $sshctl exec 192.168.x.x -- "hostname && whoami"
# 多参数会做 shell quote，可安全使用 bash -lc
& $sshctl exec 192.168.x.x -- bash -lc 'cd /tmp && pwd'
& $sshctl scp .\a.txt 192.168.x.x:C:/temp/a.txt
```

**Agent 流程：** `search -s` → 不在清单则 `add`（向用户确认凭据）→ `exec` / `scp`。首次连某主机若 host key 失败，见下节。

---

## 首次连接 / host key

默认校验 OpenSSH `~/.ssh/known_hosts`。新主机未写入时，`exec` / `scp` / `shell` 会失败（`unknown host key` / Handshake failed）。sshctl **不会**交互确认或自动写入 host key。

| 场景 | 处理 |
|------|------|
| 可信实验网 / 临时连通测试 | `& $sshctl --insecure exec <host> -- "..."`（跳过校验；**不**写入 known_hosts） |
| 正式环境 | 先用本机 OpenSSH 连一次写入 known_hosts，之后正常 `exec`，无需 `--insecure` |

```powershell
# 首次连通（仅可信环境）
& $sshctl --insecure exec 192.168.x.x -- "hostname && uname -r"
# known_hosts 已有记录后
& $sshctl exec 192.168.x.x -- "hostname"
```

`--insecure` 是跳过校验，不是「自动接受并永久信任」。

---

## 错误速查

| 情况 | 处理 |
|------|------|
| 找不到 sshctl | 构建/复制到 `$skillRoot\bin\sshctl.exe` |
| duplicate host | `add` 同 IP 覆盖，或删 JSON 重复项 |
| Windows 密码失败 | 确认密码完整；`--os Windows`；v0.2.1+ |
| `bash -lc` 路径错乱 | 需 v0.2.3+（多参数已 shell quote） |
| unknown host key / Handshake failed | 首次连接；可信环境加 `--insecure`，或先写入 known_hosts |
| 仍用旧清单目录 | `& $sshctl migrate` |

## 边界

- 远程操作只用 `$sshctl`，不用原生 ssh/scp
- **不**安装到系统 PATH（技能工作流）
- 配免密 → **ssh-key-auth-setup**
- 清单不入库
- 非可信网络避免 `--insecure`
