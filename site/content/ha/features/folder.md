---
title: Asusai bisa folda
description: ccam folder yana sa folda da duk abin da ke ƙarƙashinta su yi amfani da asusu ɗaya idan ka fara Claude Code babu sunan asusu, kamar work don lambar kamfani.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: 25f82da64580
---

Bayan `ccam folder work ~/code/company`, duk wani farawa ba tare da sunan asusu ba a `~/code/company` ko foldojin da ke cikinta yana fara Claude Code a matsayin `work`. A wajenta, asusun tsoho yana farawa kamar da. Ambaton asusu, da `ccam personal` ko gajeren umarni kamar `ccpersonal`, koyaushe yana fara wannan asusun.

## Amfani

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

Waɗannan farawa suna amfani da asusun foldar:

- `ccam -p "…"` da duk wani farawa wanda argument ɗinsa na farko zaɓi ne na Claude Code
- `ccam run` ba tare da suna ba, da `ccam --dry-run`
- `ccam --fallback -p "…"` don yunƙurinsa na farko
- dashboard, `ccam` shi kaɗai, wanda ke farawa da alamar zaɓinsa a kan asusun foldar, har ma idan ka isa foldar ta mahaɗi na alama (symbolic link)

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `<account>` | Suna ko lambar asusun da za a yi amfani da shi | |
| `[folder]` | Foldar, ta dangi (relative) ko cikakkiya (absolute) | foldar yanzu |
| `--remove` | Cire shigarwar wannan foldar daidai | a kashe |
| `--dry-run` | Ana tallafa masa: buga canjin ba tare da ajiye shi ba | a kashe |

## Yadda yake aiki

Kowace folda shigarwa ɗaya ce ta `[[folder]]` a `~/.config/ccam/config.toml`, fayil ɗin da kowane farawa ke karantawa ko ta yaya:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` yana mayar da foldar zuwa cikakkiyar hanya (absolute path) tare da warware symbolic links, kuma yana ajiye ta. Sake saita folda yana maye gurbin asusunta.
2. Farawa ba tare da suna ba yana kiran `launcher.HereProfile`. Idan babu shigarwa, yana ɗaukar asusun tsoho nan take.
3. In ba haka ba yana kwatanta foldar aiki da kowane `path` a matsayin strings: foldar kanta ko wadda ke ƙarƙashinta, ba taɓa foldar da kawai take farawa da haruffa iri ɗaya ba (`/code/company` ba ya rufe `/code/companyx`). Shigarwa mafi zurfi ce ke cin nasara, don haka `~/code` zai iya zama `personal` yayin da `~/code/company` yake `work`.
4. Sai idan foldar aiki yadda aka bayar da ita ba ta dace da komai ba ne ccam yake warware symbolic links ɗinta kuma ya sake kwatantawa.
5. A Windows da macOS kwatantawar ba ta kula da manyan haruffa da ƙanana, kamar yadda tsarukan fayilolinsu suka saba yi.

Ba a karanta wani fayil banda `config.toml`, don haka farawa a cikin folda ba ya ƙara lokacin da za a iya auna (4.0 ms da shigarwa ko babu, `hyperfine`, gudu 300).

## Tsaro da yanayi na musamman

- `ccam rename` yana ɗaukar shigarwar foldoji zuwa sabon sunan, kuma `ccam remove` yana goge shigarwar asusun da aka cire.
- Shigarwar da asusunta ba ya nan yana dakatar da farawa ba tare da suna ba a wannan foldar: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` yana ba da rahoton irin waɗannan shigarwa, da shigarwar da foldarsu ta ɓace.
- Idan ka gudanar da `ccam folder --remove` a cikin ƙaramar folda da ke gadon asusunta daga foldar da ke samanta, ccam yana ambaton foldar ta sama kuma ba ya cire komai.
- `ccam next` yana zaɓa ta juyawa ko amfani kuma ba ya duba foldoji.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Gajerun umarni](/features/shortcuts)
- [Sarrafa asusai](/features/accounts)
- [Bayanin umarnai](/docs/commands)
