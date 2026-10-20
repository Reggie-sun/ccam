---
title: Pag-share gamit ang mga group
description: Nagshe-share ang ccam group ng piniling item ng mga account nito sa iisang folder; inililipat ng ccam group ang data nang may plano, backup at conflict choice.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 0c0338362ba2
---

May sariling folder ang bawat account, pero hindi kailangang hiwalay ang karamihan ng laman nito. Sinasabi ng isang **group** kung ano ang shine-share at kung sino ang nagshe-share: nili-link ng bawat account sa group ang mga item na iyon sa iisang folder, `~/.local/share/ccam/groups/<group>/`, kaya binabasa at sinusulatan nilang lahat ang iisang kopya. Hindi kailanman shine-share ang mga login token.

Puwedeng kasapi ang isang account sa ilang group, pero iisang group lang ang pinanggagalingan ng bawat item. Karaniwang setup ang isang group para sa settings at mga extension, at pangalawang group para sa kasaysayan ng usapan sa pagitan ng dalawang account:

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

Ang unang account na idinagdag mo ang gumagawa ng `default`, at sumasali rito ang mga susunod na account maliban kung iba ang sabihin mo.

## Paggamit

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

Nagpi-print ng plano ang bawat command na naglilipat ng data, nagtatanong tungkol sa bawat file na naiiba sa kopya ng group, at nagtatanong ng `Proceed? [y/N]` bago ito magbago ng anuman:

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

Nasa tab na **Groups** ng [dashboard](/features/dashboard) din ang parehong mga pagbabago, at pinipili ng `s` sa isang account ang mga group nito.

## Mga item

| Item | Ano ang sakop |
| --- | --- |
| `mcp` | Mga MCP server, sa pamamagitan ng `mcp.json` ng group ([Mga shared na MCP server](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | ang global na `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` at `history.jsonl`, para magamit ng ibang account ang `--continue` |

Nananatili sa sariling folder ng account ang item na hindi shine-share ng alinmang group ng account.

## Mga command at option

| Command | Ano ang ginagawa |
| --- | --- |
| `ccam group` o `ccam group list` | ilista ang mga group |
| `ccam group show <group>` | kung ano ang shine-share ng group, ang mga account nito, ang folder nito at ang mga MCP server |
| `ccam group add <group> --items <items> [--members <accounts>]` | gumawa ng group; inililipat dito ang mga item ng mga account nito |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | baguhin kung ano ang shine-share ng group, para sa bawat account nito |
| `ccam group join <group> <account>...` | ilagay ang mga account sa group |
| `ccam group leave <group> <account>...` | alisin ang mga account; may naiiwang kopya sa bawat isa |
| `ccam group rename <group> <new-name>` | palitan ang pangalan sa config lang |
| `ccam group remove <group> [--leave]` | mag-alis ng group |

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: ang account na unang ipapasok ang kopya sa group | ang default na account, saka ayon sa numero |
| `--replace` | `join`: umalis muna sa mga group na nagshe-share ng parehong mga item | off |
| `--fresh` | nagsisimulang wala ang item ang mga account na tumitigil mag-share nito, sa halip na may kopya | off |
| `--leave` | `remove`: alisin muna ang mga account ng group | kailangan kapag may mga account ito |
| `--prefer <group\|mine\|merge>` | paano pagpapasyahan ang mga conflict | magtanong; sa `--yes`, ang default ng bawat file |
| `-y`, `--yes` | tumakbo nang hindi nagtatanong; sinusunod ng mga conflict ang `--prefer` | off |
| `--dry-run` (global) | i-print ang plano, walang binabago | off |

Ibinibigay ang mga account ayon sa pangalan o numero, at ang mga item bilang listahang pinaghihiwalay ng kuwit. Tumatanggap ang `ccam add` ng `--group <group>` (puwedeng ulitin) o `--no-group`.

## Paano ito gumagana

Ang mga file na kasangkot:

- `~/.config/ccam/config.toml`: ang mga entry na `[[group]]`: `name`, `items`, `members` at, kapag iba ang pangalan ng folder, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: ang `CLAUDE_CONFIG_DIR` ng account.
- `~/.local/share/ccam/groups/<folder>/`: ang mga shared na kopya.

Kapag binasa ang `config.toml`, nireresolba ang mga group ng bawat account sa iisang folder bawat item. Mga 25 µs ito para sa limang account, kaya `config.toml` pa rin lang ang binabasa ng isang launch. Error ang config kung saan makakakuha ang isang account ng item mula sa dalawang group, at walang nagsisimula hangga't hindi ito naaayos.

Gumagana ang isang pagbabago sa isang kopya ng config:

1. Binabago ng command ang kopya (`AddGroup`, `SetGroupItems`, `JoinGroup` …), na sinusuri ang mga patakaran.
2. Ikinukumpara ng `internal/group` ang luma at ang bago at, para sa bawat account na nagbabago ang mga item, bumubuo ito ng plano gamit ang `share.Build`. Isa-isang binubuo ang mga plano, at nakikita ng bawat isa kung ano ang ilalagay ng mga nauna sa mga group folder, kaya ang unang account lang ang naglilipat ng kopya nito sa group. Galing ang unang kopya sa group folder kung mayroon na ito, saka sa `--from`, saka sa default na account, saka sa iba ayon sa numero.
3. Tumatanggi ang ccam kung tumatakbo ang isa sa mga account na iyon, at sa Windows nang walang Developer Mode kung mangangailangan ng link ang isang file.
4. Pinagpapasyahan ang mga conflict, sa pamamagitan ng tanong o ng `--prefer`, ipi-print ang plano, at magtatanong ang ccam.
5. Sine-save ang bagong config bago may ilipat. Kapag huminto ang isang hakbang sa kalagitnaan, binubuo ng `ccam doctor --fix` at ng muling pagpapatakbo ang kanilang plano patungo sa mga naka-save na group at tinatapos ang pagbabago.
6. Pinapatakbo ang mga plano nang sunod-sunod. Inililipat ng bawat hakbang ang data sa group o sa backup folder; walang binubura.

Ang `share.Layout` ay listahan ng mga pinapayagang path na sakop ng bawat item. Hindi kailanman makakaabot sa group ang anumang wala rito. Pinapangalanan ng `share.Never` ang `.credentials.json` at `.claude.json`, at humihinto ang `Build` na may internal error kung lumabas man ang alinman sa layout. Walang path ang item na `mcp`: nasa loob ng `.claude.json` ang listahan ng server, kaya gumagana ang [MCP sharing](/features/mcp) sa pamamagitan ng `mcp.json` ng group.

Ito ang nangyayari sa isang folder, ang `projects/`, kapag nagsimulang mag-share ng kasaysayan ang `side` at may laman na ang group:

<!-- motion:share-on -->

Kung wala pang laman ang group sa path na iyon, ang account ang unang nagshe-share nito at inililipat doon ang folder nito nang walang pagbabago. Kung hindi, inililipat sa group ang bawat file na wala pa sa group, at nananatili sa folder ng account ang mga file na mayroon sa magkabilang panig, kaya pinapanatili ang bersyon ng group. Pagkatapos ay napupunta ang folder ng account sa `.ccam-backup/<time>/` sa loob ng folder ng account, at may link na pumapalit sa puwesto nito. Ipinapangalan ang mga file ng usapan ayon sa kanilang session ID, kaya nagsasama ang mga kasaysayan mula sa ilang account nang walang banggaan. Nag-iimbak ang mga registry ng `plugins/` ng mga absolute na path na dumadaan sa folder ng account; itinuturo ng ccam ang mga ito sa folder ng group at pinagsasama ang mga entry ng account sa registry na mayroon na ang group.

Binabaligtad ng pag-alis sa group ang link:

<!-- motion:share-off -->

Inaalis lang ng ccam ang link kung tumuturo ito sa groups folder, saka kinokopya ang kasalukuyang nilalaman ng group sa account. Sa `--fresh`, nilalaktawan ang pagkopya. Pinapanatili ng group ang nilalaman nito para sa ibang account. Ang account na naglilipat ng item sa isang group na wala pang laman ang folder ay nagdadala ng kopya ng nakita nito, kaya napapanatili nito ang mga file nito; ganoon din ang mga MCP server. Nawawala sa account na tumitigil mag-share ng MCP ang mga server na kinopya ng ccam sa `.claude.json` nito (may naiiwang backup), dahil galing ang mga iyon sa group.

Conflict ang iisang file na mayroon sa magkabilang panig na magkaiba ang nilalaman:

<!-- motion:share-conflict -->

Nag-aalok ang `history.jsonl` ng `merge by line`, `keep group` at `use mine`, na nauuna ang pagsasama. Nag-aalok ang ibang `.json` file ng `keep group`, `use mine` at `merge JSON`. Nag-aalok ang anumang iba ng `keep group` at `use mine`. Pinagsasama ng `merge JSON` ang mga object nang recursive, pinagdudugtong ang mga array nang walang duplicate, at pinapanatili ang value ng group kapag magkaiba ang itinakda ng magkabilang panig sa parehong key, kaya walang nakikitang pagbabago ang ibang account. Bawat pagpipilian maliban sa `keep group` ay sine-save muna ang kopya ng group sa `.ccam-backup/<time>/group/` ng account. Laging napupunta sa backup ang sarili mong file.

Sa bawat pagsisimula, ikinukumpara ng `share.Check` ang bawat path sa layout sa mga group ng account gamit lang ang `lstat` at `readlink`. Ginagawa agad ang nawawalang link. Pinapahinto ang pagsisimula ng totoong file kung saan dapat may link, ng link sa maling lugar, o ng lumang link, at itinuturo ang `ccam doctor --fix`.

### Pagpapalit ng pangalan, pag-alis at ang lumang format

- Config lang ang binabago ng `ccam group rename`. Pinapanatili ng folder ang pangalan nito at itinatala ito ng group bilang `dir`, kaya walang inililipat at walang puwedeng mawala.
- Tumatanggi ang `ccam group remove` sa group na may mga account maliban kung idagdag mo ang `--leave`. May naiiwang kopya sa bawat account, at napupunta ang folder sa `~/.local/share/ccam/trash/` maliban kung ginagamit ito ng ibang group.
- Bago nagkaroon ng mga group, nag-iimbak ang `config.toml` ng table na `[profile.share]` bawat account. Binabasa ng ccam ang format na iyon bilang mga group nang walang inililipat. Sa unang pag-save, isinusulat nito ang bagong format at itinatabi ang lumang file bilang `config.toml.bak`. Nagiging iisang group na ipinangalan sa folder ang mga item na shine-share ng parehong mga account; nakakakuha ang iba ng mga pangalang gaya ng `default-history` at patuloy na ginagamit ang parehong folder.

### Sa Windows

Gumagawa lang ang Windows ng mga symbolic link sa Developer Mode (Settings > System > For developers) o para sa isang administrator. Kapag wala ang karapatang iyon, nagiging directory junction ang shared na folder, na pareho ang paggana para sa Claude Code. Symbolic link lang ang puwedeng maging shared na file (`settings.json`, `CLAUDE.md`, `history.jsonl`), kaya tinatanggihan ang pagbabago sa group na nangangailangan nito bago may ilipat.

## Kaligtasan at mga espesyal na kaso

- **Mga tumatakbong account.** Tinatanggihan ang isang pagbabago habang tumatakbo ang Claude Code bilang isa sa mga account na inililipat nito, dahil patuloy na nagsusulat ang session sa mga file na iyon. Ang account lang na pinangalanan ang saklaw ng pag-alis at pagsali. Ipi-print pa rin ng `--dry-run` ang plano.
- **Mga backup.** May sariling `.ccam-backup/<time>/` folder ang bawat plano sa folder ng account. Hindi kailanman binubura ng ccam ang mga ito.
- **Sarili mong mga link.** Kung link na ginawa mo na ang isang path, halimbawa sa isang dotfiles repository, at walang laman ang group, kinokopya ng ccam ang target ng link sa group at inililipat ang link mismo sa backup. Hindi ginagalaw ang target.
- **Mga in-import na folder.** Binabago kung nasaan ito ang account na nirehistro kung nasaan ito gamit ang `ccam import` (gaya ng `~/.claude`), kaya nakikita rin ng karaniwang `claude` ang pagbabago; sinasabi iyon ng plano.
- **Dalawang account sa isang folder.** Sinusulatan ng mga account na nagshe-share ng kasaysayan ang parehong mga file. Kapag nagsimula ang pangalawa sa folder kung saan may tumatakbo na, nagbababala ang ccam at nagtatanong.
- **Mga manwal na pagbabago.** Hindi pinapansin ang miyembrong hindi account, at iniuulat ito ng `ccam doctor`. Hindi puwedeng maging pangalan ng group ang `local`, `on`, `off` at `none`.

## Kaugnay

- [Ang dashboard](/features/dashboard)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Mga shared na MCP server](/features/mcp)
- [Kopyahin ang dati nang setup](/features/migrate)
- [Lahat ng command](/docs/commands)
