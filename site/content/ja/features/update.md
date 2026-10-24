---
title: ccamのアップデート
description: ccam updateは、SHA-256チェックサムを検証してccamをGitHubの最新リリースに置き換えます。--checkは確認のみ、--versionは指定したリリースタグをインストールします。
updated: 2026-10-23
order: 10
translation: machine
sourceHash: "000d37de5f29"
---

`ccam update`は、GitHubからccamのリリースをダウンロードし、リリースの`checksums.txt`で検証してから、実行中のバイナリを1回のrenameで置き換えます。インストールスクリプトやアーカイブからccamをインストールした場合に使います。[ダッシュボード](/features/dashboard)も、開くたびにGitHubに最新のリリースを問い合わせ、`U`でアップデートを提案します。このキーは確認を取ってから同じ手順を実行します。ほかのコマンドはGitHubに接続しないので、オフラインでも同じように動作します。

## 使い方

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

新しいリリースがない場合、ccamは`ccam 0.4.0 is up to date.`と表示し、何も変更しません。`--version`に実行中のバージョンを指定した場合は、`ccam 0.4.0 is already installed.`と表示します。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--check` | 新しいリリースがあるかどうかを報告するだけです | オフ |
| `--version <tag>` | `v0.3.0`のようなリリースタグをインストールします。現在のものより新しくても古くてもかまいません | 最新 |
| `--dry-run` | 対応しています。最新のタグを調べ（`--version`で指定した場合を除く）、バイナリに対して同じ確認を行い、ダウンロードせずに、ダウンロード、チェックサム、置き換えのパスを表示します | オフ |

`--check`と`--version`は同時に使えません。タグは`vMAJOR.MINOR.PATCH`の形式で、プレリリースやビルドの接尾辞を付けてもかまいません。先頭の`v`は省略できます。それ以外は失敗します。たとえば`--version takes a release tag such as v0.3.0, not "1.2"`のようになります。

環境変数`CCAM_RELEASES_URL`は、ミラーやテストのために`https://github.com/jhyunwoo/ccam/releases`を置き換えます。構成は同じでなければなりません。`<url>/latest`が`<url>/tag/<tag>`にリダイレクトし、アセットは`<url>/download/<tag>/<asset>`にある必要があります。

## しくみ

コードは`internal/update`にあり、これをインポートするのは`ccam update`だけです。起動経路から呼び出されることはありません。

1. `cli.updateCmd`が`os.Executable`で実行中のバイナリを見つけ、`filepath.EvalSymlinks`でシンボリックリンクを解決します。置き換えられるのは、この実際のパスです。
2. `--version`がない場合、`Updater.Latest`はリダイレクトをたどらずに`<releases>/latest`にリクエストを送り、`Location`ヘッダー（`…/tag/v0.5.0`）からタグを読み取ります。GitHub APIではなくリリースのページを使うので、APIの1時間あたり60リクエストの制限は適用されません。この問い合わせの制限時間は10秒です。
3. `parseVersion`がタグを厳密なパターンで確認します。タグはダウンロードURLや表示されるシェルコマンドに入るため、`/`、`%`、空白、シェルの特殊文字を含むものは通しません。`compareVersions`が実行中のバージョンと比較します。プレリリースは、そのリリースより前として扱われます。
4. `managed()`と`writable()`が、そもそもccamがバイナリを置き換えてよいかを判断します（後述）。
5. `Updater.Fetch`が`ccam_<os>_<arch>.tar.gz`（Windowsでは`ccam_windows_<arch>.zip`）と`checksums.txt`を同時にダウンロードします。制限時間は両方合わせて2分です。各ダウンロードの上限は256 MBです。
6. ほかの処理を行う前に、アーカイブのSHA-256が`checksums.txt`の該当行と一致する必要があります。行がない場合やチェックサムが異なる場合は、`nothing was changed`で更新を停止します。
7. `extract()`がtarを読み、`ccam`という名前の最初の通常ファイルを取り出します（`extractZip`はzipから`ccam.exe`を取り出します）。ディレクトリ、リンク、その他の項目は無視し、バイナリにも256 MBの上限があります。
8. `Updater.Install`が`fsutil.AtomicWrite`で新しいバイナリを書き込みます。バイナリの隣に、古いファイルと同じ権限ビットで一時ファイル`.ccam.tmp*`を作り、ディスクにフラッシュしてから、古いファイルの上にrenameします。

タグの問い合わせからrenameまでの全体の流れは次のとおりです。

<!-- motion:update-pipeline -->

アーカイブと`checksums.txt`は並行してダウンロードされ、チェックサムが一致するまで何も書き込まれません。renameはアトミックなので、クラッシュやディスクの容量不足が起きても、残るのは古いバイナリか新しいバイナリのどちらかで、中途半端なものが残ることはありません。`~/.local/bin/ccwork`のようなショートカットコマンドはccamのバイナリへのシンボリックリンクなので、何も変更しなくても新しいバージョンが実行されます。

Windowsでは、実行中のプログラムのファイルは置き換えも削除もできず、名前の変更だけができます。`Updater.Install`は新しいバイナリを`ccam.exe.new`に書き込み、実行中の`ccam.exe`の名前を`ccam.exe.old`に変え、新しいファイルの名前を`ccam.exe`に変えます。最後の手順が失敗した場合は、古いバイナリの名前を元に戻します。次のアップデートは、削除できる`ccam.exe.old`ファイルを削除します。以前の`ccam.exe.old`からまだセッションが動いている場合は、実行中のバイナリの名前が代わりに`ccam.exe.old-<time>`になります。シンボリックリンクではなくハードリンクのショートカット（[ショートカットコマンド](/features/shortcuts)）はその直後に作り直されるので、それらも新しいバージョンを実行します。これは、置き換えたばかりのバイナリのハードリンクだけでなく、ccamが作った古いccamのハードリンクすべてが対象です。

### パッケージでインストールした場合

`.deb`からインストールしたccamはaptの管理下にあり、aptを通さずに置き換えるとパッケージのデータベースが不正確になります。ダウンロードする前に、ccamはバイナリの所有者を確認します。

<!-- motion:update-managed -->

1. `/usr/bin`、`/bin`、`/usr/sbin`、`/sbin`にあるバイナリは、パッケージで管理されているとみなします。
2. それ以外の場所でも、`dpkg`がインストールされていて`dpkg -S <path>`が成功すれば、パッケージで管理されているとみなします。
3. それ以外の場合、ccamはバイナリの隣に一時ファイル`.ccam-update-*`を作成して削除し、そこに書き込めることを確認します。

パッケージで管理されているバイナリの場合、ccamは停止し、aptでの手順を表示します。

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux以外のシステムでは、インストールスクリプトの行だけが表示されます。Windowsでは`irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`です。フォルダーに書き込めない場合、ccamはインストールスクリプトをもう一度実行するよう提案します。インストールスクリプトはccamを`~/.local/bin`にインストールします。

## 安全性と特殊なケース

- 開発版のビルドはバージョンを`dev`と報告し、比較するリリースがないので、`ccam update`と`ccam update --check`は拒否します。`ccam update --version vX.Y.Z`なら、その上にリリースをインストールできます。`go install github.com/jhyunwoo/ccam/cmd/ccam@vX.Y.Z`（または`@latest`）で作ったビルドはそのリリースを報告するので、`ccam update`とダッシュボードのアップデートの案内が使えます。`@main`や、タグより先のローカルビルドは`dev`のままです。
- ダウンロードの失敗、使っているプラットフォーム用のアセットがない場合、チェックサムの不一致では、インストール済みのバイナリはそのまま残ります。
- 存在しないリリースは`there is no release <tag>`で失敗します。GitHubがHTTP 403や429で拒否した場合、ccamはレート制限の可能性が高いことを伝え、しばらくしてから再試行するよう求めます。
- `--check`は読み取るだけです。バイナリがパッケージで管理されているか、書き込めるかは確認しません。
- 以前のバージョンに戻すには、`--version`に古いタグを指定します。

## 関連項目

- [インストール](/docs/install)
- [コマンド](/docs/commands)
- [よくある質問](/docs/faq)
