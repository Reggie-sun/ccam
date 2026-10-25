---
title: Gerenciar contas
description: Adicione, faça login, liste, renomeie, remova e importe contas do Claude Code com ccam add, ccam login, ccam list, ccam rename, ccam remove e ccam import.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "1d40ce0280e0"
---

Cada conta do ccam é uma pasta de configuração separada do Claude Code, com o próprio login. Estes comandos criam, conectam, listam, renomeiam e removem contas, e registram pastas que você já tem.

## Uso

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Sempre que um comando aceita `<name|number>`, você pode informar o nome da conta ou o número do slot dela, a coluna `#` do `ccam list`.

`ccam add` em um terminal, sem opções, abre um assistente. Ele pede o nome e como a conta faz login (um login do Claude agora, uma chave de API ou um login do Claude depois), mostra os grupos em que entrar, com `default` marcado (antes de haver grupos, os itens que o grupo `default` vai compartilhar), pergunta `Create a shortcut command for work? [y/N]` e termina com um resumo para confirmar. `Esc` volta uma etapa. Qualquer uma das opções `--group`, `--no-group`, `--shortcut` ou `--no-login`, ou um pipe no lugar de um terminal, pula o assistente; nesse caso o nome é obrigatório e nenhum atalho é criado, a menos que `--shortcut` indique um.

Sem `--group`, uma conta nova entra no grupo `default`. A primeira conta o cria, compartilhando servidores MCP, configurações, `CLAUDE.md` e extensões; veja [compartilhamento com grupos](/features/share). No Windows sem o Modo de Desenvolvedor, sem o qual não é possível criar links para arquivos, a primeira conta cria `default` só com `mcp` e `extensions` e avisa (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), e `--group` para um grupo que compartilha configurações, memória ou histórico é recusado.

`ccam list` mostra o slot, o nome, o atalho, o e-mail, os grupos e o último uso, com `*` antes da conta padrão. O e-mail é aquele com que a pasta está logada agora (do `.claude.json` dela); quando ele difere do que o ccam registrou, o registrado aparece depois como `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Com `--json`, cada conta é um objeto com `slot`, `shortcut`, `name`, `email` (o do login atual, como na tabela), `default`, `dir`, `last_used` (omitido quando a conta nunca foi iniciada), `groups` e `share`. `share` associa cada item (`mcp`, `settings`, `memory`, `extensions`, `history`) ao grupo do qual ele vem, ou a `"local"`.

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
| `ccam logout -y`, `--yes` | Faz logout sem perguntar | desativado |
| `ccam import --name <name>` | Nome da conta para a pasta informada; recusado sem uma pasta (`--name needs a folder: ccam import <folder> --name ops`) | derivado do nome da pasta |
| `ccam import -y`, `--yes` | Importa todas as pastas encontradas sem perguntar | desativado |
| `--dry-run` | Respeitado por `ccam list` e `ccam completion`, que não alteram nada. Recusado por `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, com `` --dry-run is not supported by `ccam add`; nothing was changed `` (com o nome do comando) | desativado |

`ccam login` (além de `--expired`), `ccam use` e `ccam rename` não têm opções próprias. `ccam completion` aceita `bash`, `zsh` ou `powershell`; qualquer outro valor falha com `only bash, zsh and powershell are supported`.

Os nomes de conta têm de 1 a 32 caracteres entre `a-z`, `0-9`, `-` e `_`, começando por uma letra ou um dígito. Um nome não pode ter só dígitos, porque `ccam 2` significa o slot 2. Um nome também não pode ser um comando do ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agents`, `attach`, `ccusage`), então `ccam <name>` nunca é ambíguo. A lista é `config.Reserved`, em `internal/config`.

## Como funciona

`ccam add work --shortcut ccwork` executa estas etapas em `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` verifica o nome, e `shim.Check` verifica o atalho antes que qualquer coisa seja criada.
2. `profile.Add` cria `~/.local/share/ccam/profiles/work/` com o modo `0700` e dá à conta o menor número de slot livre (`Config.NextSlot`), então o número de uma conta removida é reaproveitado. A primeira conta se torna a padrão.
3. A conta entra nos grupos dela (`default`, a menos que `--group` ou `--no-group` diga outra coisa; a primeira conta cria `default`), e `applyShare` vincula a pasta nova e vazia a eles com `share.Build` e depois `Apply`, por exemplo `settings.json` a `groups/default/settings.json`. Se isso falhar, a pasta nova é removida, porque uma pasta que sobrasse bloquearia todas as novas tentativas com `folder already exists`.
4. A conta é salva em `~/.config/ccam/config.toml` com `fsutil.AtomicWrite`: um arquivo temporário, `fsync` e depois `rename`, para que uma falha nunca deixe uma configuração pela metade.
5. `runLogin` inicia `claude auth login` com `CLAUDE_CONFIG_DIR` apontando para a pasta da conta e `--email` com o último e-mail da conta, para que a página de login venha preenchida. Se o login falhar, a conta permanece, e o ccam imprime o erro e pede que você execute `ccam login work` depois.
6. `profile.ReadIdentity` extrai `oauthAccount` de `profiles/work/.claude.json`: o e-mail, os ids da conta e da organização, o nome da organização e o plano (`claude_max`, `claude_team` …). Esse arquivo pode ter vários megabytes, então `TopLevelValue` procura só essa chave em vez de analisar o arquivo inteiro. O ccam guarda esses dados em `[profile.login]` no `config.toml`.
7. `shim.Install` cria o link de `~/.local/bin/ccwork` para o binário do ccam, que procura `argv[0]` no `config.toml` ao iniciar, e registra o link em `~/.local/state/ccam/shortcuts.json`. Se o link não puder ser criado, a conta é adicionada mesmo assim, o ccam diz por quê (`Could not create shortcuts: …`) e o resumo mostra `shortcut: ccwork (not created; see above)`.

Um login pode levar minutos, e outro ccam pode salvar a configuração nesse meio-tempo. Por isso, os comandos que esperam salvam por meio de `env.update`, que recarrega o `config.toml`, aplica a única alteração (aqui, o e-mail) e salva. `ccam remove` recarrega depois da pergunta pelo mesmo motivo.

`ccam login` executa as etapas 5 e 6 para uma conta existente. O Claude Code usa o e-mail só como sugestão para a página de login, então um login pode acabar em outra conta Claude que não aquela em que a conta do ccam estava; o ccam compara o novo login com `[profile.login]` e avisa (`ccam: this login is other@company.com, but work was dev@company.com before.`). Ele mantém o novo login e diz como desfazê-lo: ``To undo it, run `ccam login work` again and sign in as dev@company.com.`` Um login em outra conta Claude também descarta o [registro de limite](/features/usage#the-limit-record) da conta, já que aqueles limites eram da outra conta. Quando outra conta do ccam tem login na mesma conta Claude, o ccam avisa que as duas compartilham os limites de uso. `ccam login --expired` faz login, uma depois da outra, em todas as contas que nunca fizeram login ou cujo login expirou: o token de renovação expirou, ou o Claude Code apagou o login depois que o servidor o recusou. Antes de o primeiro navegador abrir, ele verifica se o `config.toml` pode ser salvo.

`ccam logout` executa `claude auth logout` na conta e apaga o e-mail salvo; o registro do login permanece, então o próximo login é comparado com ele. O Claude Code revoga o login e exclui o que guarda para a conta, inclusive os logins dos servidores MCP. Por isso, o ccam primeiro diz o que vai junto e pergunta (`--yes` pula a pergunta): quando o login tem logins de MCP, quando a conta está em execução (as sessões dela também são desconectadas) e quando a pasta de outra conta tem o mesmo login (uma pasta copiada), que o logout também revoga:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Para fazer login como outra conta Claude, basta `ccam login <name>`, e ele mantém os logins de MCP. `ccam use` apenas reescreve `default` no `config.toml`.

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

`ccam rename` e `ccam remove` atuam na pasta da conta e no que o ccam guarda com o nome da conta em outros lugares. `nameState` lista esses caminhos: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, os caches de uso `~/.local/state/ccam/usage/<name>.json` e `<name>.statusline.json`, o registro de limite `<name>.limit.json` e as notas de sessão em `~/.local/state/ccam/handoff/<name>/`. `ccam remove work` os trata nesta ordem:

<!-- motion:accounts-remove -->

`ccam remove` primeiro recusa se a conta estiver em execução, depois pergunta e verifica de novo após a resposta (uma sessão iniciada enquanto a pergunta esperava é detectada). Ele verifica se o `config.toml` pode ser salvo antes de alterar qualquer coisa. Se a conta compartilha extensões, `RewritePluginPaths` reescreve o `installed_plugins.json` e o `known_marketplaces.json` do grupo, porque o Claude Code armazena caminhos absolutos de plugins que podem passar pela pasta desta conta. Em seguida, uma conta gerenciada é desconectada: `claude auth logout` roda na pasta dela, para que nenhum login válido fique esperando na lixeira (ou, no macOS, em um item do Acesso às Chaves com o nome da pasta, que uma conta futura com esse nome herdaria). Isso é pulado para uma conta com chave de API e quando a pasta de outra conta tem o mesmo login, que o logout também revogaria; um logout que falha é informado, e a remoção continua. `profile.Trash` move a pasta para `~/.local/share/ccam/trash/<name>-<time>/`, a conta sai do `config.toml`, os caminhos de `nameState` são excluídos para que uma conta futura com o mesmo nome não os herde, e o link do atalho é removido: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`. Se, afinal, o `config.toml` não puder ser salvo, a pasta volta da lixeira. Para recuperar uma conta removida, mova a pasta de volta da lixeira e faça login de novo.

`ccam rename` renomeia a pasta, reescreve os caminhos de plugins em `plugins/` da pasta antiga para a nova, mantém a conta como padrão se ela era e move os caminhos de `nameState` para o novo nome; as notas de sessão são atualizadas para que `ccam next --continue` ainda as encontre. O atalho mantém o nome. Se a reescrita dos plugins falhar, ou se o `config.toml` não puder ser salvo, a pasta volta ao nome anterior. No macOS, a pasta mantém o nome antigo, já que o Acesso às Chaves guarda o login e os logins de MCP pelo caminho da pasta, e uma pasta movida ficaria desconectada; o mesmo vale no Windows quando o Claude Code guarda os logins no Gerenciador de Credenciais. A configuração então registra a pasta como `folder` ao lado do novo nome, e o ccam diz `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` registra uma pasta no lugar onde ela está, como `dir` no `config.toml`. Ele nunca a move nem a copia: um token de renovação copiado faria as cópias desconectarem umas às outras. Sem uma pasta como argumento, ele procura pastas `~/.claude` e `~/.claude-*` não registradas e pergunta sobre cada uma; os nomes vêm da pasta (`~/.claude` vira `main`, `~/.claude-work` vira `work`, `~/.claude-2` vira `acct2`), e uma pasta cujo nome não pode ser derivado, ou já é o de uma conta, é listada com o comando `--name` a usar. Uma pasta encontrada que falha é informada como `Skipped <path>: <reason>`, e as outras são importadas. As respostas são coletadas primeiro, e as pastas são registradas no `config.toml` como ele estiver nesse momento, então o que outros comandos do ccam salvaram nesse meio-tempo permanece; as linhas `✓` vêm depois de salvar. Sem `CLAUDE_CONFIG_DIR`, o Claude Code mantém o `.claude.json` na sua pasta pessoal, ao lado de `~/.claude`, então, quando você importa `~/.claude` e ela não tem um `.claude.json` próprio, o ccam cria o link `~/.claude/.claude.json` para `~/.claude.json`. Assim, o `claude` simples e o `ccam main` enxergam a mesma conta. Se `~/.claude` tem um `.claude.json` próprio além de `~/.claude.json`, o ccam deixa os dois e avisa: `ccam main` usaria o de dentro e o `claude` simples o de fora, mantendo separados os dados de login e os servidores MCP deles, então substitua o arquivo de dentro por um link para ficar com um só (o [doctor](/features/doctor) também informa isso). Quando a pasta tem o mesmo login que uma conta registrada (uma pasta copiada), o ccam avisa: a que renovar primeiro desconecta a outra, então execute `ccam login <name>` antes de usá-la. Contas importadas não ficam em nenhum grupo até você adicioná-las a um (`ccam group join <group> <name>`), e nenhum atalho é criado.

`ccam completion bash|zsh` imprime um script de autocompletar que também completa nomes de contas.

## Segurança e casos especiais

- `ccam rename` e `ccam remove` se recusam enquanto a conta estiver em execução (`refuseRunning`), com ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, porque a sessão continua gravando na pasta. A verificação cobre o Claude Code seja qual for a forma como foi iniciado: pelo ccam, pelo `claude` simples, por uma IDE, por sessões em segundo plano e pelo daemon delas.
- `ccam rename`, `ccam remove` e `ccam logout` se recusam enquanto outro ccam renova ou altera aquele login: ``another ccam is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Enquanto uma troca de grupos roda em outro ccam, `ccam rename` e `ccam remove` se recusam com `another ccam is changing the groups right now; try again when it is done`; depois de uma que parou no meio, com ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` em uma conta importada mantém a pasta e o login dela, para o `claude` simples. Se a conta está em grupos, ela primeiro sai deles como faz o `ccam group leave`: cada item compartilhado vira, na pasta, uma cópia do que está no grupo, e os servidores MCP do grupo vão para o `.claude.json` dela. Depois o registro é cancelado: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Os backups do `.claude.json` dela que o ccam fez vão para `<folder>/.ccam-backup/mcp-<time>`; se não puderem ser movidos, ficam em `~/.local/state/ccam/mcp/<name>/backups` e a mensagem avisa. O seu `~/.claude` não vai para a lixeira.
- Duas contas não podem ter a mesma pasta: um `config.toml` que indica uma pasta duas vezes falha ao carregar.
- `ccam add` recusa um nome cuja pasta já existe em `profiles/` e indica `ccam import`. No macOS, ele também para quando o Acesso às Chaves ainda guarda um login ou logins de MCP para o caminho da nova pasta, deixados por um ccam mais antigo, e indica o comando `security delete-generic-password` que os remove. Se o `config.toml` não puder ser salvo, a nova pasta é removida de novo.
- `ccam import` recusa uma pasta que já está registrada, um caminho que não é uma pasta e um nome que já está em uso. Uma cópia do `~/.claude` do `claude` simples também é percebida quando `~/.claude` não está registrado.
- `ccam logout` pergunta antes quando a conta está em execução; `ccam login` não verifica.
- Se o repasse não conseguir alcançar o Claude Code, por exemplo depois que o login terminou, ele imprime `ccam: could not pass the sign-in to Claude Code:` com o motivo. Um erro HTTP geralmente significa que o código expirou; comece o login de novo.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Atalhos](/features/shortcuts)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Várias contas do Claude Code](/guides/multiple-claude-code-accounts)
