---
title: Mga shared na MCP server
description: Iisang listahan ng MCP server bawat group gamit ang ccam mcp add, add-json, list, remove at import, na ipinapasa sa --mcp-config o sini-sync sa bawat account.
updated: 2026-10-23
order: 5
translation: machine
sourceHash: fefca7f8d857
---

Iniimbak ng Claude Code ang mga user MCP server sa `.claude.json`, ang parehong file na naglalaman ng login ng account, kaya hindi puwedeng i-share ang file na ito sa pagitan ng mga account. Iniimbak ng ccam ang mga shared na server sa isang group file, ang `mcp.json`, at ibinibigay ito sa bawat account ng group sa pagsisimula. Gamitin ang `ccam mcp` para magdagdag ng server nang isang beses para sa lahat ng iyong account, at ang `ccam mcp import` para ilipat sa iba ang server na idinagdag mo sa loob ng isang session.

## Paggamit

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Nalalapat ang pagbabago simula sa susunod na pagsisimula ng bawat account; pinapanatili ng mga tumatakbong session ang mga server na kasama nila nang magsimula:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Sa ilalim ng mga group, pinapangalanan ng `ccam mcp list` ang bawat server na hindi umaabot sa ibang account: `only in account work: pg` para sa server na idinagdag ng isang account para sa sarili nito, `only in account side, for /home/me/code/app only: notes` para sa server na idinagdag para sa iisang folder (karaniwang `claude mcp add`, tingnan sa ibaba), `account work has its own …, used in its sessions in place of the group's` kapag may pangalan ng isang server ng group ang sariling server ng account, ang `.mcp.json` ng folder na ito, na may `waiting for approval in` at ang mga account na hinihintayan nito ng pag-apruba, at ang mga server na `in no account, only in ~/.claude.json (plain claude's)`, kasama ang [`ccam migrate --mcp`](/features/migrate) command na kumokopya sa mga ito.

Kapag nagdagdag ka ng pangalang mayroon na ang group, pinapalitan ang definition nito at ipi-print ang `updated` sa halip na `added`. Umaabot lang ang mga server sa isang account kung kasapi ito ng group na iyon at nagshe-share ng MCP ang group (`ccam group set <group> --add mcp`); kapag hindi, sasabihin iyon ng `add`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `-g`, `--group <group>` | target na group para sa `add`, `add-json` at `remove` | ang tanging group na pinagshe-share-an ng MCP ng ilang account, kung hindi ay `default` |
| `-- <command> [args…]` (`add`) | patakbuhin ang command na ito bilang stdio server | |
| `--url <url>` (`add`) | URL ng remote server (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport ng `--url` server | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, puwedeng ulitin | wala |
| `-e`, `--env <KEY=VALUE>` (`add`) | environment variable para sa stdio server, puwedeng ulitin | wala |
| `--force` (`import`) | palitan ang group server na may parehong pangalan pero naiiba; bina-back up ang lumang `mcp.json` | off |
| `--keep-local` (`import`) | panatilihin ang mga entry sa `.claude.json` ng account | off |
| `-h`, `--help` | help para sa bawat command | |
| `--dry-run` (global) | sinusunod ng `mcp list` lang; tinatanggihan ito ng `add`, `add-json`, `remove` at `import` at walang binabago | off |

Kailangan ng `add` ang command pagkatapos ng `--` o ang `--url`; kung wala, pumapalya ito na may `a command to run (after --) or --url is required`. Nalalapat ang `--header` sa mga `--url` server at ang `--env` sa mga command server. Tumatanggap ang `add-json` ng parehong JSON ng `claude mcp add-json`, at dapat object ang definition na may `command` o `url`. Laging ginagamit ng `import` ang group na pinagshe-share-an ng MCP ng account, at ipinapakita ng `list` ang bawat group, kaya hindi binabago ng `-g` ang mga ito.

## Paano ito gumagana

Ang mga file na kasangkot:

- `~/.local/share/ccam/groups/<group>/mcp.json`: ang mga shared na server, `{"mcpServers": {…}}`, ang format na binabasa ng `--mcp-config`. Isinusulat ito gamit ang `fsutil.AtomicWrite` at mode na `0600`, dahil madalas may dalang API key ang mga server entry, at sa ilalim lang ng lock na `<group>/.mcp.json.lock`, kaya parehong napapanatili ang server ng dalawang sabay na `ccam mcp add`.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ang sariling file ng Claude Code na may mga server ng account. Binabasa lang ng ccam ang value na `mcpServers` nito gamit ang scanner na `profile.TopLevelValue` at hindi kailanman pina-parse ang buong file sa launch path.
- `~/.local/state/ccam/mcp/<name>/`: ang file na `seen`, ang `sync-stamp` (na may hash ng bawat kopyang isinulat ng ccam, para matukoy nito ang kopyang walang nagbago mula sa kopyang binago mo), ang `sync-checked` (ang huling nakita ng sync strategy), ang `kept` (ang mga server na pinanatili ng account nang umalis ito sa pag-share ng MCP) at mga backup ng `.claude.json`.

May dalawang strategy, na itinatakda ng `mcp_strategy` sa `~/.config/ccam/config.toml`. Ang default ay ang flag strategy:

<!-- motion:mcp-flag -->

1. Hinahanap ng `launcher.applyMCP` ang MCP group ng account. Sinusuri ng `mcp.FlagArgs` gamit ang isang `stat` na mayroong `mcp.json` ang group, nang hindi ito pina-parse.
2. Inilalagay ng `Plan.InsertArgs` ang `--mcp-config=<path>` kasunod mismo ng `claude`, bago ang sarili mong mga argument, kaya nalalapat pa rin sa ibabaw nito ang sarili mong `--mcp-config` o `--strict-mcp-config`.
3. Sinusuri ng `mcp.NewLocal` kung nagbago ang `.claude.json` ng account mula noong huling pagsisimula, at saka lang nito pina-parse ang mga server nito.
4. Pinapalitan ng `syscall.Exec` ang ccam ng Claude Code.

Ipinapasa ang path bilang isang argument, `--mcp-config=<path>`, hindi bilang dalawa. Tumatanggap ng ilang value ang `--mcp-config` ng Claude Code, kaya babasahin ng `claude --mcp-config a.json "prompt"` ang prompt bilang pangalawang config file at papalya na may `MCP config file not found`. Doon na humihinto ang anyong may isang argument, at gumagana ang kasunod na prompt o subcommand. Walang flag na nakukuha ang `ccam work mcp list`, `get` at `remove`, dahil tinatanggihan ito roon ng Claude Code; sa halip, pinapangalanan ng ccam ang mga server ng group sa standard error ([mga detalye](/features/launch#shared-mcp-servers)). Nilo-load ng Claude Code ang mga `--mcp-config` server bilang dagdag sa sariling server ng account; kapag nagkapareho ang pangalan, ang entry sa `.claude.json` ang nananaig. Nagsusulat lang sa `.claude.json` ang flag strategy para sa Remote Control, kapag umalis ang isang account sa pag-share ng MCP o sumali ulit dito ([sa ibaba](#leaving-and-joining)), at para sa `ccam mcp import`.

Ang `claude remote-control` (`rc`) ay hindi tumatanggap ng `--mcp-config`, at binabasa lang ng mga session nito ang mga server mula sa `.claude.json`. Para sa `ccam work rc`, hindi isinasama ng ccam ang flag at, sa isang terminal, nagtatanong kung kokopyahin ang mga server ng group sa `.claude.json`, na may backup, gaya ng ginagawa ng `sync` strategy. Nagtatanong ito hanggang may kopya na, at pagkatapos ay muli lang kapag nagbago ang `mcp.json` ng group o nawala ang isang kinopyang server ([Remote Control](/features/launch#remote-control)). Dahil nananaig ang ganoong kopya sa `--mcp-config` kapag nagkapareho ang pangalan, ina-update ng `ccam mcp add`, `add-json`, `remove` at `import` ang mga kopyang hawak ng ibang account ng group:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Lumalabas lang ang linyang `✓ Updated` kapag nagbago ang isang kopya. Sa parehong paraan ina-update ang kopyang ginawa ng mas lumang ccam, gayundin ang kopya pagkatapos ng `ccam migrate`.

Ang server na idinagdag sa loob ng session gamit ang `claude mcp add` ay napupunta lang sa `.claude.json` ng account na iyon. Ginagamit ng karaniwang `claude mcp add` ang local scope: naglo-load ang server para sa iisang account na iyon sa iisang folder na iyon (ginagawa ito ng `-s user` na pag-aari ng account sa bawat folder). Hindi ito makikita ng ccam sa oras na iyon, dahil napalitan na ito ng Claude Code, kaya sinusuri nito sa susunod na pagsisimula. Hawak ng file na `seen` ang modification time ng `.claude.json` at ang mga pangalan ng server na naiulat na. Kapag iba ang oras, ikinukumpara ng ccam ang mga pangalan sa group at nagpi-print ng pahiwatig nang isang beses para sa bawat bagong server, nang hindi naghihintay ng input, sa parehong strategy:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Pinapangalanan din ng pahiwatig ang sariling server ng account na may pangalan ng isang server ng group, dahil ang sa account ang nananaig sa mga session nito.

Sa `mcp_strategy = "sync"`, lumalabas ang mga server sa `claude mcp list` at `/mcp` bilang karaniwang user server, at isinusulat ng ccam ang mga ito sa `.claude.json` ng bawat account bago magsimula:

<!-- motion:mcp-sync -->

1. Ikinukumpara ng `mcp.Sync` ang modification time ng `mcp.json` ng group sa `sync-stamp`, at ang `.claude.json` ng account sa huli nitong nakita. Kung walang nagbago sa dalawa, wala nang ibang nangyayari. Kung hindi, sinusuri din nito kung nasa `.claude.json` pa ang bawat server ng group, at ibinabalik ang kopyang nawala (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Kinukuha nito ang sariling lock ng Claude Code, ang folder na `.claude.json.lock` katabi ng file, na ginawa gamit ang `mkdir`. Naghihintay ito nang hanggang 2 segundo at itinuturing na stale ang lock na mas matanda sa 10 segundo, gaya ng ginagawa ng lock library ng Claude Code.
3. Sine-save nito ang kasalukuyang `.claude.json` sa `backups/` sa state folder at pinapanatili ang 5 pinakabago.
4. Itinatakda nito ang bawat group server sa `mcpServers`, inaalis ang mga server na isinulat nito dati na wala na sa group, at isinusulat nang atomic ang file na may orihinal na mode nito, kapag may nagbago lang. Binabago o inaalis lang nito ang mga kopyang isinulat nito na walang nagbago mula noon. Nananatiling pag-aari ng account ang sarili nitong server na may pangalan ng isang server ng group, at sinasabi ito ng ccam nang isang beses: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. Itinatala nito ang bagong stamp kasama ang hash ng bawat kopya, at nagsisimula ang Claude Code nang walang `--mcp-config`.

Nilalaktawan ang sync, na may `ccam: account work is already running, skipping MCP sync`, kapag may tumatakbong session ng account, paano man sinimulan ang Claude Code, dahil isinusulat ulit ng Claude Code ang `.claude.json` habang tumatakbo ito. Nagpi-print ang pumalyang sync ng `ccam: MCP sync failed (launching anyway): …` at nagsisimula pa rin.

### Pag-alis at pagsali

Pinapanatili ng account na tumitigil mag-share ng MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) ang mga server ng group: kinokopya ng plano ang mga ito sa `.claude.json` nito, na may backup, at nananatili ang sarili nitong mga server na may parehong pangalan. Sa `--fresh`, wala itong nakukuha, at inililista ng plano ang mga server na hindi na maglo-load. Kapag sumali itong muli sa group na may mga server na iyon, inaalis ang mga pinanatiling kopya na walang nagbago, kaya umaabot dito ang mga susunod na pagbabago ng group.

Inililipat ng `ccam mcp import` ang sariling server ng isang account papunta sa group nito:

<!-- motion:mcp-import -->

Binabasa ng ccam ang mga server ng account, ang sarili nito at ang mga mayroon ito para sa iisang folder lang, at ikinukumpara ang bawat isa sa server ng group na may parehong pangalan gamit ang `mcp.SameDef`, na hindi pinapansin ang pagkakasunod-sunod ng key at spacing. Nagiging group server para sa bawat folder ang folder server (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); tinatanggihan ang iisang pangalan na may dalawang magkaibang folder definition. Hindi iniaalok ang mga plugin server, dahil dala ang mga ito ng plugin. Kinokopya ang server na wala sa group o mayroon ang group na may parehong definition. Nilalaktawan ang naiiba, dahil ang definition ng group ang ginagamit ng bawat ibang account:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Sa `--force`, pinapalitan ang server ng group, at sine-save muna ang dating `mcp.json` sa `.ccam-backup/<time>/mcp.json` sa loob ng group folder. Pagkatapos ay inaalis ng `mcp.Remove` ang mga kinopyang entry mula sa `.claude.json` ng account, ang sarili nito at ang sa folder sa iisang pagsulat, sa ilalim ng parehong lock at pagkatapos ng backup sa state folder. Nananatiling naka-lock ang `mcp.json` ng group mula nang i-load ito hanggang mawala ang mga entry na iyon, kaya walang makakapag-alis ng server sa group sa pagitan. Kailangan ito dahil mananaig ang entry ng account sa entry ng group kapag nagkapareho ang pangalan. Iniiwan ng `--keep-local` ang mga ito.

## Kaligtasan at mga espesyal na kaso

- **Mga sign-in.** Iniimbak ng mga server na nagsa-sign in gamit ang OAuth ang kanilang mga token sa `.credentials.json` ng bawat account (`mcpOAuth`), hindi sa listahan ng server, kaya hindi nashe-share ang sign-in kapag nag-share ka ng server. Mag-sign in nang isang beses bawat account gamit ang `/mcp`. Hindi kailanman kinokopya ng ccam ang mga sign-in, kahit ang [`ccam migrate`](/features/migrate): nagla-log out sa bawat ibang may hawak ang kinopyang refresh token. Pinapangalanan ng [`ccam ext`](/features/ext) at `ccam doctor` ang mga server na hindi pa naka-sign in ang isang account habang naka-sign in naman ang ibang account, at sinasabi ng `ccam next --continue` kung alin ang wala pa sa bagong account. Hindi kailanman nakakatanggap ng login file ang group.
- **Mga tumatakbong account.** Kumokopya ang `import` sa group pero iniiwan ang sariling entry ng tumatakbong account at sinasabi iyon; patakbuhin ulit pagkatapos ng session. Nilalaktawan ng sync ang mga tumatakbong account. Parehong binibilang din ng dalawa ang Claude Code na sinimulan nang walang ccam.
- **Mga kopya para sa Remote Control.** Hindi ginagalaw ang kopya ng tumatakbong account: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Sa `mcp_strategy = "sync"`, walang kailangang i-update, dahil nagsi-sync ang bawat pagsisimula.
- **Mga sirang file.** Kung hindi valid na JSON ang `.claude.json`, hindi ito ginagalaw ng sync at iniuulat ang error. Kung hindi makuha ang lock sa loob ng 2 segundo, iniuulat ng ccam ang `could not lock … (Claude Code is writing the file)`.
- **Mga account na hindi nagshe-share ng MCP.** Tumatanggi ang `import` at sinasabing ilagay muna ang account sa isang group na nagshe-share ng MCP (`ccam group join <group> side`).
- **Mga pangalan ng group.** Pangalan ng group lang ang tinatanggap ng `-g`; nagiging folder ito sa ilalim ng `groups/`, kaya tinatanggihan ang mga pangalang gaya ng `../x`.

## Kaugnay

- [Sharing sa pagitan ng mga account](/features/share)
- [Paggamit ng Remote Control sa ccam](/guides/remote-control)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Kopyahin ang dati nang setup](/features/migrate)
- [Lahat ng command](/docs/commands)
