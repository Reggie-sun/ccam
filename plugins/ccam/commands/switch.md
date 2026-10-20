---
description: Carry this conversation on as another ccam account (the next one, or the one you name)
argument-hint: "[account]"
---

The ccam plugin's hook handles this command before it reaches you: it notes this conversation for `ccam next --continue` and ends the turn. If you are reading this, the hook did not run (Claude Code before 2.1.78, the ccam plugin's hooks turned off, or `ccam` not on PATH).

Tell the user, briefly: type `/exit`, then run `ccam next --continue` in this folder (add `--auto` to take the account with the most room). ccam then opens this folder's last conversation as another account; Claude Code copies it over when that account does not share the history.
