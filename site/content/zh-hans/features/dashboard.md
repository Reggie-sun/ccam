---
title: 仪表盘
description: 在终端中运行 ccam 会打开全屏仪表盘，显示账号、组、文件夹规则、扩展和检查结果；每项更改都会先显示计划。
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` 不带参数运行时，如果标准输入和标准输出都是终端，就会打开全屏仪表盘；否则打印 `ccam list` 的输出。仪表盘显示每个账号及其用量和组，并能完成各个命令所做的事：启动、添加、登录、更改组、设置文件夹规则、检查和修复。启动 Claude Code 时仪表盘会关闭；添加账号或登录时会暂时离开仪表盘回到终端，完成后再返回。

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

第一帧只根据 `config.toml` 绘制，用时远低于 100 ms；用量会在每个账号返回结果后陆续填入。百分比从 80% 起变为黄色，达到 100% 时变为红色；查询失败后取自缓存的值会以淡色显示。

## 选项卡

| 选项卡 | 显示内容 | 按键 |
| --- | --- | --- |
| Accounts | 编号、名称、邮箱、5 小时用量条、7 天用量、上次使用时间、组；`✔` 表示默认账号，`▸` 表示当前文件夹的账号，`•` 表示正在运行 | 见下文 |
| Groups | 每个组、它共享的内容和它的账号，以及不在任何组中的账号 | `enter` 编辑，`n` 新建，`e` 重命名，`x` 删除 |
| Folders | [`ccam folder`](/features/folder) 的规则，以及在当前位置生效的那一条 | `n` 在此处使用某个账号，`enter` 更改，`x` 删除 |
| Extensions | [`ccam ext`](/features/ext) 显示的内容，每个扩展只列一次，并注明它来自哪些组 | `r` 重新加载 |
| Health | [`ccam doctor`](/features/doctor) 报告的内容 | `f` 修复（`doctor --fix`） |

`←`/`→` 或 `Tab` 切换选项卡，`?` 列出当前选项卡的按键，`r` 跳过短期缓存重新查询用量，`q` 或 `Esc` 退出。

## Accounts 选项卡中的按键

| 按键 | 作用 |
| --- | --- |
| `Enter` | 以所选账号启动 Claude Code |
| `1`–`9` | 移到该编号的账号；快速输入 `1` `2` 可移到 12 号。只有 `Enter` 才会启动 |
| `↑`/`↓`、`k`/`j` | 移动 |
| `n` | 启动 [`ccam next --auto`](/features/next) 会选择的账号 |
| `o` | 带选项启动：`--continue`、`--resume`、其他参数，或试运行 |
| `u` | 用量详情：5 小时、7 天和各模型的时间窗口，以及它们的重置时间和数字的来源 |
| `a` | 添加账号（见下文的向导） |
| `l` | 重新登录 |
| `s` | 选择账号所属的组 |
| `d` | 设为默认账号 |
| `f` | 在当前文件夹中使用它 |
| `e` | 重命名 |
| `c` | 设置或删除它的快捷命令 |
| `x` | 删除它（先询问，默认选中 No） |

## 更改组

在 Groups 选项卡中，`Enter` 会打开编辑器：左边是项，右边是账号，用 `Tab` 切换，用 `Space` 勾选或取消。如果某个账号已经从另一个组获得某一项，该项会变灰并显示原因；会重复获得某一项的账号也同样变灰，因为一项只能来自一个组。`Enter` 显示计划：

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

与组中副本不同的文件会先询问（保留组中的、使用我的，或者合并）。然后更改在后台运行，使用的代码与 [`ccam group`](/features/share) 相同：先保存配置，再按顺序移动每个账号的数据，原始文件放入 `.ccam-backup/`。会移动正在运行的账号数据的更改会被拒绝，并显示该账号的名称。

## 添加账号

`a` 会打开 `ccam add` 在终端中同样使用的向导：

1. 账号名称。
2. 登录方式：现在登录 Claude 订阅（在浏览器中）、[API 密钥](/features/api)，或稍后再登录 Claude 订阅。
3. 使用 API 密钥时：ccam 每次启动时从哪里读取密钥（`env:NAME` 或 `command:COMMAND`）、端点、密钥的发送方式以及模型。
4. 要加入的组，默认勾选 `default`；还没有任何组时，则选择 `default` 组共享哪些内容。
5. 快捷命令（如果需要）。
6. 供确认的摘要。

`Esc` 返回上一步；在第一步时则关闭向导。

## 颜色

颜色沿用 Claude Code 终端界面中各角色的配色，因此 ccam 放在它旁边也显得协调；组名的颜色由名称决定。ccam 会向终端查询背景色，并选择浅色或深色主题。你也可以在 `config.toml` 中自行设置：

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 主题适用于红绿色盲，`-ansi` 主题使用终端自己的 16 种颜色。`NO_COLOR` 会关闭颜色。

## 安全性与边界情况

- 仪表盘调用的函数与各个命令相同，因此所有规则都同样适用：移动数据前先显示计划、拒绝更改正在运行的账号、保留备份，以及登录令牌永远不会离开账号。
- 仪表盘没有试运行：不带名称的 `ccam --dry-run` 会打印此文件夹的账号将如何启动，在某个账号上按 `o` 则可以只对这一次启动进行试运行。
- 在 Claude Code 启动或登录页面打开之前，终端会先恢复原状。
- 屏幕太小、放不下对话框时，会先省略选项卡中的行；过长的行用 `…` 截断，而不是换行。

## 相关内容

- [通过组共享](/features/share)
- [启动 Claude Code](/features/launch)
- [管理账号](/features/accounts)
- [用量限制](/features/usage)
