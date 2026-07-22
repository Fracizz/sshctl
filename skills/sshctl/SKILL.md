---
name: sshctl
description: |
  sshctl 远程主机 CLI（清单 search/exec/shell/scp/add/install）。优先 sshctl，尽量不用原生 ssh/scp。
  触发词：sshctl、search -s、exec、shell、scp、install、servers.json、SSHCTL、.sshctl。
---

# sshctl · 远程主机 CLI（Windows）

**默认用 sshctl 做所有远程操作**；不要手写 `ssh` / `scp` / `sshpass`（配免密见 `ssh-key-auth-setup` 技能）。

| 操作 | 命令 |
|------|------|
| 搜主机 | `sshctl search -s <关键词>` |
| 远程执行 | `sshctl exec <host> -- <cmd>` |
| 交互 shell | `sshctl shell <host>` |
| 传文件 | `sshctl scp <src> <dst>` |
| 写入清单 | `sshctl add --host ... --user ... --password '...'` |
| 系统安装 | `sshctl install`（需管理员） |

---

## 安装到系统目录（推荐）

安装后 **`sshctl` 在 PATH**，不必再写全路径。

### 方式 1：sshctl install（管理员 PowerShell）

```powershell
cd D:\WorkSpace\code\sshctl
$env:GOROOT = 'C:\Program Files\Go'
& 'C:\Program Files\Go\bin\go.exe' build -o bin\sshctl.exe .
# 以管理员运行：
.\bin\sshctl.exe install
# → C:\Program Files\sshctl\sshctl.exe + 写入机器 PATH
```

### 方式 2：安装脚本

```powershell
# 管理员 PowerShell
cd D:\WorkSpace\code\sshctl
.\scripts\install.ps1
```

### 方式 3：发布 zip

从 [Releases](https://github.com/Fracizz/sshctl/releases) 下载 `sshctl-windows-amd64.zip`，解压后以管理员运行：

```powershell
.\sshctl.exe install
```

### 验证

```powershell
# 新开终端
sshctl version    # 0.2.0+
sshctl list
```

| 路径 | 说明 |
|------|------|
| 系统安装 | `C:\Program Files\sshctl\sshctl.exe` |
| 开发构建 | `D:\WorkSpace\code\sshctl\bin\sshctl.exe` |
| go install | `%USERPROFILE%\go\bin\sshctl.exe` |

---

## 配置

| 项 | 路径 / 变量 |
|----|-------------|
| 清单（默认） | `%USERPROFILE%\.sshctl\servers.json` |
| 从 sshfrac 迁移 | `sshctl migrate`（或任意命令自动迁移） |
| Legacy 备份 | 迁移后 `~/.sshfrac/servers.json.bak` |
| 覆盖 | `$SSHCTL_CONFIG` |
| Legacy 环境变量 | `$SSHFRAC_CONFIG` 仍可读（显式指定时） |
| 主密码 | `$SSHCTL_MASTER_PASSWORD`、`$SSHCTL_BIND_MACHINE=1` |

**规则：** 每个 IP 仅一条；`add` 同 IP 覆盖；密码含特殊字符须完整引号包裹。

---

## 常用命令

```powershell
sshctl migrate              # ~/.sshfrac → ~/.sshctl（legacy 重命名为 .bak）
sshctl init
sshctl list
sshctl search -s 192.168
sshctl add --host 192.168.x.x --user administrator --password '...' --os Windows --desc "说明"
sshctl exec 192.168.x.x -- "hostname && whoami"
sshctl scp .\a.txt 192.168.x.x:C:/temp/a.txt
```

**Agent 流程：** `search -s` → 不在清单则 `add`（密码/账号向用户确认）→ `exec` / `scp`。

---

## 打包 / 开发

```powershell
cd D:\WorkSpace\code\sshctl
$env:VERSION = '0.2.0'
.\scripts\build.ps1          # dist\sshctl-*.zip
go build -o bin\sshctl.exe . # 需 GOROOT 指向完整 Go 安装
```

仓库：https://github.com/Fracizz/sshctl（模块名 `github.com/Fracizz/sshctl`）

---

## 错误速查

| 情况 | 处理 |
|------|------|
| 找不到 sshctl | 运行 `sshctl install` 或用 `bin\sshctl.exe` 全路径 |
| install 失败 | 管理员 PowerShell |
| duplicate host | `add` 同 IP 覆盖，或删 JSON 重复项 |
| Windows 密码失败 | 确认密码完整；`--os Windows`；v0.2.0+ |
| 仍用旧名 sshfrac | 运行 `sshctl migrate`；legacy 备份为 `servers.json.bak` |

## 边界

- 远程操作只用 sshctl，不用原生 ssh/scp
- 配免密 / 查免密 → 使用 **ssh-key-auth-setup** 技能
- 清单不入库；不擅自改他人 authorized_keys
