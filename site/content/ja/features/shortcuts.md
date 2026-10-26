---
title: ショートカットコマンド
description: ccam shortcutで、アカウントにccworkのような専用のコマンドを付け、ccam shims installで~/.local/binのリンクを作り直します。
updated: 2026-10-26
order: 3
translation: machine
sourceHash: 066e7e622517
---

ショートカットは、1つのアカウントとしてClaude Codeを起動する、自分で名前を付けたコマンドです。`ccwork --resume`は`ccam work --resume`と同じ動作をします。実体は`~/.local/bin`にあるccamへのリンクなので、bash、zsh、スクリプトなど、`PATH`からプログラムを実行するあらゆる場所で動作します。ショートカットは任意です。よく起動するアカウントに作成してください。

## 使い方

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccamは変更のたびに結果を表示します。

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

1つのアカウントが持てるショートカットは最大1つで、`~/.config/ccam/config.toml`に`shortcut`として保存されます。`ccam add`のウィザードは`Create a shortcut command for <name>? [y/N]`と尋ね、続いて`Shortcut name [cc<name>]:`と尋ねます。ウィザード以外では、`ccam add`は`--shortcut`を付けた場合だけショートカットを作成します。`ccam list`は各アカウントのショートカットを表示し、ない場合は`-`と表示します。

ショートカットは、`ccam <name>`が名前より後の引数を扱うのと同じように、すべての引数をClaude Codeに渡します。`--dry-run`や`--help`も含まれます。起動をプレビューするには、`ccam --dry-run work`を使ってください。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--remove` | アカウントのショートカットとそのリンクを削除します | オフ |
| `-h`、`--help` | コマンドのヘルプ | |
| `--dry-run` | `ccam shortcut`は対応しています。名前を確認し、保存もリンクもせずに`Would set the shortcut of work to ccwork.`または`Would remove ccwork, the shortcut of work.`と表示します | オフ |

`ccam shims install`は引数を取らず、`--dry-run`を`` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``で拒否します。

名前の規則は、次の順序で確認されます。

- `a-z`、`0-9`、`.`、`-`、`_`を使った1～64文字で、先頭は英字か数字であること
- `ccam`、`claude`、`run`や`usage`などのccamのコマンドではないこと
- ほかのアカウントのショートカットではないこと
- `~/.local/bin`にある既存のファイルと同じ名前ではないこと。ただし、ccamが作成したリンク（このアカウント自身のショートカットなど）は除きます
- `PATH`上にあるほかの実行ファイルの名前ではないこと

`config.toml`も、読み込むときに同じ規則で確認されます。ショートカットを手で編集して、パターンに合わなくなった場合、`claude`か`ccam`になった場合、または2つのアカウントで使われるようになった場合は、直すまで、どのコマンドもそのアカウントの名前を示すメッセージとともに失敗します。

コマンドと`--remove`を同時に指定すると、`give a command or --remove, not both`で停止します。アカウントがすでに持っている名前を設定すると`work: ccwork (unchanged)`と表示し、リンクがなければ作り直します。

## しくみ

名前を設定すると`checkShortcut`が実行され、binフォルダーとccamのバイナリの実際のパスを指定して`shim.Check`を呼び出します。

<!-- motion:shortcuts-check -->

名前は確認を順に通っていきます。各確認の横には、その確認で拒否される名前を示しています。`Work`はパターンに合わず、`run`はccamのコマンドで、`ccside`はほかのアカウントのもので、`ccfile`は`~/.local/bin`にある既存のファイルで、`git`は`/usr/bin/git`を隠してしまいます。すべての確認に通った名前だけが保存され、リンクされます。

1. ccamは、名前を`^[a-z0-9][a-z0-9._-]{0,63}$`、`config.Reserved`、`cfg.FindShortcut`の順に照合します（`shim.Check`）。`~/.local/bin/<name>`に対して`os.Lstat`を呼び出し、`PATH`上のそれ以外の各フォルダーに対しては`os.Stat`を呼び出して、ccamへのリンクではない実行ファイルがあれば拒否します。
2. ccamは`fsutil.AtomicWrite`（一時ファイル、`fsync`、rename）で`config.toml`を保存します。
3. ccamは古いショートカットのリンクを削除します。ただし、それがccamへのリンクである場合だけです（`shim.Remove`）。
4. ccamは、設定されているすべてのショートカットを名前の順にたどり、足りないものについて`os.Symlink`でccamのバイナリへのリンクを作成します（`shim.Install`）。バイナリのパスは`os.Executable`と`filepath.EvalSymlinks`で求めます。1つのショートカットの場所にファイルがあっても、ほかのショートカットは止まりません。エラーにはそれぞれの名前が示されます。

ccamは、作成したすべてのリンクと、各リンクが実行するバイナリを`~/.local/state/ccam/shortcuts.json`に記録します。`os.Readlink`がccamのパスを返す場合、またはリンクとccamのパスが同じファイルに解決される場合、そのリンクはccam自身のものとみなされます。記録されたリンクのうち、その後バイナリが移動したり置き換えられたりしたもの（別のフォルダーに移したccamのバイナリ、再インストール）は古くなったリンクで、ccamがそれをこのccamのバイナリに向け直します（`shim.Install`）。インストールされている別のccamのコピーへの記録されたリンクは、そのまま残します。削除と上書きの確認はすべてこれらの判定を通るので、自分で置いたccamのコピーや自分で作ったリンクのように、ccamが作成していないファイルに触れることはありません。

ショートカットを実行したとき、ccamが必要とするのは入力された名前だけです。

<!-- motion:shortcuts-link -->

シェルが`PATH`上で`ccwork`を見つけ、カーネルがリンクをたどってccamのバイナリを起動します。このとき`argv[0]`は`ccwork`のままです。`app.Main`では、`ccam`以外の名前の場合、`config.toml`を一度読み込んで`shim.Lookup`を呼び出し、`shortcut`が`filepath.Base(argv[0])`と一致するアカウントを見つけます。続いてccamは、読み込み済みの設定を使って[`ccam <name>`](/features/launch)と同じ手順をたどり、`argv[1:]`をそのまま`claude`に渡します（`launcher.LaunchWith`）。

シェルのエイリアスやラッパースクリプトではなくシンボリックリンクを使う理由は、すべてのシェルとスクリプトで動作し、余分なシェルの起動（約1～2 msかかります）を避けられ、起動時に読むファイルも1つのままだからです。`scripts/perf-check.sh`は、`ccam <name>`と同じ10 ms（p95）の上限で、`hyperfine`を使ってショートカットを計測します。

### Windowsの場合

ショートカットは`%USERPROFILE%\.local\bin`にある`ccwork.exe`で、ccamは検索の前に`.exe`を（大文字小文字を区別せずに）取り除きます（`shim.Name`）。そのため、`ccwork`、`ccwork.exe`、`CCWORK.EXE`はどれも`work`を起動します。ccamはまずシンボリックリンクを作成します。Windowsでこれができるのは、開発者モードのときか管理者の場合です。その権限がない場合はハードリンクを作成します。これは同じ`ccam.exe`ファイルの2つ目の名前で、権限は必要ありません。ハードリンクは、`ccam.exe`と同じファイルであればccam自身のものと見なされます。

ハードリンクは常に作成元のファイルを実行するので、`ccam update`は`ccam.exe`の場所に新しいファイルを置くと、`shim.Relink`でハードリンクのショートカットをすべて作り直し、`✓ shortcut ccwork renewed`と表示します。`install.ps1`をもう一度実行した場合も`ccam shims install`を通じて同じことが行われ（`OK 2 shortcuts now run the new ccam`）、古いバージョンのccamのコピーになっているショートカットについては、`ccam shims install`と`ccam doctor --fix`も同じことをします。名前の確認では、ほかの`PATH`フォルダーで、`PATHEXT`の各拡張子（`.exe`、`.cmd`など）を付けた`<name>`を探します。

ショートカットはシェル（PowerShell、cmd.exe、Git Bash）で入力してください。エクスプローラー、［ファイル名を指定して実行］ダイアログ、`Start-Process`から開いたシンボリックリンクは、ただの`ccam`として起動するので、アカウントではなくダッシュボードが開きます。

## 安全性と特殊なケース

- **どのアカウントのものでもないリンク。** `config.toml`の`shortcut`を手動で編集した後などに、`~/.local/bin`にあるccamへのリンクの名前をどのアカウントも使っていない場合、ccamはデフォルトのアカウントにフォールバックしません。`` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``で停止します。`ccam shortcut`で名前を設定し直すか、リンクを削除してください。
- **既存のファイル。** ccamは自分で作成していないファイルを置き換えることはありません。`ccam shims install`は`<path> already exists and is not a ccam link, so it was left alone`と報告し、ほかのリンクは作成します。その場合、`ccam add`と`ccam import`は`shortcut: ccwork (not created; see above)`と表示します。
- **ccamのバイナリを移動した場合。** ccamのバイナリを移動すると、ccamが作成したリンクは存在しないパスを指すことになります。`ccam shims install`（または`ccam doctor --fix`）が、それらを実行に使ったccamのバイナリに向け直します。
- **名前の変更と削除。** `ccam rename`ではショートカットの名前は変わらず、リンクも動作し続けます。リンクはccamを指しており、ccamが`config.toml`で名前を調べるためです。`ccam remove`はアカウントのリンクを削除します。`ccam shortcut <account> <new>`は、新しいリンクを作成する前に古いリンクを削除します。
- **PATH。** ccamがリンクを作成したときに`~/.local/bin`が`PATH`にない場合、`Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`に続けて、追加すべき`export PATH=...`の行を表示します。Windowsでは、代わりにユーザーの`PATH`にフォルダーを追加するPowerShellの行を表示します。
- **あとから確認する。** `ccam doctor`は`missing shortcuts in /home/me/.local/bin: ccwork`と、古いバージョンのccamを実行するショートカットを報告し、フォルダーが`PATH`にない場合は警告します。`ccam doctor --fix`は足りないリンクを作成して古くなったリンクを作り直し、自分で直したものについてだけ`→ fixed`と表示します。まだインストールされている別のccamを実行するショートカットは動作するので、`ccam doctor`は情報としてだけ示します。`ccam doctor --fix`はそれを実行に使ったccamのバイナリに向け直し、`ccam add`と`ccam import`は``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``と表示します。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [アカウントの管理](/features/accounts)
- [コマンドリファレンス](/docs/commands#shortcuts)
- [ccamのインストール](/docs/install)
