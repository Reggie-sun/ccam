---
title: よくある質問
description: ccam、Claude Codeのアカウント、使用量の上限、共有、プライバシー、対応システムについての回答です。
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 5ff6e94786a2
---

## ccamはAnthropicの公式ツールですか？

いいえ。ccamは非公式のオープンソースプロジェクトであり、Anthropicとは提携しておらず、Anthropicの承認も受けていません。ccamは、複数のアカウントを並行して使うためにClaude Codeがドキュメントで説明している`CLAUDE_CONFIG_DIR`設定を利用しています。

## 複数のClaude Codeアカウントを使ってもよいのですか？

それは各アカウントの規約によります。仕事用のアカウントの場合は、組織のポリシーにもよります。複数のアカウントを使う前に確認してください。ccamは、すでに持っているアカウント間の切り替えを簡単にするだけです。

## ccamで使用量の上限を回避できますか？

いいえ。各アカウントにはそれぞれの上限がそのまま適用されます。ccamは上限を表示し、あるアカウントの上限に達したときに、自分が所有する別のアカウントで作業を続けられるようにします。

## 1台のコンピューターで2つのClaude Codeアカウントを使うにはどうすればよいですか？

`CLAUDE_CONFIG_DIR`を使って、各アカウントに専用の設定フォルダーを割り当てます。ccamはこれを自動で行います。アカウントごとに`ccam add`を実行し、`ccam personal`と`ccam work`、または`ccwork`のように自分で名前を付けたショートカットで起動してください。手動で設定する方法は[ガイド](/guides/multiple-claude-code-accounts)でも説明しています。

## 別のアカウントで会話を続けるにはどうすればよいですか？

`ccam share <name> --history on`で2つのアカウント間で会話履歴を共有してから、`ccam next --continue`を実行します。Claude Codeが、現在のフォルダーで最も新しい会話を次のアカウントで開きます。

## ccamはログイン情報を見たり送信したりしますか？

ccamがClaudeのログインをアカウント間で共有することはありません。アカウントのアクセストークンを読み取るのは、`ccam usage`、`ccam next --auto`、`ccam doctor`を実行したときに、Anthropicの使用量エンドポイントにそのアカウントの使用量を問い合わせるためだけです。`ccam migrate`は、MCPサーバーの保存済みのサインイン情報を`~/.claude`から、同じコンピューター上にある各アカウント自身のログインファイルにコピーします。ccamにはテレメトリーがなく、それ以外の場所には何も送信しません。

## リモートサーバーでサインインするにはどうすればよいですか？

SSHセッションで、いつもどおり`ccam add`または`ccam login <name>`を実行します。Claude Codeはサーバー上でサインインを待つため、手元のコンピューターのブラウザーは最後に`http://localhost:<port>/callback?code=…`のページに移動しますが、このページは読み込まれません。そのアドレスをブラウザーからコピーしてターミナルに貼り付けてください。ccamがそれをサーバー上のClaude Codeに渡し、Claude Codeがサインインを完了します。アドレスには一度だけ使えるコードが含まれているので、他人と共有しないでください。

## アカウントはどこに保存されますか？

`~/.local/share/ccam/profiles/<name>/`に、アカウントごとに1つのフォルダーとして保存されます。macOSでは、Claude Codeはログイン情報そのものをキーチェーンに、そのフォルダーから導出した名前で保存します。

## 既存の~/.claudeはどうなりますか？

`ccam import`か`ccam migrate`を実行しない限り、何も変わりません。インポートはフォルダーをその場所のまま登録するので、通常の`claude`も`ccam <name>`も、引き続き同じアカウントと設定を使います。`ccam migrate`は、フォルダーを読み取るだけで、その設定、`CLAUDE.md`、拡張機能、MCPサーバーを共有設定にコピーします。

## ccamはどのシステムに対応していますか？

Linux（Ubuntu、Debian、Fedora、Arch、Alpineなど）とmacOSです。Windowsでは、WSLの中で使ってください。

## MCPサーバーが1つのアカウントにしか表示されないのはなぜですか？

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントにだけ保存されます。ccamは次にそのアカウントを起動したときにそのことを知らせます。`ccam mcp import <account> <server>`でサーバーを共有グループに移動してください。

## ccamを使うとClaude Codeは遅くなりますか？

ccam経由でClaude Codeを起動すると、起動時間が約4ミリ秒増えます。これはノートパソコンで`scripts/perf-check.sh`を使って測定した値です。ccamは小さなファイルを1つ読み、共有リンクを確認してから自身を`claude`に置き換えるため、作業中にccamが動作していることはありません。

## ccamをアップデートするにはどうすればよいですか？

`ccam update`を実行します。最新のリリースをダウンロードしてチェックサムを検証し、`ccam`ファイルを置き換えます。`ccam update --check`は、新しいリリースがあるかどうかを知らせるだけです。ccamが自動で更新されることはありません。`.deb`パッケージでインストールした場合は、代わりに新しい`.deb`をaptでインストールしてください。詳しくは[アップデート](/docs/install#update)を参照してください。

## ccamをアンインストールするにはどうすればよいですか？

各アカウントのショートカットを`ccam shortcut <account> --remove`で削除し、`~/.local/bin`から`ccam`を削除します。アカウントが不要になった場合は、続けて`~/.config/ccam`、`~/.local/share/ccam`、`~/.local/state/ccam`も削除します。詳しくは[インストール](/docs/install#uninstall)を参照してください。
