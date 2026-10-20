---
title: フォルダーごとのアカウント
description: ccam folderを使うと、アカウント名なしでClaude Codeを起動したとき、フォルダーとその下すべてで1つのアカウントが使われます。たとえば会社のコードではworkです。
updated: 2026-10-20
order: 11
translation: machine
sourceHash: "46d1410a3d81"
---

`ccam folder work ~/code/company`を実行すると、`~/code/company`とその下のどのフォルダーでも、アカウント名を指定しない起動はすべて`work`としてClaude Codeを起動します。その外では、これまでどおりデフォルトのアカウントが起動します。`ccam personal`や`ccpersonal`のようなショートカットでアカウントを指定した場合は、常にそのアカウントが起動します。

## 使い方

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

次の起動では、フォルダーのアカウントが使われます。

- `ccam -p "…"`、および最初の引数がClaude Codeのオプションであるそのほかの起動
- 名前なしの`ccam run`と、`ccam --dry-run`
- `ccam --fallback -p "…"`の最初の試行
- ダッシュボード（`ccam`単体）。フォルダーのアカウントにカーソルを置いた状態で開きます

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `<account>` | 使うアカウントの名前か番号 | |
| `[folder]` | フォルダー（相対パスまたは絶対パス） | 現在のフォルダー |
| `--remove` | そのフォルダーちょうどのエントリーを削除します | オフ |
| `--dry-run` | 対応しています。変更を保存せずに表示します | オフ |

## しくみ

各フォルダーは、`~/.config/ccam/config.toml`の1つの`[[folder]]`エントリーです。このファイルは、どの起動でもいずれにせよ読み込まれます。

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder`は、フォルダーをシンボリックリンクを解決した絶対パスにして保存します。同じフォルダーをもう一度設定すると、そのアカウントが置き換わります。
2. 名前なしの起動は`launcher.HereProfile`を呼び出します。エントリーがなければ、すぐにデフォルトのアカウントを使います。
3. エントリーがある場合は、作業ディレクトリを各`path`と文字列として比較します。一致するのはそのフォルダー自体かその下のフォルダーだけで、同じ文字で始まるだけのフォルダーは一致しません（`/code/company`は`/code/companyx`を含みません）。最も深いエントリーが優先されるので、`~/code`を`personal`にしつつ、`~/code/company`を`work`にできます。
4. 作業ディレクトリがそのままではどれにも一致しない場合に限り、ccamはそのシンボリックリンクを解決して、もう一度比較します。
5. Windowsでは、Windowsと同じように、比較で大文字と小文字を区別しません。

`config.toml`以外のファイルは読まないので、フォルダー内での起動にかかる時間は測定できるほど増えません（エントリーの有無にかかわらず4.0 ms、`hyperfine`、300回）。

## 安全性と特殊なケース

- `ccam rename`はフォルダーのエントリーを新しい名前に引き継ぎ、`ccam remove`は削除したアカウントのエントリーを削除します。
- アカウントがもう存在しないエントリーがあると、そのフォルダーでの名前なしの起動は停止します: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`。`ccam doctor`は、そのようなエントリーと、フォルダーがなくなったエントリーを報告します。
- 上のフォルダーからアカウントを受け継いでいるフォルダーの中で`ccam folder --remove`を実行すると、何も削除せずに、その上のフォルダーを示します。
- `ccam next`はローテーションか使用量でアカウントを選び、フォルダーは考慮しません。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [ショートカットコマンド](/features/shortcuts)
- [アカウントの管理](/features/accounts)
- [コマンドリファレンス](/docs/commands)
