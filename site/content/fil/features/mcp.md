---
title: Mga shared na MCP server
description: Magkaroon ng iisang listahan ng mga MCP server bawat group gamit ang ccam mcp add, add-json, list, remove at import, na ipinapasa sa --mcp-config o sini-sync sa bawat account.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: f1de7db4b7be
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

Kapag nagdagdag ka ng pangalang mayroon na ang group, pinapalitan ang definition nito at ipi-print ang `updated` sa halip na `added`. Umaabot lang ang mga server sa isang account kung nagshe-share ito ng MCP sa group na iyon (`ccam share <name> --mcp on`); kapag wala pang account na nagshe-share ng MCP, sasabihin iyon ng `add`.

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: ang mga shared na server, `{"mcpServers": {…}}`, ang format na binabasa ng `--mcp-config`. Isinusulat ito gamit ang `fsutil.AtomicWrite` at mode na `0600`, dahil madalas may dalang API key ang mga server entry.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ang sariling file ng Claude Code na may mga server ng account. Binabasa lang ng ccam ang value na `mcpServers` nito gamit ang scanner na `profile.TopLevelValue` at hindi kailanman pina-parse ang buong file sa launch path.
- `~/.local/state/ccam/mcp/<name>/`: ang file na `seen`, ang `sync-stamp` at mga backup ng `.claude.json`.

May dalawang strategy, na itinatakda ng `mcp_strategy` sa `~/.config/ccam/config.toml`. Ang default ay ang flag strategy:

<!-- motion:mcp-flag -->

1. Hinahanap ng `launcher.applyMCP` ang MCP group ng account. Sinusuri ng `mcp.FlagArgs` gamit ang isang `stat` na mayroong `mcp.json` ang group, nang hindi ito pina-parse.
2. Inilalagay ng `Plan.InsertArgs` ang `--mcp-config=<path>` kasunod mismo ng `claude`, bago ang sarili mong mga argument, kaya nalalapat pa rin sa ibabaw nito ang sarili mong `--mcp-config` o `--strict-mcp-config`.
3. Sinusuri ng `mcp.NewLocal` kung nagbago ang `.claude.json` ng account mula noong huling pagsisimula, at saka lang nito pina-parse ang mga server nito.
4. Pinapalitan ng `syscall.Exec` ang ccam ng Claude Code.

Ipinapasa ang path bilang isang argument, `--mcp-config=<path>`, hindi bilang dalawa. Tumatanggap ng ilang value ang `--mcp-config` ng Claude Code, kaya babasahin ng `claude --mcp-config a.json "prompt"` ang prompt bilang pangalawang config file at papalya na may `MCP config file not found`. Doon na humihinto ang anyong may isang argument, at gumagana ang kasunod na prompt o subcommand gaya ng `mcp list`. Nilo-load ng Claude Code ang mga `--mcp-config` server bilang dagdag sa sariling server ng account; kapag nagkapareho ang pangalan, ang entry sa `.claude.json` ang nananaig. Sa Remote Control lang nagsusulat sa `.claude.json` ang flag strategy.

Ang `claude remote-control` (`rc`) ay hindi tumatanggap ng `--mcp-config`, at binabasa lang ng mga session nito ang mga server mula sa `.claude.json`. Para sa `ccam work rc`, hindi isinasama ng ccam ang flag at, sa isang terminal, nagtatanong kung kokopyahin ang mga server ng group sa `.claude.json`, na may backup, gaya ng ginagawa ng `sync` strategy. Nagtatanong ito hanggang may kopya na, at pagkatapos ay muli lang kapag nagbago ang `mcp.json` ng group ([Remote Control](/features/launch#remote-control)). Dahil nananaig ang ganoong kopya sa `--mcp-config` kapag nagkapareho ang pangalan, ina-update ng `ccam mcp add`, `add-json`, `remove` at `import` ang mga kopyang hawak ng ibang account ng group:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Ang server na idinagdag sa loob ng session gamit ang `claude mcp add` ay napupunta lang sa `.claude.json` ng account na iyon. Hindi ito makikita ng ccam sa oras na iyon, dahil napalitan na ito ng Claude Code, kaya sinusuri nito sa susunod na pagsisimula. Hawak ng file na `seen` ang modification time ng `.claude.json` at ang mga pangalan ng server na naiulat na. Kapag iba ang oras, ikinukumpara ng ccam ang mga pangalan sa group at nagpi-print ng pahiwatig nang isang beses para sa bawat bagong server, nang hindi naghihintay ng input:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Sa `mcp_strategy = "sync"`, lumalabas ang mga server sa `claude mcp list` at `/mcp` bilang karaniwang user server, at isinusulat ng ccam ang mga ito sa `.claude.json` ng bawat account bago magsimula:

<!-- motion:mcp-sync -->

1. Ikinukumpara ng `mcp.Sync` ang modification time ng `mcp.json` ng group sa `sync-stamp`. Kung hindi ito nagbago, hindi binubuksan ang `.claude.json` ng account.
2. Kinukuha nito ang sariling lock ng Claude Code, ang folder na `.claude.json.lock` katabi ng file, na ginawa gamit ang `mkdir`. Naghihintay ito nang hanggang 2 segundo at itinuturing na stale ang lock na mas matanda sa 10 segundo, gaya ng ginagawa ng lock library ng Claude Code.
3. Sine-save nito ang kasalukuyang `.claude.json` sa `backups/` sa state folder at pinapanatili ang 5 pinakabago.
4. Itinatakda nito ang bawat group server sa `mcpServers`, inaalis ang mga server na isinulat nito dati na wala na sa group, pinapanatili ang sariling server ng account, at isinusulat nang atomic ang file na may orihinal na mode nito.
5. Itinatala nito ang bagong stamp at ang mga pangalang pinamamahalaan nito, at nagsisimula ang Claude Code nang walang `--mcp-config`.

Nilalaktawan ang sync, na may `ccam: account work is already running, skipping MCP sync`, kapag may tumatakbong session ng account, dahil isinusulat ulit ng Claude Code ang `.claude.json` habang tumatakbo ito. Nagpi-print ang pumalyang sync ng `ccam: MCP sync failed (launching anyway): …` at nagsisimula pa rin. Sa flag strategy lang ang pahiwatig para sa mga server na idinagdag sa loob ng session.

Inililipat ng `ccam mcp import` ang sariling server ng isang account papunta sa group nito:

<!-- motion:mcp-import -->

Binabasa ng ccam ang mga server ng account gamit ang `mcp.LocalServers` at ikinukumpara ang bawat isa sa server ng group na may parehong pangalan gamit ang `mcp.SameDef`, na hindi pinapansin ang pagkakasunod-sunod ng key at spacing. Kinokopya ang server na wala sa group o mayroon ang group na may parehong definition. Nilalaktawan ang naiiba, dahil ang definition ng group ang ginagamit ng bawat ibang account:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Sa `--force`, pinapalitan ang server ng group, at sine-save muna ang dating `mcp.json` sa `.ccam-backup/<time>/mcp.json` sa loob ng group folder. Pagkatapos ay inaalis ng `mcp.RemoveLocal` ang mga kinopyang entry mula sa `.claude.json` ng account, sa ilalim ng parehong lock at pagkatapos ng backup sa state folder. Kailangan ito dahil mananaig ang entry ng account sa entry ng group kapag nagkapareho ang pangalan. Iniiwan ng `--keep-local` ang mga ito.

## Kaligtasan at mga espesyal na kaso

- **Mga sign-in.** Iniimbak ng mga server na nagsa-sign in gamit ang OAuth ang kanilang mga token sa `.credentials.json` ng bawat account (`mcpOAuth`), hindi sa listahan ng server, kaya hindi nashe-share ang sign-in kapag nag-share ka ng server. Mag-sign in nang isang beses bawat account gamit ang `/mcp`, o kopyahin ang mga sign-in mula sa `~/.claude` gamit ang [`ccam migrate`](/features/migrate). Hindi kailanman nakakatanggap ng login file ang group.
- **Mga tumatakbong account.** Kumokopya ang `import` sa group pero iniiwan ang sariling entry ng tumatakbong account at sinasabi iyon; patakbuhin ulit pagkatapos ng session. Nilalaktawan ng sync ang mga tumatakbong account.
- **Mga kopya para sa Remote Control.** Hindi ginagalaw ang kopya ng tumatakbong account: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Sa `mcp_strategy = "sync"`, walang kailangang i-update, dahil nagsi-sync ang bawat pagsisimula.
- **Mga sirang file.** Kung hindi valid na JSON ang `.claude.json`, hindi ito ginagalaw ng sync at iniuulat ang error. Kung hindi makuha ang lock sa loob ng 2 segundo, iniuulat ng ccam ang `could not lock … (Claude Code is writing the file)`.
- **Mga account na hindi nagshe-share ng MCP.** Tumatanggi ang `import` at sinasabing patakbuhin muna ang `ccam share side --mcp on`.
- **Mga pangalan ng group.** Pangalan ng group lang ang tinatanggap ng `-g`; nagiging folder ito sa ilalim ng `groups/`, kaya tinatanggihan ang mga pangalang gaya ng `../x`.

## Kaugnay

- [Sharing sa pagitan ng mga account](/features/share)
- [Paggamit ng Remote Control sa ccam](/guides/remote-control)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Kopyahin ang dati nang setup](/features/migrate)
- [Lahat ng command](/docs/commands)
