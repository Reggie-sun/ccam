---
title: Usar várias contas do Claude Code no mesmo computador
description: Execute várias contas do Claude Code lado a lado com CLAUDE_CONFIG_DIR, manualmente com aliases de shell ou automaticamente com o ccam, e troque entre elas.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "5c5dbd9fcd23"
---

Dê a cada conta a própria pasta de configuração. O Claude Code lê a pasta da variável de ambiente `CLAUDE_CONFIG_DIR`, e cada pasta guarda o próprio login, as configurações e o histórico, então todas as contas continuam logadas e você pode executá-las ao mesmo tempo. Você pode configurar isso manualmente com um alias de shell ou deixar o ccam fazer isso.

## Como configurar manualmente?

Esta é a abordagem que [a documentação do Claude Code descreve](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Adicione uma linha para cada conta extra ao `~/.bashrc` ou ao `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Abra um novo terminal e execute `claude-work`. O Claude Code pede que você faça login, porque a pasta é nova. A partir daí, `claude` é sua primeira conta e `claude-work` a segunda. Adicione um alias por conta, cada um com a própria pasta.

No Windows sem WSL, uma função do PowerShell no seu perfil (`notepad $PROFILE`) faz o mesmo:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## O que fica de fora com aliases?

Funciona bem para duas contas. O que você precisa saber:

- Cada pasta começa vazia: sem configurações, sem `CLAUDE.md`, sem servidores MCP, sem agentes nem skills. Você copia ou cria links para o que precisar e mantém as cópias sincronizadas por conta própria.
- Aliases só existem em shells interativos. Scripts, Makefiles e outros programas não os enxergam.
- Os servidores MCP ficam no `.claude.json`, junto com o login, então não dá para simplesmente criar links deles entre pastas.
- As conversas ficam na conta em que começaram, então você não pode continuar uma delas em outra conta.
- Não há nada que mostre o uso de todas as contas de uma vez.

## Posso trocar de conta com /login?

Sim, se você só troca de vez em quando. Digitar `/login` dentro do Claude Code [troca a conta](https://code.claude.com/docs/en/quickstart) da pasta em que você está. Só uma conta fica logada por vez, então você faz login de novo a cada troca e não pode executar duas contas lado a lado. Pastas separadas evitam as duas coisas.

## Como configurar com o ccam?

O ccam automatiza a mesma ideia e acrescenta as partes que faltam:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Agora `ccam personal` e `ccam work` iniciam as duas contas, assim como `ccwork`, o atalho que você deu à conta work. Como são programas de verdade e não aliases, eles também funcionam em scripts. Execute `ccam` sozinho para escolher uma conta em uma lista, e `ccam use work` para tornar uma delas a padrão.

Configurações, `CLAUDE.md`, servidores MCP e extensões são compartilhados desde o início. `ccam next --continue` leva uma conversa para a outra conta, e o histórico de conversas compartilhado, que você pode ativar, a mantém como uma só conversa. `ccam usage` mostra o uso de 5 horas e semanal de todas as contas de uma vez.

Se você já usa `~/.claude`, execute `ccam import` em vez de adicionar essa conta de novo, ou `ccam migrate` para copiar as configurações, os servidores MCP e as extensões dela para as configurações compartilhadas.

## Como saber qual conta está em execução?

Dentro do Claude Code, `/status` mostra a conta em que você está logado. Com o ccam, execute `ccam statusline install` uma vez e a linha de status passa a mostrar o tempo todo o nome da conta e o uso dela, por exemplo `[work] · Opus · 5h 24% · 7d 41%`.

## Funciona no Windows e no WSL?

O Claude Code roda [nativamente no Windows e dentro do WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` funciona nos dois, então no Windows nativo o alias ou a função do PowerShell acima é suficiente. O ccam também roda nativamente no Windows: instale no PowerShell com `irm https://ccam.dev/install.ps1 | iex`, ou dentro do WSL com o comando do Linux. Os atalhos viram `ccwork.exe`, e compartilhar arquivos como `settings.json` exige o Modo de Desenvolvedor ([Instalar o ccam](/docs/install)).

## Solução de problemas

**O Claude Code pede login toda vez.** A primeira execução em uma pasta nova sempre pede. Se continuar pedindo, o login expirou ou foi removido: execute `ccam login work` (ou inicie o alias uma vez e faça login).

**Um servidor MCP aparece em só uma conta.** Servidores adicionados dentro de uma sessão com `claude mcp add` pertencem àquela conta e, sem `-s user`, só àquela pasta. Leve um deles para todas as contas com `ccam mcp import work <server>`.

**As configurações são diferentes entre as contas.** `ccam share work` mostra o que a conta compartilha e por meio de qual grupo; `ccam group join default work` faz com que ela compartilhe o que o `default` compartilha, depois de mostrar o que vai mudar. `ccam doctor --fix` corrige links quebrados e guarda um backup de cada arquivo que altera.

**Duas contas abrem a mesma conversa.** Evite trabalhar em uma conversa em duas contas ao mesmo tempo. O ccam avisa se você tentar.

## Qual escolher?

Para duas contas sem nada compartilhado, um alias basta. Quando você quer configurações e servidores MCP compartilhados, mais contas, ver o uso de relance ou continuar conversas, o ccam poupa o trabalho de manutenção. A [comparação com aliases](/compare/aliases) lista as diferenças lado a lado, e [o ccam e outras ferramentas para trocar de conta](/compare/account-switchers) trata das ferramentas que trocam um único login.

Se você pode usar várias contas depende dos termos de cada conta. Verifique-os antes.

**Assentos dos planos Team e Enterprise.** A escolha automática (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e o [`auto_switch`](/features/next#auto-switch), quando você o ativa) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Usá-la nesse caso é uma decisão sua: consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados na nota acima.
