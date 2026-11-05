# Using ccam with oh-my-claudecode and Claude Unlimited

Notes from running ccam next to [oh-my-claudecode](https://github.com/Yeachan-Heo/oh-my-claudecode) (OMC) and
[Claude Unlimited](https://github.com/DevDock-AI/claude-unlimited) (`cu`). Everything here was tried on Linux.

## Subagents cannot use another account

OMC's built-in subagents (`executor`, `explore`, ...) run inside the current Claude Code process and use that
session's login. ccam picks an account when `claude` starts, so it cannot move a subagent to another account.
Only work that runs in a separate `claude` process can use another account:

- **A dispatcher subagent.** `ccam-worker.agent.md` in this folder: copy it to your shared `agents/` directory.
  It runs a self-contained task as `ccam --fallback "$(ccam pick --strategy headroom)" -p "..."`, so a limit hit
  mid-run reruns on the account with the most room. The headless session sees none of the parent conversation, so
  the task prompt must stand alone.
- **OMC `/team` workers.** They start as separate `claude` processes with `OMC_TEAM_WORKER` set. A wrapper placed
  first on `PATH` can route them through ccam:

  ```sh
  #!/bin/sh
  if [ -n "$OMC_TEAM_WORKER" ] && [ -z "$CCAM_POOL_PICKED" ] && acct="$(ccam pick --strategy headroom 2>/dev/null)" && [ -n "$acct" ]; then
    CCAM_POOL_PICKED=1 exec ccam "$acct" "$@"
  fi
  exec <your original claude launcher> "$@"
  ```

  The account is chosen once, at worker start. `ccam pick` deliberately spreads calls made within two minutes, so
  it is not strictly "emptiest first".

## `auto_strategy`

`weekly-reset` takes the account whose week resets first among those with at least 60% of the five-hour window
left, so a 39%-used account can win over a 0%-used one. Use `headroom` to always take the lowest five-hour usage.

## Claude Unlimited next to ccam

- **Do not copy ccam's credentials into `cu`.** Both tools would hold the same OAuth refresh token. Whichever
  refreshes first rotates it, and the other copy is rejected with `invalid_grant`. This was observed: the imported
  `cu` profiles died within hours while ccam's logins kept working. Log each account into `cu` separately with
  `cu add-account`, which creates an independent session.
- **A launcher that strips `ANTHROPIC_*` breaks `cu code`.** `cu code` routes `claude` through its local proxy with
  `ANTHROPIC_BASE_URL` and a placeholder token. If the `claude` on `PATH` is a wrapper that removes `ANTHROPIC_*`
  variables, the session silently talks to Anthropic directly and `cu` serves nothing (every profile's
  `last_used_at` stays empty). Put the real Claude Code binary first on `PATH` for `cu code`:

  ```sh
  mkdir -p ~/tools/cu-bin
  ln -s /path/to/@anthropic-ai/claude-code/bin/claude.exe ~/tools/cu-bin/claude
  alias cud='PATH="$HOME/tools/cu-bin:$PATH" cu code --distribute'
  ```

- With `--distribute`, the main agent and each subagent start on the least-busy account. Four parallel subagents
  landed on four different accounts in a test; one account that was at 100% of its five-hour window was skipped.
- Usage adds up on Anthropic's side but not in each tool's view, so `ccam usage` under-reports an account that is
  also used through `cu`.
- `cu` stores and refreshes subscription OAuth tokens itself. Its own README says this may conflict with
  Anthropic's terms for subscription accounts. Read it before adding an account.
