---
title: 次のアカウントへの切り替え
description: ccam nextは最後に使ったアカウントの次のアカウントでClaude Codeを開きます。--autoは5時間の使用量が最も少ないアカウントを選び、--continueは会話を再開します。
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next`は、アカウント名を指定しなくても、別のアカウントとしてClaude Codeを開きます。アカウントが使用量の上限に達したときに使います。オプションなしの`ccam next`は最後に使ったアカウントの次のアカウントを、`--auto`は5時間の枠に最も余裕があるアカウントを選び、`--continue`は新しいアカウントで最後の会話を続けます。

## 使い方

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccamは切り替えを`→ <from> → <to>`として標準エラー出力に表示し、`ccam <name>`と同じ方法で新しいアカウントとしてClaude Codeを起動します。`<from>`は最後に使ったアカウントです。`--`の後の引数は、そのまま`claude`に渡されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--auto` | ログイン済みのアカウントのうち、5時間の使用量が最も少ないものを選びます。100%のアカウントはスキップします | オフ |
| `-c`、`--continue` | 最後に使ったアカウントと会話履歴を共有しているアカウントだけを候補にし、`claude --continue`を起動します | オフ |
| `-- <claude args>` | ccam自身の引数の後に`claude`に渡す引数 | なし |
| `--dry-run` | 対応しています。`→ <from> → <to>`と実行するコマンドラインを表示し、Claude Codeの起動も起動の記録も行いません | オフ |

`--auto`と`--continue`は一緒に使えます。`--continue`を付けると、`--continue`は自分で指定した引数より前に置かれます。

## しくみ

ccamを通じて起動するたびに、その時刻が`~/.local/state/ccam/last-used/<name>`に書き込まれます。`ccam next`はこれらのファイルを読み、どのアカウントがログイン済みかを確認して、アカウントを選びます。

1. `cli.candidates`が、`config.toml`からアカウントごとに1つの`next.Candidate`を作ります。内容は、アカウント、最後に起動した時刻（`launcher.LastUsed`）、ログインがあるかどうか（`creds.Has`: 読み取り可能な`.credentials.json`、またはmacOSではキーチェーンの項目）です。
2. `--auto`がない場合、`next.Rotate`はアカウントをスロット番号順に並べ、最も最近起動したアカウントを見つけます。それが切り替え元のアカウントです。まだどのアカウントも起動していない場合は、最も小さいスロットになります。
3. そこからスロット順に進み、最後のスロットの次は最初のスロットに戻って、最初に見つかったログイン済みのアカウントを返します。
4. `launcher.LaunchWith`が、`ccam <name>`とまったく同じように、選んだアカウントを準備して起動します。アカウントのリンクを確認し、MCPの共有を適用し、セッションを登録して、`syscall.Exec`でccamを`claude`に置き換えます。

探索は最後に使ったアカウントから始まり、その後で最初に見つかったログイン済みのアカウントで止まります。

<!-- motion:next-rotate -->

`work`が最も最近起動されたので、探索はその次から始まります。`side`はログインしていないのでスキップされ、`lab`が選ばれます。`lab`もログアウトしていた場合は、探索は`personal`に戻ります（図の`% len`）。ローテーションが開始元のアカウントを返すことはありません。ほかにログイン済みのアカウントがない場合は、エラーで停止します。

### --auto

`--auto`は、スロット順ではなく使用量で決めます。`cli.pickAuto`はログイン済みのアカウントを残し、`ccam usage`と同じ`Fetcher.FetchAll`で、それらすべての使用量を一度に問い合わせ、`usage.Best()`に選ばせます。

<!-- motion:next-auto -->

`usage.Best()`は5時間の使用率が最も低いアカウントを選び、100%のアカウントや5時間の使用量が不明なアカウントはスキップします。問い合わせは`ccam usage`の規則に従います。取得してから60秒未満の数値は再利用し、各アカウントの制限時間は3秒で、応答がないアカウントやログインの期限が切れたアカウントは保存済みの数値で判断します。保存済みの時間枠のリセット時刻が過ぎている場合は0%として扱うので、その後リセットされたアカウントが見送られることはありません。`--auto`がトークンを更新することはありません。通常のローテーションとは異なり、`--auto`は、最後に使ったアカウントに最も余裕が残っていれば、そのアカウントを選ぶこともあります。

### --continue

Claude Codeの`--continue`は、現在のフォルダーで最も新しい会話を開きます。これが機能するのは、新しいアカウントからその会話が見える場合だけです。そのため、`--continue`を付けると、ccamは最後に使ったアカウントと同じ`history`グループにあるアカウントだけを候補にし、自分で指定した引数の前に`--continue`を置きます。履歴が`local`のアカウントは誰とも共有していないので、そのようなアカウントからの`--continue`は常に失敗します。先に`ccam share <name> --history on`で共有を有効にしてください（[共有](/docs/sharing)を参照）。

## 安全性と特殊なケース

- ログインしていないアカウントが選ばれることはありません。ほかにログイン済みのアカウントがない場合、オプションなしの`ccam next`は`no logged-in account other than <name>`で停止します。
- `--continue`を付けて、履歴を共有しているアカウントがない場合: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``
- `--auto`を付けて、ログイン済みのアカウントがない場合: `no logged-in account to pick`。すべてのアカウントが上限に達しているか使用量が不明な場合: `` no account with known usage and headroom (see `ccam usage`) ``
- `--auto`は`--dry-run`を付けてもネットワークリクエストを行います。読み取るだけで、トークンファイルに書き込むことはありません。
- 選ばれたアカウントは`ccam <name>`と同じ経路で起動するので、同じ確認が行われます。同じ履歴を共有する別のClaude Codeがすでにこのフォルダーで実行中の場合の警告も含まれます。
- `last-used`は、Claude Codeが実際に起動したときにだけ書き込まれます。`--dry-run`はこれを変更しないので、ドライランでローテーションが進むことはありません。

## 関連項目

- [使用量の上限](/features/usage)
- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [使用量の上限に達したときのワークフロー](/guides/usage-limit-workflow)
