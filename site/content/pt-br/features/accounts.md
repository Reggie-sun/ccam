---
title: Gerenciar contas
description: Adicione, faça login, liste, renomeie, remova e importe contas do Claude Code com ccam add, ccam login, ccam list, ccam rename, ccam remove e ccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

Cada conta do ccam é uma pasta de configuração separada do Claude Code, com o próprio login. Estes comandos criam, conectam, listam, renomeiam e removem contas, e registram pastas que você já tem.

## Uso

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

Sempre que um comando aceita `<name|number>`, você pode informar o nome da conta ou o número do slot dela, a coluna `#` do `ccam list`.

`ccam add` em um terminal, sem opções, abre um assistente. Ele pede o nome e como a conta faz login (um login do Claude agora, uma chave de API ou um login do Claude depois), mostra os grupos em que entrar, com `default` marcado (antes de haver grupos, os itens que o grupo `default` vai compartilhar), pergunta `Create a shortcut command for work? [y/N]` e termina com um resumo para confirmar. `Esc` volta uma etapa. Qualquer uma das opções `--group`, `--no-group`, `--shortcut` ou `--no-login`, ou um pipe no lugar de um terminal, pula o assistente; nesse caso o nome é obrigatório e nenhum atalho é criado, a menos que `--shortcut` indique um.

Sem `--group`, uma conta nova entra no grupo `default`. A primeira conta o cria, compartilhando servidores MCP, configurações, `CLAUDE.md` e extensões; veja [compartilhamento com grupos](/features/share).

`ccam list` mostra o slot, o nome, o atalho, o e-mail, os grupos e o último uso, com `*` antes da conta padrão:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

Com `--json`, cada conta é um objeto com `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (omitido quando a conta nunca foi iniciada), `groups` e `share`. `share` associa cada item (`mcp`, `settings`, `memory`, `extensions`, `history`) ao grupo do qual ele vem, ou a `"local"`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `ccam add --no-login` | Cria a conta sem fazer login; faça login depois com `ccam login <name>` | desativado |
| `ccam add --group <group>` | Entra neste grupo em vez de `default`. Pode ser repetida; `default` é criado quando não existe | `default` |
| `ccam add --no-group` | Não entra em nenhum grupo: a conta não compartilha nada | desativado |
| `ccam add --shortcut <command>` | Também cria este atalho em `~/.local/bin` | nenhum |
| `ccam add --api --key-from <source>` | Uma conta que usa uma chave de API em vez de um login; veja [Contas com chave de API](/features/api) | desativado |
| `ccam list --json` | Imprime as contas em JSON | desativado |
| `ccam remove -y`, `--yes` | Remove sem perguntar | desativado |
| `ccam import --name <name>` | Nome da conta para a pasta informada | derivado do nome da pasta |
| `ccam import -y`, `--yes` | Importa todas as pastas encontradas sem perguntar | desativado |
| `--dry-run` | Respeitado por `ccam list` e `ccam completion`, que não alteram nada. Recusado por `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, com `` --dry-run is not supported by `ccam add`; nothing was changed `` (com o nome do comando) | desativado |

`ccam login`, `ccam logout`, `ccam use` e `ccam rename` não têm opções próprias. `ccam completion` aceita `bash`, `zsh` ou `powershell`; qualquer outro valor falha com `only bash, zsh and powershell are supported`.

Os nomes de conta têm de 1 a 32 caracteres entre `a-z`, `0-9`, `-` e `_`, começando por uma letra ou um dígito. Um nome não pode ter só dígitos, porque `ccam 2` significa o slot 2. Um nome também não pode ser um comando do ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), então `ccam <name>` nunca é ambíguo. A lista é `config.Reserved`, em `internal/config`.

## Como funciona

`ccam add work --shortcut ccwork` executa estas etapas em `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` verifica o nome, e `shim.Check` verifica o atalho antes que qualquer coisa seja criada.
2. `profile.Add` cria `~/.local/share/ccam/profiles/work/` com o modo `0700` e dá à conta o menor número de slot livre (`Config.NextSlot`), então o número de uma conta removida é reaproveitado. A primeira conta se torna a padrão.
3. A conta entra nos grupos dela (`default`, a menos que `--group` ou `--no-group` diga outra coisa; a primeira conta cria `default`), e `applyShare` vincula a pasta nova e vazia a eles com `share.Build` e depois `Apply`, por exemplo `settings.json` a `groups/default/settings.json`. Se isso falhar, a pasta nova é removida, porque uma pasta que sobrasse bloquearia todas as novas tentativas com `folder already exists`.
4. A conta é salva em `~/.config/ccam/config.toml` com `fsutil.AtomicWrite`: um arquivo temporário, `fsync` e depois `rename`, para que uma falha nunca deixe uma configuração pela metade.
5. `runLogin` inicia `claude auth login` com `CLAUDE_CONFIG_DIR` apontando para a pasta da conta e `--email` com o último e-mail da conta, para que a página de login venha preenchida. Se o login falhar, a conta permanece, e o ccam imprime o erro e pede que você execute `ccam login work` depois.
6. `profile.ReadIdentity` extrai `oauthAccount` de `profiles/work/.claude.json`: o e-mail, os ids da conta e da organização, o nome da organização e o plano (`claude_max`, `claude_team` …). Esse arquivo pode ter vários megabytes, então `TopLevelValue` procura só essa chave em vez de analisar o arquivo inteiro. O ccam guarda esses dados em `[profile.login]` no `config.toml`.
7. `shim.Install` cria o link de `~/.local/bin/ccwork` para o binário do ccam, que procura `argv[0]` no `config.toml` ao iniciar.

Um login pode levar minutos, e outro ccam pode salvar a configuração nesse meio-tempo. Por isso, os comandos que esperam salvam por meio de `env.update`, que recarrega o `config.toml`, aplica a única alteração (aqui, o e-mail) e salva. `ccam remove` recarrega depois da pergunta pelo mesmo motivo.

`ccam login` executa as etapas 5 e 6 para uma conta existente. O Claude Code usa o e-mail só como sugestão para a página de login, então um login pode acabar em outra conta Claude que não aquela em que a conta do ccam estava; o ccam compara o novo login com `[profile.login]` e avisa (`ccam: this login is other@company.com, but work was dev@company.com before.`). Ele mantém o novo login e diz como desfazê-lo: `ccam logout work`, depois `ccam login work` e o login com o e-mail certo. Quando outra conta do ccam tem login na mesma conta Claude, o ccam avisa que as duas compartilham os limites de uso. `ccam login --expired` faz login, uma depois da outra, em todas as contas que nunca fizeram login ou cujo login (o token de renovação) expirou.

`ccam logout` executa `claude auth logout` na conta e apaga o e-mail salvo; o registro do login permanece, então o próximo login é comparado com ele. `ccam use` apenas reescreve `default` no `config.toml`.

`ccam list` acrescenta uma coluna `Organization` quando as contas pertencem a mais de uma organização (planos pessoais mostram `personal`), e `ccam list --json` tem `plan` e `org` para cada conta.

### Fazer login em uma máquina remota

Por SSH, a parte do navegador de `claude auth login` pode falhar. O Claude Code espera o login em um servidor de retorno em `http://localhost:<port>/callback` no servidor, mas o navegador que conclui o login roda no seu computador, onde esse endereço não leva a lugar nenhum, e o navegador mostra um erro como `ERR_CONNECTION_REFUSED`. O endereço ainda contém o código de uso único, então o ccam pode entregá-lo:

<!-- motion:accounts-relay -->

`runLogin` não entrega o seu terminal ao Claude Code. Ele conecta a entrada padrão do Claude Code a um pipe, e `relayLogin` copia para ele cada linha que você digita, então um código colado no prompt `Paste code here` do Claude Code continua chegando. Quando uma linha é um endereço de retorno, `callbackURL` só a aceita se o host for `localhost`, `127.0.0.1` ou `[::1]`, houver uma porta, o caminho for `/callback` e `code` e `state` estiverem presentes; `http://` pode ser omitido. Em seguida, `deliverCallback` solicita esse endereço no servidor com um limite de 10 segundos e sem seguir redirecionamentos, como o navegador teria feito. A solicitação chega ao servidor de retorno do `claude auth login` que está esperando (a linha tracejada no diagrama). O Claude Code recebe o código e salva o login, no Linux em `.credentials.json` na pasta da conta. Uma linha colada nunca faz o ccam contatar outro host.

O ccam imprime isto antes de o login começar:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Se der certo, ele imprime `→ Passed the sign-in to Claude Code on this machine.`

Enquanto o login espera uma entrada em um terminal, o ccam desliga o eco e o modo de linha do terminal e mostra um `*` para cada caractere que você digita, então um código ou endereço de retorno colado não aparece na tela. `readLine` faz a edição da linha por conta própria: `Enter` envia a linha, `Backspace` apaga um caractere e `Ctrl+D` em uma linha vazia encerra a entrada. O terminal é restaurado quando o login termina, e também com `Ctrl+C`, `SIGTERM` ou `SIGHUP`, que continuam encerrando o ccam. Sem um terminal na entrada padrão e na saída de erro padrão, as linhas são copiadas conforme chegam.

### Renomear, remover e importar

`ccam rename` e `ccam remove` atuam na pasta da conta e no que o ccam guarda com o nome da conta em outros lugares. `nameState` lista esses caminhos: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` e os caches de uso `~/.local/state/ccam/usage/<name>.json` e `<name>.statusline.json`. `ccam remove work` os trata nesta ordem:

<!-- motion:accounts-remove -->

`ccam remove` primeiro recusa se a conta estiver em execução e depois pergunta. Se a conta compartilha extensões, `RewritePluginPaths` reescreve o `installed_plugins.json` e o `known_marketplaces.json` do grupo, porque o Claude Code armazena caminhos absolutos de plugins que podem passar pela pasta desta conta. `profile.Trash` move a pasta para `~/.local/share/ccam/trash/<name>-<time>/`, a conta sai do `config.toml`, os caminhos de `nameState` são excluídos para que uma conta futura com o mesmo nome não os herde, e o link do atalho é removido. Para recuperar, mova a pasta de volta da lixeira.

`ccam rename` renomeia a pasta, reescreve os caminhos de plugins em `plugins/` da pasta antiga para a nova, mantém a conta como padrão se ela era e move os caminhos de `nameState` para o novo nome. O atalho mantém o nome. Se a reescrita dos plugins falhar, a pasta volta ao nome anterior, já que a configuração não é salva.

`ccam import` registra uma pasta no lugar onde ela está, como `dir` no `config.toml`. Ele nunca a move nem a copia: um token de renovação copiado faria as cópias desconectarem umas às outras. Sem uma pasta como argumento, ele procura pastas `~/.claude` e `~/.claude-*` não registradas e pergunta sobre cada uma; os nomes vêm da pasta (`~/.claude` vira `main`, `~/.claude-work` vira `work`, `~/.claude-2` vira `acct2`), e uma pasta cujo nome não pode ser derivado é listada com o comando `--name` a usar. Sem `CLAUDE_CONFIG_DIR`, o Claude Code mantém o `.claude.json` na sua pasta pessoal, ao lado de `~/.claude`, então, quando você importa `~/.claude` e ela não tem um `.claude.json` próprio, o ccam cria o link `~/.claude/.claude.json` para `~/.claude.json`. Assim, o `claude` simples e o `ccam main` enxergam a mesma conta. Contas importadas não ficam em nenhum grupo até você adicioná-las a um (`ccam group join <group> <name>`), e nenhum atalho é criado.

`ccam completion bash|zsh` imprime um script de autocompletar que também completa nomes de contas.

## Segurança e casos especiais

- `ccam rename` e `ccam remove` se recusam enquanto a conta estiver em execução (`refuseRunning`), com `work is running; quit its Claude Code session first`, porque a sessão continua gravando na pasta.
- `ccam remove` em uma conta importada apenas cancela o registro: `✓ main unregistered (folder ... left in place)`. O seu `~/.claude` não vai para a lixeira.
- `ccam add` recusa um nome cuja pasta já existe em `profiles/` e indica `ccam import`.
- `ccam import` recusa uma pasta que já está registrada, um caminho que não é uma pasta e um nome que já está em uso.
- `ccam logout` e `ccam login` não verificam se a conta está em execução.
- Se o repasse não conseguir alcançar o Claude Code, por exemplo depois que o login terminou, ele imprime `ccam: could not pass the sign-in to Claude Code:` com o motivo. Um erro HTTP geralmente significa que o código expirou; comece o login de novo.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Atalhos](/features/shortcuts)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Várias contas do Claude Code](/guides/multiple-claude-code-accounts)
