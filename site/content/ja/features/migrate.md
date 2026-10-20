---
title: 既存の環境のコピー
description: ccam migrateで、既存の~/.claudeの設定、CLAUDE.md、拡張機能、MCPサーバー、MCPのサインイン情報をccamのグループにコピーします。
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

ccamを使う前からClaude Codeを使っていた場合、環境は`~/.claude`と`~/.claude.json`にあります。`ccam migrate`はそれを共有グループにコピーするので、グループを通じてそれらの項目を共有するすべてのアカウントが、これまでと同じ環境で使い始められます。コピー元のフォルダーは読み取るだけで、通常の`claude`もこれまでどおり動作します。最初の`ccam add`の後に一度実行し、アカウントがMCPサーバーの共有を始めるたびにもう一度実行してください。

## 使い方

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

ターミナルで、項目のフラグも`--yes`も付けずに実行すると、ccamは最初にチェックリストを表示します。

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Spaceキー（または`x`）でカーソル位置の項目を切り替え、矢印キーか`j`/`k`で移動し、Enterキーで計画を表示します。Esc、`q`、Ctrl+Cを押すと`Cancelled.`と表示されます。コピー元に該当するものがない項目には`not found`と表示されます。すべてのチェックを外すと、ccamは`Nothing selected.`と表示して終了します。

続いてccamは計画を表示し、コピーする前に確認を求めます。

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

グループの内容が届くのはそのグループのアカウントだけで、しかもグループが共有する項目だけなので、概要ではグループのアカウントを一覧表示します。コピーした項目のうちグループが共有していないものがある場合や、グループに入っていないアカウントがある場合は、それを変更する`ccam group set <group> --add <items>`と`ccam group join <group> <accounts>`のコマンドを表示します。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `[folder]` | コピー元のClaude Codeのフォルダー | `~/.claude` |
| `--settings` | `settings.json`をコピーします | オフ |
| `--memory` | グローバルの`CLAUDE.md`をコピーします | オフ |
| `--extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/`をコピーします | オフ |
| `--mcp` | ユーザーのMCPサーバー（`.claude.json`の`mcpServers`）をコピーします | オフ |
| `--history` | 会話履歴（`projects/`、`history.jsonl`）をコピーします | オフ |
| `-g`、`--group string` | コピー先の共有グループ | `default` |
| `-y`、`--yes` | 確認なしでコピーします。競合した場合はグループのコピーを残します | オフ |
| `--dry-run`（グローバル） | 計画を表示し、確認やコピーの前に終了します | オフ |

項目のフラグを付けない場合、ターミナルではチェックリストで決めます。ターミナル以外の場合や`--yes`を付けた場合、ccamは会話履歴以外のすべてをコピーします。`projects/`は大きくなることがあるため、履歴はデフォルトで除外されます。

`--group on`は`default`を意味します。`--group off`と`--group local`は`--group takes a group name, not "off"`で拒否されます。`ccam migrate`は`--dry-run`に対応しています。ドライランでもチェックリストは表示されるので、選択内容をプレビューできます。

## しくみ

`internal/cli/migrate.go`が項目を選び、`migrate.Build`がコピー元とグループを読み取って計画を返します。確認した後の`Plan.Apply`まで、何も書き込まれません。

1. **コピー元を見つける。** `checkSource`は、存在しないパス、フォルダーではないパス、Claude Codeの既知の項目が1つもないパスを拒否します。`migrate.ClaudeJSON`は、そのフォルダーに対応する`.claude.json`を見つけます。フォルダー内にあればそれを使い、`~/.claude`の中にない場合は、その隣の`~/.claude.json`を使います。`CLAUDE_CONFIG_DIR`が設定されていない場合、Claude Codeはそこにこのファイルを置きます。
2. **項目をパスに対応付ける。** すべてのパスは、グループが使うのと同じ許可リストである`share.Layout`から取得します。MCPサーバーは`mcp.LocalServers`で`.claude.json`から読み取ります。これはファイル全体を解析せず、`profile.TopLevelValue`スキャナーを使います。
3. **移動せず、コピーする。** ファイルは`fsutil.AtomicWrite`か`fsutil.CopyFile`で`~/.local/share/ccam/groups/<group>/`にコピーし、MCPサーバーはグループの`mcp.json`に入れます。

各項目は、1つのコピー元のパスを1つのグループのパスに対応付けます。

<!-- motion:migrate-items -->

履歴の2行が破線になっているのは、`--history`を付けた場合だけコピーされるためです。`.credentials.json`と`.claude.json`は`share.Never`に含まれており、グループに届くことはありません。万一これらがレイアウトに現れた場合、`Build`は内部エラーを返します。フォルダーをたどる処理でも、名前でこれらをスキップします。

`buildEntry`はすべてのパスに状態を割り当てます。計画の各行が説明しているのはこの状態です。

<!-- motion:migrate-states -->

- **Missing**: コピー元にないか、フォルダーがあるべき場所にファイルがあります。項目ごとに1回、`not found, skipped`として表示されます。
- **Linked**: コピー元のパスがすでにグループのコピーを指しています。たとえば、コピー元がccamのプロファイルフォルダーの場合です。
- **Same**: グループにすでに同じ内容があります。
- **New**: グループにないので、そのまま全体をコピーします。
- **Merge**: 両方にフォルダーがあります。`copyMissing`はグループにないファイルだけを追加します。内容が異なるファイルはグループのまま残り、計画には`group copies kept`と表示されます。`plugins/`の場合は、レジストリの`installed_plugins.json`と`known_marketplaces.json`をマージし、プラグインの絶対パスをグループのフォルダーを指すように書き換えます。
- **Conflict**: 両方に内容の異なるファイルがあります。ファイルごとに`keep group`か`use mine`を選びます。`.json`ファイルでは`merge JSON`（値が衝突した場合はグループが優先されるディープマージ）、`history.jsonl`では`merge by line`も選べます。定義が異なるMCPサーバーも競合として扱われ、`keep group`か`use mine`を選びます。Enterキーを押したときのデフォルトは、最初に表示される選択肢です。`--yes`では常にグループのコピーを残します。

チェックリスト自体は`tui.MigratePicker`で、`ccam add`が使う`checklist`コンポーネントで作られています。`migrate.Found`は、lstatの呼び出しとMCPサーバーの一覧を使って、見つからない項目に印を付けます。

### MCPのサインイン情報

OAuthでサインインするMCPサーバーは、トークンをサーバーの一覧ではなく、ログインファイルである`.credentials.json`の`mcpOAuth`に保存します。Claude Codeは各項目のキーを、`<name>|`に、サーバーの`type`、`url`、`headers`の`sha256`ハッシュの先頭16文字を続けたものにしています。キーはフォルダーに依存しないので、サーバーの定義が同じであれば、コピーした項目は別のアカウントでも機能します。そのため、`internal/migrate/logins.go`の`buildLogins`は、最終的にグループがコピー元の定義で持つことになるサーバー、つまり新規のもの、すでに同じもの、競合して`use mine`を選んだものの項目だけをコピーします。プラグインのサーバーのサインイン情報（キーが`plugin:`で始まるもの）は、拡張機能をコピーし、コピー元に`plugins/`がある場合に一緒にコピーされます。

<!-- motion:migrate-signins -->

項目は、受け取る各アカウント自身の`.credentials.json`に、`fsutil.AtomicWrite`を使ってモード`0600`で書き込まれます。グループに入ることはありません。アカウントが受け取るのは、対象のグループを通じて`mcp`（プラグインのサーバーの場合は`extensions`）を共有している場合だけで、追加されるのはそのアカウントにないキーだけです。すでにある項目は、より新しいものかもしれないので残します。計画の`sign-ins`の行には、サーバー、受け取るアカウント、スキップしたアカウントが表示されます。ccamは、次のアカウントをスキップし、その名前と理由を表示します。

- `ccam import`でその場所のまま登録したアカウント。通常の`claude`がそのログインファイルを使っている可能性があるためです。
- 実行中のアカウント。Claude Codeがログインファイルを書き換えるためです（書き込む直前にもう一度確認し、`started meanwhile`として報告します）。
- まだログインしていないアカウント。
- ログインをmacOSのキーチェーンに保存しているアカウント。ccamはキーチェーンに書き込みません。

スキップしたアカウントでは、`/mcp`でサーバーにサインインしてください。あとからMCPサーバーの共有を始めたアカウントには、`ccam migrate`をもう一度実行するとサインイン情報がコピーされます。概要には``Then run `ccam migrate` again to give them the MCP sign-ins too.``と表示されます。

## 安全性と特殊なケース

- コピー元は読み取るだけです。中のファイルを移動、リンク、変更することはなく、コピー元のログインに書き込むこともありません。キーチェーンから読み取ることがあるmacOSでも同様です。
- `use mine`、マージ、既存の`mcp.json`への新しいMCPサーバーの追加でグループのファイルを置き換える前に、グループのコピーを`<group>/.ccam-backup/<time>/`に保存します。フォルダー名は実行ごとに一意なので、同じ秒に2回実行してもバックアップが上書きされることはありません。概要にそのパスが表示されます。
- `--yes`ではすべての競合でグループのコピーを残し、代わりに選ぶ方法として、同じコマンドを`--yes`なしで実行するよう案内します。
- 何も変わらない場合は`Nothing to copy.`と表示し、すべての競合でグループのコピーを残す場合は`Nothing to copy: the group keeps its own copies.`と表示します。
- コピーが途中で失敗した場合、ccamは`The source folder was not changed. Group files replaced so far are saved in`に続けてバックアップフォルダーを表示し、終了します。
- コピー元がない場合は`no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`で終了します。

## 関連項目

- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [点検と修復](/features/doctor)
- [Claude Codeの複数アカウント](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
