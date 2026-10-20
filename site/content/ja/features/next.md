---
title: 次のアカウントへの切り替え
description: ccam nextは前回の次のアカウントでClaude Codeを開きます。--autoは余裕が最も多いものを選び、--continueは再開し、--fallbackは-pのジョブを再実行します。
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next`は、アカウント名を指定しなくても、別のアカウントとしてClaude Codeを開きます。アカウントが使用量の上限に達したときに使います。オプションなしの`ccam next`は最後に使ったアカウントの次のアカウントを、`--auto`は5時間の枠に最も余裕があるアカウントを選び、`--continue`は新しいアカウントで最後の会話を続けます。`claude -p`を実行するスクリプトでは、`ccam --fallback`がこの切り替えを自動で行います（[上限に達したら再実行](#run-again-at-the-limit)を参照）。

## 使い方

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
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
| `--auto` | ログイン済みのアカウントを使用量で選びます（下記参照）。100%のアカウントはスキップします | オフ |
| `--max-usage <1-100>` | `--auto`と併用。使用量がこの値未満のアカウントを優先します | `auto_max_usage`、なければ100 |
| `--strategy headroom\|expiring` | `--auto`と併用。5時間の使用量が最も少ないもの、または5時間の枠のリセットが最も早いもの | `auto_strategy`、なければ`headroom` |
| `--wait[=<duration>]` | `--auto`と併用。余裕のあるアカウントがないとき、最初のアカウントがリセットされるまで最長でこの時間だけ待ちます | オフ。値なしで指定すると`5h` |
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

`--auto`は、スロット順ではなく使用量で決めます。`cli.pickAuto`はログイン済みのアカウントを残し、`ccam usage`と同じ`Fetcher.FetchAll`で、それらすべての使用量を一度に問い合わせ、`usage.Policy.Choose()`に選ばせます。

<!-- motion:next-auto -->

アカウントは5時間と週間の枠で判断され、Claude Codeへの引数に`--model`がある場合は、そのモデルの週間の枠でも判断されます（名前に`opus`、`sonnet`、`haiku`、`fable`を含むもの）。`Policy.Choose()`はアカウントを3つに分けます。

1. **余裕あり**: すべての枠が`auto_max_usage`（`config.toml`、デフォルトは100）未満。この中から戦略で選びます。`headroom`（デフォルト）は5時間の使用量が最も少ないものを、`expiring`は5時間の枠のリセットが最も早いものを選び、残りのクォータが失われる前に使い切ります。
2. **超過**: 100%未満だが`auto_max_usage`以上。余裕のあるアカウントがないときだけ使い、5時間の使用量が最も少ないものから選びます。
3. **満杯**: 100%の枠がある、または5時間の使用量が不明。選ばれることはありません。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

選べるアカウントがない場合、ccamは終了コード75（`EX_TEMPFAIL`）で終了するので、スクリプトは「あとで再試行」と失敗を区別できます。`--wait`を付けると、代わりに満杯のアカウントのうち、上限になっている枠が最も早くリセットされるまで（エンドポイントに新しい枠が表示される、リセット時刻の30秒後まで）待ち、`ccam: no account has room; waiting until 15:10`と表示してから、もう一度選びます。どのアカウントもリセット時刻を示さない場合は、10分後にもう一度確認します。待機時間が終わるまでに余裕ができなければ、75で終了します。問い合わせは`ccam usage`の規則に従います。取得してから60秒未満の数値は再利用し、各アカウントの制限時間は3秒で、応答がないアカウントやログインの期限が切れたアカウントは保存済みの数値で判断します。保存済みの時間枠のリセット時刻が過ぎている場合は0%として扱うので、その後リセットされたアカウントが見送られることはありません。`--auto`がトークンを更新することはありません。通常のローテーションとは異なり、`--auto`は、最後に使ったアカウントに最も余裕が残っていれば、そのアカウントを選ぶこともあります。

### --continue

Claude Codeの`--continue`は、現在のフォルダーで最も新しい会話を開きます。これが機能するのは、新しいアカウントからその会話が見える場合だけです。そのため、`--continue`を付けると、ccamは最後に使ったアカウントと同じ`history`グループにあるアカウントだけを候補にし、自分で指定した引数の前に`--continue`を置きます。履歴が`local`のアカウントは誰とも共有していないので、そのようなアカウントからの`--continue`は常に失敗します。先に`ccam share <name> --history on`で共有を有効にしてください（[共有](/docs/sharing)を参照）。

## 上限に達したら再実行

`ccam next`は、止まった対話型セッションのためのものです。非対話型の実行、つまりスクリプトや定期実行ジョブの中の`claude -p`はあなたを待てないので、アカウント名の前に`--fallback`を付けると、ccamがその実行のために切り替えます。

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback`は`-p`か`--print`と一緒の場合にだけ動作します。そうでない場合、ccamは停止して`ccam next --auto --continue`を案内します。
2. 標準入力は一度だけ読み込まれ、すべての試行に渡されます。
3. ccamは`exec`ではなく子プロセスとして`claude`を起動し、標準出力を保留して、標準エラー出力は届いたとおりに流します。
4. `claude`が0で終了した場合や、別の理由で失敗した場合は、保留した出力と終了コードをそのまま返します。
5. Claude Codeのメッセージが使用量の上限を示している場合（`You've hit your limit`、`usage limit reached`）や、アカウントの使用量の問い合わせで5時間または週間の枠が100%になっている場合、その失敗は使用量の上限と見なされます。
6. その場合、その試行の出力は破棄され、同じコマンドが別のログイン済みのアカウントで実行されます。順序は、5時間の使用量が最も少ないアカウント、次に使用量を問い合わせられないアカウントを番号順です。100%のアカウントはスキップされ、各アカウントは最大1回だけ試されます。
7. 引数に`-c`、`--continue`、`-r`、`--resume`がある場合は、会話を見つけられるように、同じ履歴グループのアカウントだけを試します。

次のアカウントは`--auto`と同じ規則（`auto_max_usage`、`auto_strategy`、`--model`の枠）で選ばれます。試せるアカウントが残っていない場合、ccamは最後の試行の出力を表示し、余裕のあるログイン済みのアカウントがないことを伝えて、終了コード75で終了します。アカウント名の前に`--wait`を付けると（`ccam --fallback --wait=2h work -p …`）、代わりに最初のアカウントがリセットされるまで待ち、そのあとすべてのアカウントをもう一度試すことがあります。必ず先に一度は待つので、上限に達したばかりのアカウントがすぐに再実行されることはありません。`ccam --dry-run --fallback work -p …`は、最初の試行のコマンドを表示します。

対話型セッションが自動で切り替わることはありません。代わりに、`ccam statusline install`を実行しておくと、5時間または週間の使用量が90%に達したときに、ステータスラインに`→ ccam next --auto --continue`が追加されます（[使用量の上限](/features/usage)）。

## 安全性と特殊なケース

- ログインしていないアカウントが選ばれることはありません。ほかにログイン済みのアカウントがない場合、オプションなしの`ccam next`は`no logged-in account other than <name>`で停止します。
- `--continue`を付けて、履歴を共有しているアカウントがない場合: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``
- `--auto`を付けて、ログイン済みのアカウントがない場合: `no logged-in account to pick`。すべてのアカウントが上限に達しているか使用量が不明な場合: `` no account with known usage and headroom (see `ccam usage`) ``。どちらも終了コード75で終了します。
- `--max-usage`、`--strategy`、`--wait`は`--auto`と一緒にしか使えず、1-100以外の値やほかの戦略は拒否されます。
- `--auto`は`--dry-run`を付けてもネットワークリクエストを行います。読み取るだけで、トークンファイルに書き込むことはありません。
- 選ばれたアカウントは`ccam <name>`と同じ経路で起動するので、同じ確認が行われます。同じ履歴を共有する別のClaude Codeがすでにこのフォルダーで実行中の場合の警告も含まれます。
- `last-used`は、Claude Codeが実際に起動したときにだけ書き込まれます。`--dry-run`はこれを変更しないので、ドライランでローテーションが進むことはありません。
- `--fallback`は試行が終わるまで標準出力を保留するので、`--output-format stream-json`のようなストリーミング出力は最後にまとめて届きます。最初の試行が上限に達する前に変更したファイルは変更されたままで、次の試行はプロンプトを最初からやり直します。

## 関連項目

- [使用量の上限](/features/usage)
- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [使用量の上限に達したときのワークフロー](/guides/usage-limit-workflow)
