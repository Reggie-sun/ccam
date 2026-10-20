---
title: "Trocar de conta no Claude Code: ccam e alternativas"
description: Como o ccam se compara a claude-swap, CCS, clauth, CAAM, cux e outras ferramentas para trocar de conta no Claude Code, e qual abordagem serve para cada uso.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

As ferramentas para trocar de conta no Claude Code seguem uma de duas abordagens. A maioria troca o login dentro da única pasta `~/.claude`, então só uma conta fica ativa por vez. Outras, entre elas o ccam, dão a cada conta a própria pasta com `CLAUDE_CONFIG_DIR`, então todas as contas continuam logadas e várias podem rodar ao mesmo tempo. O ccam acrescenta configurações e servidores MCP compartilhados e continua uma conversa na pasta de outra conta.

## Trocar um login ou usar pastas separadas?

**Trocar um login** mantém um único conjunto de configurações, servidores MCP e histórico, porque só existe uma pasta. A troca é rápida, e uma sessão em andamento pode passar a usar o novo login. Mas só uma conta fica ativa por vez, para todos os terminais e editores ao mesmo tempo.

**Pastas separadas** é o método que [a documentação do Claude Code descreve](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Todas as contas continuam logadas, e você pode executar uma sessão de trabalho e uma pessoal lado a lado. O preço é que cada pasta começa com as próprias configurações, servidores MCP e histórico, a menos que alguma ferramenta os compartilhe.

## Comparação das ferramentas

Retirado da documentação de cada projeto em 2026-10-19. "Não documentado" significa que o projeto não diz, não que o recurso não exista.

| Ferramenta | Abordagem | Contas ao mesmo tempo | Uso de 5 horas e semanal | Quando um limite é atingido | Plataformas |
| --- | --- | --- | --- | --- | --- |
| ccam | Pastas separadas | Sim | Sim, `ccam usage` | Você executa `ccam next --auto --continue` para continuar a conversa na conta com mais espaço | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Troca o login; pastas separadas com `cswap run` | Sim, com `cswap run` | Sim | `cswap auto` troca o login ao passar de um limiar; a sessão em andamento continua | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Pastas separadas | Sim | Em parte, em um app da barra de menus do macOS | Pools de contas passam as requisições para a próxima conta por um proxy local | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Troca o login; pastas separadas com `clauth start` | Sim, com `clauth start` | Sim | Troca seguindo uma cadeia de fallback | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Troca o login ou usa perfis isolados | Sim, com perfis isolados | Sim | `caam run` tenta de novo comandos não interativos em outra conta | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Troca o login e encapsula o `claude` | Não documentado | Sim | Troca o login e retoma a mesma conversa | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Troca o login | Não | Sim | Não documentado | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Troca o login | Não | Sim | Troca opcional em 100% | macOS |

O ccswap é um fork do claude-swap que também gerencia logins do Codex. Várias dessas ferramentas também lidam com outros agentes de programação, como Codex ou Gemini.

## No que o ccam se concentra

- Todas as contas continuam logadas. Trocar nunca desconecta uma conta, e `ccam work` e `ccam personal` podem rodar em dois terminais ao mesmo tempo.
- Servidores MCP, `settings.json`, `CLAUDE.md`, agentes, skills e plugins são compartilhados entre as pastas por meio de links, item por item e conta por conta. Logins nunca são compartilhados.
- Com o histórico compartilhado, um único comando, `ccam next --auto --continue`, abre a mesma conversa na pasta da conta com mais espaço.
- O ccam inicia o Claude Code e se substitui por ele, então nada fica rodando em segundo plano. Ele não monitora seu uso nem troca de conta sozinho; você decide quando.

## Onde outras ferramentas fazem mais

- claude-swap, CCS, clauth, CAAM e cux rodam no Windows sem WSL. O ccam não.
- claude-swap, clauth e cux trocam sozinhos quando uma conta passa de um limiar de uso. O ccam deixa a troca com você.
- CCSwitcher, Claude Switcher e claude-swap mostram o uso na barra de menus do macOS.
- CCS, CAAM, ccswap e Claude Switcher também gerenciam o Codex e outros agentes de programação.

## Qual escolher?

Escolha uma ferramenta que troca o login se você usa uma conta por vez, quer que a troca aconteça automaticamente ou está no Windows sem WSL. Escolha o ccam se você quer contas lado a lado, uma única configuração compartilhada e a mesma conversa na próxima conta. Se você não precisa compartilhar nada, [aliases de shell simples](/compare/aliases) bastam.

Se você pode usar várias contas depende dos termos de cada conta. Verifique-os antes.
