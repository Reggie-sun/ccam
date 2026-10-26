---
title: 更新 ccam
description: ccam update 在通过 SHA-256 校验后，用最新的 GitHub 发布版本替换 ccam；--check 只报告，--version 安装指定的发布标签。
updated: 2026-10-26
order: 10
translation: machine
sourceHash: "598931e4801f"
---

`ccam update` 从 GitHub 下载一个 ccam 发布版本，用该版本的 `checksums.txt` 进行校验，然后通过一次 rename 替换正在运行的可执行文件。如果你是用安装脚本或从压缩包安装的 ccam，就用它来更新。[仪表盘](/features/dashboard)也会在每次打开时向 GitHub 询问最新版本，并通过 `U` 提供更新，该键会先询问，再执行相同的步骤。其他命令从不连接 GitHub，因此它们在离线时照常工作。

## 用法

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

没有更新的版本时，ccam 打印 `ccam 0.4.0 is up to date.`，不做任何更改。如果 `--version` 指定的是你当前已有的版本，它打印 `ccam 0.4.0 is already installed.`

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--check` | 只报告是否有更新的版本 | 关闭 |
| `--version <tag>` | 安装这个发布标签（如 `v0.3.0`），可以比当前版本新或旧 | 最新 |
| `--dry-run` | 支持：查询最新标签（除非 `--version` 指定了标签），对可执行文件执行相同的检查，并打印下载、校验和与替换的路径，但不下载 | 关闭 |

`--check` 和 `--version` 不能同时使用。标签必须形如 `vMAJOR.MINOR.PATCH`，可以带预发布或构建后缀；开头的 `v` 可以省略。其他形式都会失败，例如 `--version takes a release tag such as v0.3.0, not "1.2"`。

环境变量 `CCAM_RELEASES_URL` 可以替换 `https://github.com/jhyunwoo/ccam/releases`，用于镜像或测试。它的结构必须相同：`<url>/latest` 重定向到 `<url>/tag/<tag>`，资源文件位于 `<url>/download/<tag>/<asset>`。

## 工作原理

ccam 在改动已安装的可执行文件之前会先校验下载内容，并通过一次 rename 替换可执行文件。代码位于 `internal/update` 中，只有 `ccam update` 会导入它。启动路径从不调用它。

1. ccam 用 `os.Executable` 找到正在运行的可执行文件，并用 `filepath.EvalSymlinks` 解析符号链接（`cli.updateCmd`）。被替换的就是这个真实路径。
2. 不带 `--version` 时，ccam 请求 `<releases>/latest` 但不跟随重定向（`Updater.Latest`），并从 `Location` 请求头（`…/tag/v0.5.0`）中读取标签。这里使用的是发布页面而不是 GitHub API，因此不受 API 每小时 60 次请求的限制。查询在 10 秒后超时。
3. ccam 用严格的格式检查标签（`parseVersion`）。标签会出现在下载 URL 和打印出的 shell 命令中，因此带有 `/`、`%`、空格或 shell 字符的内容都无法通过。然后 ccam 把它与正在运行的版本进行比较（`compareVersions`）；预发布版本排在正式版本之前。
4. ccam 决定它是否可以替换该可执行文件（`managed()`、`writable()`；见下文）。
5. ccam 同时下载 `ccam_<os>_<arch>.tar.gz`（在 Windows 上是 `ccam_windows_<arch>.zip`）和 `checksums.txt`，两者共用 2 分钟时限（`Updater.Fetch`）。每个下载的上限为 256 MB。
6. 在进行任何其他操作之前，压缩包的 SHA-256 必须与 `checksums.txt` 中对应的那一行一致。缺少该行或校验和不同，都会以 `nothing was changed` 中止更新。
7. ccam 读取 tar 并取出第一个名为 `ccam` 的普通文件（`extract()`；`extractZip` 从 zip 中取出 `ccam.exe`）。目录、链接和其他条目会被忽略，可执行文件的上限同样为 256 MB。
8. ccam 用 `fsutil.AtomicWrite` 写入新的可执行文件（`Updater.Install`）：在可执行文件旁边创建临时文件 `.ccam.tmp*`，沿用旧文件的权限位，刷新到磁盘后 rename 覆盖旧文件。

从查询标签到 rename 的完整流程：

<!-- motion:update-pipeline -->

压缩包和 `checksums.txt` 并行下载，在校验和一致之前不会写入任何内容。rename 是原子操作，因此即使崩溃或磁盘已满，留下的也只会是旧的或新的可执行文件，绝不会是半个。`~/.local/bin/ccwork` 这样的快捷命令是指向 ccam 可执行文件的符号链接，因此无需任何更改就会运行新版本。

在 Windows 上，正在运行的程序的文件无法被替换或删除，只能重命名。ccam 把新的可执行文件写入 `ccam.exe.new`，把正在运行的 `ccam.exe` 重命名为 `ccam.exe.old`，再把新文件重命名为 `ccam.exe`（`Updater.Install`）。如果最后一步失败，旧的可执行文件会被改回原名。下一次更新会删除它能删除的 `ccam.exe.old` 文件；如果仍有会话在使用更早的 `ccam.exe.old`，正在运行的可执行文件会改为重命名为 `ccam.exe.old-<time>`。使用硬链接而不是符号链接的快捷命令（[快捷命令](/features/shortcuts)）会紧接着被更新，因此它们也会运行新版本；这包括 ccam 为任何旧版本 ccam 创建的硬链接，而不只是刚被替换的那个可执行文件的硬链接。

### 由软件包管理的安装

从 `.deb` 安装的 ccam 副本归 apt 管理，绕过 apt 替换它会让软件包数据库与文件不一致。下载之前，ccam 会检查可执行文件归谁所有：

<!-- motion:update-managed -->

1. 位于 `/usr/bin`、`/bin`、`/usr/sbin` 或 `/sbin` 中的可执行文件视为由软件包管理。
2. 在其他位置，如果安装了 `dpkg` 且 `dpkg -S <path>` 成功，也视为由软件包管理。
3. 否则，ccam 会在可执行文件旁边创建并删除一个临时的 `.ccam-update-*` 文件，以确认可以在那里写入。

对于由软件包管理的可执行文件，ccam 会停止并打印 apt 的操作步骤：

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

在 Linux 以外的系统上，只打印安装脚本那一行：在 Windows 上是 `irm https://ccam.dev/install.ps1 | iex`。当文件夹不可写时，ccam 建议重新运行安装脚本，它会把 ccam 安装到 `~/.local/bin`。

## 安全性与边界情况

- 开发构建报告的版本为 `dev`，没有可以比较的发布版本，因此 `ccam update` 和 `ccam update --check` 会拒绝它。`ccam update --version vX.Y.Z` 可以在它之上安装一个发布版本。用 `go install` 构建的版本请参见[安装 ccam](/docs/install#other-ways-to-install)。
- 下载失败、缺少适用于你平台的资源文件，或者校验和不一致，都会让已安装的可执行文件保持原样。
- 不存在的发布版本会以 `there is no release <tag>` 失败。当 GitHub 以 HTTP 403 或 429 拒绝时，ccam 会提示这可能是频率限制，并请你稍后再试。
- `--check` 只读取。它不检查可执行文件是否由软件包管理或是否可写。
- 用 `--version` 指定较旧的标签，就可以回退版本。

## 相关内容

- [安装 ccam](/docs/install)
- [命令参考](/docs/commands)
- [常见问题](/docs/faq)
