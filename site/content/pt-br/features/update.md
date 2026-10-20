---
title: Atualizar o ccam
description: ccam update troca o ccam pela versão mais recente do GitHub após verificar o SHA-256; --check só informa e --version instala uma tag de versão específica.
updated: 2026-10-19
order: 10
translation: machine
sourceHash: 83fc20191328
---

`ccam update` baixa uma versão do ccam do GitHub, confere-a com o `checksums.txt` da versão e substitui o binário em execução com um único rename. Use-o quando você instalou o ccam com o script de instalação ou a partir de um arquivo compactado. O ccam nunca procura atualizações por conta própria: só este comando contata o GitHub, então todos os outros comandos funcionam do mesmo jeito sem conexão.

## Uso

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Quando não existe uma versão mais nova, o ccam imprime `ccam 0.4.0 is up to date.` e não altera nada. Com `--version` e a versão que você já usa, ele imprime `ccam 0.4.0 is already installed.`

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--check` | Só informa se existe uma versão mais nova | desativado |
| `--version <tag>` | Instala esta tag de versão, como `v0.3.0`, mais nova ou mais antiga que a atual | a mais recente |
| `--dry-run` | Respeitado: consulta a tag mais recente (a menos que `--version` indique uma), faz as mesmas verificações no binário e imprime os caminhos de download, checksum e substituição sem baixar nada | desativado |

`--check` e `--version` não podem ser usados juntos. A tag precisa ter o formato `vMAJOR.MINOR.PATCH`, com um sufixo opcional de pré-lançamento ou de build; o `v` inicial pode ser omitido. Qualquer outra coisa falha, por exemplo `--version takes a release tag such as v0.3.0, not "1.2"`.

A variável de ambiente `CCAM_RELEASES_URL` substitui `https://github.com/jhyunwoo/ccam/releases`, para um espelho ou para testes. Ela precisa ter a mesma estrutura: `<url>/latest` redireciona para `<url>/tag/<tag>`, e os arquivos ficam em `<url>/download/<tag>/<asset>`.

## Como funciona

O código fica em `internal/update`, que só o `ccam update` importa. O caminho de inicialização nunca o chama.

1. `cli.updateCmd` encontra o binário em execução com `os.Executable` e resolve os links simbólicos com `filepath.EvalSymlinks`. Esse caminho real é o que será substituído.
2. Sem `--version`, `Updater.Latest` solicita `<releases>/latest` sem seguir o redirecionamento e lê a tag do cabeçalho `Location` (`…/tag/v0.5.0`). Isso usa a página de versões em vez da API do GitHub, então o limite de 60 solicitações por hora da API não se aplica. A consulta tem 10 segundos.
3. `parseVersion` confere a tag com um padrão estrito. Uma tag acaba em URLs de download e em um comando de shell impresso, então nada com `/`, `%`, espaços ou caracteres de shell passa. `compareVersions` a compara com a versão em execução; um pré-lançamento vem antes da versão final correspondente.
4. `managed()` e `writable()` decidem se o ccam pode substituir o binário (veja abaixo).
5. `Updater.Fetch` baixa `ccam_<os>_<arch>.tar.gz` e `checksums.txt` ao mesmo tempo, com 2 minutos para os dois. Cada download é limitado a 256 MB.
6. O SHA-256 do arquivo compactado precisa corresponder à linha dele no `checksums.txt` antes de qualquer outra coisa. Uma linha ausente ou um checksum diferente interrompe a atualização com `nothing was changed`.
7. `extract()` lê o tar e pega o primeiro arquivo comum chamado `ccam`. Pastas, links e outras entradas são ignorados, e o binário também é limitado a 256 MB.
8. `Updater.Install` grava o novo binário com `fsutil.AtomicWrite`: um arquivo temporário `.ccam.tmp*` ao lado do binário, com os bits de permissão do arquivo antigo, gravado no disco e renomeado por cima do antigo.

O processo inteiro, da consulta da tag até o rename:

<!-- motion:update-pipeline -->

O arquivo compactado e o `checksums.txt` são baixados lado a lado, e nada é gravado até que o checksum corresponda. O rename é atômico, então uma falha ou um disco cheio deixa o binário antigo ou o novo, nunca a metade de um. Atalhos como `~/.local/bin/ccwork` são links simbólicos para o binário do ccam, então executam a nova versão sem nenhuma alteração.

### Instalações gerenciadas por pacote

Um ccam instalado pelo `.deb` pertence ao apt, e substituí-lo sem o apt saber deixaria o banco de dados de pacotes incorreto. Antes de baixar, o ccam verifica quem é o dono do binário:

<!-- motion:update-managed -->

1. Um binário em `/usr/bin`, `/bin`, `/usr/sbin` ou `/sbin` é considerado gerenciado por pacote.
2. Em outros lugares, se o `dpkg` estiver instalado e `dpkg -S <path>` funcionar, ele também é considerado gerenciado por pacote.
3. Caso contrário, o ccam cria e remove um arquivo temporário `.ccam-update-*` ao lado do binário para ter certeza de que pode gravar ali.

Para um binário gerenciado por pacote, o ccam para e imprime as etapas do apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Em sistemas que não são Linux, só a linha do script de instalação é impressa. Quando a pasta não permite gravação, o ccam sugere executar de novo o script de instalação, que instala o ccam em `~/.local/bin`.

## Segurança e casos especiais

- Um build de desenvolvimento informa a versão como `dev` e não tem uma versão publicada para comparar, então `ccam update` e `ccam update --check` o recusam. `ccam update --version vX.Y.Z` instala uma versão por cima dele.
- Um download que falha, um arquivo ausente para a sua plataforma ou um checksum que não corresponde deixa o binário instalado como estava.
- Uma versão que não existe falha com `there is no release <tag>`. Quando o GitHub recusa com HTTP 403 ou 429, o ccam diz que provavelmente é um limite de taxa e pede que você tente de novo mais tarde.
- `--check` só lê. Ele não verifica se o binário é gerenciado por pacote nem se permite gravação.
- `--version` com uma tag mais antiga é a forma de voltar a uma versão anterior.

## Relacionados

- [Instalação](/docs/install)
- [Comandos](/docs/commands)
- [Perguntas frequentes](/docs/faq)
