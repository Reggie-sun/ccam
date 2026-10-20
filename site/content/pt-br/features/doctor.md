---
title: Verificar e reparar
description: Verifique links de compartilhamento, permissões de arquivos, logins e a versão do Claude Code com ccam doctor, e corrija o que é seguro com ccam doctor --fix.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: "6ae82bbf5585"
---

`ccam doctor` verifica a sua instalação do ccam: o binário do Claude Code, a pasta, o login e os links de compartilhamento de cada conta, os servidores MCP compartilhados, o endpoint de uso e os seus atalhos. `ccam doctor --fix` repara o que consegue sem perder dados e, antes, move para uma pasta de backup qualquer arquivo que substituir. Execute-o quando `ccam <name>` parar com um problema de link de compartilhamento, depois de restaurar arquivos manualmente ou depois de uma atualização do Claude Code.

## Uso

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Cada resultado é uma linha com uma marca: `✓` OK, `·` informação, `!` um aviso, `✗` um erro. Os resultados de contas começam com o nome da conta.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Com `--fix`, cada resultado reparado termina em `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--fix` | Repara o que é seguro reparar (os arquivos alterados ganham backup) | desativado |
| `--dry-run` (global) | Executa todas as verificações, mas não repara nada, mesmo com `--fix` | desativado |

`ccam doctor` não aceita argumentos. Ele respeita `--dry-run`: `ccam --dry-run doctor --fix` imprime os mesmos resultados que `ccam doctor` e não altera nada.

O código de saída é `1` quando sobra pelo menos um resultado `✗` sem reparo, e `0` nos outros casos. Avisos e informações nunca fazem a execução falhar. Quando a execução falha sem `--fix`, o ccam termina com uma linha em branco e `` `ccam doctor --fix` repairs what it can. ``. Uma configuração que não pode ser carregada também termina com `1` e um erro `ccam:`.

## Como funciona

`internal/cli/doctor.go` carrega `~/.config/ccam/config.toml` e chama `doctor.Run` em `internal/doctor/doctor.go`, que executa as verificações em uma ordem fixa e retorna os resultados. `--fix` é repassado como `Options.Fix`, definido como false em uma simulação.

<!-- motion:doctor-checks -->

O bloco do meio é executado uma vez por conta. As verificações, e o que `--fix` faz em cada uma:

1. **Versão do Claude Code** (`checkClaude`). Localiza o binário `claude` como o inicializador faz e executa `claude --version`. Se a versão não começar com `2.1.302`, a versão contra a qual o layout de pastas do ccam foi verificado (`doctor.LayoutVersion`), ele adiciona `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Um binário ausente ou um `--version` que falha é um erro. Sem reparo.
2. **Credenciais no ambiente** (`checkAuthEnv`). Avisa quando `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN` está definida: toda conta iniciada nesse shell a herda, e o Claude Code a prefere ao login da própria conta. Sem reparo; remova a variável você mesmo.
3. **Pasta** (`checkProfile`). Uma pasta de conta ausente é um erro e encerra as verificações dessa conta. Permissões com qualquer bit de grupo ou de outros definido são um aviso; `--fix` executa `chmod 0700`. O Windows não tem bits de permissão, então lá esta verificação e a de `0600` abaixo são puladas; a pasta no seu perfil de usuário é protegida pela lista de controle de acesso dela.
4. **Arquivo de login** (`checkCredentials`). `.credentials.json` é verificado com `lstat`. Um link simbólico é um erro, porque cada conta precisa ter o próprio token, e não é reparado. Um modo diferente de `0600` é um aviso; `--fix` executa `chmod 0600`. Depois o login é lido (no macOS, primeiro do Acesso às Chaves). A falta de login é um aviso que indica `ccam login <name>`. Um token de renovação expirado é o erro `login expired`; um token de acesso expirado é apenas informação, já que o Claude Code o renova na próxima inicialização. Um login válido mostra o e-mail e o plano (`login valid (dev@company.com, Claude Team · Acme)`). Nenhum desses é reparado; `ccam login --expired` faz login em todas as contas que precisam.
5. **Links de compartilhamento** (`share.Check`). Compara cada caminho de `share.Layout` com os grupos da conta, usando só `lstat` e `readlink`. Ele informa quatro tipos: um link ausente, um link que aponta para outro lugar, um arquivo comum onde deveria haver um link e um link que sobrou de um item que agora é local. Um link ausente é um aviso; os outros são erros. Os caminhos adicionados a um item em uma versão posterior do ccam (`file-history/`, `tasks/` e `paste-cache/` no histórico, `keybindings.json` nas configurações) são só um aviso enquanto a conta ainda tem a própria cópia: `is not shared with the group yet (ccam doctor --fix moves it in)`. Uma inicialização nunca para por causa deles.
6. **Entradas desconhecidas** (`checkUnknownEntries`). Lista os nomes de primeiro nível na pasta da conta que o ccam não conhece (`share.Known`), como um indício de que o layout do Claude Code mudou. Apenas informação.
7. **Caminhos de plugins** (`checkPlugins`). Os registros de plugins (`installed_plugins.json`, `known_marketplaces.json`) armazenam caminhos absolutos. No `plugins/` de cada grupo, caminhos que passam pela pasta de uma conta deixam de funcionar quando essa conta é removida ou renomeada. `--fix` os reescreve para a pasta do grupo com `profile.RewritePluginPaths`, uma gravação atômica do arquivo de registro.
8. **Servidores MCP** (`checkMCP`). Carrega o `mcp.json` de cada grupo uma vez. Um arquivo que não pode ser lido é um erro; um servidor que não passa em `mcp.Validate` é um aviso. Sem reparo.
9. **Formato da resposta de uso** (`checkUsageFormat`). Consulta o endpoint de uso, sem o cache e com o limite habitual de 3 segundos, para uma conta após a outra, até que uma responda. Se a resposta não puder mais ser interpretada, ele avisa que `ccam usage` passa a usar o cache da linha de status. Esta é a única verificação que usa a rede.
10. **Atalhos** (`checkShims`). Só quando uma conta tem um atalho. Links ausentes em `~/.local/bin` são um aviso; `--fix` chama `shim.Install`, que só cria os links ausentes e nunca mexe em outros arquivos. Uma pasta bin que não está no `PATH` é um aviso e não é reparada.
11. **Pastas** (`checkFolders`). Só quando existem entradas de [`ccam folder`](/features/folder). Uma entrada cuja conta não existe é um erro, e uma entrada cuja pasta sumiu é um aviso. Nenhuma das duas é reparada, porque remover uma entrada mudaria qual conta inicia ali.
12. **Grupos e configurações** (`checkGroups`). Um membro de grupo que não é uma conta, uma conta cujo nome virou depois um comando do ccam (inicie-a com `ccam run <name>` ou renomeie-a) e um `theme` ou `usage_source` desconhecido no `config.toml` são avisos. Sem reparo.
13. **Logins** (`checkIdentities`). Lê, no `.claude.json` de cada pasta, com qual conta ela está logada. Uma pasta logada agora em uma conta diferente da que o `ccam login` registrou (por exemplo, depois de um `/login` dentro de uma sessão) é um aviso, e também uma pasta logada na mesma conta Claude que outra: as duas dividem os limites de uso. Sem reparo.

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

Aqui, `settings.json` é um arquivo comum. Como o grupo já tem um `settings.json` diferente, o conflito usa a escolha padrão, `KeepGroup`, e a cópia do grupo permanece. O arquivo da conta vai para `.ccam-backup/<time>/settings.json` dentro da pasta da conta, e `os.Symlink` vincula o caminho ao grupo. Outros casos:

- Se o grupo ainda não tem uma cópia, o arquivo da conta vai para o grupo e se torna o primeiro conteúdo dele.
- No caso de uma pasta, os arquivos que o grupo não tem vão para o grupo; o resto vai para o backup.
- Para `history.jsonl`, o padrão é `merge by line`: as linhas do grupo e depois as suas novas. A cópia anterior do grupo é salva em `.ccam-backup/<time>/group/`.
- Um link para outro grupo é substituído. Um link para uma pasta fora do ccam vai para o backup, e o conteúdo dela inicia o grupo se ele estiver vazio.
- Um link que sobrou de um item que agora é local é removido, e a conta recebe uma cópia do conteúdo do grupo.

Exceto pela remoção de links, cada etapa move dados para o grupo ou para a pasta de backup. Nada é excluído.

## Segurança e casos especiais

- Um `ccam doctor` simples nunca altera nada. Em scripts, use `--dry-run` com `--fix` para ter certeza.
- `doctor --fix` não verifica se uma conta está em execução. Encerre o Claude Code de uma conta antes de reparar os links de compartilhamento dela, porque a sessão continua gravando nesses arquivos.
- Se um reparo de compartilhamento falhar, o ccam informa `could not repair share links: <error>` para essa conta e pula a verificação de entradas desconhecidas dela. As outras contas continuam sendo verificadas.
- Se uma alteração com `ccam group` falhar no meio, os novos grupos já estão salvos, e o ccam pede que você conclua com `ccam doctor --fix`.
- Um `.credentials.json` com link simbólico, um login expirado e credenciais no ambiente são informados, mas nunca alterados.
- As reescritas dos registros de plugins são gravadas de forma atômica, mas não são copiadas para `.ccam-backup/`.

## Relacionados

- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Copiar uma configuração existente](/features/migrate)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Comandos](/docs/commands)
