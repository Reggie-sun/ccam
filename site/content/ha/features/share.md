---
title: Rabawa tsakanin asusai
description: Yadda ccam share ke haɗa saituna, CLAUDE.md, extensions da tarihi da foldar rukuni, tare da shiri, kwafin ajiya da zaɓin rikici, da yadda off ke mayar da kwafi.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Kowane asusu yana da foldarsa, amma yawancin abin da ke cikinta ba ya buƙatar zama daban. `ccam share` yana kunna ko kashe rabawa ga asusu ɗaya da abu ɗaya a lokaci ɗaya: ana matsar da fayil ko foldar asusun cikin foldar rukuni kuma ana maye gurbinsa da mahaɗi (symbolic link), don haka kowane asusu a rukunin yana karantawa da rubutawa a kwafi ɗaya. Yi amfani da shi bayan `ccam add` don canza abin da asusu ke rabawa, ko don sanya wasu asusai a rukunin nasu.

## Amfani

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Ba tare da zaɓuɓɓukan abubuwa ba, `ccam share` yana buga yanayin yanzu:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Da zaɓuɓɓuka, ccam yana buga shiri, yana tambaya game da kowane fayil da ya bambanta da kwafin rukunin, kuma yana tambaya `Proceed? [y/N]` kafin ya canza komai. Sabon asusu yana raba komai sai tarihin tattaunawa; [Raba sabobin MCP, saituna da tarihi](/docs/sharing) yana bayyana kowane abu daga ɓangaren mai amfani.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Sabobin MCP | ba a canza ba |
| `--settings <on\|off\|group>` | `settings.json` | ba a canza ba |
| `--memory <on\|off\|group>` | `CLAUDE.md` na gaba ɗaya | ba a canza ba |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | ba a canza ba |
| `--history <on\|off\|group>` | `projects/` da `history.jsonl` | ba a canza ba |
| `--fresh` | idan ana kashe rabawa, fara babu komai maimakon kwafar abin da ke cikin rukunin | a kashe |
| `--prefer <group\|mine\|merge>` | yadda za a yanke hukuncin rikici | tambaya; da `--yes`, tsohon zaɓin kowane abu |
| `-y`, `--yes` | gudana ba tare da tambaya ba; rikici yana bin `--prefer` | a kashe |
| `-h`, `--help` | taimako ga `share` | |
| `--dry-run` (na gaba ɗaya) | ana karɓar sa: yana tambaya game da rikici, yana buga shirin, ba ya canza komai | a kashe |

Ƙimar `on` tana nufin rukunin `default`, `off` tana nufin na gida (ba a raba ba), kuma kowace ƙima dabam sunan rukuni ce. Sunayen rukunoni suna zama sunayen foldoji ƙarƙashin `~/.local/share/ccam/groups/`, don haka `ccam share side --history ../x` yana gazawa da `invalid group name: "../x"`. Ana ba da asusun ta suna ko ta lamba, kuma abubuwan da ka ambata ne kawai ke canzawa.

## Yadda yake aiki

Fayilolin da abin ya shafa:

- `~/.config/ccam/config.toml`: teburin `[profile.share]` na kowane asusu, ƙima ɗaya ga kowane abu.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` na asusun.
- `~/.local/share/ccam/groups/<group>/`: kwafin da aka raba.

1. `internal/cli/share.go` yana karanta saitin, yana aiwatar da kowace ƙimar `--<item>` a kan kwafin saitunan asusun da `config.Share.Set`, kuma yana ƙi idan asusun yana aiki.
2. `share.Build` yana kwatanta saitunan yanzu da waɗanda ake so kuma yana lissafa matakan. Yana karanta faifai ne kawai: `os.Lstat` da `os.Readlink` ga kowace hanya a `share.Layout`. Shirin jerin ayyuka ne: matsawa cikin gurbin rukuni mara komai (adopt), matsar da fayilolin da rukunin ba shi da su (merge dir), rubuta rikicin da aka warware cikin rukunin, ajiya, cire mahaɗi, haɗawa, kwafa baya, sake rubuta hanyoyin plugin da haɗa rajistar plugin.
3. Ana yanke hukuncin rikici, ta tambaya ko ta `--prefer`.
4. `Plan.Describe` yana buga shirin, kuma ccam yana tambaya.
5. Ana rubuta saitunan da ake so a `config.toml` kafin komai ya motsa. Idan `Plan.Apply` ya tsaya a tsakiya, `ccam doctor --fix` da sake gudanarwa suna gina shirinsu zuwa ga abin da aka ajiye kuma suna kammala canjin.
6. `Plan.Apply` yana gudanar da matakan a tsari. Kowane mataki yana matsar da bayanai cikin rukunin ko cikin foldar ajiya; ba a goge komai.

`share.Layout` jerin izini ne na hanyoyin da kowane abu ya ƙunsa. Duk abin da ba ya cikinsa ba zai taɓa kai wa rukuni ba. `share.Never` yana ambaton `.credentials.json` da `.claude.json`, kuma `Build` yana tsayawa da kuskure na ciki idan ɗaya daga cikinsu ya taɓa bayyana a tsarin. Abin `mcp` ba shi da hanyoyi: jerin sabobin yana cikin `.claude.json`, don haka [rabawar MCP](/features/mcp) tana aiki ta `mcp.json` na rukunin maimakon haka, kuma `--mcp` yana canza saitin ne kawai.

Ga abin da ke faruwa ga folda ɗaya, `projects/`, idan `side` ya kunna tarihi kuma rukunin ya riga yana da abubuwa:

<!-- motion:share-on -->

Idan rukunin bai da komai a wannan hanyar tukuna, asusun shi ne na farko da ya raba ta kuma ana matsar da foldarsa can kamar yadda take, ba tare da ajiya ba. In ba haka ba `mergeDir` yana matsar da kowane fayil da rukunin ba shi da shi cikin rukunin, kuma fayilolin da ke ɓangarorin biyu suna zama a foldar asusun, don haka ana riƙe sigar rukunin. Sannan foldar asusun tana zuwa `.ccam-backup/<time>/` a cikin foldar asusun, kuma `os.Symlink` yana saka mahaɗi a wurinta. Ana sanya wa fayilolin tattaunawa suna bisa session ID ɗinsu, don haka tarihi daga asusai da dama yana haɗuwa ba tare da karo ba.

Matsawa tana amfani da `rename(2)` kuma tana komawa ga kwafa da gogewa tsakanin tsarin fayiloli (file systems). Mahaɗin dangi a cikin bishiyar da aka matsar wanda ke nuna wajenta ana sake rubuta shi a matsayin cikakken mahaɗi (`fsutil.PortableLink`), don ya ci gaba da aiki daga sabon wurin. Rajistocin `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) suna ajiye cikakkun hanyoyi ta foldar asusun; ccam yana nuna su zuwa foldar rukunin kuma, idan rukunin ya riga yana da rajista, yana haɗa shigarwar asusun cikinta.

Kashe abu yana juya mahaɗin:

<!-- motion:share-off -->

ccam yana cire mahaɗin ne kawai idan yana nuna cikin foldar groups, sannan yana kwafar abin da ke cikin rukunin a yanzu zuwa asusun da `CopyTree`. Da `--fresh` ana tsallake kwafin kuma asusun yana farawa babu komai a wannan hanyar. Rukunin yana riƙe abin da ke cikinsa don sauran asusai. Matsar da abu daga rukuni ɗaya zuwa wani yana cire mahaɗinsa kuma yana haɗa shi da sabon rukunin, yana barin tsohon rukunin yadda yake; shirin yana faɗin haka a layin `Note:`.

Fayil ɗaya da ke ɓangarorin biyu mai abubuwa daban rikici ne:

<!-- motion:share-conflict -->

`share.NewConflict` yana lissafa zaɓuɓɓukan da suka dace da fayil ɗin. `history.jsonl` yana ba da `merge by line`, `keep group` da `use mine`, tare da haɗawa a farko. Sauran fayilolin `.json` suna ba da `keep group`, `use mine` da `merge JSON`. Duk wani abu dabam yana ba da `keep group` da `use mine`. Zaɓi na farko shi ne tsoho:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) yana haɗa abubuwa (objects) a maimaice, yana haɗa arrays ba tare da maimaitawa ba, kuma yana riƙe ƙimar rukunin idan ɓangarorin biyu sun saita maɓalli ɗaya daban-daban, don haka sauran asusai ba sa ganin canji a ɗabi'a. `merge by line` yana riƙe layukan rukunin, sannan yana ƙara layukanka da ba ya da su. Kowane zaɓi banda `keep group` yana fara ajiye kwafin rukunin a `.ccam-backup/<time>/group/` a cikin foldar asusun, sannan yana rubuta sakamakon da `fsutil.AtomicWrite`. Fayil ɗinka koyaushe yana zuwa ajiya.

`--prefer group`, `mine` ko `merge` yana yanke hukuncin kowane rikici ba tare da tambaya ba; `merge` yana nufin `merge JSON` ko `merge by line`, duk wanda fayil ɗin ke bayarwa, kuma fayil da ba ya bayar da ko ɗaya yana riƙe tsoho. `--yes` shi kaɗai yana ɗaukar tsohon zaɓin kowane fayil.

A kowane farawa, `share.Check` yana kwatanta kowace hanya a tsarin da saitin ta amfani da `lstat` da `readlink` kawai, wanda ke sa ya yi sauƙi isa ga hanyar farawa. Ana ƙirƙirar mahaɗin da ya ɓace nan take, domin babu abin da za a rasa. Fayil na gaske inda ya kamata mahaɗi ya kasance, mahaɗi zuwa wurin da bai dace ba, ko tsohon mahaɗi yana dakatar da farawar:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` yana gudanar da `share.Build` da saitunan da aka ajiye a ɓangarorin biyu, wanda ke samar da kowane mataki daga abin da ke kan faifai, yana riƙe abin da ke cikin rukunin idan an yi rikici kuma yana ajiye kwafin duk abin da ya maye gurbinsa.

## Tsaro da yanayi na musamman

- **Asusan da ke aiki.** ccam yana ƙin canza rabawa yayin da Claude Code ke gudana a matsayin asusun (`account work is running; quit it first, then change its sharing`), domin zaman yana ci gaba da rubutu a cikin fayilolin da ake matsarwa kuma layukan da yake ƙarawa a wannan lokacin za su ɓace. `--dry-run` har yanzu yana buga shirin.
- **Kwafin ajiya.** Kowane shiri yana samun foldar `.ccam-backup/<time>/` nasa a cikin foldar asusun, tare da ƙarin lamba idan shirye-shirye biyu suka fara a daƙiƙa ɗaya. ccam ba ya taɓa goge su.
- **Gazawa a tsakiya.** Idan mataki ya gaza, ccam yana buga kuskuren da inda ainihin fayilolin suke, kuma yana gaya maka ka kammala da `ccam doctor --fix`. An riga an ajiye abin da ake so, don haka gyaran yana ci gaba zuwa gare shi.
- **Mahaɗanka na kanka.** Idan hanya ta riga ta zama mahaɗin da ka yi, misali zuwa ma'ajiyar dotfiles, kuma rukunin babu komai a ciki, ccam yana kwafar abin da mahaɗin ke nunawa cikin rukunin kuma yana matsar da mahaɗin kansa zuwa ajiya. Ba a taɓa abin da ake nunawa ba.
- **Fayilolin shiga.** Ba a taɓa haɗa, matsar ko kwafar `.credentials.json` da `.claude.json` cikin rukuni ba.
- **Asusai biyu a folda ɗaya.** Asusan da ke raba tarihi suna rubutu a fayiloli iri ɗaya. Idan asusu na biyu na rukunin tarihi ɗaya ya fara a foldar da ɗaya ke aiki tuni, ccam yana gargaɗi kuma yana tambaya kafin ya ci gaba.

## Masu alaƙa

- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Sabobin MCP da aka raba](/features/mcp)
- [Kwafar saitin da ke akwai](/features/migrate)
- [Dukkan umarnai](/docs/commands)
