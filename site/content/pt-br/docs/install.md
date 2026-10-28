---
title: Instalar o ccam
metaTitle: Instalar o ccam no Linux, no macOS e no Windows
description: Instale o ccam no Linux, macOS ou Windows com um comando, um pacote .deb ou go install, e atualize com ccam update. O Claude Code precisa estar instalado.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

O ccam inicia o Claude Code para você, então instale o [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) primeiro e confirme que `claude --version` funciona.

Execute isto em um terminal no Linux ou no macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

O script baixa a versão para o seu sistema, verifica o checksum SHA-256 e coloca o `ccam` em `~/.local/bin`, onde o ccam também cria os atalhos que você escolher depois, como `ccwork`. Se essa pasta não estiver no seu `PATH`, o script imprime a linha que você deve adicionar ao `~/.bashrc` ou ao `~/.zshrc`.

No Windows, execute isto no PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Ele baixa `ccam_windows_<arch>.zip`, verifica o checksum SHA-256 e coloca o `ccam.exe` em `%USERPROFILE%\.local\bin`, a pasta onde o instalador do Claude Code coloca o `claude.exe`. Se essa pasta ainda não estiver no `PATH` do seu usuário, o script a adiciona e mantém entradas como `%USERPROFILE%\bin` como estão. Depois, abra um novo terminal.

## Sistemas suportados

- Linux em x86-64 e ARM64: Ubuntu, Debian, Fedora, Arch, Alpine e outras distribuições. O binário é estático e não tem dependências.
- macOS em Apple silicon e Intel.
- Windows 10 e 11 em x86-64 e ARM64, no PowerShell, no Windows Terminal ou no Git Bash:
  - Pastas compartilhadas funcionam sem o Modo de Desenvolvedor.
  - Arquivos compartilhados (`settings.json`, `CLAUDE.md`, `history.jsonl`) exigem o Modo de Desenvolvedor (Configurações > Sistema > Para desenvolvedores), porque o Windows só permite criar links para arquivos nesse modo ou a um administrador; consulte [compartilhamento no Windows](/features/share#on-windows).
  - Sem o Modo de Desenvolvedor, o primeiro `ccam add` cria o grupo `default` só com servidores MCP e extensões e explica como adicionar o resto depois que o Modo de Desenvolvedor estiver ativado.
- Dentro do WSL, o ccam funciona como em qualquer outro sistema Linux.

## Outras formas de instalar

**Pacote para Debian e Ubuntu.** Baixe `ccam_<version>_linux_amd64.deb` da [página de versões](https://github.com/jhyunwoo/ccam/releases) e instale-o:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

O pacote instala `/usr/bin/ccam`.

**Pelo código-fonte.** Com Go 1.26 ou mais recente:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

Se você instalar uma tag de versão (`@latest` ou `@vX.Y.Z`), o ccam informa essa versão, e o `ccam update` e a oferta de atualização do painel funcionam. Um build de `@main`, ou de um checkout local com commits depois de uma tag, informa `dev`, e o `ccam update` o recusa, a menos que você passe `--version`.

**Uma versão específica.** Defina `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Preenchimento automático no shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

O preenchimento automático funciona para todos os comandos e opções do `ccam`, incluindo os nomes das contas.

## Plugin do Claude Code

O plugin do ccam ensina o próprio Claude Code a usar o ccam: a verificar qual conta tem espaço, a iniciar workers `claude -p` em paralelo em contas diferentes com `ccam pick` e a explicar uma skill ou um servidor MCP ausente com `ccam ext`. Instale-o dentro do Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Ele adiciona:

- uma skill que sugere comandos que fazem login ou movem arquivos, mas nunca os executa;
- os comandos `/ccam:usage`, `/ccam:ext` e `/ccam:switch`;
- dois hooks ([detalhes](/features/next#inside-claude-code)). Um mostra uma notificação no terminal quando uma sessão para em um limite de uso; no Linux e no macOS, com o [`auto_switch`](/features/next#auto-switch) ativado, o ccam então move a sessão para outra conta com espaço. O outro permite que `/ccam:switch [account]` passe a conversa para outra conta.

Os hooks executam `ccam`, então ele precisa estar no `PATH` que o Claude Code enxerga.

## Verificar a instalação

```sh
ccam version
ccam doctor
```

`ccam doctor` verifica se o Claude Code foi encontrado e se todas as pastas de conta, links e logins estão em ordem. Continue com o [início rápido](/docs/quick-start).

## Atualizar

```sh
ccam update --check
ccam update
```

`ccam update --check` apenas informa se existe uma versão mais recente. `ccam update` baixa a versão para o seu sistema, confere o checksum SHA-256 com o `checksums.txt` da versão e então substitui o arquivo `ccam` de uma só vez. Se o checksum não bater, nada é alterado. Atalhos como `ccwork` são links para o ccam, então já executam a nova versão na hora.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

O painel (`ccam` sem argumentos, em um terminal) também procura uma nova versão toda vez que abre e, quando há uma, a oferece: pressione `U` e confirme. Nenhum outro comando procura atualizações, e o ccam nunca se atualiza sem perguntar. Para instalar uma versão específica, inclusive uma mais antiga, passe `--version`. `ccam --dry-run update` mostra o que seria baixado e substituído sem alterar nada.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Executar o script de instalação de novo também atualiza a cópia em `~/.local/bin`. Se você compilou o ccam com `go install` a partir de uma tag de versão, você pode executar `ccam update` ou o mesmo comando `go install` outra vez.

**Windows.** `ccam update` baixa o `.zip` para Windows. No Windows, um programa em execução não pode ser substituído, então o binário antigo é renomeado para `ccam.exe.old`, que a próxima atualização exclui; enquanto uma sessão ainda roda a partir de um `ccam.exe.old` anterior, o nome é `ccam.exe.old-<time>`. Os atalhos que são hard links são renovados logo em seguida. Executar o `install.ps1` de novo funciona do mesmo jeito, mesmo enquanto uma sessão do Claude Code ainda roda a partir do binário antigo, e imprime `OK 2 shortcuts now run the new ccam` quando renova atalhos que são hard links.

**Pacote para Debian e Ubuntu.** `ccam update` não substitui um `/usr/bin/ccam` instalado pelo apt e, em vez disso, mostra estes passos. Baixe o novo `ccam_<version>_linux_amd64.deb` da [página de versões](https://github.com/jhyunwoo/ccam/releases) e instale-o por cima do antigo:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Desinstalar

`ccam uninstall` remove o ccam e tudo o que ele criou. Primeiro ele lista o que vai ser removido e pergunta:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Cada conta é removida como `ccam remove` faz. Uma conta nas pastas do ccam é desconectada primeiro e excluída junto com a pasta dela, incluindo o histórico de conversas e as configurações.
- Uma pasta `~/.claude` que você [importou](/docs/commands) permanece onde está, para o `claude` comum. Se ela estiver em um grupo, recebe de volta uma cópia de tudo o que compartilhava. A linha de status do ccam no `settings.json` dela dá lugar à linha de status que o ccam envolvia, se havia uma. Se ela tiver o plugin do ccam, o ccam imprime o comando `/plugin uninstall` para executar no Claude Code.
- Os atalhos, `~/.config/ccam`, `~/.local/share/ccam` (com a lixeira do ccam) e `~/.local/state/ccam` são excluídos, e depois o próprio comando `ccam`.

O ccam se recusa a continuar enquanto uma conta estiver em execução; encerre antes as sessões do Claude Code dela. Um `/usr/bin/ccam` instalado pelo apt permanece; remova-o depois com `sudo apt remove ccam`. No Windows, um programa em execução não consegue excluir a si mesmo, então `ccam uninstall` termina imprimindo o comando `Remove-Item` para o `ccam.exe`.
