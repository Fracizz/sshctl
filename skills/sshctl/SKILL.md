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

---

## 二进制（与 skill 同目录）

| 项 | 路径 |
|----|------|
| skillRoot | 本 `SKILL.md` 所在文件夹 |
| 二进制 | `$skillRoot\bin\sshctl.exe` |

### 构建 / 更新

在仓库根目录（仅 Windows amd64；会同步到仓库 skill 与已存在的 `.claude` / `.codex` skill `bin\`）：

```powershell
$env:VERSION = '0.2.1'
.\scripts\build.ps1
```

或：

```powershell
go build -o skills\sshctl\bin\sshctl.exe .
```

`bin/` 下的 exe **不入库**。也可从 [Releases](https://github.com/Fracizz/sshctl/releases) 取 `sshctl-windows-amd64.zip`，把 `sshctl.exe` 放到 `$skillRoot\bin\`。

多平台发布包：推送 `v*` 标签，由 GitHub Actions（`release.yml`）打 zip 并创建 Release；**不要**本机打发布包。

### 验证

```powershell
& $sshctl version    # 0.2.1+
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
& $sshctl scp .\a.txt 192.168.x.x:C:/temp/a.txt
```

**Agent 流程：** `search -s` → 不在清单则 `add`（向用户确认凭据）→ `exec` / `scp`。

---

## 错误速查

| 情况 | 处理 |
|------|------|
| 找不到 sshctl | 构建/复制到 `$skillRoot\bin\sshctl.exe` |
| duplicate host | `add` 同 IP 覆盖，或删 JSON 重复项 |
| Windows 密码失败 | 确认密码完整；`--os Windows`；v0.2.1+ |
| 仍用旧清单目录 | `& $sshctl migrate` |

## 边界

- 远程操作只用 `$sshctl`，不用原生 ssh/scp
- **不**安装到系统 PATH（技能工作流）
- 配免密 → **ssh-key-auth-setup**
- 清单不入库
