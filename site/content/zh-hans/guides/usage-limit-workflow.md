---
title: Claude Code 达到用量限制时该怎么办
description: 等待额度重置，或立即在另一个 Claude Code 账号上继续同一个对话。逐步说明 5 小时限制和每周限制的工作方式。
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "18ef61a68c55"
---

当一个账号达到用量限制（常被称为 rate limit）时，你有两个选择：等到限制重置，或者现在就在你拥有的另一个账号上继续同一个对话。使用 ccam 时，第二种做法只需一条命令：`ccam next --auto --continue`。

## Claude Code 有哪些限制？

在 Pro 和 Max 套餐中，Claude Code 有[5 小时会话限制和每周限制](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)。这两个限制与网页版和应用中的 Claude 共用，用量会同时计入两者。其中一个用完时，Claude Code 会[显示一条消息](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)，例如 "You've hit your session limit · resets 3:45pm"。切换模型没有用，因为这些限制覆盖所有模型。

## 方法 1：等待重置

如果你只有一个账号，就只能等待。较新版本的 Claude Code 会[在打开的会话中等待](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)，限制重置后自动继续任务。在 Claude Code 中运行 `/usage` 可以查看你已用掉多少套餐额度。

## 方法 2：现在就在另一个账号上继续

如果你有第二个账号，就可以立即接着同一个对话继续工作。只需设置一次：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

现在，工作时状态栏会显示当前账号及其 5 小时用量，达到限制时就不会措手不及。

达到限制后，先退出 Claude Code（`/exit`），然后运行：

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

在余量最多的账号中打开这个对话：

```sh
ccam next --auto --continue
```

Claude Code 会以 `work` 账号启动，并打开当前文件夹中最近的对话，包括之前的全部内容。共享历史（上面的 `ccam group set default --add history`）能让它保持为同一个对话；不共享时，Claude Code 会把对话复制到 `work` 中，ccam 会告诉你第一次回复要重新读取多少 token。

在 Claude Code 中安装 [ccam 插件](/docs/install#claude-code-plugin)后，因达到限制而停止的会话会显示一条通知，并记录已满的时间窗口，这样 `ccam next --auto` 在其重置之前不会再次选中这个账号；`/ccam:switch`（或 `/ccam:switch work`）会为 `ccam next --continue` 记下这个对话，不留任何需要猜测的地方。

**Team 和 Enterprise 席位。** 自动选择（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换。Anthropic 帮助中心关于 [Team 套餐](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的说明写道："Each seat on a Team plan is assigned to one named individual, and only that person may use it"（Team 套餐的每个席位都分配给一位指定的个人，只有此人可以使用），以及 "Users can't route their work through another user's account to use that user's remaining usage"（用户不能通过另一位用户的账号完成自己的工作，以使用该用户剩余的用量）；关于 [Enterprise 套餐](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的说明也是如此。因此，在席位之间使用自动选择可能违反 Anthropic 的条款。在这些席位之间是否使用，由你自己决定：在依赖这项功能之前，请先向你所在的组织和 Anthropic 确认。个人套餐（Pro、Max）的情况请参见[常见问题中关于条款的说明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)。

## ccam 会自动切换账号吗？

只有在你开启后才会。默认情况下，你工作时 ccam 并没有在运行：它启动 Claude Code 后就用 Claude Code 替换了自身，所以何时切换由你决定，`ccam next --auto` 会选择 5 小时用量最低的账号，跳过用量已达 100% 或未登录的账号，并且只在同一个 Team 或 Enterprise 组织内，或在个人套餐之间选择。在 `config.toml` 中设置 `relay = true` 后，ccam 会改为在会话旁边等待，这样 `/ccam:switch` 就能自行继续，并沿用会话的选项、权限模式和推理强度。设置 `auto_switch = true`（仅限 Linux 和 macOS；实验性功能）后，ccam 还会自己完成切换：当 Claude Code 因用量限制停止会话，并且另一个同类账号还有余量时，ccam 会结束该会话，并在那个账号上继续这个对话（[条件与代价](/features/next#auto-switch)）。另一些工具则是替换正在运行的会话的登录信息；[ccam 与其他账号切换工具](/compare/account-switchers)解释了两者的区别。

## 注意事项

- 在你工作时所在的同一个文件夹中运行 `ccam next --continue`；对话是按项目文件夹保存的。
- 避免让两个账号同时处理同一个对话。你这样做时，ccam 会发出警告。
- 不共享历史时，对话会被复制到另一个账号，两份副本此后各自发展；`/rewind` 用来恢复代码的检查点留在第一个账号中。
- 新账号上的第一条消息可能比平时消耗更多额度，因为该账号还没有缓存这个对话。`ccam next --continue --brief` 会改为以一份简短的摘要开始。
- 是否可以这样在账号之间切换，取决于每个账号的条款。

`ccam usage` 和 `ccam next` 的所有选项请参见[用量限制](/docs/usage-limits)。
