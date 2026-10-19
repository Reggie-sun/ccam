---
title: Compartilhamento entre contas
description: Como o ccam share vincula configurações, CLAUDE.md, extensões e histórico a uma pasta de grupo, com um plano, backups e escolhas de conflito, e como off copia de volta.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Cada conta tem a própria pasta, mas a maior parte do que há nela não precisa ficar separada. O `ccam share` ativa ou desativa o compartilhamento de uma conta e de um item por vez: o arquivo ou a pasta da conta é movido para uma pasta de grupo e substituído por um link simbólico, então todas as contas do grupo leem e gravam a mesma cópia. Use-o depois do `ccam add` para mudar o que uma conta compartilha, ou para colocar algumas contas em um grupo só delas.

## Uso

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Sem opções de itens, o `ccam share` imprime o estado atual:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Com opções, o ccam imprime um plano, pergunta sobre cada arquivo que difere da cópia do grupo e pergunta `Proceed? [y/N]` antes de alterar qualquer coisa. Uma conta nova compartilha tudo, exceto o histórico de conversas; [Compartilhar servidores MCP, configurações e histórico](/docs/sharing) descreve cada item do ponto de vista do usuário.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Servidores MCP | inalterado |
| `--settings <on\|off\|group>` | `settings.json` | inalterado |
| `--memory <on\|off\|group>` | o `CLAUDE.md` global | inalterado |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | inalterado |
| `--history <on\|off\|group>` | `projects/` e `history.jsonl` | inalterado |
| `--fresh` | ao desativar o compartilhamento, começa vazio em vez de copiar o conteúdo do grupo | desativado |
| `--prefer <group\|mine\|merge>` | como decidir os conflitos | pergunta; com `--yes`, o padrão de cada item |
| `-y`, `--yes` | executa sem perguntar; os conflitos seguem `--prefer` | desativado |
| `-h`, `--help` | ajuda do `share` | |
| `--dry-run` (global) | respeitado: pergunta sobre os conflitos, imprime o plano e não altera nada | desativado |

O valor `on` significa o grupo `default`, `off` significa local (não compartilhado), e qualquer outro valor é um nome de grupo. Os nomes de grupos viram nomes de pastas em `~/.local/share/ccam/groups/`, então `ccam share side --history ../x` falha com `invalid group name: "../x"`. A conta é indicada pelo nome ou pelo número, e só os itens que você indicar mudam.

## Como funciona

Os arquivos envolvidos:

- `~/.config/ccam/config.toml`: a tabela `[profile.share]` de cada conta, um valor por item.
- `~/.local/share/ccam/profiles/<name>/`: o `CLAUDE_CONFIG_DIR` da conta.
- `~/.local/share/ccam/groups/<group>/`: as cópias compartilhadas.

1. `internal/cli/share.go` lê a configuração, aplica cada valor `--<item>` a uma cópia das configurações da conta com `config.Share.Set` e se recusa se a conta estiver em execução.
2. `share.Build` compara as configurações atuais com o destino e lista as etapas. Ele só lê o disco: `os.Lstat` e `os.Readlink` para cada caminho de `share.Layout`. O plano é uma lista de operações: mover para um lugar vazio no grupo (adotar), mover os arquivos que o grupo não tem (mesclar pasta), gravar um conflito resolvido no grupo, fazer backup, remover o link, criar o link, copiar de volta, reescrever caminhos de plugins e mesclar o registro de plugins.
3. Os conflitos são decididos, pela pergunta ou por `--prefer`.
4. `Plan.Describe` imprime o plano, e o ccam pergunta.
5. As configurações de destino são gravadas no `config.toml` antes que qualquer coisa seja movida. Se `Plan.Apply` parar no meio, `ccam doctor --fix` e uma nova execução montam o plano em direção ao destino salvo e concluem a alteração.
6. `Plan.Apply` executa as etapas em ordem. Cada etapa move dados para o grupo ou para a pasta de backup; nada é excluído.

`share.Layout` é uma lista de permissões dos caminhos que cada item abrange. Qualquer coisa fora dela nunca chega a um grupo. `share.Never` indica `.credentials.json` e `.claude.json`, e `Build` para com um erro interno se algum deles aparecer no layout. O item `mcp` não tem caminhos: a lista de servidores fica dentro do `.claude.json`, então o [compartilhamento de MCP](/features/mcp) funciona pelo `mcp.json` do grupo, e `--mcp` só altera a configuração.

Isto é o que acontece com uma pasta, `projects/`, quando `side` ativa o histórico e o grupo já tem conteúdo:

<!-- motion:share-on -->

Se o grupo ainda não tem nada nesse caminho, a conta é a primeira a compartilhá-lo, e a pasta dela é movida para lá como está, sem backup. Caso contrário, `mergeDir` move para o grupo cada arquivo que ele não tem, e os arquivos que existem dos dois lados ficam na pasta da conta, então a versão do grupo é mantida. Depois, a pasta da conta vai para `.ccam-backup/<time>/` dentro da pasta da conta, e `os.Symlink` coloca um link no lugar dela. Os arquivos de conversa recebem o nome do ID da sessão, então os históricos de várias contas se combinam sem conflitos.

A movimentação usa `rename(2)` e, entre sistemas de arquivos diferentes, recorre a copiar e excluir. Um link relativo dentro da árvore movida que apontava para fora dela é reescrito como um link absoluto (`fsutil.PortableLink`), para que continue funcionando a partir do novo lugar. Os registros de `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) armazenam caminhos absolutos que passam pela pasta da conta; o ccam os aponta para a pasta do grupo e, quando o grupo já tem um registro, mescla nele as entradas da conta.

Desativar um item inverte o link:

<!-- motion:share-off -->

O ccam só remove o link se ele apontar para dentro da pasta de grupos e, em seguida, copia o conteúdo atual do grupo para a conta com `CopyTree`. Com `--fresh`, a cópia é pulada e a conta começa sem nada nesse caminho. O grupo mantém o conteúdo para as outras contas. Mover um item de um grupo para outro remove o link e o vincula ao novo grupo, deixando o grupo antigo como está; o plano avisa isso em uma linha `Note:`.

Um único arquivo que existe dos dois lados com conteúdo diferente é um conflito:

<!-- motion:share-conflict -->

`share.NewConflict` lista as opções adequadas ao arquivo. `history.jsonl` oferece `merge by line`, `keep group` e `use mine`, com a mesclagem primeiro. Outros arquivos `.json` oferecem `keep group`, `use mine` e `merge JSON`. Qualquer outro arquivo oferece `keep group` e `use mine`. A primeira opção é a padrão:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) mescla objetos recursivamente, junta arrays sem duplicatas e mantém o valor do grupo quando os dois lados definem a mesma chave de forma diferente, para que as outras contas não vejam mudança de comportamento. `merge by line` mantém as linhas do grupo e depois adiciona as suas linhas que ele não tem. Toda escolha, exceto `keep group`, salva antes a cópia do grupo em `.ccam-backup/<time>/group/` na pasta da conta e depois grava o resultado com `fsutil.AtomicWrite`. O seu próprio arquivo sempre vai para o backup.

`--prefer group`, `mine` ou `merge` decide todos os conflitos sem perguntar; `merge` significa `merge JSON` ou `merge by line`, o que o arquivo oferecer, e um arquivo que não oferece nenhum dos dois fica com o padrão. `--yes` sozinho usa o padrão de cada arquivo.

A cada inicialização, `share.Check` compara cada caminho do layout com a configuração usando apenas `lstat` e `readlink`, o que o mantém leve o suficiente para o caminho de inicialização. Um link ausente é criado na hora, porque não há nada a perder. Um arquivo comum onde deveria haver um link, um link para o lugar errado ou um link obsoleto interrompe a inicialização:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` executa `share.Build` com as configurações salvas dos dois lados, o que deriva cada etapa a partir do que está no disco, mantém o conteúdo do grupo nos conflitos e faz backup de tudo o que substitui.

## Segurança e casos especiais

- **Contas em execução.** O ccam se recusa a alterar o compartilhamento enquanto o Claude Code estiver em execução com a conta (`account work is running; quit it first, then change its sharing`), porque a sessão continua gravando nos arquivos que estão sendo movidos, e as linhas que ela acrescentasse nesse meio-tempo se perderiam. `--dry-run` ainda imprime o plano.
- **Backups.** Cada plano recebe a própria pasta `.ccam-backup/<time>/` na pasta da conta, com um sufixo numerado quando dois planos começam no mesmo segundo. O ccam nunca as exclui.
- **Falhas parciais.** Se uma etapa falhar, o ccam imprime o erro e onde estão os originais, e pede que você conclua com `ccam doctor --fix`. O destino já está salvo, então o reparo continua em direção a ele.
- **Seus próprios links.** Se um caminho já é um link que você criou, por exemplo para um repositório de dotfiles, e o grupo está vazio, o ccam copia o destino do link para o grupo e move o próprio link para o backup. O destino não é alterado.
- **Arquivos de login.** `.credentials.json` e `.claude.json` nunca são vinculados, movidos nem copiados para um grupo.
- **Duas contas na mesma pasta.** Contas que compartilham o histórico gravam nos mesmos arquivos. Quando uma segunda conta do mesmo grupo de histórico inicia em uma pasta onde outra já está em execução, o ccam avisa e pergunta antes de continuar.

## Relacionados

- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Servidores MCP compartilhados](/features/mcp)
- [Copiar uma configuração existente](/features/migrate)
- [Todos os comandos](/docs/commands)
