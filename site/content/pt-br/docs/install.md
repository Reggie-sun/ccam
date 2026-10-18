---
title: Instalar o ccam
description: Instale o ccam no Linux ou no macOS com um comando, por um pacote .deb ou pelo código-fonte. O Claude Code precisa ser instalado antes.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Execute isto em um terminal no Linux ou no macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

O script baixa a versão para o seu sistema, verifica o checksum SHA-256 e coloca dois comandos em `~/.local/bin`: `ccam` e seu alias curto `ccs`. Se essa pasta não estiver no seu `PATH`, o script imprime a linha que você deve adicionar ao `~/.bashrc` ou ao `~/.zshrc`.

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

O pacote instala `/usr/bin/ccam` e o link `/usr/bin/ccs`.

**Pelo código-fonte.** Com Go 1.26 ou mais recente:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

O preenchimento automático funciona para `ccam` e `ccs`, incluindo os nomes das contas.

## Verificar a instalação

```sh
ccam version
ccam doctor
```

`ccam doctor` verifica se o Claude Code foi encontrado e se todas as pastas de conta, links e logins estão em ordem. Continue com o [início rápido](/docs/quick-start).

## Desinstalar

Remova os dois comandos e os atalhos das contas e depois as pastas de dados, se você não precisar mais das contas:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Excluir `~/.local/share/ccam` também exclui os logins e o histórico de conversas armazenados ali. Uma pasta `~/.claude` que você [importou](/docs/commands) permanece onde está.
