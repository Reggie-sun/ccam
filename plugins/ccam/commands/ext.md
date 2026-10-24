---
description: Show which skills, commands, MCP servers and hooks each ccam account loads, and any problems
allowed-tools: Bash(ccam ext:*)
---

!`ccam ext`

Point out the problems listed at the end, if any, and what to run for each (`ccam doctor --fix` for broken share links, `ccam mcp import <account> <server>` to share a server, `/mcp` in a session of the account for a missing sign-in, `ccam migrate --mcp` for servers only `~/.claude.json` has). Mention items that only some accounts have.
