---
title: ccamとほかのClaude Codeアカウント切り替えツール
description: ccamとclaude-swap、CCS、clausona、clauth、CAAM、cuxなど、Claude Codeのアカウントを切り替えるツールを比較し、用途に合う方式を説明します。
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "6d08203f1ea6"
---

Claude Codeのアカウントを切り替えるツールの方式は2つに分かれます。多くは1つの`~/.claude`フォルダーの中でログインを入れ替えるので、使えるアカウントは一度に1つです。ccamを含むほかのツールは、`CLAUDE_CONFIG_DIR`でアカウントごとに専用のフォルダーを割り当てるので、どのアカウントもログインしたままで、複数を同時に実行できます。ccamはさらに設定とMCPサーバーを共有し、別のアカウントのフォルダーで会話を続けられます。

## ログインを入れ替えるか、フォルダーを分けるか

**ログインを入れ替える方式**では、フォルダーが1つしかないため、設定、MCPサーバー、履歴も1組だけです。切り替えは速く、実行中のセッションが新しいログインを引き継ぐこともできます。ただし、使えるアカウントは一度に1つで、すべてのターミナルとエディターに同時に適用されます。

入れ替えるために、ツールはほかのすべてのアカウントのログインのコピーを持っておき、それを所定の場所に書き込みます。Claude Codeはログインを更新するたびにその中のリフレッシュトークンを置き換えるので、保存したコピーはどれも、実行中のセッションと常に同期させておく必要があります。切り替えた後にアカウントがログアウトされる、あるいは`/login`をやり直す必要がある、といった問題は、こうしたプロジェクトが最もよく修正している問題の1つです。

**フォルダーを分ける方式**は、[Claude Codeのドキュメントで説明されている方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)です。どのアカウントもログインしたままで、仕事用と個人用のセッションを並行して実行できます。各ログインはそれぞれのフォルダーに置かれたままで、Claude Codeがいつもどおりそこで更新します。その代わり、ツールで共有しない限り、各フォルダーは設定、MCPサーバー、履歴をそれぞれ別に持つ状態から始まります。

## ツールの比較

2026-10-20時点の各プロジェクトのドキュメントとソースコードに基づきます（clausona: v0.12.0-beta）。「記載なし」は、プロジェクトが言及していないという意味で、機能がないという意味ではありません。

| ツール | 方式 | 複数アカウントの同時使用 | 5時間と週間の使用量 | 上限に達したとき | プラットフォーム |
| --- | --- | --- | --- | --- | --- |
| ccam | フォルダーを分ける | 可 | 可、`ccam usage` | セッション内で通知する。`ccam next --auto --continue`か`/ccam:switch`が、履歴の共有の有無にかかわらず、最も余裕があるアカウントに会話を持ち越す。`-p`のジョブは`ccam --fallback`がそのアカウントで自動的に再実行し、`ccam pick`はワーカーを分散する | Linux、macOS、Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | ログインを入れ替える。`cswap run`ではフォルダーを分ける | 可（`cswap run`） | 可 | `cswap auto`がしきい値でログインを切り替え、実行中のセッションはそのまま続く | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | フォルダーを分ける | 可 | 一部 | アカウントプールがローカルプロキシ経由でリクエストを次のアカウントに回す | Linux、macOS、Windows |
| [clausona](https://github.com/larcane97/clausona) | フォルダーを分ける | 可 | 可 | 各実行の前にルートがクォータの残っているアカウントを選び、リセットを待つこともでき、ワーカーには別々のアカウントを割り当てる | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | ログインを入れ替える。`clauth start`ではフォルダーを分ける | 可（`clauth start`） | 可 | フォールバックチェーンに沿って切り替える | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | ログインを入れ替える、または分離したプロファイル | 可（分離したプロファイル） | 可 | `caam run`が非対話型のコマンドを別のアカウントで再試行する | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | ログインを入れ替え、`claude`をラップする | 不可、使えるアカウントは1つ | 可 | ログインを入れ替えて同じ会話を再開する | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | ログインを入れ替える | 不可 | 可 | しきい値で切り替え（オプション） | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | ログインを入れ替える | 不可 | 可 | 100%で切り替え（オプション） | macOS |

ccswapはclaude-swapのフォークで、Codexのログインも管理します。[claude-account](https://github.com/hamzarehmandeveloper/claude-account)、[claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile)、[claude-code-profiles](https://github.com/quinnjr/claude-code-profiles)も、アカウントごとに専用のフォルダーを割り当てます。これらのツールのいくつかは、CodexやGeminiなどほかのコーディングエージェントにも対応しています。

## 各ツールのつくり

| ツール | 実装言語 | Claude Codeの実行中 | ログインを別の場所にコピーするか | アカウント間で共有できるもの |
| --- | --- | --- | --- | --- |
| ccam | Go、静的バイナリ1つ | 何もない。ccamは自身を`claude`に置き換える（置き換えられないWindowsでは終了を待つ。オプションのリレーはそばに留まる） | しない | グループ。各グループが共有する項目とアカウントを決め、1つのアカウントは複数のグループに入れる |
| claude-swap | Python 3.12以上、uvまたはpipxでインストール | `cswap run`では何もない。`cswap auto`とメニューバーアプリは動き続ける | する（ログインを入れ替えるとき） | `cswap run`の場合、設定、`CLAUDE.md`、拡張機能をまとめて、履歴は任意。MCPサーバーは起動のたびにコピー |
| CCS | Node.js 18以上またはBun | CCSが親プロセスとして残る。プールはローカルプロキシを追加 | アカウントフォルダーではしない。プールは独自のログインを持つ | アカウントごとに共有かローカル、履歴はコンテキストグループで。MCPサーバーはコピー |
| clausona | TypeScript、Node.js 20以上 | 素の`claude`はシェル関数を通して実行。`clausona run`ではNodeが親プロセスとして残る。終了のたびにバックグラウンドのNodeプロセス | しない。ただし使用量を読むため、使っていないアカウントのトークンを自分で更新する | アカウントごとのファイル以外すべてを1つのプライマリから。履歴は任意。`claude mcp add`で追加したサーバーは共有されない |
| clauth | Rust | clauthが親プロセスとして残る | する（ログインを入れ替えるとき） | `~/.claude`全体か、何も共有しないか |
| CAAM | Go | モードによる。デーモンは任意 | する（ボールトモード） | モードで決まる。MCPサーバーは起動のたびにコピー |
| cux | Go | ラッパーが`claude`を包んだまま残る | する | すべて。フォルダーは1つ |
| CCSwitcher, Claude Switcher | macOSアプリ（Swift、Python） | メニューバーアプリ | する | すべて。フォルダーは1つ |

## ccamの強み

### ログインはClaude Codeが置いた場所から動かない

ccamは、Claudeのログインをある場所から別の場所へコピーすることはありません。各アカウントのログインはそれぞれのフォルダー、macOSではそのフォルダー専用のキーチェーン項目に置かれたままで、更新するのはClaude Codeです。ずれてしまう保存済みのコピーはありません。`ccam import`で登録した既存の`~/.claude`は、通常の`claude`が使うキーチェーン項目をそのまま使うので、どちらから起動してもログインしたままです。

ccamが期限切れのトークンを自分で更新するのは、`ccam usage --refresh`を実行したときだけです。先にロックを取り、そのアカウントが実行中であれば、Claude Codeがどう起動されたかにかかわらず（バックグラウンドセッションとそのデーモンも含みます）、またClaude Code自身がそのログインを更新している最中も、更新を拒否します。そのため、動いているセッションの下でトークンを入れ替えることはありません。

### あなたとClaude Codeの間では何も動かない

ccamはファイルを1つ読み、リンクを確認し、自身を`claude`に置き換えます（`exec`）。作業中はccamのプロセスも、ラッパーも、プロキシも、デーモンもありません。標準入力、シグナル、終了コードはClaude Codeのもので、リクエストはAnthropicへ直接送られます。Windowsではプロセスを置き換えられないため、ccamは同じコンソールでClaude Codeを起動して終了コードを待つだけで、ほかには何もしません。

この方法でアカウントを起動するときの追加時間は約4 ms、95パーセンタイルでも6 ms未満で、ピークメモリーは8 MBです（hyperfineで[`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh)を300回実行）。同じノートパソコンで、Pythonインタープリターは起動するだけで13 ms、Node.jsは31 msかかります。プログラムが何かをする前の時点でです。

### 何を共有するかはグループが決める

MCPサーバー、`settings.json`、`CLAUDE.md`、拡張機能（エージェント、コマンド、スキル、プラグイン）、会話履歴は、それぞれ独立した5つの項目です。グループは共有する項目と所属するアカウントを決めるので、2つの仕事用アカウントは`work`グループで履歴を共有し、個人用アカウントは自分の履歴を別に持ちながら、3つとも`default`を通じて同じMCPサーバーを使う、といったことができます。`ccam group set work --add memory`は、グループのすべてのアカウントを1つの計画で移行します。[ダッシュボード](/features/dashboard)ではチェックボックスで同じことができ、アカウントがすでに別のグループから受け取っている項目はグレー表示になります。

- MCPサーバーは起動時に`--mcp-config`で渡すので、既定では起動しても、Claude Codeが実行中に書き換えるファイルである`.claude.json`には手を付けません。ccamがそこに書き込むのは、そう明示した手順だけで、Claude Codeのロックを取り、バックアップを残してから行います。対象は、確認を取ったうえでのRemote Control、アカウントがMCPの共有から抜けるときや戻るとき、そして`ccam mcp import`です。`ccam mcp add`は、グループ全体にサーバーを1回で追加します。
- ほかの項目はグループのフォルダーへのリンクなので、1つのアカウントでの変更はすぐにほかのアカウントにも反映されます。起動時にコピーするものはありません。
- ログイン、`.claude.json`、実行中の状態は決して共有しません。

### すべての変更を先に見せ、バックアップする

`ccam group`と`ccam migrate`は、何かを移動する前に計画を表示して確認を求めます。アカウントとグループでファイルの内容が異なる場合は、グループのものを残すか、自分のものを使うか、JSONをマージするかを選べます。フォルダーの中ではグループのコピーが残り、計画に内容の異なるファイルの名前が表示されます。元のファイルは`.ccam-backup/<time>/`に移ります。削除したアカウントはログアウトされ、そのフォルダーはごみ箱フォルダーに移ります。また、ファイルは1回のリネームで置き換えるので、書き込みが途中で止まっても古いファイルが残ります。`--dry-run`は、起動、共有、移行、ショートカットコマンド、アップデート、修復が何をするかを表示します。`ccam doctor --fix`は、壊れたリンクと権限を修復します。

### 次のアカウントで同じ会話を

`ccam next --auto --continue`がすべてのアカウントの使用量を一度に問い合わせ、最も余裕があるアカウントを開いて会話を続けます。履歴を共有していれば同じ会話のままです。共有していなければ、Claude Codeが会話を新しいアカウントにコピーし（`--resume <file> --fork-session`）、ccamは最初の応答で読み直すトークン数を表示します。上限でセッションが止まると、ccamプラグインのフックが通知を表示し、`/ccam:switch`が会話を引き渡します。切り替えるタイミングはあなたが決めます。バックグラウンドで使用量を監視するものはなく、自動的な選択は1つの組織の中で行われるので、会社の会話がうっかり個人のアカウントに移ることはありません。アカウントが90%に達すると、ステータスラインがこのコマンドを提案します。

スクリプトはあなたを待てないので、`ccam --fallback work -p "…"`がスクリプトのために切り替えます。アカウントが上限に達すると、同じコマンドを最も余裕があるアカウントで再実行し、最後まで完了した実行の出力だけを表示します。使用量の上限以外の失敗は、そのまま返します。並列のワーカーには、`ccam pick`が余裕のある別々のアカウントを1つずつ渡します。どのアカウントにも余裕がなければ、どちらも終了コード75で終了します。

**TeamとEnterpriseのシート。** 自動的な選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。そのため、シート間で自動的な選択を使うと、Anthropicの規約に違反するおそれがあります。自動的な選択に頼る前に、所属する組織とAnthropicに確認してください。個人プラン（Pro、Max）については、[FAQの規約に関する説明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を参照してください。

### フォルダーごとに適切なアカウント

`ccam folder work ~/code/company`を実行すると、そのフォルダー内でアカウント名を指定せずに起動したときは常に`work`が使われ、ほかのプロジェクトは既定のアカウントのままです。この設定は、起動のたびに読む同じ`config.toml`に保存されるので、時間はかかりません。

### バイナリ1つ、検証済みのダウンロード

ccamは約10 MBの静的バイナリで、インストールが必要なランタイムはなく、Linux、macOS、Windowsに対応しています。インストールスクリプト（`install.sh`、`install.ps1`）と`ccam update`は、すべてのダウンロードをリリースのSHA-256チェックサムと照合します。ccamが確認なしにアップデートすることはありません。テレメトリーもありません。ccamが送るリクエストは、あなたのアカウントの使用量をAnthropicに問い合わせるもの（並列、3秒の制限時間、60秒のキャッシュ）、ダッシュボードを開いたときに最新のリリースをGitHubに問い合わせるもの、あなたが指示したダウンロードだけです。`usage_source = "claude"`を設定すると、その問い合わせもClaude Codeを通じて行われ、ccamがトークンを扱うことは一切ありません。

## ほかのツールのほうが得意なこと

- claude-swap、clauth、cux、CAAM、CCSwitcherは、アカウントが使用量のしきい値を超えると、実行中の対話型セッションを自動で切り替えられます。ccamは`-p`の実行を`--fallback`で切り替え、対話型セッションは、上限での通知と`/ccam:switch`を用意したうえで、あなたに任せます。
- 開発者モードがオフのWindowsでは、ccamはフォルダーを共有できますが、`settings.json`のような単一のファイルは共有できません（その場合、最初の`ccam add`はMCPサーバーと拡張機能だけを共有します）。claude-swapとCCSはその場合ファイルをコピーします。
- CCSwitcher、Claude Switcher、claude-swapは、macOSのメニューバーに使用量を表示します。
- CCS、clausona、CAAM、clauth、ccswap、Claude Switcherは、Codexやほかのコーディングエージェントも管理します。
- clausonaは、アカウントごとに拡張機能をオン・オフでき（ccamの`ccam ext`は表示するだけです）、ローカルのトークンコストを記録し、tmuxなどのターミナルで多数のワーカーを動かすClaude Codeプラグインも提供しています。
- CAAM、CCSwitcher、Claude SwitcherはHomebrewでインストールでき、clauthとCAAMはリリースのチェックサムに署名しています。ccamはスクリプト、`.deb`、`go install`でインストールし、署名なしでチェックサムを確認します。

## どれが合うか

一度に1つのアカウントしか使わない場合や、実行中のセッションを自動で切り替えたい場合は、ログインを入れ替えるツールを選んでください。Linux、macOS、Windowsで、複数のアカウントを並行して使い、1つの共通の環境を保ち、次のアカウントで同じ会話を続け、Claude Codeの横で何も動かしたくない場合は、ccamを選んでください。何も共有する必要がなければ、[普通のシェルエイリアス](/compare/aliases)で十分です。

複数のアカウントを使ってよいかどうかは、各アカウントの規約によります。使う前に確認してください。
