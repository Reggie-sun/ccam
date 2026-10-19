---
title: ccamのインストール
description: LinuxまたはmacOSにccamを1つのコマンドで、.debパッケージから、またはソースからインストールします。先にClaude Codeのインストールが必要です。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

LinuxまたはmacOSのターミナルで、次のコマンドを実行します。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

このスクリプトは、お使いのシステム用のリリースをダウンロードしてSHA-256チェックサムを検証し、`~/.local/bin`に`ccam`を配置します。アカウントを追加すると、`ccam add`が同じフォルダーに`claude1`、`claude2`…のショートカットを作成します。そのフォルダーが`PATH`に含まれていない場合は、`~/.bashrc`または`~/.zshrc`に追加する行をスクリプトが表示します。

ccamがClaude Codeを起動するので、先に[Claude Code](https://docs.claude.com/en/docs/claude-code/overview)をインストールし、`claude --version`が動作することを確認してください。

## 対応システム

- x86-64およびARM64のLinux: Ubuntu、Debian、Fedora、Arch、Alpine、その他のディストリビューション。バイナリは静的リンクで、依存関係はありません。
- Appleシリコン搭載およびIntel搭載のmacOS。
- Windowsには直接対応していません。WSLの中では、ccamはほかのLinuxシステムと同じように動作します。

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

**特定のバージョン。** `CCAM_VERSION`を設定します。

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## シェル補完

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

補完は`ccam`のすべてのコマンドとオプションに対応し、アカウント名も補完します。

## インストールの確認

```sh
ccam version
ccam doctor
```

`ccam doctor`は、Claude Codeが見つかるか、すべてのアカウントのフォルダー、リンク、ログインに問題がないかを確認します。次は[クイックスタート](/docs/quick-start)に進んでください。

## アンインストール

コマンドとアカウントのショートカットを削除します。アカウントが不要になった場合は、データフォルダーも削除します。

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam`を削除すると、そこに保存されているログインと会話履歴も削除されます。[インポートした](/docs/commands)`~/.claude`フォルダーは、その場所に残ります。
