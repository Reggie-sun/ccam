---
title: Claude Code 达到用量限制时该怎么办
description: 查看每个 Claude Code 账号还剩多少余量，并在某个账号达到用量限制时，在另一个账号上继续同一段对话。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

当一个账号达到用量限制时，在你拥有的另一个账号中继续同一段对话：先共享一次对话历史，然后运行 `ccam next --continue`。如果你只有一个账号，请等到 Claude Code 显示的重置时间。

## 提前做好准备

把账号设置好，让切换只需一条命令：

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

现在，状态栏会在你工作时显示账号及其 5 小时用量，因此达到用量限制时你不会措手不及。

## 达到用量限制时

退出 Claude Code（`/exit`），然后运行：

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

在余量最多的账号中打开这段对话：

```sh
ccam next --auto --continue
```

Claude Code 会以 `work` 账号启动，并打开当前文件夹最近的一次对话，包括到目前为止的全部内容。

## 注意事项

- 在你之前工作的同一个文件夹中运行 `ccam next --continue`；对话是按项目文件夹存储的。
- 避免让两个账号同时处理同一段对话。如果你这样做，ccam 会发出警告。
- 另一个账号也必须共享对话历史，并且在同一个组中。
- 你能否像这样在账号之间切换，取决于每个账号适用的条款。

关于 `ccam usage` 和 `ccam next` 的所有选项，参见[用量限制](/docs/usage-limits)。
