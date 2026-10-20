---
title: "Trocar de conta no Claude Code: ccam e alternativas"
description: Como o ccam se compara a claude-swap, CCS, clausona, clauth, CAAM, cux e outras ferramentas de troca de conta do Claude Code, e qual abordagem escolher.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "a90ede56b767"
---

As ferramentas para trocar de conta no Claude Code seguem uma de duas abordagens. A maioria troca o login dentro da única pasta `~/.claude`, então só uma conta fica ativa por vez. Outras, entre elas o ccam, dão a cada conta a própria pasta com `CLAUDE_CONFIG_DIR`, então todas as contas continuam logadas e várias podem rodar ao mesmo tempo. O ccam acrescenta configurações e servidores MCP compartilhados e continua uma conversa na pasta de outra conta.

## Trocar um login ou usar pastas separadas?

**Trocar um login** mantém um único conjunto de configurações, servidores MCP e histórico, porque só existe uma pasta. A troca é rápida, e uma sessão em andamento pode passar a usar o novo login. Mas só uma conta fica ativa por vez, para todos os terminais e editores ao mesmo tempo.

Para fazer a troca, a ferramenta guarda uma cópia do login de cada uma das outras contas e a coloca de volta no lugar. O Claude Code substitui o refresh token de um login sempre que o renova, então cada cópia salva precisa acompanhar as sessões em andamento. Contas que ficam deslogadas ou que pedem `/login` de novo depois de uma troca estão entre os problemas que esses projetos mais corrigem.

**Pastas separadas** é o método que [a documentação do Claude Code descreve](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Todas as contas continuam logadas, e você pode executar uma sessão de trabalho e uma pessoal lado a lado. Cada login fica na própria pasta, onde o Claude Code o renova normalmente. O preço é que cada pasta começa com as próprias configurações, servidores MCP e histórico, a menos que alguma ferramenta os compartilhe.

## Comparação das ferramentas

Retirado da documentação e do código-fonte de cada projeto em 2026-10-20 (clausona: v0.12.0-beta). "Não documentado" significa que o projeto não diz, não que o recurso não exista.

| Ferramenta | Abordagem | Contas ao mesmo tempo | Uso de 5 horas e semanal | Quando um limite é atingido | Plataformas |
| --- | --- | --- | --- | --- | --- |
| ccam | Pastas separadas | Sim | Sim, `ccam usage` | Uma notificação na sessão; `ccam next --auto --continue` ou `/ccam:switch` leva a conversa para a conta com mais espaço, com ou sem histórico compartilhado; `ccam --fallback` executa de novo ali, sozinho, os jobs `-p`, e `ccam pick` distribui os workers | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Troca o login; pastas separadas com `cswap run` | Sim, com `cswap run` | Sim | `cswap auto` troca o login ao passar de um limiar; a sessão em andamento continua | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Pastas separadas | Sim | Em parte | Pools de contas passam as requisições para a próxima conta por um proxy local | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Pastas separadas | Sim | Sim | Rotas escolhem uma conta com cota antes de cada execução, podem esperar um reset e distribuem contas diferentes aos workers | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Troca o login; pastas separadas com `clauth start` | Sim, com `clauth start` | Sim | Troca seguindo uma cadeia de fallback | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Troca o login ou usa perfis isolados | Sim, com perfis isolados | Sim | `caam run` tenta de novo comandos não interativos em outra conta | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Troca o login e encapsula o `claude` | Não, uma conta ativa | Sim | Troca o login e retoma a mesma conversa | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Troca o login | Não | Sim | Troca opcional ao passar de um limiar | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Troca o login | Não | Sim | Troca opcional em 100% | macOS |

O ccswap é um fork do claude-swap que também gerencia logins do Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) e [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) também dão a cada conta a própria pasta. Várias dessas ferramentas também lidam com outros agentes de programação, como Codex ou Gemini.

## Como elas são construídas

| Ferramenta | Escrita em | Enquanto o Claude Code roda | Copia logins de um lugar para outro | O que as contas podem compartilhar |
| --- | --- | --- | --- | --- |
| ccam | Go, um único binário estático | Nada: o ccam se substitui pelo `claude` (no Windows, que não permite isso, ele espera o processo terminar; um relay opcional fica ao lado dele) | Não | Grupos: cada um indica os itens que compartilha e as contas dele; uma conta pode estar em vários |
| claude-swap | Python 3.12+, instalado com uv ou pipx | Nada com `cswap run`; `cswap auto` e o app da barra de menus continuam rodando | Sim, quando troca o login | Com `cswap run`: configurações, `CLAUDE.md` e extensões juntas, histórico opcional; servidores MCP copiados a cada início |
| CCS | Node.js 18+ ou Bun | O CCS fica como processo pai; os pools adicionam um proxy local | Não para as pastas de conta; os pools guardam os próprios logins | Compartilhado ou local por conta, histórico em grupos de contexto; servidores MCP copiados |
| clausona | TypeScript, Node.js 20+ | O `claude` simples passa por uma função do shell; `clausona run` mantém o Node como processo pai; um processo Node em segundo plano depois de cada saída | Não, mas renova sozinho os tokens das contas ociosas para ler o uso delas | Tudo exceto os arquivos próprios de cada conta, a partir de uma conta principal; histórico opcional; servidores adicionados com `claude mcp add` não são compartilhados |
| clauth | Rust | O clauth fica como processo pai | Sim, quando troca o login | Todo o `~/.claude`, ou nada |
| CAAM | Go | Depende do modo; um daemon opcional | Sim, no modo cofre | Definido pelo modo; servidores MCP copiados a cada início |
| cux | Go | Um wrapper fica em volta do `claude` | Sim | Tudo: uma pasta |
| CCSwitcher, Claude Switcher | Apps de macOS (Swift, Python) | O app da barra de menus | Sim | Tudo: uma pasta |

## O que o ccam faz bem

### Os logins ficam onde o Claude Code os colocou

O ccam nunca copia um login do Claude de um lugar para outro. O login de cada conta fica na própria pasta, ou, no macOS, no item das Chaves próprio dessa pasta, e o Claude Code é o programa que o renova. Não existe cópia salva que possa ficar dessincronizada. Um `~/.claude` existente registrado com `ccam import` mantém o item das Chaves que o `claude` puro usa, então continua logado nos dois casos.

O ccam só renova por conta própria um token expirado quando você executa `ccam usage --refresh`. Ele obtém um bloqueio antes e se recusa enquanto aquela conta estiver rodando, então nunca troca um token debaixo de uma sessão ativa.

### Nada roda entre você e o Claude Code

O ccam lê um arquivo, verifica seus links e se substitui pelo `claude` (`exec`). Enquanto você trabalha, não há processo do ccam, wrapper, proxy nem daemon: a entrada padrão, os sinais e o código de saída pertencem ao Claude Code, e as requisições dele vão direto para a Anthropic. O Windows não consegue substituir um processo, então lá o ccam inicia o Claude Code no mesmo console e espera o código de saída dele, sem fazer mais nada.

Iniciar uma conta assim acrescenta cerca de 4 ms, menos de 6 ms no percentil 95, com 8 MB de pico de memória ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) com hyperfine, 300 execuções). No mesmo notebook, só iniciar o interpretador Python leva 13 ms e o Node.js 31 ms, antes de um programa fazer qualquer coisa.

### Os grupos decidem o que é compartilhado

Servidores MCP, `settings.json`, `CLAUDE.md`, extensões (agentes, comandos, skills, plugins) e histórico de conversas são cinco itens separados. Um grupo indica os itens que compartilha e as contas que fazem parte dele, então duas contas de trabalho podem compartilhar o histórico em um grupo `work` enquanto uma conta pessoal mantém o seu, e as três compartilham os mesmos servidores MCP por meio do `default`. `ccam group set work --add memory` move todas as contas do grupo em um único plano; o [painel](/features/dashboard) faz o mesmo com caixas de seleção e deixa em cinza um item que uma conta já recebe de outro grupo.

- Os servidores MCP são passados com `--mcp-config` na inicialização, então por padrão o ccam nunca edita o `.claude.json`, o arquivo que o Claude Code reescreve enquanto roda. `ccam mcp add` adiciona um servidor uma única vez para o grupo inteiro.
- Os outros itens são links para a pasta do grupo, então uma mudança feita em uma conta aparece nas outras na hora. Nada é copiado na inicialização.
- Logins, `.claude.json` e o estado de execução nunca são compartilhados.

### Toda mudança é mostrada antes e tem backup

`ccam group` e `ccam migrate` mostram um plano e perguntam antes de mover qualquer coisa. Quando um arquivo difere entre uma conta e o grupo dela, você escolhe: manter o do grupo, usar o seu ou mesclar o JSON. Os originais vão para `.ccam-backup/<time>/`, a pasta de uma conta removida vai para uma pasta de lixeira, e os arquivos são substituídos com uma única renomeação, então uma gravação interrompida deixa o arquivo antigo. `--dry-run` mostra o que a inicialização, o compartilhamento, a migração, os atalhos, as atualizações e os reparos fariam. `ccam doctor --fix` corrige links quebrados e permissões.

### A mesma conversa na próxima conta

`ccam next --auto --continue` consulta o uso de todas as contas ao mesmo tempo, abre a que tem mais espaço e continua a conversa. Com o histórico compartilhado, é a mesma conversa; sem ele, o Claude Code copia a conversa para a nova conta (`--resume <file> --fork-session`), e o ccam informa quantos tokens a primeira resposta lê de novo. Quando um limite interrompe uma sessão, o hook do plugin do ccam mostra uma notificação, e `/ccam:switch` passa a conversa adiante. Você decide quando trocar; nada monitora seu uso em segundo plano, e as escolhas automáticas ficam dentro de uma mesma organização, então uma conversa da empresa nunca vai parar por acidente em uma conta pessoal. A linha de status sugere o comando quando uma conta chega a 90%.

Scripts não podem esperar por você, então `ccam --fallback work -p "…"` faz a troca por eles: quando a conta atinge o limite, o mesmo comando roda de novo com a conta que tem mais espaço, e só a saída da execução que terminou é impressa. Uma falha que não seja limite de uso é repassada como está. Para workers em paralelo, `ccam pick` entrega a cada um uma conta diferente com espaço, e os dois saem com código 75 quando nenhuma conta tem mais espaço.

### A conta certa em cada pasta

`ccam folder work ~/code/company` faz com que toda inicialização sem nome de conta dentro dessa pasta use `work`, enquanto o resto dos seus projetos continua com a conta padrão. A configuração fica no mesmo `config.toml` que toda inicialização lê, então não custa tempo.

### Um binário, downloads verificados

O ccam é um binário estático de cerca de 10 MB, sem runtime para instalar, para Linux, macOS e Windows. Os scripts de instalação (`install.sh`, `install.ps1`) e o `ccam update` verificam cada download com os checksums SHA-256 da versão, e o ccam nunca procura atualizações por conta própria. Não há telemetria: as únicas requisições que ele faz são consultas de uso das suas contas à Anthropic, em paralelo com limite de 3 segundos e cache de 60 segundos, e os downloads que você pede. Com `usage_source = "claude"`, até essas consultas passam pelo Claude Code, e o ccam nunca lida com um token.

## Onde outras ferramentas fazem mais

- claude-swap, clauth, cux, CAAM e CCSwitcher podem trocar sozinhos a conta de uma sessão interativa em andamento quando uma conta passa de um limiar de uso. O ccam troca execuções `-p` com `--fallback` e deixa as sessões interativas com você, com uma notificação no limite e `/ccam:switch`.
- No Windows sem o Modo de Desenvolvedor, o ccam compartilha pastas, mas não arquivos avulsos como `settings.json`; claude-swap e CCS os copiam nesse caso.
- CCSwitcher, Claude Switcher e claude-swap mostram o uso na barra de menus do macOS.
- CCS, clausona, CAAM, clauth, ccswap e Claude Switcher também gerenciam o Codex e outros agentes de programação.
- O clausona também liga e desliga extensões por conta (o `ccam ext` do ccam só as mostra), acompanha o custo local de tokens e traz plugins do Claude Code que rodam frotas de workers no tmux e em outros terminais.
- CAAM, CCSwitcher e Claude Switcher são instalados com Homebrew, e clauth e CAAM assinam os checksums das suas versões. O ccam é instalado com o próprio script, um `.deb` ou `go install`, e verifica os checksums sem assinatura.

## Qual escolher?

Escolha uma ferramenta que troca o login se você usa uma conta por vez ou quer que sessões em andamento troquem de conta automaticamente. Escolha o ccam se você quer contas lado a lado, uma única configuração compartilhada, a mesma conversa na próxima conta e nada rodando ao lado do Claude Code, no Linux, macOS ou Windows. Se você não precisa compartilhar nada, [aliases de shell simples](/compare/aliases) bastam.

Se você pode usar várias contas depende dos termos de cada conta. Verifique-os antes.
