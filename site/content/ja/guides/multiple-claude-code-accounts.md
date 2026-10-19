---
title: 1台のPCで複数のClaude Codeアカウントを使う
description: CLAUDE_CONFIG_DIRで複数のClaude Codeアカウントを並行して使い、切り替える方法。シェルエイリアスによる手動設定とccamによる自動設定を説明します。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2f4deb40c81
---

アカウントごとに専用の設定フォルダーを割り当てます。Claude Codeは`CLAUDE_CONFIG_DIR`環境変数からフォルダーを読み取り、各フォルダーにはそれぞれのログイン、設定、履歴が保存されます。そのため、どのアカウントもログインしたままになり、同時に実行することもできます。シェルエイリアスを使って手動で設定することも、ccamに任せることもできます。

## 手動で設定するには？

これは[Claude Codeのドキュメントで説明されている方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)です。追加するアカウントごとに1行を`~/.bashrc`または`~/.zshrc`に追加します。

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

新しいターミナルを開き、`claude-work`を実行します。フォルダーが新しいので、Claude Codeからログインを求められます。それ以降は、`claude`が1つ目のアカウント、`claude-work`が2つ目のアカウントになります。アカウントごとにエイリアスを1つずつ、それぞれ別のフォルダーで追加してください。

WSLを使わないWindowsでは、PowerShellのプロファイル（`notepad $PROFILE`）に関数を追加すると同じことができます。

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## エイリアスの方法で足りないものは？

アカウントが2つならこの方法で十分です。知っておくべき点は次のとおりです。

- 各フォルダーは空の状態から始まります。設定も、`CLAUDE.md`も、MCPサーバーも、エージェントやスキルもありません。必要なものをコピーまたはリンクし、コピーの内容を自分でそろえ続ける必要があります。
- エイリアスは対話型シェルの中にしか存在しません。スクリプト、Makefile、その他のプログラムからは見えません。
- MCPサーバーはログインと同じ`.claude.json`に保存されるため、単純にフォルダー間でリンクすることはできません。
- 会話は開始したアカウントに残るため、別のアカウントで続けることはできません。
- すべてのアカウントの使用量を一度に表示する手段はありません。

## 代わりに/loginで切り替えられますか？

たまに切り替えるだけなら可能です。Claude Codeの中で`/login`と入力すると、今いるフォルダーの[アカウントが切り替わります](https://code.claude.com/docs/en/quickstart)。ログインできるのは一度に1つのアカウントだけなので、切り替えるたびにログインし直す必要があり、2つのアカウントを並行して実行することもできません。フォルダーを分ければ、どちらの問題も避けられます。

## ccamで設定するには？

ccamは同じ考え方を自動化し、足りない部分を補います。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

これで、`ccam personal`と`ccam work`で2つのアカウントを起動できます。仕事用のアカウントに付けたショートカット`ccwork`でも起動できます。これらはエイリアスではなく実際のプログラムなので、スクリプトの中でも動作します。`ccam`だけを実行すると一覧からアカウントを選べ、`ccam use work`でそのアカウントをデフォルトにできます。

設定、`CLAUDE.md`、MCPサーバー、拡張機能は最初から共有されます。会話履歴の共有を有効にすれば、`ccam next --continue`でもう一方のアカウントに切り替えて会話を続けられます。`ccam usage`は、すべてのアカウントの5時間と週間の使用量を一度に表示します。

すでに`~/.claude`を使っている場合は、そのアカウントを改めて追加する代わりに`ccam import`を実行してください。または`ccam migrate`で、その設定、MCPサーバー、拡張機能を共有設定にコピーできます。

## どのアカウントで実行中かを確認するには？

Claude Codeの中では、`/status`でログイン中のアカウントを確認できます。ccamでは、`ccam statusline install`を一度実行しておくと、ステータスラインにアカウント名と使用量が常に表示されます。たとえば`[work] · Opus · 5h 24% · 7d 41%`のように表示されます。

## WindowsやWSLでも使えますか？

Claude Codeは[Windowsでネイティブに、またはWSLの中で動作します](https://code.claude.com/docs/en/setup)。`CLAUDE_CONFIG_DIR`はどちらでも機能するため、ネイティブのWindowsでは上記のエイリアスやPowerShell関数で十分です。ccamはLinuxとmacOSで動作します。Windowsでは、WSLの中にインストールして使ってください。Linuxと同じように動作します。

## トラブルシューティング

**毎回ログインを求められる。** 新しいフォルダーで初めて起動したときは、必ずログインを求められます。それ以降も求められ続ける場合は、ログインの有効期限が切れたか、削除されています。`ccam login work`を実行してください（またはエイリアスで一度起動してログインします）。

**MCPサーバーが1つのアカウントにしか表示されない。** セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントに属します。`ccam mcp import work <server>`で、すべてのアカウントに移せます。

**アカウントごとに設定が異なる。** `ccam share work`で、そのアカウントが共有している内容を表示できます。`ccam share work --settings on`は、変更内容を表示したうえでそのアカウントの設定を共有します。`ccam doctor --fix`は壊れたリンクを修復し、変更するすべてのファイルのバックアップを残します。

**2つのアカウントで同じ会話を開いている。** 1つの会話を2つのアカウントで同時に進めるのは避けてください。そうしようとすると、ccamが警告します。

## どちらを選ぶか

アカウントが2つで何も共有しないなら、エイリアスで十分です。設定やMCPサーバーを共有したい場合、アカウントがもっと多い場合、使用量をひと目で確認したい場合、会話を続けたい場合は、ccamを使えば管理の手間が省けます。[エイリアスとの比較](/compare/aliases)で違いを並べて紹介しています。1つのログインを入れ替える方式のツールについては、[ccamとほかのアカウント切り替えツール](/compare/account-switchers)を参照してください。

複数のアカウントを使ってよいかどうかは、各アカウントの規約によります。使う前に確認してください。
