---
title: Limites de uso
description: ccam usage mostra os limites de 5 horas e semanais de todas as contas de uma vez, e ccam statusline coloca o nome da conta e o uso no Claude Code.
updated: 2026-10-20
order: 7
translation: machine
sourceHash: f08a2fdc6f0b
---

`ccam usage` mostra quanto dos limites de 5 horas e semanal cada conta já usou, com o tempo até cada janela reiniciar. Use-o antes de começar uma tarefa longa, ou quando o Claude Code avisar que um limite está próximo, para ver qual conta ainda tem espaço. `ccam statusline` é o comando de linha de status do Claude Code: ele mostra o nome da conta e o uso dentro do Claude Code e salva os números para o `ccam usage`.

## Uso

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Em um terminal, a tabela aparece imediatamente com `fetching…` em cada linha, e cada linha é redesenhada quando a conta dela responde. `As of` indica a idade dos números; `(cached)` significa que eles vêm de um arquivo salvo, não do endpoint. A última linha indica a conta com o menor uso de 5 horas abaixo de 100% e só aparece quando você tem mais de uma conta. Uma conta sem login mostra `not logged in → ccam login <name>`, e uma conta com chave de API mostra `no limits (API)`.

Depois de `ccam statusline install`, o Claude Code mostra uma linha como `[work] · Opus · 5h 24% · 7d 41%`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--json` | Imprime os relatórios em JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | desativado |
| `--watch[=N]` | Redesenha a cada `N` segundos até você pressionar Ctrl+C. `--watch` sozinho, ou um valor menor ou igual a 0, significa 60 | desativado |
| `-v`, `--verbose` | Mostra os limites semanais por modelo e uma linha `source:` (`live`, `cache` ou `statusline`) | desativado |
| `--no-cache` | Ignora o cache de 60 segundos e consulta o endpoint de novo | desativado |
| `--source live\|cache` | `cache` nunca consulta o endpoint e mostra apenas os números salvos; `live` é o comportamento padrão | live |
| `--refresh` | Renova os tokens de acesso expirados, só das contas que não estão em execução | desativado |
| `--dry-run` | Aceito: o `ccam usage` não altera nada além do próprio cache. Recusado junto com `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | desativado |

Qualquer outro valor de `--source` falha com `--source must be live or cache`. `ccam usage` não aceita argumentos.

`ccam statusline` lê da entrada padrão o JSON da linha de status do Claude Code. A única opção dele é `--wrap <command>`: um comando de linha de status existente, cuja primeira linha de saída recebe o nome da conta na frente. `ccam statusline install` não aceita opções. Os dois recusam `--dry-run` com `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Como funciona

Os números vêm de dois lugares: o endpoint de uso não documentado da Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, chamado com o login de cada conta, e os valores que o Claude Code passa ao comando de linha de status. Tudo isso fica em `internal/usage`.

1. `cli.usageCmd` carrega o `config.toml` e monta um `usage.Account` por conta: a pasta dela `~/.local/share/ccam/profiles/<name>/`, a pasta de sessões (para saber se ela está em execução) e se ela foi importada no lugar.
2. `Fetcher.FetchAll` inicia uma goroutine por conta. Todas compartilham um único `http.Client`, então as conexões com o endpoint são reaproveitadas. A conexão pode levar 2 s, e cada conta tem 3 s no total. O cliente não segue redirecionamentos, já que um redirecionamento reenviaria o token para outro lugar, e lê no máximo 1 MB da resposta.
3. À medida que cada conta termina, a linha dela é redesenhada. Uma conta que falha ou não responde a tempo passa a usar os números salvos, então uma conta lenta nunca atrasa as outras.

Cada conta consulta o endpoint na própria goroutine; a que estoura o tempo mostra os números salvos.

<!-- motion:usage-fetch -->

`personal` e `work` respondem a tempo e mostram números atuais. `side` leva mais de 3 s, então `fallback()` lê os dois arquivos de cache dela e usa o gravado por último. `usage.Table()` desenha as linhas e `usage.Best()` escolhe a conta indicada na última linha.

Para uma conta, `Fetcher.Fetch` faz estas verificações em ordem:

<!-- motion:usage-cache -->

1. Com `--source cache`, vai direto para `fallback()`.
2. A menos que `--no-cache` esteja definido, um resultado atual em `~/.local/state/ccam/usage/<name>.json` com menos de 60 segundos (`CacheTTL`) é retornado como está. Executar `ccam usage` duas vezes seguidas custa uma solicitação por conta.
3. `readCredentials()` lê o `.credentials.json` (ou o Acesso às Chaves do macOS) sem alterá-lo, e só a chave `claudeAiOauth` é decodificada. Sem um login de assinatura, o resultado é `not logged in`, ou `no limits (API)` quando o `.claude.json` contém uma chave de API.
4. Um token de acesso expirado não é usado. Sem `--refresh`, a conta mostra os números salvos marcados como expirados.
5. Caso contrário, o ccam chama o endpoint. Um 401 ou 403 conta como login expirado, qualquer outro status ou um corpo que ele não consegue ler conta como erro, e nos dois casos ele usa os números salvos.
6. Uma resposta válida é gravada em `<name>.json` com o horário em que foi medida.

`fallback()` lê `<name>.json` e `<name>.statusline.json` e fica com o mais recente. Antes de mostrar números salvos, ele verifica o horário de reinício de cada janela: uma janela que reiniciou depois que os números foram salvos mostra 0%, porque a porcentagem antiga não vale mais. É por isso que `personal`, acima, mostra `0%` sem horário de reinício.

O endpoint tem dois formatos de resposta. O mais novo tem um array `limits` com as entradas `session` (5 horas), `weekly_all` e `weekly_scoped` (por modelo). O mais antigo tem os objetos `five_hour`, `seven_day`, `seven_day_opus` e `seven_day_sonnet`. `ParseLive` lê os dois: o array `limits` prevalece, e os objetos antigos completam o que ele não cobre. Um corpo sem nenhum limite é informado como mudança de formato, e o `ccam doctor` avisa sobre isso.

### Renovar tokens

Os tokens de renovação são rotativos: renovar um gera um novo token de renovação e invalida o antigo. Se o ccam renovasse o login enquanto o Claude Code estivesse em execução com a mesma conta, essa sessão seria desconectada. Por isso, o `ccam usage` lê os tokens e nunca os grava, a menos que você passe `--refresh`. Nesse caso, `Fetcher.refresh`:

1. recusa as contas registradas no lugar com `ccam import`, já que o `claude` simples pode usar esse login sem o ccam saber;
2. recusa os logins armazenados no Acesso às Chaves do macOS;
3. obtém um arquivo de bloqueio, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (em `~/.local/state/ccam/locks/` quando essa variável não está definida), para que dois processos do ccam nunca renovem o mesmo token duas vezes;
4. com o bloqueio, recusa se uma sessão da conta estiver em execução e pula a solicitação se outro processo já tiver renovado o token;
5. envia o token de renovação ao endpoint de tokens e grava os novos tokens de volta no `.credentials.json` com `fsutil.AtomicWrite`, mantendo todas as outras chaves do arquivo.

### O coletor da linha de status

O Claude Code executa o comando de linha de status e envia a ele um documento JSON por pipe. O ccam lê dele o nome do modelo e `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` encontra a conta por `CCAM_PROFILE`, que o ccam define ao iniciar o Claude Code, ou então comparando `CLAUDE_CONFIG_DIR` com as pastas das contas.
2. `ParseStatusLine` lê `rate_limits.five_hour` e `rate_limits.seven_day`. O JSON não tem limites de taxa para contas com chave de API, nem antes da primeira resposta da sessão; nesses casos, nada é salvo.
3. `SaveStatusLine` grava `~/.local/state/ccam/usage/<name>.statusline.json`. O `ccam usage` o usa quando o endpoint não responde.
4. O ccam imprime `[name] · model · 5h N% · 7d N%`. Com `--wrap`, ele executa o seu comando com `sh -c` e a mesma entrada e, em vez disso, imprime `[name]` seguido da primeira linha do comando. No Windows, ele executa o comando como o Claude Code faz: com o Git Bash, se estiver instalado, senão com o PowerShell.
5. Quando o uso de 5 horas ou o semanal está em 90% ou mais e existe outra conta, a linha termina com `→ ccam next --auto --continue`, por exemplo `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

O comando encapsulado tem 2 segundos. Depois disso, ele é encerrado e a linha mostra só o nome da conta, então um comando lento não consegue travar a linha de status do Claude Code.

`ccam statusline install` grava `statusLine` no `settings.json` de cada conta. Quando uma conta já tem um comando de linha de status, o ccam o mantém como `--wrap '<old command>'`. No Windows, o caminho do ccam é gravado com barras normais e `~` para a pasta pessoal (`~/.local/bin/ccam.exe statusline`), que tanto o Git Bash quanto o PowerShell entendem. Um `settings.json` compartilhado por um grupo é gravado uma única vez, pelo link, então todas as contas do grupo o recebem. Executá-lo de novo imprime `(already set)` para os arquivos que já usam o ccam.

## Segurança e casos especiais

- O `ccam usage` nunca altera um arquivo de token, a menos que você passe `--refresh`. Quando um login expirar, inicie a conta uma vez com `ccam <name>` e o Claude Code o renova.
- `--refresh` nunca mexe em contas em execução, contas importadas nem logins do Acesso às Chaves. Uma renovação recusada mostra os números salvos e o motivo.
- O endpoint de uso não é documentado publicamente e pode mudar. Quando isso acontecer, você continua vendo os números salvos, e o `ccam doctor` informa a mudança de formato.
- Sem números salvos e sem resposta, a linha mostra `lookup failed:` e o motivo, por exemplo `no saved usage` com `--source cache`.
- Os arquivos de cache guardam só porcentagens e horários de reinício, nunca tokens.

## Relacionados

- [Limites de uso e troca de contas](/docs/usage-limits)
- [Trocar para a próxima conta](/features/next)
- [Fluxo de trabalho para limites de uso](/guides/usage-limit-workflow)
- [Comandos](/docs/commands)
