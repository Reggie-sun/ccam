---
title: Perguntas frequentes
metaTitle: "Perguntas frequentes: várias contas do Claude Code"
description: Respostas sobre o ccam, contas do Claude Code, limites de uso, compartilhamento, privacidade e sistemas suportados.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: "c4e3b89bb4ea"
---

## O ccam é uma ferramenta oficial da Anthropic?

Não. O ccam é um projeto de código aberto não oficial e não é afiliado à Anthropic nem endossado por ela. Ele usa a configuração `CLAUDE_CONFIG_DIR`, que o [Claude Code documenta](https://code.claude.com/docs/en/env-vars) para executar várias contas lado a lado.

## Posso usar várias contas do Claude Code?

Isso depende dos termos de cada conta e, no caso de contas de trabalho, da política da sua organização. Verifique-os antes de usar várias contas. O ccam apenas facilita a troca entre contas que você já tem.

**Assentos dos planos Team e Enterprise.** A escolha automática (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados no parágrafo acima.

## O ccam contorna os limites de uso?

Não. Cada conta mantém os próprios limites. O ccam mostra esses limites e permite que você continue em outra conta sua quando uma delas se esgota.

## Como usar duas contas do Claude Code no mesmo computador?

Dê a cada conta a própria pasta de configuração com `CLAUDE_CONFIG_DIR`. O ccam faz isso por você: execute `ccam add` para cada conta e depois inicie-as com `ccam personal` e `ccam work`, ou com atalhos que você mesmo nomeia, como `ccwork`. O [guia](/guides/multiple-claude-code-accounts) também mostra como fazer isso manualmente.

## Como continuar uma conversa em outra conta?

Execute `ccam next --continue` na pasta: ele leva a última conversa da pasta para a próxima conta (`--auto` para a que tem mais espaço). Dentro do Claude Code, o `/ccam:switch` do [plugin do ccam](/docs/install#claude-code-plugin) faz o mesmo com a sessão em que você está. Quando as duas contas compartilham o histórico de conversas (`ccam group set default --add history`, ou um grupo só delas: `ccam group add pair --items history --members work,side`), a nova conta abre essa mesma conversa; caso contrário, o Claude Code a copia para o histórico da nova conta.

## O ccam vê ou envia meu login?

O ccam nunca compartilha nem copia seu login do Claude entre contas, e nunca copia os logins de servidores MCP. Ele lê o token de acesso de uma conta apenas para consultar o uso dessa conta no endpoint de uso da Anthropic, quando você executa `ccam usage`, `ccam next --auto` ou `ccam doctor`, e só renova um token expirado com `ccam usage --refresh`; com `usage_source = "claude"`, ele nem lê o token. `ccam ext` e `ccam doctor` leem os nomes dos logins de MCP de uma conta, nunca os tokens deles. O ccam não tem telemetria. Além dessas consultas, a única requisição que ele faz por conta própria é a pergunta do painel ao GitHub sobre a versão mais recente do ccam, que não leva nada sobre você ou suas contas.

## Como fazer login em um servidor remoto?

Execute `ccam add` ou `ccam login <name>` na sessão SSH, como de costume. O Claude Code espera o login no servidor, então o navegador do seu computador termina em uma página `http://localhost:<port>/callback?code=…` que não carrega. Copie esse endereço do navegador e cole no terminal: o ccam o repassa ao Claude Code no servidor, que conclui o login. O endereço contém um código de uso único, então não o compartilhe.

## Onde minhas contas ficam armazenadas?

Em `~/.local/share/ccam/profiles/<name>/`, uma pasta por conta. No macOS, o Claude Code guarda o próprio login no Acesso às Chaves, com um nome derivado dessa pasta, então `ccam rename` mantém ali o nome antigo da pasta.

## O que acontece com a minha pasta ~/.claude existente?

Nada, a menos que você execute `ccam import` ou `ccam migrate`. A importação registra a pasta no lugar onde ela está, então o `claude` simples e o `ccam <name>` continuam usando a mesma conta e as mesmas configurações. O `ccam migrate` copia as configurações, o `CLAUDE.md`, as extensões e os servidores MCP dela para as configurações compartilhadas e só lê a pasta.

## Quais sistemas o ccam suporta?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e outras distribuições), macOS e Windows 10 e 11, nativamente ou dentro do WSL. No Windows, compartilhar arquivos como `settings.json` exige o Modo de Desenvolvedor; pastas compartilhadas funcionam sem ele, e, sem ele, o primeiro `ccam add` cria o grupo `default` só com servidores MCP e extensões ([compartilhamento no Windows](/features/share#on-windows)).

## Por que um servidor MCP aparece em apenas uma conta?

Servidores adicionados dentro de uma sessão com `claude mcp add` ficam armazenados apenas para essa conta e, sem `-s user`, apenas para a pasta atual. O ccam avisa você na próxima vez que você iniciar a conta, e `ccam mcp list` e `ccam ext` os mostram; mova o servidor para o grupo compartilhado com `ccam mcp import <account> <server>`.

## O ccam deixa o Claude Code mais lento?

Iniciar o Claude Code pelo ccam acrescenta cerca de 4 milissegundos, medidos com `scripts/perf-check.sh` em um notebook. O ccam lê um arquivo pequeno, verifica os links compartilhados e depois se substitui pelo `claude`, então ele não fica em execução enquanto você trabalha. O Windows não consegue substituir um processo, então lá o ccam inicia o Claude Code e espera por ele, sem fazer mais nada.

## Como atualizar o ccam?

Execute `ccam update`. Ele baixa a versão mais recente, verifica o checksum e substitui o arquivo `ccam`; `ccam update --check` apenas informa se existe uma versão mais recente. O painel também avisa quando existe uma versão mais recente, e `U` a instala depois de perguntar; o ccam nunca se atualiza sozinho. Se você instalou o pacote `.deb`, instale o novo `.deb` com o apt. Consulte [atualização](/docs/install#update).

## Como desinstalar o ccam?

Remova o atalho de cada conta com `ccam shortcut <account> --remove`. Se você importou o `~/.claude` e ele está em um grupo, execute `ccam remove <name>` para ele antes, para que ele receba de volta uma cópia do que compartilha. Depois, exclua o `ccam` de `~/.local/bin` e remova `~/.config/ccam`, `~/.local/share/ccam` e `~/.local/state/ccam` se você não precisar mais das contas. Consulte [instalação](/docs/install#uninstall).
