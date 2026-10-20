---
title: Rabawa da rukunoni
description: Rukunin ccam yana raba abubuwan da aka zaɓa tsakanin asusansa ta folda ɗaya; ccam group yana matsar da bayanai tare da shiri, kwafin ajiya da zaɓin rikici.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 0c0338362ba2
---

Kowane asusu yana da foldarsa, amma yawancin abin da ke cikinta ba ya buƙatar zama daban. **Rukuni** yana faɗin abin da ake rabawa da kuma su wa ke rabawa: kowane asusu a rukunin yana haɗa waɗannan abubuwa da folda ɗaya ta mahaɗi, `~/.local/share/ccam/groups/<group>/`, don haka dukansu suna karantawa da rubutawa a kwafi ɗaya. Ba a taɓa raba tokens na shiga ba.

Asusu zai iya kasancewa a rukunoni da dama, amma kowane abu yana fitowa daga rukuni ɗaya kawai. Saitin da aka saba shi ne rukuni ɗaya don saituna da extensions, da na biyu don tarihin tattaunawa tsakanin asusai biyu:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Asusun farko da ka ƙara yana ƙirƙirar `default`, kuma asusan da ke biye suna shiga cikinsa sai dai idan ka faɗi wani abu dabam.

## Amfani

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Kowane umarni da ke matsar da bayanai yana buga shiri, yana tambaya game da kowane fayil da ya bambanta da kwafin rukunin, kuma yana tambaya `Proceed? [y/N]` kafin ya canza komai:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Canje-canje iri ɗaya suna nan a shafin **Groups** na [dashboard](/features/dashboard), kuma `s` a kan asusu yana zaɓar rukunoninsa.

## Abubuwa

| Abu | Abin da ya ƙunsa |
| --- | --- |
| `mcp` | Sabobin MCP, ta `mcp.json` na rukunin ([Sabobin MCP da aka raba](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` na gaba ɗaya |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` da `history.jsonl`, don wani asusu ya iya `--continue` |

Abin da babu rukunin asusun da ke raba shi yana zama a foldar asusun kansa.

## Umarnai da zaɓuɓɓuka

| Umarni | Abin da yake yi |
| --- | --- |
| `ccam group` ko `ccam group list` | lissafa rukunonin |
| `ccam group show <group>` | abin da rukunin ke rabawa, asusansa, foldarsa da sabobin MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | ƙirƙiri rukuni; abubuwan asusansa suna motsawa cikinsa |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | canza abin da rukunin ke rabawa, ga kowane asusu a cikinsa |
| `ccam group join <group> <account>...` | saka asusai a rukunin |
| `ccam group leave <group> <account>...` | fitar da asusai; kowanne yana riƙe kwafi |
| `ccam group rename <group> <new-name>` | sake suna a cikin saitin kawai |
| `ccam group remove <group> [--leave]` | cire rukuni |

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: asusun da kwafinsa ke fara shiga rukunin | asusun tsoho, sannan bisa lamba |
| `--replace` | `join`: fara fita daga rukunonin da ke raba abubuwa iri ɗaya | a kashe |
| `--fresh` | asusan da suka daina raba wani abu suna farawa ba tare da shi ba maimakon da kwafi | a kashe |
| `--leave` | `remove`: fara fitar da asusan rukunin | ana buƙata idan yana da asusai |
| `--prefer <group\|mine\|merge>` | yadda za a yanke hukuncin rikici | tambaya; da `--yes`, tsohon zaɓin kowane fayil |
| `-y`, `--yes` | gudana ba tare da tambaya ba; rikici yana bin `--prefer` | a kashe |
| `--dry-run` (na gaba ɗaya) | buga shirin, kada a canza komai | a kashe |

Ana ba da asusai ta suna ko ta lamba; abubuwa kuma a matsayin jeri da aka raba da waƙafi. `ccam add` yana karɓar `--group <group>` (ana iya maimaitawa) ko `--no-group`.

## Yadda yake aiki

Fayilolin da abin ya shafa:

- `~/.config/ccam/config.toml`: shigarwar `[[group]]`: `name`, `items`, `members` kuma, idan sunan foldar ya bambanta, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` na asusun.
- `~/.local/share/ccam/groups/<folder>/`: kwafin da aka raba.

Idan aka karanta `config.toml`, ana warware rukunonin kowane asusu zuwa folda ɗaya ga kowane abu. Wannan yana ɗaukar kusan 25 µs ga asusai biyar, don haka farawa har yanzu yana karanta `config.toml` kawai. Saitin da asusu zai sami abu daga rukunoni biyu a cikinsa kuskure ne, kuma babu abin da ke farawa har sai an gyara shi.

Canji yana aiki a kan kwafin saitin:

1. Umarnin yana gyara kwafin (`AddGroup`, `SetGroupItems`, `JoinGroup` …), wanda ke duba ƙa'idodin.
2. `internal/group` yana kwatanta tsoho da sabo kuma, ga kowane asusu da abubuwansa ke canzawa, yana gina shiri da `share.Build`. Ana gina shirye-shiryen ɗaya bayan ɗaya, kuma kowanne yana ganin abin da na farkon za su saka a foldojin rukuni, don haka asusun farko ne kaɗai ke matsar da kwafinsa ciki. Kwafin farko yana fitowa daga foldar rukunin idan tana da shi tuni, sannan `--from`, sannan asusun tsoho, sannan sauran bisa lamba.
3. ccam yana ƙi idan ɗaya daga cikin waɗannan asusai yana aiki, kuma a Windows ba tare da Developer Mode ba idan fayil zai buƙaci mahaɗi.
4. Ana yanke hukuncin rikici, ta tambaya ko ta `--prefer`, ana buga shirin, kuma ccam yana tambaya.
5. Ana ajiye sabon saitin kafin komai ya motsa. Idan mataki ya tsaya a tsakiya, `ccam doctor --fix` da sake gudanarwa suna gina shirinsu zuwa ga rukunonin da aka ajiye kuma suna kammala canjin.
6. Shirye-shiryen suna gudana a tsari. Kowane mataki yana matsar da bayanai cikin rukunin ko cikin foldar ajiya; ba a goge komai.

`share.Layout` jerin izini ne na hanyoyin da kowane abu ya ƙunsa. Duk abin da ba ya cikinsa ba zai taɓa kai wa rukuni ba. `share.Never` yana ambaton `.credentials.json` da `.claude.json`, kuma `Build` yana tsayawa da kuskure na ciki idan ɗaya daga cikinsu ya taɓa bayyana a tsarin. Abin `mcp` ba shi da hanyoyi: jerin sabobin yana cikin `.claude.json`, don haka [rabawar MCP](/features/mcp) tana aiki ta `mcp.json` na rukunin.

Ga abin da ke faruwa ga folda ɗaya, `projects/`, idan `side` ya fara raba tarihi kuma rukunin ya riga yana da abubuwa:

<!-- motion:share-on -->

Idan rukunin bai da komai a wannan hanyar tukuna, asusun shi ne na farko da ya raba ta kuma ana matsar da foldarsa can kamar yadda take. In ba haka ba ana matsar da kowane fayil da rukunin ba shi da shi cikin rukunin, kuma fayilolin da ke ɓangarorin biyu suna zama a foldar asusun, don haka ana riƙe sigar rukunin. Sannan foldar asusun tana zuwa `.ccam-backup/<time>/` a cikin foldar asusun, kuma mahaɗi yana ɗaukar wurinta. Ana sanya wa fayilolin tattaunawa suna bisa session ID ɗinsu, don haka tarihi daga asusai da dama yana haɗuwa ba tare da karo ba. Rajistocin `plugins/` suna ajiye cikakkun hanyoyi ta foldar asusun; ccam yana nuna su zuwa foldar rukunin kuma yana haɗa shigarwar asusun cikin rajistar da rukunin ke da ita tuni.

Fita tana juya mahaɗin:

<!-- motion:share-off -->

ccam yana cire mahaɗin ne kawai idan yana nuna cikin foldar groups, sannan yana kwafar abin da ke cikin rukunin a yanzu zuwa asusun. Da `--fresh` ana tsallake kwafin. Rukunin yana riƙe abin da ke cikinsa don sauran asusai. Asusun da ke matsar da abu zuwa rukunin da foldarsa har yanzu babu komai a ciki yana kawo kwafin abin da ya gani, don haka yana riƙe fayilolinsa; sabobin MCP suna zuwa ta hanya ɗaya. Asusun da ya daina raba MCP yana rasa sabobin da ccam ya kwafa cikin `.claude.json` ɗinsa (ana ajiye kwafin ajiya), tunda daga rukunin suka fito.

Fayil ɗaya da ke ɓangarorin biyu mai abubuwa daban rikici ne:

<!-- motion:share-conflict -->

`history.jsonl` yana ba da `merge by line`, `keep group` da `use mine`, tare da haɗawa a farko. Sauran fayilolin `.json` suna ba da `keep group`, `use mine` da `merge JSON`. Duk wani abu dabam yana ba da `keep group` da `use mine`. `merge JSON` yana haɗa abubuwa (objects) a maimaice, yana haɗa arrays ba tare da maimaitawa ba, kuma yana riƙe ƙimar rukunin idan ɓangarorin biyu sun saita maɓalli ɗaya daban-daban, don haka sauran asusai ba sa ganin canji. Kowane zaɓi banda `keep group` yana fara ajiye kwafin rukunin a `.ccam-backup/<time>/group/` na asusun. Fayil ɗinka koyaushe yana zuwa ajiya.

A kowane farawa, `share.Check` yana kwatanta kowace hanya a tsarin da rukunonin asusun ta amfani da `lstat` da `readlink` kawai. Ana ƙirƙirar mahaɗin da ya ɓace nan take. Fayil na gaske inda ya kamata mahaɗi ya kasance, mahaɗi zuwa wurin da bai dace ba, ko tsohon mahaɗi yana dakatar da farawar kuma yana nuna `ccam doctor --fix`.

### Sake suna, cirewa da tsohon tsari

- `ccam group rename` yana canza saitin ne kawai. Foldar tana riƙe sunanta kuma rukunin yana rubuta ta a matsayin `dir`, don haka babu abin da ke motsawa kuma babu abin da zai iya ɓacewa.
- `ccam group remove` yana ƙin rukunin da ke da asusai sai dai idan ka ƙara `--leave`. Sannan kowane asusu yana riƙe kwafi, kuma foldar tana zuwa `~/.local/share/ccam/trash/` sai dai idan wani rukuni yana amfani da ita.
- Kafin rukunoni, `config.toml` yana ajiye teburin `[profile.share]` ga kowane asusu. ccam yana karanta wannan tsari a matsayin rukunoni ba tare da matsar da komai ba. A ajiyewa ta farko yana rubuta sabon tsarin kuma yana riƙe tsohon fayil a matsayin `config.toml.bak`. Abubuwan da asusai iri ɗaya ke rabawa suna zama rukuni ɗaya mai sunan foldar; sauran suna samun sunaye kamar `default-history` kuma suna ci gaba da amfani da folda ɗaya.

### A Windows

Windows yana ƙirƙirar symbolic links ne kawai a Developer Mode (Settings > System > For developers) ko ga administrator. Ba tare da wannan izinin ba, foldar da aka raba tana zama directory junction, wanda ke aiki iri ɗaya ga Claude Code. Fayil da aka raba (`settings.json`, `CLAUDE.md`, `history.jsonl`) zai iya zama symbolic link ne kawai, don haka ana ƙin canjin rukuni da ke buƙatar ɗaya kafin komai ya motsa.

## Tsaro da yanayi na musamman

- **Asusan da ke aiki.** Ana ƙin canji yayin da Claude Code ke gudana a matsayin ɗaya daga cikin asusan da yake matsarwa, domin zaman yana ci gaba da rubutu a cikin waɗannan fayiloli. Fita da shiga sun shafi asusun da aka ambata ne kawai. `--dry-run` har yanzu yana buga shirin.
- **Kwafin ajiya.** Kowane shiri yana samun foldar `.ccam-backup/<time>/` nasa a cikin foldar asusun. ccam ba ya taɓa goge su.
- **Mahaɗanka na kanka.** Idan hanya ta riga ta zama mahaɗin da ka yi, misali zuwa ma'ajiyar dotfiles, kuma rukunin babu komai a ciki, ccam yana kwafar abin da mahaɗin ke nunawa cikin rukunin kuma yana matsar da mahaɗin kansa zuwa ajiya. Ba a taɓa abin da ake nunawa ba.
- **Foldojin da aka shigo da su.** Asusun da aka yi wa rajista a wurinsa da `ccam import` (kamar `~/.claude`) yana canzawa a wurinsa, don haka `claude` na yau da kullum ma yana ganin canjin; shirin yana faɗin haka.
- **Asusai biyu a folda ɗaya.** Asusan da ke raba tarihi suna rubutu a fayiloli iri ɗaya. Idan na biyu ya fara a foldar da ɗaya ke aiki tuni, ccam yana gargaɗi kuma yana tambaya.
- **Gyara da hannu.** Memba da ba asusu ba ne ana watsi da shi kuma `ccam doctor` yana bayar da rahotonsa. Sunayen `local`, `on`, `off` da `none` ba za su iya zama sunayen rukuni ba.

## Masu alaƙa

- [Dashboard](/features/dashboard)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Sabobin MCP da aka raba](/features/mcp)
- [Kwafar saitin da ke akwai](/features/migrate)
- [Dukkan umarnai](/docs/commands)
