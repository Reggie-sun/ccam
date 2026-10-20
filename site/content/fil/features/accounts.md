---
title: Pamamahala ng mga account
description: Magdagdag, mag-log in, maglista, magpalit ng pangalan, mag-alis at mag-import ng Claude Code account gamit ang ccam add, login, list, rename, remove at import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

Ang bawat ccam account ay hiwalay na Claude Code config folder na may sariling login. Gumagawa, nagsa-sign in, naglilista, nagpapalit ng pangalan at nag-aalis ng mga account ang mga command na ito, at nagrerehistro rin ng mga folder na mayroon ka na.

## Paggamit

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Saanman tumatanggap ang isang command ng `<name|number>`, puwede mong ibigay ang pangalan ng account o ang slot number nito, ang column na `#` ng `ccam list`.

Kapag pinatakbo sa terminal nang walang option, nagbubukas ang `ccam add` ng wizard. Hinihingi nito ang pangalan at kung paano magsa-sign in ang account (Claude login ngayon, API key, o Claude login mamaya), saka ipinapakita ang mga group na sasalihan na naka-check ang `default` (kapag wala pang group, ang mga item na ishe-share ng group na `default`), saka nagtatanong ng `Create a shortcut command for work? [y/N]`, at nagtatapos sa buod na kukumpirmahin mo. Bumabalik nang isang hakbang ang `Esc`. Nilalaktawan ang wizard kapag may `--group`, `--no-group`, `--shortcut` o `--no-login`, o kapag pipe ang input sa halip na terminal; sa ganitong kaso kailangan ang pangalan at walang shortcut na gagawin maliban kung may pinangalanan ang `--shortcut`.

Kapag walang `--group`, sumasali ang bagong account sa group na `default`. Ang unang account ang gumagawa nito, na nagshe-share ng mga MCP server, settings, `CLAUDE.md` at mga extension; tingnan ang [pag-share gamit ang mga group](/features/share).

Ipinapakita ng `ccam list` ang slot, pangalan, shortcut, email, mga group at huling paggamit, na may `*` bago ang default na account:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

Sa `--json`, object ang bawat account na may `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (wala kapag hindi pa kailanman sinimulan ang account), `groups` at `share`. Itinutugma ng `share` ang bawat item (`mcp`, `settings`, `memory`, `extensions`, `history`) sa group na pinanggagalingan nito, o sa `"local"`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `ccam add --no-login` | Gumawa ng account nang hindi nagla-log in; mag-log in mamaya gamit ang `ccam login <name>` | off |
| `ccam add --group <group>` | Sumali sa group na ito sa halip na sa `default`. Puwedeng ulitin; ginagawa ang `default` kapag wala pa ito | `default` |
| `ccam add --no-group` | Huwag sumali sa anumang group: walang shine-share ang account | off |
| `ccam add --shortcut <command>` | Gumawa rin ng shortcut command na ito sa `~/.local/bin` | wala |
| `ccam add --api --key-from <source>` | Account na gumagamit ng API key sa halip na login; tingnan ang [Mga API key account](/features/api) | off |
| `ccam list --json` | I-print ang mga account bilang JSON | off |
| `ccam remove -y`, `--yes` | Mag-alis nang hindi nagtatanong | off |
| `ccam import --name <name>` | Pangalan ng account para sa folder na ibinigay mo | kinukuha mula sa pangalan ng folder |
| `ccam import -y`, `--yes` | I-import ang bawat nahanap na folder nang hindi nagtatanong | off |
| `--dry-run` | Sinusunod ng `ccam list` at `ccam completion`, na walang binabago. Tinatanggihan ng `add`, `login`, `logout`, `use`, `rename`, `remove` at `import`, na may `` --dry-run is not supported by `ccam add`; nothing was changed `` (kasama ang pangalan ng command) | off |

Walang sariling option ang `ccam login`, `ccam logout`, `ccam use` at `ccam rename`. Tumatanggap ang `ccam completion` ng `bash`, `zsh` o `powershell`; pumapalya ang iba na may `only bash, zsh and powershell are supported`.

Ang mga pangalan ng account ay 1 hanggang 32 character ng `a-z`, `0-9`, `-` at `_`, na nagsisimula sa titik o numero. Hindi puwedeng puro numero ang pangalan, dahil ang ibig sabihin ng `ccam 2` ay slot 2. Hindi rin puwedeng ccam command ang pangalan (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), kaya hindi kailanman malabo ang `ccam <name>`. Ang listahan ay `config.Reserved` sa `internal/config`.

## Paano ito gumagana

Pinapatakbo ng `ccam add work --shortcut ccwork` ang mga hakbang na ito sa `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. Sinusuri ng `config.ValidateName` ang pangalan, at sinusuri ng `shim.Check` ang shortcut bago may gawin.
2. Ginagawa ng `profile.Add` ang `~/.local/share/ccam/profiles/work/` na may mode na `0700` at ibinibigay sa account ang pinakamababang bakanteng slot number (`Config.NextSlot`), kaya muling ginagamit ang numero ng inalis na account. Ang unang account ang nagiging default.
3. Sumasali ang account sa mga group nito (`default` maliban kung iba ang sinasabi ng `--group` o `--no-group`; ang unang account ang gumagawa ng `default`), at nili-link ng `applyShare` ang bago at walang lamang folder sa mga iyon gamit ang `share.Build` at saka ang `Apply`, halimbawa ang `settings.json` sa `groups/default/settings.json`. Kapag pumalya ito, inaalis ang bagong folder, dahil haharangin ng natirang folder ang bawat pagsubok ulit na may `folder already exists`.
4. Sine-save ang account sa `~/.config/ccam/config.toml` sa pamamagitan ng `fsutil.AtomicWrite`: isang pansamantalang file, `fsync`, saka `rename`, kaya hindi kailanman mag-iiwan ng kalahating config ang isang crash.
5. Sinisimulan ng `runLogin` ang `claude auth login` na nakatakda ang `CLAUDE_CONFIG_DIR` sa folder ng account, at ang `--email` na may huling email ng account para napunan na ang login page. Kapag pumalya ang login, nananatili ang account, at ipi-print ng ccam ang error at sasabihing patakbuhin mo ang `ccam login work` mamaya.
6. Kinukuha ng `profile.ReadIdentity` ang `oauthAccount` mula sa `profiles/work/.claude.json`: ang email, ang mga id ng account at organisasyon, ang pangalan ng organisasyon at ang plan (`claude_max`, `claude_team` …). Puwedeng ilang megabyte ang file na iyon, kaya hinahanap ng `TopLevelValue` ang iisang key sa halip na i-parse ang buong file. Itinatabi ng ccam ang mga ito sa `[profile.login]` sa `config.toml`.
7. Nili-link ng `shim.Install` ang `~/.local/bin/ccwork` sa ccam binary, na hinahanap ang `argv[0]` sa `config.toml` kapag nagsisimula ito.

Puwedeng umabot nang ilang minuto ang login, at baka may ibang ccam na mag-save ng config habang naghihintay. Kaya ang mga command na naghihintay ay nagse-save sa pamamagitan ng `env.update`, na binabasa ulit ang `config.toml`, inilalapat ang iisang pagbabago (dito, ang email) at nagse-save. Binabasa ulit ng `ccam remove` ang config pagkatapos ng tanong nito sa parehong dahilan.

Pinapatakbo ng `ccam login` ang hakbang 5 at 6 para sa isang umiiral na account. Pahiwatig lang para sa login page ang email para sa Claude Code, kaya puwedeng mapunta ang isang login sa ibang Claude account kaysa sa dati ng ccam account; ikinukumpara ng ccam ang bagong login sa `[profile.login]` at sinasabi iyon (`ccam: this login is other@company.com, but work was dev@company.com before.`). Pinapanatili nito ang bagong login at sinasabi kung paano ito babawiin: `ccam logout work`, saka `ccam login work` at mag-sign in gamit ang tamang email. Kapag naka-log in ang isa pang ccam account sa parehong Claude account, nagbababala ang ccam na iisa ang usage limit ng dalawa. Isa-isang ina-log in ng `ccam login --expired` ang bawat account na hindi pa kailanman naka-log in o nag-expire na ang login (ang refresh token nito).

Pinapatakbo ng `ccam logout` ang `claude auth logout` sa account at binubura ang naka-save na email; nananatili ang tala ng login, kaya sinusuri laban dito ang susunod na login. Isinusulat lang ulit ng `ccam use` ang `default` sa `config.toml`.

Nagdadagdag ang `ccam list` ng column na `Organization` kapag kabilang ang mga account sa higit sa isang organisasyon (`personal` ang ipinapakita ng mga personal na plan), at may `plan` at `org` ang `ccam list --json` para sa bawat account.

### Pag-sign in sa isang remote na makina

Sa SSH, puwedeng masira ang bahagi ng `claude auth login` na nasa browser. Naghihintay ang Claude Code sa sign-in sa isang callback server sa `http://localhost:<port>/callback` sa server, pero ang browser na tumatapos sa sign-in ay tumatakbo sa iyong computer, kung saan walang pinupuntahan ang address na iyon at nagpapakita ang browser ng error gaya ng `ERR_CONNECTION_REFUSED`. Dala pa rin ng address ang one-time code, kaya kaya itong ihatid ng ccam:

<!-- motion:accounts-relay -->

Hindi ibinibigay ng `runLogin` ang iyong terminal sa Claude Code. Ikinokonekta nito ang standard input ng Claude Code sa isang pipe at kinokopya ng `relayLogin` doon ang bawat linyang tina-type mo, kaya dumarating pa rin ang code na idinikit sa prompt na `Paste code here` ng Claude Code. Kapag callback address ang isang linya, tinatanggap lang ito ng `callbackURL` kung `localhost`, `127.0.0.1` o `[::1]` ang host, may port, `/callback` ang path at parehong naroon ang `code` at `state`; puwedeng wala ang `http://`. Pagkatapos ay hinihiling ng `deliverCallback` ang address na iyon sa server na may 10 segundong limit at nang hindi sumusunod sa mga redirect, gaya ng gagawin sana ng browser. Umaabot ang request sa callback server ng naghihintay na `claude auth login` (ang putol-putol na linya sa diagram). Natatanggap ng Claude Code ang code at sine-save ang login, sa Linux sa `.credentials.json` sa folder ng account. Hindi kailanman mapapakontak ng isang idinikit na linya ang ccam sa ibang host.

Ipi-print ito ng ccam bago magsimula ang login:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Kapag nagtagumpay, ipi-print nito ang `→ Passed the sign-in to Claude Code on this machine.`

Habang naghihintay ng input ang login sa isang terminal, pinapatay ng ccam ang echo at line mode ng terminal at nagpapakita ng `*` para sa bawat character na tina-type mo, kaya hindi lumalabas sa screen ang idinikit na code o callback address. Ang `readLine` mismo ang gumagawa ng line editing: ipinapadala ng `Enter` ang linya, binubura ng `Backspace` ang isang character, at tinatapos ng `Ctrl+D` sa walang lamang linya ang input. Ibinabalik sa dati ang terminal kapag natapos ang login, at gayundin sa `Ctrl+C`, `SIGTERM` o `SIGHUP`, na tumatapos pa rin sa ccam. Kapag walang terminal sa standard input at standard error, kinokopya ang mga linya habang dumarating.

### Pagpapalit ng pangalan, pag-alis at pag-import

Kumikilos ang `ccam rename` at `ccam remove` sa folder ng account at sa mga iniimbak ng ccam sa ibang lugar sa ilalim ng pangalan ng account. Inililista ng `nameState` ang mga path na iyon: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` at ang mga usage cache na `~/.local/state/ccam/usage/<name>.json` at `<name>.statusline.json`. Hinaharap ng `ccam remove work` ang mga ito sa ganitong pagkakasunod-sunod:

<!-- motion:accounts-remove -->

Tumatanggi muna ang `ccam remove` kung tumatakbo ang account, saka nagtatanong. Kung nagshe-share ng mga extension ang account, isinusulat ulit ng `RewritePluginPaths` ang `installed_plugins.json` at `known_marketplaces.json` ng group, dahil nag-iimbak ang Claude Code ng mga absolute na plugin path na puwedeng dumaan sa folder ng account na ito. Inililipat ng `profile.Trash` ang folder sa `~/.local/share/ccam/trash/<name>-<time>/`, inaalis ang account sa `config.toml`, binubura ang mga path ng `nameState` para hindi ito mamana ng susunod na account na may parehong pangalan, at inaalis ang link ng shortcut. Para maibalik, ilipat pabalik ang folder mula sa trash.

Pinapalitan ng `ccam rename` ang pangalan ng folder, isinusulat ulit ang mga plugin path sa `plugins/` nito mula sa lumang folder papunta sa bago, dinadala ang default, at inililipat ang mga path ng `nameState` sa bagong pangalan. Pinapanatili ng shortcut command ang pangalan nito. Kapag pumalya ang pagsulat ulit ng plugin, ibinabalik ang lumang pangalan ng folder, dahil hindi naka-save ang config.

Nirerehistro ng `ccam import` ang isang folder kung nasaan ito, bilang `dir` sa `config.toml`. Hindi nito kailanman inililipat o kinokopya ito: kapag kinopya ang refresh token, magla-log out sa isa't isa ang mga kopya. Kapag walang folder argument, hinahanap nito ang mga hindi pa nakarehistrong `~/.claude` at `~/.claude-*` folder at nagtatanong tungkol sa bawat isa; kinukuha ang mga pangalan mula sa folder (`main` ang `~/.claude`, `work` ang `~/.claude-work`, `acct2` ang `~/.claude-2`), at ang folder na hindi makuhanan ng pangalan ay inililista kasama ang `--name` command na gagamitin. Kapag walang `CLAUDE_CONFIG_DIR`, iniimbak ng Claude Code ang `.claude.json` sa iyong home folder katabi ng `~/.claude`, kaya kapag in-import mo ang `~/.claude` at wala itong sariling `.claude.json`, gagawa ang ccam ng link na `~/.claude/.claude.json` papunta sa `~/.claude.json`. Sa ganitong paraan iisang account ang nakikita ng karaniwang `claude` at ng `ccam main`. Hindi kasapi sa anumang group ang mga in-import na account hanggang idagdag mo sila sa isa (`ccam group join <group> <name>`), at walang shortcut na ginagawa.

Nagpi-print ang `ccam completion bash|zsh` ng completion script na nagko-complete rin ng mga pangalan ng account.

## Kaligtasan at mga espesyal na kaso

- Tumatanggi ang `ccam rename` at `ccam remove` habang tumatakbo ang account (`refuseRunning`), na may `work is running; quit its Claude Code session first`, dahil patuloy na nagsusulat ang session sa folder.
- Sa isang in-import na account, inaalis lang ng `ccam remove` ang rehistro nito: `✓ main unregistered (folder ... left in place)`. Hindi inililipat sa trash ang iyong `~/.claude`.
- Tinatanggihan ng `ccam add` ang pangalan na ang folder ay mayroon na sa `profiles/`, at itinuturo ang `ccam import`.
- Tinatanggihan ng `ccam import` ang folder na nakarehistro na, ang path na hindi folder, at ang pangalan na gamit na.
- Hindi sinusuri ng `ccam logout` at `ccam login` kung tumatakbo ang account.
- Kapag hindi maabot ng relay ang Claude Code, halimbawa pagkatapos matapos ang login, ipi-print nito ang `ccam: could not pass the sign-in to Claude Code:` kasama ang dahilan. Kadalasan, ang HTTP error ay nangangahulugang nag-expire na ang code; simulan ulit ang login.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Mga shortcut command](/features/shortcuts)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Maraming Claude Code account](/guides/multiple-claude-code-accounts)
