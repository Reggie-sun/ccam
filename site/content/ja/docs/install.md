---
title: ccamのインストール
metaTitle: "Linux、macOS、Windowsにccamをインストール"
description: Linux、macOS、Windowsにccamを1つのコマンド、.debパッケージ、go installでインストールし、ccam updateで更新します。Claude Codeを先に入れてください。
updated: 2026-10-26
order: 1
translation: machine
sourceHash: 2ada7b2650d6
---

ccamがClaude Codeを起動するので、先に[Claude Code](https://docs.claude.com/en/docs/claude-code/overview)をインストールし、`claude --version`が動作することを確認してください。

LinuxまたはmacOSのターミナルで、次のコマンドを実行します。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

このスクリプトは、お使いのシステム用のリリースをダウンロードしてSHA-256チェックサムを検証し、`~/.local/bin`に`ccam`を配置します。後から`ccwork`のようなショートカットコマンドを選ぶと、ccamは同じフォルダーにそれを作成します。そのフォルダーが`PATH`に含まれていない場合は、`~/.bashrc`または`~/.zshrc`に追加する行をスクリプトが表示します。

Windowsでは、PowerShellで次を実行します。

```powershell
irm https://ccam.dev/install.ps1 | iex
```

このスクリプトは`ccam_windows_<arch>.zip`をダウンロードしてSHA-256チェックサムを検証し、`ccam.exe`を`%USERPROFILE%\.local\bin`（Claude Codeのインストーラーが`claude.exe`を置くフォルダー）に配置します。そのフォルダーがまだユーザーの`PATH`になければ、スクリプトが追加し、`%USERPROFILE%\bin`のような項目はそのまま残します。その後、新しいターミナルを開いてください。

## 対応システム

- x86-64およびARM64のLinux: Ubuntu、Debian、Fedora、Arch、Alpine、その他のディストリビューション。バイナリは静的リンクで、依存関係はありません。
- Appleシリコン搭載およびIntel搭載のmacOS。
- x86-64およびARM64のWindows 10と11。PowerShell、Windows Terminal、Git Bashで使えます。
  - 共有フォルダーは開発者モードなしで動作します。
  - 共有ファイル（`settings.json`、`CLAUDE.md`、`history.jsonl`）には開発者モード（設定 > システム > 開発者向け）が必要です。Windowsでは、このモードか管理者でなければファイルへのリンクを作成できないためです。[Windowsでの共有](/features/share#on-windows)を参照してください。
  - 開発者モードがない場合、最初の`ccam add`はMCPサーバーと拡張機能だけを含む`default`グループを作成し、開発者モードを有効にした後で残りを追加する方法を表示します。
- WSLの中では、ccamはほかのLinuxシステムと同じように動作します。

## そのほかのインストール方法

**DebianおよびUbuntuのパッケージ。** [リリースページ](https://github.com/jhyunwoo/ccam/releases)から`ccam_<version>_linux_amd64.deb`をダウンロードし、インストールします。

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

このパッケージは`/usr/bin/ccam`をインストールします。

**ソースから。** Go 1.26以降を使います。

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

リリースタグ（`@latest`または`@vX.Y.Z`）をインストールした場合、ccamはそのバージョンを報告し、`ccam update`とダッシュボードのアップデート提案が使えます。`@main`や、タグの後にコミットがあるローカルのチェックアウトからのビルドは`dev`と報告し、`--version`を指定しない限り`ccam update`は拒否します。

**特定のバージョン。** `CCAM_VERSION`を設定します。

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## シェル補完

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

補完は`ccam`のすべてのコマンドとオプションに対応し、アカウント名も補完します。

## Claude Codeプラグイン

ccamプラグインを使うと、Claude Code自身がccamを使えるようになります。どのアカウントに余裕があるかを確認し、`ccam pick`で別々のアカウントに並列の`claude -p`ワーカーを起動し、スキルやMCPサーバーが見つからない理由を`ccam ext`で説明します。Claude Codeの中でインストールします。

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

次のものが追加されます。

- ログインやファイルの移動を行うコマンドを提案はするが、実行はしないスキル
- コマンド`/ccam:usage`、`/ccam:ext`、`/ccam:switch`
- 2つのフック（[詳細](/features/next#inside-claude-code)）。1つは、セッションが使用量の上限で止まるとターミナルに通知を表示します。LinuxとmacOSでは、[`auto_switch`](/features/next#auto-switch)を有効にしていれば、続いてccamが余裕のある別のアカウントにセッションを移します。もう1つは、`/ccam:switch [account]`で会話を別のアカウントに引き渡せるようにします。

フックは`ccam`を実行するので、ccamはClaude Codeから見える`PATH`上にある必要があります。

## インストールの確認

```sh
ccam version
ccam doctor
```

`ccam doctor`は、Claude Codeが見つかるか、すべてのアカウントのフォルダー、リンク、ログインに問題がないかを確認します。次は[クイックスタート](/docs/quick-start)に進んでください。

## アップデート

```sh
ccam update --check
ccam update
```

`ccam update --check`は、新しいリリースがあるかどうかを知らせるだけです。`ccam update`は、お使いのシステム用のリリースをダウンロードし、リリースの`checksums.txt`と照合してSHA-256チェックサムを検証してから、`ccam`ファイルを一度に置き換えます。チェックサムが一致しない場合は何も変更しません。`ccwork`のようなショートカットはccamへのリンクなので、すぐに新しいバージョンで動作します。

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ダッシュボード（ターミナルで引数なしで実行した`ccam`）も、開くたびに新しいリリースがないか確認し、見つかればアップデートを提案します。`U`を押して確認すれば更新されます。ほかのコマンドがアップデートを確認することはなく、ccamが確認なしに更新することもありません。古いリリースも含めて特定のリリースをインストールするには、`--version`を指定します。`ccam --dry-run update`は、何も変更せずに、何をダウンロードして何を置き換えるかを表示します。

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

インストールスクリプトをもう一度実行しても、`~/.local/bin`にあるccamが更新されます。リリースタグから`go install`でccamをビルドした場合は、`ccam update`を実行しても、同じ`go install`コマンドをもう一度実行してもかまいません。

**Windows。** `ccam update`はWindows用の`.zip`をダウンロードします。Windowsでは実行中のプログラムを置き換えられないため、古いバイナリは`ccam.exe.old`に名前を変更され、次のアップデートで削除されます。以前の`ccam.exe.old`からまだセッションが動いている場合、名前は`ccam.exe.old-<time>`になります。ハードリンクのショートカットはその直後に作り直されます。`install.ps1`をもう一度実行した場合も同じように動作し、古いバイナリからClaude Codeのセッションがまだ動いていても置き換えられます。ハードリンクのショートカットを作り直すと、`OK 2 shortcuts now run the new ccam`と表示します。

**DebianおよびUbuntuのパッケージ。** aptでインストールした`/usr/bin/ccam`は`ccam update`では置き換えず、代わりに次の手順を表示します。[リリースページ](https://github.com/jhyunwoo/ccam/releases)から新しい`ccam_<version>_linux_amd64.deb`をダウンロードし、古いパッケージの上からインストールします。

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## アンインストール

`~/.local/share/ccam`を削除すると、そこに保存されているログインと会話履歴も削除されます。[インポートした](/docs/commands)`~/.claude`フォルダーは、その場所に残ります。ただし、グループに入っている間は、共有しているファイルが`~/.local/share/ccam/groups/`にあります。何かを削除する前に、そのアカウントに対して`ccam remove <name>`（または`ccam group leave`）を実行してください。これでアカウントがグループから外れ、共有していたものすべてのコピーがフォルダーに戻ります。

まず、ccamがインストールされているうちに、アカウントのショートカットを削除します。`ccam list`のShortcut列にショートカットが表示されます。それぞれ`ccam shortcut <account> --remove`で削除するか、`~/.local/bin`からそれらのリンクを自分で削除してください。続いてコマンドを削除し、アカウントが不要になった場合はデータフォルダーも削除します。

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windowsでは、PowerShellで次を実行します。

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```
