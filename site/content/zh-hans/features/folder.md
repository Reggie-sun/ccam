---
title: 按文件夹使用账号
description: ccam folder 让一个文件夹及其下所有文件夹在不指定账号名启动 Claude Code 时使用同一个账号，例如公司代码使用 work。
updated: 2026-10-23
order: 11
translation: machine
sourceHash: "48de98d42b99"
---

`ccam folder work ~/code/company` 让 `~/code/company` 及其下任何文件夹中所有不指定账号名的启动都以 `work` 启动 Claude Code。在这个文件夹之外，照旧启动默认账号。指定账号名时，无论是用 `ccam personal` 还是 `ccpersonal` 这样的快捷命令，总是启动该账号。

## 用法

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

以下启动会使用文件夹的账号：

- `ccam -p "…"`，以及其他第一个参数是 Claude Code 选项的启动
- 不带账号名的 `ccam run`，以及 `ccam --dry-run`
- `ccam --fallback -p "…"` 的第一次尝试
- 仪表盘，即单独的 `ccam`，打开时光标位于文件夹的账号上，通过符号链接进入该文件夹时也一样

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `<account>` | 要使用的账号的名称或编号 | |
| `[folder]` | 文件夹，相对路径或绝对路径 | 当前文件夹 |
| `--remove` | 删除恰好对应该文件夹的条目 | 关闭 |
| `--dry-run` | 支持：打印更改而不保存 | 关闭 |

## 工作原理

每个文件夹是 `~/.config/ccam/config.toml` 中的一个 `[[folder]]` 条目，而每次启动本来就会读取这个文件：

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` 把文件夹解析为绝对路径（解析掉符号链接）并保存。再次设置同一个文件夹会替换它的账号。
2. 不指定账号名的启动会调用 `launcher.HereProfile`。没有任何条目时，它会立即取默认账号。
3. 否则，它把工作目录与每个 `path` 按字符串比较：匹配文件夹本身或其下的文件夹，绝不匹配只是开头字母相同的文件夹（`/code/company` 不包括 `/code/companyx`）。最深的条目优先，所以 `~/code` 可以是 `personal`，而 `~/code/company` 是 `work`。
4. 只有当原样的工作目录什么都没有匹配上时，ccam 才会解析它的符号链接并再比较一次。
5. 在 Windows 和 macOS 上，比较时不区分大小写，与它们的文件系统通常的行为一致。

除了 `config.toml` 之外不读取任何文件，所以在文件夹中启动不会增加可测量的时间（无论有无条目都是 4.0 ms，`hyperfine`，300 次运行）。

## 安全性与边界情况

- `ccam rename` 会把文件夹条目带到新名称，`ccam remove` 会删除被删除账号的条目。
- 如果条目对应的账号已不存在，在该文件夹中不指定账号名的启动会停止：`<folder> is set to account "<name>", which does not exist. Change it with ccam folder`。`ccam doctor` 会报告这类条目，以及文件夹已不存在的条目。
- 如果当前文件夹的账号来自上层文件夹的设置，在其中运行 `ccam folder --remove` 不会删除任何内容，而是指出那个上层文件夹。
- `ccam next` 按轮换或用量选择账号，不考虑文件夹。

## 相关内容

- [启动 Claude Code](/features/launch)
- [快捷命令](/features/shortcuts)
- [管理账号](/features/accounts)
- [命令参考](/docs/commands)
