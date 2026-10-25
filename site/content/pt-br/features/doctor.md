---
title: Verificar e reparar
description: Verifique links de compartilhamento, permissões de arquivos, logins e a versão do Claude Code com ccam doctor, e corrija o que é seguro com ccam doctor --fix.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` verifica a sua instalação do ccam: o binário do Claude Code, a pasta, o login e os links de compartilhamento de cada conta, os servidores MCP compartilhados, o endpoint de uso e os seus atalhos. `ccam doctor --fix` repara o que consegue sem perder dados e, antes, move para uma pasta de backup qualquer arquivo que substituir. Execute-o quando `ccam <name>` parar com um problema de link de compartilhamento, depois de restaurar arquivos manualmente ou depois de uma atualização do Claude Code.

## Uso

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

A primeira linha marca todas as contas de uma vez, para você ver onde estão os problemas antes de ler qualquer outra coisa. Abaixo dela, cada conta ganha uma caixa, e os resultados que não pertencem a nenhuma conta (Claude Code, variáveis de ambiente, plugins, servidores MCP, atalhos, regras de pasta, grupos e configurações, `auto_switch`) vão para uma caixa `General`. O título da caixa mostra a pior marca da conta e quantos problemas ela tem. Cada linha dentro dela é uma área verificada, sempre `Folder`, `Login` e `Share links` nessa ordem, depois as demais, com uma marca: `✔` OK, `ℹ` informação, `⚠` um aviso, `✘` um erro. Uma pasta sem nenhum problema mostra `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Com `--fix`, cada resultado reparado termina em `→ fixed` e deixa de contar como problema:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Em um terminal, as caixas se ajustam à largura dele e as mensagens longas quebram de linha, nas cores do próprio terminal (`NO_COLOR` as desativa). Uma caixa nunca tem menos de 29 colunas, então nenhuma palavra é cortada em um terminal estreito. Enviada para um pipe ou um arquivo, a saída não tem cores e nenhuma mensagem é quebrada ou cortada, então `ccam doctor | grep settings.json` encontra a linha inteira; cada caixa fica tão larga quanto a linha mais longa. A aba Health do painel mostra as mesmas caixas.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--fix` | Repara o que é seguro reparar (os arquivos alterados ganham backup) | desativado |
| `--dry-run` (global) | Executa todas as verificações, mas não repara nada, mesmo com `--fix` | desativado |

`ccam doctor` não aceita argumentos. Ele respeita `--dry-run`: `ccam --dry-run doctor --fix` imprime os resultados do `ccam doctor`, acrescenta o que `--fix` deixaria de lado e não altera nada:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

O código de saída é `1` quando sobra pelo menos um resultado `✘` sem reparo, e `0` nos outros casos. Avisos e informações nunca fazem a execução falhar. Quando a execução falha sem `--fix`, o ccam termina com uma linha em branco e `` `ccam doctor --fix` repairs what it can. ``. Uma configuração que não pode ser carregada também termina com `1` e um erro `ccam:`.

## Como funciona

`internal/cli/doctor.go` carrega `~/.config/ccam/config.toml` e chama `doctor.Run` em `internal/doctor/doctor.go`, que executa as verificações em uma ordem fixa e retorna os resultados, cada um com a conta a que pertence e a área que verificou. `tui.HealthReport` os agrupa nas caixas. `--fix` é repassado como `Options.Fix`, definido como false em uma simulação.

<!-- motion:doctor-checks -->

O bloco do meio é executado uma vez por conta. As verificações, e o que `--fix` faz em cada uma:

1. **Versão do Claude Code** (`checkClaude`). Localiza o binário `claude` como o inicializador faz e executa `claude --version`. Um binário ausente, um `claude_bin` que não é um arquivo executável ou um `--version` que falha é um erro. Sem reparo.
2. **Credenciais no ambiente** (`checkAuthEnv`). Avisa quando `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN` está definida: toda conta iniciada nesse shell a herda, e o Claude Code a prefere ao login da própria conta. Sem reparo; remova a variável você mesmo.
3. **Pasta** (`checkProfile`). Uma pasta de conta ausente é um erro e encerra as verificações dessa conta. Permissões com qualquer bit de grupo ou de outros definido são um aviso; `--fix` executa `chmod 0700`. Um `~/.claude` importado com um `.claude.json` próprio além de `~/.claude.json` é informação: `ccam main` e o `claude` simples passam a manter dados de login e servidores MCP separados. O Windows não tem bits de permissão, então lá esta verificação e a de `0600` abaixo são puladas; a pasta no seu perfil de usuário é protegida pela lista de controle de acesso dela.
4. **Arquivo de login** (`checkCredentials`). `.credentials.json` é verificado com `lstat`. Um link simbólico é um erro, porque cada conta precisa ter o próprio token, e não é reparado. Um modo diferente de `0600` é um aviso; `--fix` executa `chmod 0600`. Depois o login é lido (no macOS, primeiro do Acesso às Chaves). A falta de login, ou um arquivo de login sem login, é um aviso que indica `ccam login <name>`. Um token de renovação expirado, ou um login que o Claude Code apagou depois que o servidor o recusou, é o erro `` login expired → `ccam login <name>` ``; um token de acesso expirado é apenas informação, já que o Claude Code o renova na próxima inicialização. Um login guardado no Gerenciador de Credenciais do Windows é informação, já que o ccam não consegue lê-lo. Um login válido mostra o e-mail com que a pasta está logada agora e o plano (`login valid (dev@company.com, Claude Team · Acme)`). Nenhum desses é reparado; `ccam login --expired` faz login em todas as contas que precisam.
5. **Links de compartilhamento** (`share.Check`). Compara cada caminho de `share.Layout` com os grupos da conta, usando só `lstat` e `readlink`. Ele informa cinco tipos: um link ausente, um link que aponta para outro lugar, um arquivo comum onde deveria haver um link, um link que sobrou de um item que agora é local e um link correto cuja pasta de grupo está ausente (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; uma inicialização não verifica isso). Um link ausente e uma pasta de grupo ausente são avisos; os outros são erros. Uma cópia de volta do grupo que parou no meio também é informada: `--fix` coloca no lugar uma cópia preparada completa, mas, quando o caminho foi criado de novo desde então (por uma inicialização), ele indica a pasta com a cópia completa para você resolver. Os caminhos adicionados a um item em uma versão posterior do ccam (`file-history/`, `tasks/` e `paste-cache/` no histórico, `keybindings.json` nas configurações) são só um aviso enquanto a conta ainda tem a própria cópia: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Uma inicialização nunca para por causa deles.
6. **Entradas desconhecidas** (`checkUnknownEntries`). Lista os nomes de primeiro nível na pasta da conta que o ccam não conhece (`share.Known`), como um indício de que o layout do Claude Code mudou. Apenas informação.
7. **Caminhos de plugins** (`checkPlugins`). Os registros de plugins (`installed_plugins.json`, `known_marketplaces.json`) armazenam caminhos absolutos. No `plugins/` de cada grupo, caminhos que passam pela pasta de uma conta deixam de funcionar quando essa conta é removida ou renomeada. `--fix` os reescreve para a pasta do grupo com `profile.RewritePluginPaths`, uma gravação atômica do arquivo de registro.
8. **Servidores MCP** (`checkMCP`). Carrega o `mcp.json` de cada grupo uma vez. Um arquivo que não pode ser lido é um erro; um servidor que não passa em `mcp.Validate` é um aviso. Servidores no `~/.claude.json` do `claude` simples que nenhuma conta carrega são um aviso que indica o comando `ccam migrate --mcp` que os copia (com `--group <a|b>` quando vários grupos compartilham MCP). Servidores que só uma conta de um grupo tem são informação, com o comando `ccam mcp import` que os compartilha, e uma conta sem login em servidores MCP em que outras contas têm login recebe um aviso para entrar com `/mcp` (o ccam nunca copia logins). Sem reparo.
9. **Formato da resposta de uso** (`checkUsageFormat`). Consulta o endpoint de uso, sem o cache, para uma conta após a outra, até que uma responda, tudo dentro de um único limite de 4 segundos. Com `usage_source = "claude"`, ele pergunta ao Claude Code de uma conta em vez disso. Se a resposta não puder mais ser interpretada, ele avisa que `ccam usage` passa a usar o cache da linha de status; quando nenhuma conta respondeu, ele diz isso como informação (`could not check the usage response format: …`). Esta é a única verificação que usa a rede.
10. **Atalhos** (`checkShims`). Só quando uma conta tem um atalho. Links ausentes em `~/.local/bin` são um aviso, assim como atalhos criados pelo ccam que ainda executam um ccam mais antigo, que foi substituído ou movido (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` chama `shim.Install`, que cria e renova só os links que o ccam criou e registrou em `~/.local/state/ccam/shortcuts.json`, e nunca mexe em outros arquivos. O aviso termina em `→ fixed` só quando cada um desses links foi criado; se algum não pôde ser, o motivo é acrescentado. Um atalho que executa outro ccam que ainda está instalado funciona e é informação; `--fix` o aponta para este ccam. Uma pasta bin que não está no `PATH` é um aviso e não é reparada.
11. **Pastas** (`checkFolders`). Só quando existem entradas de [`ccam folder`](/features/folder). Uma entrada cuja conta não existe é um erro, e uma entrada cuja pasta sumiu é um aviso. Nenhuma das duas é reparada, porque remover uma entrada mudaria qual conta inicia ali.
12. **Grupos e configurações** (`checkGroups`). Um membro de grupo que não é uma conta, uma conta cujo nome virou depois um comando do ccam (inicie-a com `ccam run <name>` ou renomeie-a), um `theme` ou `usage_source` desconhecido e linhas do `config.toml` que o ccam não lê são avisos. As últimas geralmente são uma configuração escrita abaixo de uma tabela `[[group]]`, `[[profile]]` ou `[[folder]]`, onde ela pertence a essa tabela (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); os comandos que alteram a configuração se recusam até que isso seja corrigido. Sem reparo.
13. **Logins** (`checkIdentities`). Lê, no `.claude.json` de cada pasta, com qual conta ela está logada. Uma pasta logada agora em uma conta diferente da que o `ccam login` registrou (por exemplo, depois de um `/login` dentro de uma sessão) é um aviso; uma conta sem registro de login, como uma importada ou criada por um ccam mais antigo, é comparada pelo e-mail registrado. Uma pasta logada na mesma conta Claude que outra também é um aviso: as duas dividem os limites de uso. Uma pasta que tem o mesmo token de login que outra, uma pasta copiada (inclusive o `~/.claude` do `claude` simples), é um erro, já que a que renovar primeiro desconecta a outra; execute `ccam login <name>` para dar a ela um login próprio. Sem reparo.
14. **`auto_switch`** (`checkAutoSwitch`). Só quando `auto_switch = true`. Uma informação na caixa `General` diz o que o `auto_switch` faz e leva à [descrição dele](/features/next#auto-switch). Para cada plano com assentos entre as contas (Claude Team, Claude Enterprise), uma linha cita essas contas, diz que o `auto_switch` move uma sessão entre assentos de uma mesma organização, cita a regra de assentos da Anthropic com o link da central de ajuda e deixa a escolha com você. Um Claude Code anterior à 2.1.214, que não informa o próprio processo aos hooks (`CLAUDE_PID`), gera um aviso nessa caixa, já que nenhuma sessão pode ser movida. Uma conta sem o plugin do ccam, ou com o plugin desativado, recebe um aviso na própria caixa com o comando que corrige isso no Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (ou `/plugin enable ccam@ccam`). Contas com chave de API ficam de fora. No Windows, a única linha é `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` Sem reparo.

Antes das contas, uma **troca de grupos** que parou no meio é um erro na caixa `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` a conclui na ordem em que foi planejada, com as respostas dadas aos conflitos dela; até lá, nenhuma outra troca de grupos, renomeação ou remoção começa. As contas envolvidas só ganham os links ausentes até que ela termine. Se ela não puder ser concluída, o erro diz que remover aquele arquivo abre mão do resto da troca. `ccam --dry-run doctor --fix` mostra a ordem (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) e o que a parte de cada conta deixaria de lado. Uma troca que outro ccam está fazendo agora é informação, e, enquanto ela roda, `--fix` não repara links de compartilhamento (`its share links are not repaired while another ccam changes the groups; …`). `--fix` mantém `~/.local/state/ccam/group-change.lock` enquanto repara e lê o `config.toml` de novo depois de concluir uma troca.

### Reparar links de compartilhamento

O caminho de inicialização cria na hora um link ausente, mas para diante de qualquer coisa que possa perder dados e indica este comando:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Para qualquer problema de link de compartilhamento, `--fix` chama `share.Build` com as configurações de compartilhamento da conta como estado atual e como estado de destino, e depois `Plan.Apply`. `Build` deriva de novo cada etapa a partir do que está no disco, o mesmo plano que uma alteração com `ccam group` usa, então o reparo segue as mesmas regras de entrar em um grupo.

<!-- motion:doctor-repair -->

Aqui, `settings.json` é um arquivo comum. Como o grupo já tem um `settings.json` diferente, o conflito usa a escolha padrão, `KeepGroup`, e a cópia do grupo permanece. O arquivo da conta vai para `.ccam-backup/<time>/settings.json` dentro da pasta da conta, e `os.Symlink` vincula o caminho ao grupo. `--fix` avisa: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Outros casos:

- Se o grupo ainda não tem uma cópia, o arquivo da conta vai para o grupo e se torna o primeiro conteúdo dele.
- No caso de uma pasta, os arquivos que o grupo não tem vão para o grupo. Os arquivos que o grupo tem idênticos, byte a byte, são descartados, e os que diferem vão para o backup, indicados em uma nota (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Uma pasta de grupo ausente é criada primeiro.
- Para `history.jsonl`, o padrão é `merge by line`: as linhas do grupo e depois as suas novas. A cópia anterior do grupo é salva em `.ccam-backup/<time>/group/`.
- Um link para outro grupo é substituído. Um link para uma pasta fora do ccam vai para o backup, e o conteúdo dela inicia o grupo se ele estiver vazio.
- Um link que sobrou de um item que agora é local é removido, e a conta recebe uma cópia do conteúdo do grupo.

Exceto pela remoção de links e pelo descarte de arquivos que o grupo já tem idênticos, byte a byte, cada etapa move dados para o grupo ou para a pasta de backup. Nada se perde.

## Segurança e casos especiais

- Um `ccam doctor` simples nunca altera nada. Em scripts, use `--dry-run` com `--fix` para ter certeza.
- Para uma conta em execução, seja qual for a forma como o Claude Code foi iniciado, `doctor --fix` só cria os links ausentes e as pastas de grupo ausentes, o que não perde nada. Quando sobra alguma outra coisa, ele adiciona `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, porque a sessão continua gravando nos arquivos que um reparo moveria.
- Se um reparo de compartilhamento falhar, o ccam informa `could not repair share links: <error>` para essa conta e pula a verificação de entradas desconhecidas dela. As outras contas continuam sendo verificadas.
- Se uma alteração com `ccam group` falhar no meio, os novos grupos já estão salvos e a troca fica anotada em `~/.local/state/ccam/group-change.json`; o ccam pede que você conclua com `ccam doctor --fix`, que faz o resto dela, inclusive as etapas de MCP.
- Um `.credentials.json` com link simbólico, um login expirado e credenciais no ambiente são informados, mas nunca alterados.
- As reescritas dos registros de plugins são gravadas de forma atômica, mas não são copiadas para `.ccam-backup/`.

## Relacionados

- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Copiar uma configuração existente](/features/migrate)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Comandos](/docs/commands)
