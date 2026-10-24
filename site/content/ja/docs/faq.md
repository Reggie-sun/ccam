---
title: よくある質問
metaTitle: "ccamのよくある質問：Claude Codeの複数アカウント"
description: ccam、Claude Codeのアカウント、使用量の上限、共有、プライバシー、対応システムについての回答です。
updated: 2026-10-24
order: 6
translation: machine
sourceHash: "c4e3b89bb4ea"
---

## ccamはAnthropicの公式ツールですか？

いいえ。ccamは非公式のオープンソースプロジェクトであり、Anthropicとは提携しておらず、Anthropicの承認も受けていません。ccamは、複数のアカウントを並行して使うために[Claude Codeがドキュメントで説明している](https://code.claude.com/docs/en/env-vars)`CLAUDE_CONFIG_DIR`設定を利用しています。

## 複数のClaude Codeアカウントを使ってもよいのですか？

それは各アカウントの規約によります。仕事用のアカウントの場合は、組織のポリシーにもよります。複数のアカウントを使う前に確認してください。ccamは、すでに持っているアカウント間の切り替えを簡単にするだけです。

**TeamとEnterpriseのシート。** 自動的な選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。そのため、シート間で自動的な選択を使うと、Anthropicの規約に違反するおそれがあります。自動的な選択に頼る前に、所属する組織とAnthropicに確認してください。個人プラン（Pro、Max）については、上の段落を参照してください。

## ccamで使用量の上限を回避できますか？

いいえ。各アカウントにはそれぞれの上限がそのまま適用されます。ccamは上限を表示し、あるアカウントの上限に達したときに、自分が所有する別のアカウントで作業を続けられるようにします。

## 1台のコンピューターで2つのClaude Codeアカウントを使うにはどうすればよいですか？

`CLAUDE_CONFIG_DIR`を使って、各アカウントに専用の設定フォルダーを割り当てます。ccamはこれを自動で行います。アカウントごとに`ccam add`を実行し、`ccam personal`と`ccam work`、または`ccwork`のように自分で名前を付けたショートカットで起動してください。手動で設定する方法は[ガイド](/guides/multiple-claude-code-accounts)でも説明しています。

## 別のアカウントで会話を続けるにはどうすればよいですか？

そのフォルダーで`ccam next --continue`を実行します。フォルダーの最後の会話が次のアカウントに引き継がれます（最も余裕があるアカウントにするには`--auto`）。Claude Codeの中では、[ccamプラグイン](/docs/install#claude-code-plugin)の`/ccam:switch`が、いま使っているセッションに対して同じことをします。2つのアカウントが会話履歴を共有していれば（`ccam group set default --add history`、またはそのアカウント専用のグループ: `ccam group add pair --items history --members work,side`）、新しいアカウントはまさにその会話を開きます。共有していなければ、Claude Codeが会話を新しいアカウントの履歴にコピーします。

## ccamはログイン情報を見たり送信したりしますか？

ccamがClaudeのログインをアカウント間で共有したりコピーしたりすることはなく、MCPサーバーのサインイン情報をコピーすることもありません。アカウントのアクセストークンを読み取るのは、`ccam usage`、`ccam next --auto`、`ccam doctor`を実行したときに、Anthropicの使用量エンドポイントにそのアカウントの使用量を問い合わせるためだけです。期限切れのトークンを更新するのは`ccam usage --refresh`を実行したときだけで、`usage_source = "claude"`ではトークンをまったく読み取りません。`ccam ext`と`ccam doctor`が読み取るのはアカウントのMCPサインインの名前だけで、そのトークンは読み取りません。ccamにはテレメトリーがありません。これらの問い合わせ以外にccamが自分から送るリクエストは、ダッシュボードが最新のccamリリースをGitHubに問い合わせるものだけで、そこにはあなたやあなたのアカウントに関する情報は一切含まれません。

## リモートサーバーでサインインするにはどうすればよいですか？

SSHセッションで、いつもどおり`ccam add`または`ccam login <name>`を実行します。Claude Codeはサーバー上でサインインを待つため、手元のコンピューターのブラウザーは最後に`http://localhost:<port>/callback?code=…`のページに移動しますが、このページは読み込まれません。そのアドレスをブラウザーからコピーしてターミナルに貼り付けてください。ccamがそれをサーバー上のClaude Codeに渡し、Claude Codeがサインインを完了します。アドレスには一度だけ使えるコードが含まれているので、他人と共有しないでください。

## アカウントはどこに保存されますか？

`~/.local/share/ccam/profiles/<name>/`に、アカウントごとに1つのフォルダーとして保存されます。macOSでは、Claude Codeはログイン情報そのものをキーチェーンに、そのフォルダーから導出した名前で保存します。そのため、`ccam rename`はそこでフォルダーの古い名前をそのまま使います。

## 既存の~/.claudeはどうなりますか？

`ccam import`か`ccam migrate`を実行しない限り、何も変わりません。インポートはフォルダーをその場所のまま登録するので、通常の`claude`も`ccam <name>`も、引き続き同じアカウントと設定を使います。`ccam migrate`は、フォルダーを読み取るだけで、その設定、`CLAUDE.md`、拡張機能、MCPサーバーを共有設定にコピーします。

## ccamはどのシステムに対応していますか？

Linux（Ubuntu、Debian、Fedora、Arch、Alpineなど）、macOS、Windows 10と11（ネイティブ、またはWSLの中）です。Windowsで`settings.json`のようなファイルを共有するには開発者モードが必要です。共有フォルダーは開発者モードなしで動作し、開発者モードがなければ、最初の`ccam add`はMCPサーバーと拡張機能だけを含むグループ`default`を作成します（[Windowsでの共有](/features/share#on-windows)）。

## MCPサーバーが1つのアカウントにしか表示されないのはなぜですか？

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントにだけ保存され、`-s user`を付けなければ現在のフォルダーにだけ保存されます。ccamは次にそのアカウントを起動したときにそのことを知らせ、`ccam mcp list`と`ccam ext`にもそれらが表示されます。`ccam mcp import <account> <server>`でサーバーを共有グループに移動してください。

## ccamを使うとClaude Codeは遅くなりますか？

ccam経由でClaude Codeを起動すると、起動時間が約4ミリ秒増えます。これはノートパソコンで`scripts/perf-check.sh`を使って測定した値です。ccamは小さなファイルを1つ読み、共有リンクを確認してから自身を`claude`に置き換えるため、作業中にccamが動作していることはありません。Windowsではプロセスを置き換えられないため、ccamはClaude Codeを起動して終了を待つだけで、ほかには何もしません。

## ccamをアップデートするにはどうすればよいですか？

`ccam update`を実行します。最新のリリースをダウンロードしてチェックサムを検証し、`ccam`ファイルを置き換えます。`ccam update --check`は、新しいリリースがあるかどうかを知らせるだけです。ダッシュボードも新しいリリースがあれば知らせ、`U`を押すと確認を取ってからインストールします。ccamが自動で更新されることはありません。`.deb`パッケージでインストールした場合は、代わりに新しい`.deb`をaptでインストールしてください。詳しくは[アップデート](/docs/install#update)を参照してください。

## ccamをアンインストールするにはどうすればよいですか？

各アカウントのショートカットを`ccam shortcut <account> --remove`で削除します。`~/.claude`をインポートしていて、それがグループに入っている場合は、先にそのアカウントに対して`ccam remove <name>`を実行してください。共有しているもののコピーがそのアカウントに戻ります。次に`~/.local/bin`から`ccam`を削除し、アカウントが不要になった場合は`~/.config/ccam`、`~/.local/share/ccam`、`~/.local/state/ccam`も削除します。詳しくは[インストール](/docs/install#uninstall)を参照してください。
