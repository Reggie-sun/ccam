---
title: 用量限制与切换账号
description: 用 ccam usage 查看每个 Claude Code 账号的 5 小时用量和每周用量，并用 ccam next --auto 打开余量最多的账号。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: d288177682d7
---

`ccam usage` 显示每个账号已用掉多少 5 小时限额和每周限额。`ccam next --auto` 打开余量最多的账号，`--continue` 会把当前对话带过去。

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

各列依次为：账号、5 小时用量、距离重置的时间、每周用量、距离重置的时间，以及这些数字是多久之前获取的。

## 数字从哪里来

ccam 使用每个账号自己的登录，同时向 Anthropic 的用量接口查询所有账号。每个账号最多等待 3 秒，某个账号一返回结果就立即显示它。

这个接口没有公开文档，可能会发生变化。请求失败时，ccam 会显示最后一次获得的数字，并将其标记为缓存。在这些数字保存之后已经重置的时间窗口会显示 0%。如果响应格式发生了变化，`ccam doctor` 会告诉你。

60 秒内获取的数字会被复用。选项：

| 选项 | 作用 |
| --- | --- |
| `--watch [seconds]` | 持续刷新，默认每 60 秒一次 |
| `--json` | 机器可读的输出 |
| `-v` | 同时显示按模型区分的每周限额，以及每个数字的来源 |
| `--no-cache` | 总是重新查询 |
| `--source cache` | 从不查询；只显示已保存的数字 |

## 过期的登录

当某个账号的访问令牌已过期时，ccam 不会为它续期，因为如果某个 Claude Code 会话正在使用同一登录，续期可能会让该会话退出登录。用 `ccam <name>` 启动一次该账号，Claude Code 就会为它续期。

`ccam usage --refresh` 会自行为过期的令牌续期，但仅限于未在运行的账号。它从不为存储在 macOS 钥匙串中的登录或已导入的 `~/.claude` 中的登录续期。

## 在状态栏中显示用量

```sh
ccam statusline install
```

之后 Claude Code 会在状态栏中显示账号名及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。如果你已经有状态栏命令，ccam 会保留它，并在前面加上账号名。状态栏也会保存它的数字，当接口没有响应时，`ccam usage` 会使用这些数字。

## 切换

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` 会跳过未登录的账号，`--auto` 会跳过用量已达 100% 的账号。要继续对话，两个账号都必须共享对话历史；参见[共享](/docs/sharing)。

[用量限制工作流程指南](/guides/usage-limit-workflow)通过一个完整示例进行了讲解。
