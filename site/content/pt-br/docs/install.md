---
title: Instalar o ccam
description: Instale o ccam no Linux ou no macOS com um comando, por um pacote .deb ou pelo código-fonte, e mantenha-o atualizado com ccam update.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 4d7291e0f86e
---

Execute isto em um terminal no Linux ou no macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

O script baixa a versão para o seu sistema, verifica o checksum SHA-256 e coloca o `ccam` em `~/.local/bin`, onde o ccam também cria os atalhos que você escolher depois, como `ccwork`. Se essa pasta não estiver no seu `PATH`, o script imprime a linha que você deve adicionar ao `~/.bashrc` ou ao `~/.zshrc`.

O ccam inicia o Claude Code para você, então instale o [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) primeiro e confirme que `claude --version` funciona.

## Sistemas suportados

- Linux em x86-64 e ARM64: Ubuntu, Debian, Fedora, Arch, Alpine e outras distribuições. O binário é estático e não tem dependências.
- macOS em Apple silicon e Intel.
- O Windows não é suportado diretamente. Dentro do WSL, o ccam funciona como em qualquer sistema Linux.

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

**Uma versão específica.** Defina `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Preenchimento automático no shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

O preenchimento automático funciona para todos os comandos e opções do `ccam`, incluindo os nomes das contas.

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

O ccam nunca procura atualizações por conta própria: ele só busca uma nova versão quando você executa este comando. Para instalar uma versão específica, inclusive uma mais antiga, passe `--version`. `ccam --dry-run update` mostra o que seria baixado e substituído sem alterar nada.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Executar o script de instalação de novo também atualiza a cópia em `~/.local/bin`. Se você compilou o ccam com `go install`, execute o mesmo comando outra vez.

**Pacote para Debian e Ubuntu.** `ccam update` não substitui um `/usr/bin/ccam` instalado pelo apt e, em vez disso, mostra estes passos. Baixe o novo `ccam_<version>_linux_amd64.deb` da [página de versões](https://github.com/jhyunwoo/ccam/releases) e instale-o por cima do antigo:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Desinstalar

Remova primeiro os atalhos das contas, enquanto o ccam ainda está instalado. `ccam list` os mostra na coluna Shortcut; remova cada um com `ccam shortcut <account> --remove` ou exclua você mesmo esses links de `~/.local/bin`. Depois remova o comando e, se você não precisar mais das contas, as pastas de dados:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Excluir `~/.local/share/ccam` também exclui os logins e o histórico de conversas armazenados ali. Uma pasta `~/.claude` que você [importou](/docs/commands) permanece onde está.
