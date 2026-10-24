---
title: ข้อมูลอ้างอิงคำสั่ง
metaTitle: "ข้อมูลอ้างอิงคำสั่ง ccam"
description: "ทุกคำสั่งและตัวเลือกของ ccam ตั้งแต่ ccam add และ ccam <name> ไปจนถึง ccam usage, ccam next และ ccam doctor"
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "36af7e209cd9"
---

ใส่ `--dry-run` ไว้หน้าคำสั่งเพื่อดูว่าคำสั่งนั้นจะทำอะไรโดยไม่ลงมือทำจริง เช่น `ccam --dry-run work` ตัวเลือกนี้ใช้ได้เมื่อเริ่ม Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) และกับ `group`, `migrate`, `shortcut`, `update` และ `doctor --fix` ส่วนคำสั่งอื่นที่เปลี่ยนแปลงบางอย่างจะปฏิเสธตัวเลือกนี้แทนที่จะทำงาน

<!-- commands -->

## การเริ่ม Claude Code

`ccam <name|number> [arguments]` เริ่ม Claude Code ในชื่อบัญชีนั้น ccam ไม่อ่านอาร์กิวเมนต์ที่อยู่หลังชื่อบัญชี แต่ส่งต่อให้ `claude` ตามลำดับเดิม ส่วนอินพุตมาตรฐาน รหัสออก (exit code) และสัญญาณ (signal) ทำงานเหมือนกับการใช้ `claude` โดยตรงทุกประการ จะใส่ `--` หลังชื่อบัญชีหรือไม่ก็ได้

`ccam` ที่มีอาร์กิวเมนต์แต่ไม่มีชื่อบัญชี เช่น `ccam -p "summarise"` จะใช้บัญชีที่ตั้งไว้สำหรับโฟลเดอร์ปัจจุบัน หากไม่มีจะใช้บัญชีเริ่มต้น ตั้งค่าทั้งสองได้ด้วย `ccam folder <name> [folder]` และ `ccam use <name>`

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` ที่วางไว้หน้าชื่อบัญชีใช้กับการรันแบบ `-p`/`--print` เมื่อบัญชีถึงขีดจำกัดการใช้งาน ccam จะรันคำสั่งเดิมอีกครั้งด้วยบัญชีที่เข้าสู่ระบบอยู่และเหลือโควตามากที่สุด แล้วพิมพ์เฉพาะผลลัพธ์ของการรันที่ทำงานจนจบ ดู [บัญชีตามโฟลเดอร์](/features/folder) และ [รันอีกครั้งเมื่อถึงขีดจำกัด](/features/next#run-again-at-the-limit)

**ที่นั่ง (seat) ของ Team และ Enterprise** การเลือกอัตโนมัติ (`ccam next`, `--fallback`, `ccam pick` และปุ่ม `n` ในแดชบอร์ด) จะสลับระหว่างที่นั่งต่าง ๆ ในองค์กร Claude Team หรือ Enterprise เดียวกันด้วย ศูนย์ช่วยเหลือของ Anthropic ระบุว่า "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (แต่ละที่นั่งในแพ็กเกจ Team จะกำหนดให้กับบุคคลที่ระบุชื่อเพียงคนเดียว และมีเพียงบุคคลนั้นที่ใช้ได้) และ "Users can't route their work through another user's account to use that user's remaining usage" (ผู้ใช้ไม่สามารถส่งงานของตนผ่านบัญชีของผู้ใช้คนอื่นเพื่อใช้ปริมาณการใช้งานที่เหลืออยู่ของผู้ใช้คนนั้นได้) ([แพ็กเกจ Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ส่วน [แพ็กเกจ Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ก็เช่นเดียวกัน) ดังนั้น การใช้การเลือกอัตโนมัติระหว่างที่นั่งต่าง ๆ จึงอาจละเมิดข้อกำหนดของ Anthropic โปรดสอบถามองค์กรของคุณและ Anthropic ก่อนที่จะพึ่งพาการเลือกอัตโนมัติ สำหรับแพ็กเกจส่วนตัว (Pro, Max) ดูที่ [หมายเหตุเรื่องข้อกำหนดในคำถามที่พบบ่อย](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)

คำสั่งลัดอย่าง `ccwork` เริ่มบัญชีของตัวเองในลักษณะเดียวกัน: `ccwork --resume` ทำงานเหมือนกับ `ccam work --resume` ดูที่ [คำสั่งลัด](/docs/commands#shortcuts)

## การเพิ่มและลบบัญชี

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ที่ไม่ระบุโฟลเดอร์จะค้นหา `~/.claude` และ `~/.claude-*` แล้วถามคุณทีละโฟลเดอร์ โฟลเดอร์ที่ไม่สำเร็จจะถูกรายงานเป็น `Skipped <path>: <reason>` และโฟลเดอร์อื่นจะถูกลงทะเบียนตามปกติ ส่วนโฟลเดอร์ที่ชื่อซ้ำกับบัญชีที่มีอยู่แล้วต้องใช้ `--name` ซึ่งระบุพร้อมกับโฟลเดอร์ โฟลเดอร์ที่นำเข้าจะถูกลงทะเบียนไว้ในตำแหน่งเดิม ccam ไม่ย้ายหรือคัดลอกการเข้าสู่ระบบใด ๆ

`ccam login --expired` ยังเข้าสู่ระบบให้บัญชีที่ Claude Code ล้างการเข้าสู่ระบบไปแล้วด้วย ส่วน `ccam remove` จะออกจากระบบบัญชีที่ ccam จัดการก่อนที่โฟลเดอร์จะไปอยู่ในถังขยะ (ยกเว้นเมื่อบัญชีอื่นถือการเข้าสู่ระบบเดียวกันอยู่) และ `ccam remove` ของบัญชีที่นำเข้าจะทิ้งสำเนาของสิ่งที่ใช้ร่วมกันไว้ให้กลุ่มของบัญชีนั้น บน macOS `ccam rename` จะคงชื่อโฟลเดอร์ไว้ เพราะพวงกุญแจ (Keychain) เก็บการเข้าสู่ระบบไว้ภายใต้พาธของโฟลเดอร์

`ccam add --api` เพิ่มบัญชีที่ใช้ API key แทนการเข้าสู่ระบบ Claude ccam อ่านคีย์ทุกครั้งที่เริ่มและไม่เคยเก็บไว้ ดู[บัญชี API key](/features/api)

`ccam rename`, `ccam remove` และคำสั่ง `ccam group` จะไม่ยอมเปลี่ยนแปลงบัญชีในขณะที่ Claude Code กำลังทำงานด้วยบัญชีนั้น ไม่ว่าจะเริ่มด้วยวิธีใด (รวมถึงเซสชันเบื้องหลัง) ให้ออกจากเซสชันนั้นก่อน คำสั่งเหล่านี้ยังปฏิเสธเมื่อการเปลี่ยนกลุ่มยังทำไม่เสร็จ: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``

`ccam list --json` จะแสดงทุกบัญชีพร้อมสล็อต คำสั่งลัด อีเมล และโฟลเดอร์ อาร์เรย์ `groups` ในผลลัพธ์จะแสดงกลุ่มของบัญชี และออบเจ็กต์ `share` จะจับคู่แต่ละรายการ (`mcp`, `settings`, `memory`, `extensions`, `history`) กับกลุ่มที่รายการนั้นมาจาก หรือ `local`

## คำสั่งลัด

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

คำสั่งลัดคือคำสั่งที่คุณตั้งชื่อเองเพื่อเริ่ม Claude Code ในชื่อบัญชีหนึ่ง เช่นเดียวกับ `ccam work` คำสั่งลัดจะส่งทุกอาร์กิวเมนต์ให้ Claude Code โดยไม่เปลี่ยนแปลง ดังนั้น `ccwork --resume` และ `ccwork -p "hi"` จึงทำงานตามที่คุณคาดไว้ คำสั่งลัดเป็นลิงก์ไปยัง ccam ใน `~/.local/bin`

คำสั่งลัดเป็นตัวเลือกเสริม ตัวช่วยตั้งค่าของ `ccam add` จะถามว่าจะสร้างหรือไม่ ส่วน `ccam add <name> --shortcut <command>` จะสร้างให้โดยไม่ถาม และ `ccam list` จะแสดงคำสั่งลัดของแต่ละบัญชี หรือแสดง `-` หากไม่มี

ชื่อใช้ได้เฉพาะอักขระ a-z, 0-9, `.`, `-` และ `_` ชื่อต้องไม่ซ้ำกับคำสั่งของ ccam คำสั่งลัดของบัญชีอื่น หรือชื่อของโปรแกรมอื่นใน `PATH` ของคุณ และต้องไม่ใช่ไฟล์ใน `~/.local/bin` ที่ ccam ไม่ได้สร้าง เช่น สำเนาหรือลิงก์ของ ccam ที่คุณสร้างเอง คำสั่งลัดที่แก้ไขด้วยมือใน `config.toml` แล้วผิดกฎเหล่านี้ เป็น `claude` หรือ `ccam` หรือถูกกำหนดให้สองบัญชี จะทำให้โหลด `config.toml` ไม่สำเร็จ โดยระบุชื่อบัญชีนั้น

`ccam rename` จะคงชื่อคำสั่งลัดไว้ และ `ccam remove` จะลบคำสั่งลัดไปพร้อมกับบัญชี `ccam shims install` จะสร้างลิงก์ของทุกบัญชีขึ้นใหม่ เช่น หลังจากที่คุณย้าย ccam ไปที่อื่น ไฟล์ที่ขวางคำสั่งลัดหนึ่งอยู่จะไม่ทำให้คำสั่งลัดอื่นหยุดไปด้วย

## การใช้ร่วมกัน

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

กลุ่มหนึ่งใช้รายการร่วมกันระหว่างบัญชีในกลุ่ม บัญชีหนึ่งอยู่ได้หลายกลุ่ม แต่แต่ละรายการมาจากกลุ่มเดียวเท่านั้น ดูที่ [การใช้ร่วมกันด้วยกลุ่ม](/features/share)

## การคัดลอกการตั้งค่าที่มีอยู่

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` คัดลอก `settings.json`, `CLAUDE.md` ส่วนขยาย และเซิร์ฟเวอร์ MCP ของผู้ใช้จากโฟลเดอร์ Claude Code ไปยังกลุ่ม เพื่อให้ทุกบัญชีที่ใช้รายการเหล่านั้นร่วมกันได้รับด้วย โฟลเดอร์ต้นทางจะถูกอ่านอย่างเดียว ตัวเลือก `--settings`, `--memory`, `--extensions`, `--mcp` และ `--history` จะคัดลอกเฉพาะรายการที่คุณระบุ หากไม่ระบุ บนเทอร์มินัลจะมีรายการตรวจสอบ (checklist) ให้คุณเลือก โดยเลือกทุกอย่างไว้ยกเว้นประวัติการสนทนา ส่วนในกรณีอื่นจะคัดลอกทุกอย่างยกเว้นประวัติการสนทนา `-g` หรือ `--group` ใช้เลือกกลุ่ม และ `-y` หรือ `--yes` จะคัดลอกโดยไม่ถาม และเก็บสำเนาของกลุ่มไว้เมื่อไฟล์ไม่ตรงกัน การลงชื่อเข้าใช้ MCP จะไม่ถูกคัดลอก เพราะการลงชื่อเข้าใช้ที่ถูกคัดลอกจะทำให้ผู้ถือรายอื่นหลุดจากระบบ แผนจะระบุชื่อบัญชีที่ต้องลงชื่อเข้าใช้ด้วย `/mcp` นอกจากนี้ `--mcp` ยังใช้ได้เมื่อ `~/.claude` ว่างเปล่าหรือไม่มีอยู่ แต่ `~/.claude.json` มีเซิร์ฟเวอร์อยู่ ดูที่ [การคัดลอกการตั้งค่าที่มีอยู่](/docs/sharing#copy-an-existing-setup)

## เซิร์ฟเวอร์ MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` จะย้ายเซิร์ฟเวอร์ที่มีเฉพาะในบัญชีนั้นไปยังกลุ่มของบัญชี รวมถึงเซิร์ฟเวอร์ที่บัญชีนั้นมีไว้สำหรับโฟลเดอร์เดียว (`claude mcp add` แบบธรรมดา) ซึ่งหลังจากนั้นจะโหลดในทุกโฟลเดอร์ เซิร์ฟเวอร์ที่ใช้ชื่อซึ่งกลุ่มใช้อยู่แล้วกับนิยามอื่นจะถูกข้ามไป เพิ่ม `--force` เพื่อแทนที่สำเนาของกลุ่ม โดย ccam จะสำรองสำเนานั้นไว้ก่อน

## การใช้งานและการสลับบัญชี

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` และ `ccam pick` เลือกด้วยวิธีเดียวกัน: บัญชีที่ต่ำกว่า `auto_max_usage` (`config.toml` ค่าเริ่มต้น 100) มาก่อน ตาม `auto_strategy` (`headroom`, `expiring` หรือ `balance`) โดยนับช่วงรายสัปดาห์ของ `--model` ด้วย ไม่เลือกบัญชีที่อยู่ที่ 100% เด็ดขาด (ช่วงที่ Claude Code รายงานว่าเต็มจะนับเป็น 100% จนกว่าจะรีเซ็ต) ไม่เลือกบัญชีที่เข้าสู่ระบบด้วย API key และเลือกเฉพาะภายในองค์กรของบัญชีที่ใช้เป็นจุดเริ่มต้น (ใช้ `--cross-org` หรือ `cross_org = true` เพื่อมองไกลกว่านั้น) `--wait` รับค่าหลัง `=` (`--wait=2h`) หากเขียนโดยเว้นวรรคจะถูกปฏิเสธพร้อมคำแนะนำ เมื่อเลือกบัญชีใดไม่ได้ จะออกด้วยรหัส 75 หรือเมื่อใช้ `--wait` จะรอให้บัญชีแรกรีเซ็ต `ccam next --continue` พาการสนทนาไปต่อได้แม้ระหว่างบัญชีที่ไม่ได้ใช้ประวัติร่วมกัน โดย Claude Code จะคัดลอกการสนทนาไปยังบัญชีใหม่ ต่างจาก `ccam <name>` ตรงที่ `ccam next` อ่านตัวเลือกของตัวเอง ดังนั้นตัวเลือกของ Claude Code ต้องวางไว้หลัง `--` หากไม่มี ccam จะหยุดพร้อม `unknown flag` ภายใน Claude Code `/ccam:switch [account]` ของปลั๊กอินจะส่งต่อการสนทนา และ `ccam hook` คือสิ่งที่ hook ของปลั๊กอินรัน ดู[การสลับไปยังบัญชีถัดไป](/features/next) [การเลือกบัญชีสำหรับสคริปต์](/features/pick) และ[เซสชันเบื้องหลังของทุกบัญชี](/features/agents)

## การบำรุงรักษา

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` จะตรวจสอบ checksum แบบ SHA-256 ของรุ่นนั้นก่อนแทนที่ ccam และ `--version v0.2.0` จะติดตั้งรุ่นที่ระบุ ส่วน ccam ที่ติดตั้งจากแพ็กเกจ `.deb` ให้อัปเดตด้วย apt แทน ดูที่ [การอัปเดต](/docs/install#update)

## ไฟล์

| พาธ | เนื้อหา |
| --- | --- |
| `~/.config/ccam/config.toml` | บัญชี กลุ่ม และการตั้งค่า (`theme`, `usage_source`, `relay`, `cross_org` …) การตั้งค่าต้องอยู่ด้านบนสุด เหนือ `[[group]]`, `[[profile]]` หรือ `[[folder]]` ตัวแรก: การตั้งค่าที่อยู่ใต้ตารางจะเป็นของตารางนั้น ccam จึงระบุถึงการตั้งค่านั้นทุกครั้งที่เปิดและใน `ccam doctor` และคำสั่งที่บันทึกการตั้งค่าจะปฏิเสธจนกว่าจะย้ายมันขึ้นไป ความคิดเห็น (comment) จะไม่ถูกเก็บไว้เมื่อ ccam บันทึกไฟล์ |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ของแต่ละบัญชี |
| `~/.local/share/ccam/groups/<group>/` | การตั้งค่า เซิร์ฟเวอร์ MCP ส่วนขยาย และประวัติที่ใช้ร่วมกัน |
| `~/.local/state/ccam/` | แคชการใช้งานและบันทึกขีดจำกัด (`usage/<name>.json`, `usage/<name>.limit.json`) เวลาที่ใช้ล่าสุด ตัวบ่งชี้เซสชันและล็อก (`locks/`) คำสั่งลัดที่ ccam สร้าง (`shortcuts.json`) การเปลี่ยนกลุ่มที่ยังไม่เสร็จ (`group-change.json`) และบันทึกที่ `/ccam:switch` และ hook ขีดจำกัดทิ้งไว้ (`handoff/<account>/<session>.json`) |
| `~/.local/bin/<shortcut>` | คำสั่งลัดอย่าง `ccwork` ซึ่งเป็นลิงก์ไปยัง ccam |

ตัวแปร `XDG_CONFIG_HOME`, `XDG_DATA_HOME` และ `XDG_STATE_HOME` ใช้เปลี่ยนตำแหน่งของโฟลเดอร์เหล่านี้ บน Windows `~` คือ `%USERPROFILE%` และคำสั่งลัดเป็น `<shortcut>.exe`
