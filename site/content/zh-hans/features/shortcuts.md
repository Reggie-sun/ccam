---
title: 快捷命令
description: 用 ccam shortcut 为账号创建 ccwork 这样的专属命令，并用 ccam shims install 重新创建 ~/.local/bin 中的链接。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

快捷命令是你自己命名的命令，用来以某个账号启动 Claude Code：`ccwork --resume` 的效果与 `ccam work --resume` 相同。它是 `~/.local/bin` 中指向 ccam 的链接，因此在 bash、zsh、脚本以及其他任何从 `PATH` 运行程序的地方都能使用。快捷命令是可选的；为你经常启动的账号创建一个即可。

## 用法

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam 会确认每次更改：

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

一个账号最多有一个快捷命令，以 `shortcut` 保存在 `~/.config/ccam/config.toml` 中。`ccam add` 向导会询问 `Create a shortcut command for <name>? [y/N]`，然后询问 `Shortcut name [cc<name>]:`。在向导之外，`ccam add` 只有在带 `--shortcut` 时才会创建快捷命令。`ccam list` 显示每个账号的快捷命令，没有时显示 `-`。

快捷命令会把所有参数传给 Claude Code，与 `ccam <name>` 处理账号名之后的参数的方式相同。其中也包括 `--dry-run` 和 `--help`：要预览一次启动，请使用 `ccam --dry-run work`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--remove` | 删除账号的快捷命令及其链接 | 关闭 |
| `-h`、`--help` | 命令的帮助 | |
| `--dry-run` | `ccam shortcut` 支持它：检查名称并打印 `Would set the shortcut of work to ccwork.` 或 `Would remove ccwork, the shortcut of work.`，不保存也不创建链接 | 关闭 |

`ccam shims install` 不接受参数，并会拒绝 `--dry-run`，显示 `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``。

名称规则，按以下顺序检查：

- 由 1 到 64 个 `a-z`、`0-9`、`.`、`-` 和 `_` 字符组成，以字母或数字开头；
- 不能是 `ccam`、`claude`，也不能是 `run` 或 `usage` 这样的 ccam 命令；
- 不能是其他账号的快捷命令；
- 不能是 `~/.local/bin` 中已有的、且并非指向 ccam 的链接的文件；
- 不能与你 `PATH` 上的其他可执行文件同名。

同时给出命令名和 `--remove` 会以 `give a command or --remove, not both` 停止。设置账号已有的名称会打印 `work: ccwork (unchanged)`，并在链接缺失时重新创建它。

## 工作原理

设置名称时会运行 `checkShortcut`，它以 bin 文件夹和 ccam 可执行文件的真实路径调用 `shim.Check`：

<!-- motion:shortcuts-check -->

名称按顺序依次通过各项检查。每项检查旁边是一个会被它拒绝的名称：`Work` 不符合格式，`run` 是 ccam 命令，`ccside` 属于另一个账号，`ccfile` 是 `~/.local/bin` 中已有的文件，`git` 会遮住 `/usr/bin/git`。只有通过所有检查的名称才会被保存并创建链接。

1. `shim.Check` 依次用 `^[a-z0-9][a-z0-9._-]{0,63}$`、`config.Reserved` 和 `cfg.FindShortcut` 检查名称。它对 `~/.local/bin/<name>` 调用 `os.Lstat`，对 `PATH` 上的其他每个文件夹调用 `os.Stat`，并拒绝不是指向 ccam 的链接的可执行文件。
2. ccam 用 `fsutil.AtomicWrite` 保存 `config.toml`（临时文件、`fsync`、rename）。
3. `shim.Remove` 删除旧快捷命令的链接，但仅当它是指向 ccam 的链接时。
4. `shim.Install` 为每个已配置但缺少链接的快捷命令创建链接，用 `os.Symlink` 指向 ccam 可执行文件，其路径由 `os.Executable` 和 `filepath.EvalSymlinks` 解析。它从不删除任何内容。

当 `os.Readlink` 返回 ccam 的路径，或者链接与 ccam 路径解析到同一个文件时，该链接就被视为 ccam 自己的链接。每次删除和每次覆盖检查都要经过这项判断，因此 ccam 不会触碰不是它创建的文件。

运行快捷命令时，ccam 只需要你输入的那个名称：

<!-- motion:shortcuts-link -->

shell 在 `PATH` 上找到 `ccwork`，内核跟随链接启动 ccam 可执行文件，而 `argv[0]` 仍然是 `ccwork`。在 `app.Main` 中，只要名称不是 `ccam`，就会加载一次 `config.toml` 并调用 `shim.Lookup`，后者找到 `shortcut` 等于 `filepath.Base(argv[0])` 的账号。然后 `launcher.LaunchWith` 在配置已加载的情况下执行与 [`ccam <name>`](/features/launch) 相同的步骤，并把 `argv[1:]` 原样传给 `claude`。

为什么用符号链接，而不是 shell 别名或包装脚本：符号链接在每种 shell 和脚本中都能工作，不会多启动一次 shell（约 1 到 2 ms），启动时也仍然只读取一个文件。`scripts/perf-check.sh` 用 `hyperfine` 测量快捷命令，其上限与 `ccam <name>` 相同，为 10 ms（p95）。

## 安全性与边界情况

- **没有账号认领的链接。** 如果 `~/.local/bin` 中某个指向 ccam 的链接的名称没有任何账号在使用（例如你手动编辑了 `config.toml` 中的 `shortcut` 之后），ccam 不会退回到默认账号，而是以 `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` 停止。请用 `ccam shortcut` 重新设置名称，或者删除该链接。
- **已有文件。** ccam 从不替换不是它创建的文件。`ccam shims install` 会以 `<path> already exists and is not a ccam link, so it was left alone` 停止。
- **移动了 ccam 可执行文件。** 旧链接会指向一个已不存在的路径，ccam 不会把它们视为自己的链接。请删除它们，然后运行 `ccam shims install`。
- **重命名和删除。** `ccam rename` 保留快捷命令的名称，链接也继续有效，因为它指向 ccam，而 ccam 会在 `config.toml` 中查找该名称。`ccam remove` 会删除账号的链接。`ccam shortcut <account> <new>` 会先删除旧链接，再创建新链接。
- **PATH。** 当 ccam 创建链接而 `~/.local/bin` 不在你的 `PATH` 上时，它会打印 `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`，后面跟着需要添加的 `export PATH=...` 一行。
- **之后检查。** `ccam doctor` 会报告 `missing shortcuts in /home/me/.local/bin: ccwork`，并在该文件夹不在 `PATH` 上时发出警告。`ccam doctor --fix` 会创建缺少的链接。

## 相关内容

- [启动 Claude Code](/features/launch)
- [管理账号](/features/accounts)
- [命令参考](/docs/commands#shortcuts)
- [安装 ccam](/docs/install)
