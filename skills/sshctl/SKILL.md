---
name: sshctl
description: |
  sshctl 远程主机 CLI（清单 search/exec/shell/scp/add/migrate）。优先 sshctl，尽量不用原生 ssh/scp。
  二进制与技能同目录（SKILL.md 所在文件夹下的 bin/sshctl.exe），不安装到系统 PATH。
  触发词：sshctl、search -s、exec、shell、scp、servers.json、SSHCTL、.sshctl。
---

# sshctl · 远程主机 CLI（Windows）

**默认用 sshctl 做所有远程操作**；不要手写 `ssh` / `scp` / `sshpass`（配免密见 `ssh-key-auth-setup` 技能）。

**Agent 必须**通过技能目录下的二进制调用，**不要**假设 `sshctl` 在系统 PATH 中。

```powershell
# skillRoot = 本 SKILL.md 所在目录（Agent 从附带的技能文件路径解析，勿用仓库根或硬编码绝对路径）
$skillRoot = '...'  # 本技能：skills/sshctl/（相对仓库）或 ~/.claude/skills/sshctl/ 等
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
| 清单迁移 | `& $sshctl migrate`（或任意命令自动迁移） |

---

## 二进制位置（推荐）

技能与二进制同目录，**不写入系统 PATH**：

| 项 | 路径 |
|----|------|
| skillRoot | 本 `SKILL.md` 所在文件夹 |
| 二进制 | `$skillRoot\bin\sshctl.exe` |
| 仓库内示例 | `skills/sshctl/bin/sshctl.exe`（克隆后本地构建） |
| Claude 技能 | `%USERPROFILE%\.claude\skills\sshctl\bin\sshctl.exe` |
| Codex 技能 | `%USERPROFILE%\.codex\skills\sshctl\bin\sshctl.exe` |

### 构建 / 更新二进制

在仓库根目录：

```powershell
$env:VERSION = '0.2.0'
.\scripts\build.ps1
# → skills/sshctl/bin/sshctl.exe
# → 若存在 ~/.claude/skills/sshctl/ 或 ~/.codex/skills/sshctl/，同步复制 bin/sshctl.exe
```

或仅构建技能二进制：

```powershell
go build -o skills\sshctl\bin\sshctl.exe .
```

`skills/sshctl/bin/sshctl.exe` 已加入 `.gitignore`，**不入库**；克隆后需本地构建或从 [Releases](https://github.com/Fracizz/sshctl/releases) 解压 `sshctl-windows-amd64.zip` 中的 `sshctl.exe` 到 `$skillRoot\bin\`。

### 验证

```powershell
& $sshctl version    # 0.2.0+
& $sshctl list
```

### 可选：系统 PATH 安装（高级）

仅当需要全局 `sshctl` 命令时，以**管理员 PowerShell**运行：

```powershell
& $sshctl install
# 或 .\scripts\install.ps1
# → C:\Program Files\sshctl\sshctl.exe + 机器 PATH
```

技能工作流**默认不用**此方式。

---

## 配置

| 项 | 路径 / 变量 |
|----|-------------|
| 清单（默认） | `%USERPROFILE%\.sshctl\servers.json` |
| 从 sshfrac 迁移 | `& $sshctl migrate`（或任意命令自动迁移） |
| Legacy 备份 | 迁移后 `~/.sshfrac/servers.json.bak` |
| 覆盖 | `$SSHCTL_CONFIG` |
| Legacy 环境变量 | `$SSHFRAC_CONFIG` 仍可读（显式指定时） |
| 主密码 | `$SSHCTL_MASTER_PASSWORD`、`$SSHCTL_BIND_MACHINE=1` |

**规则：** 每个 IP 仅一条；`add` 同 IP 覆盖；密码含特殊字符须完整引号包裹。

---

## 常用命令

```powershell
& $sshctl migrate
& $sshctl init
& $sshctl list
& $sshctl search -s 192.168
& $sshctl add --host 192.168.x.x --user administrator --password '...' --os Windows --desc "说明"
& $sshctl exec 192.168.x.x -- "hostname && whoami"
& $sshctl scp .\a.txt 192.168.x.x:C:/temp/a.txt
```

**Agent 流程：** `search -s` → 不在清单则 `add`（密码/账号向用户确认）→ `exec` / `scp`。

---

## 错误速查

| 情况 | 处理 |
|------|------|
| 找不到 sshctl | 构建到 `$skillRoot\bin\sshctl.exe`，或从 Release zip 复制到该目录 |
| duplicate host | `add` 同 IP 覆盖，或删 JSON 重复项 |
| Windows 密码失败 | 确认密码完整；`--os Windows`；v0.2.0+ |
| 仍用旧名 sshfrac | 运行 `& $sshctl migrate`；legacy 备份为 `servers.json.bak` |

## 边界

- 远程操作只用 sshctl，不用原生 ssh/scp
- 不安装到系统 PATH（技能工作流）
- 配免密 / 查免密 → 使用 **ssh-key-auth-setup** 技能
- 清单不入库；不擅自改他人 authorized_keys
