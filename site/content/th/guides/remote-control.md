---
title: การใช้ Remote Control กับ ccam
description: เริ่ม Remote Control ของ Claude Code ในชื่อบัญชีใดก็ได้ด้วย ccam <name> rc ดูตัวเลือกที่มีผล เหตุที่ ccam เสนอคัดลอกเซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน และการแก้ปัญหา
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc` เริ่ม Remote Control ของ Claude Code ในชื่อบัญชี `work` คุณจึงควบคุมเซสชันบนคอมพิวเตอร์เครื่องนี้ได้จาก claude.ai/code หรือแอป Claude บนมือถือ คำสั่งนี้ทำงานเหมือน `claude rc` โดยต่างกันเพียงข้อเดียว: หากบัญชีใช้เซิร์ฟเวอร์ MCP ร่วมกัน ccam จะถามหนึ่งครั้งว่าจะคัดลอกเซิร์ฟเวอร์เหล่านั้นลงใน `.claude.json` ของบัญชีหรือไม่ เพราะเซสชันของ Remote Control อ่านเซิร์ฟเวอร์จากที่นั่นเท่านั้น

## การเริ่มใช้งาน

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[คำสั่งลัด](/features/shortcuts)ใช้ได้แบบเดียวกัน: `ccwork rc` คือ `ccam work rc` ให้รันคำสั่งในโฟลเดอร์ที่คุณต้องการทำงาน แล้วเปิด claude.ai/code หรือแอป Claude ที่ลงชื่อเข้าใช้ด้วยบัญชี Claude เดียวกับ `work`

## ตัวเลือก

ตัวเลือกสองประเภทมาบรรจบกันที่นี่: ตัวเลือกของ ccam เอง ซึ่งอยู่ก่อนชื่อบัญชีหรืออยู่ในการตั้งค่าของ ccam และตัวเลือกของ Remote Control ซึ่งอยู่หลัง `rc` และส่งถึง Claude Code โดยไม่เปลี่ยนแปลง

### ตัวเลือกและการตั้งค่าของ ccam

| ตัวเลือกหรือการตั้งค่า | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--dry-run` (ก่อนชื่อบัญชี) | พิมพ์ตัวแปรสภาพแวดล้อมและคำสั่ง รวมถึงข้อความแจ้งเรื่อง MCP หากมี ไม่ถามและไม่คัดลอก | ปิด |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | ถามเมื่อบัญชีใช้เซิร์ฟเวอร์ MCP ร่วมกันและ `.claude.json` ของบัญชียังไม่มีรายการปัจจุบันของกลุ่ม `y` คัดลอกเซิร์ฟเวอร์ ส่วน `N` หรือ Enter จะเริ่มโดยไม่มีเซิร์ฟเวอร์เหล่านั้น | `N` |
| `mcp_strategy = "sync"` ใน `~/.config/ccam/config.toml` | คัดลอกเซิร์ฟเวอร์ของกลุ่มลงใน `.claude.json` ก่อนการเริ่มทุกครั้ง ไม่ว่าจะเป็น Remote Control หรือไม่ จึงไม่มีการถามเลย | `"flag"` |
| `ccam share work --mcp off` | เลิกใช้เซิร์ฟเวอร์ MCP ร่วมกับบัญชีนี้ จากนั้น Remote Control จะเริ่มด้วยเซิร์ฟเวอร์ของบัญชีเองเท่านั้นและไม่มีการถาม | เปิด |

### ตัวเลือกของ Remote Control

ทุกอย่างหลัง `rc` เป็นของ Claude Code ตารางนี้คือตัวเลือกของ `claude remote-control` ใน Claude Code 2.1.304 ส่วน `ccam work rc --help` จะพิมพ์รายการสำหรับเวอร์ชันของคุณ

| ตัวเลือก | สิ่งที่ทำ |
| --- | --- |
| `--name <name>` | ชื่อของเซสชันที่แสดงใน claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | คำนำหน้าของชื่อเซสชันที่สร้างขึ้นอัตโนมัติ (ค่าเริ่มต้น: ชื่อโฮสต์) |
| `-c`, `--continue` | กลับไปเชื่อมต่อกับเซสชันที่ Remote Control บันทึกไว้ล่าสุดสำหรับโฟลเดอร์นี้ ภายในประมาณ 4 ชั่วโมง |
| `--session-id <id>` | กลับไปเชื่อมต่อกับเซสชันที่ระบุ |
| `--permission-mode <mode>` | โหมดสิทธิ์ของเซสชันใหม่: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` หรือ `plan` |
| `--spawn <mode>` | `same-dir` (ค่าเริ่มต้น), `worktree` (git worktree หนึ่งชุดต่อเซสชัน) หรือ `session` (เซสชันเดียว และออกเมื่อเซสชันจบ) |
| `--capacity <N>` | จำนวนเซสชันสูงสุดพร้อมกันในโหมด `same-dir` และ `worktree` (ค่าเริ่มต้น 32) |
| `--[no-]create-session-in-dir` | สร้างเซสชันแรกในโฟลเดอร์ปัจจุบันตอนเริ่ม (ค่าเริ่มต้น: เปิด) |
| `--[no-]chrome` | Claude in Chrome สำหรับเซสชันใหม่ (ค่าเริ่มต้น: ปิด) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | เอาต์พุตสำหรับดีบักและแบบละเอียด |

ให้วาง `rc` ต่อจากชื่อบัญชีทันที ccam ถือว่าอาร์กิวเมนต์แรกที่ไม่ขึ้นต้นด้วย `-` คือคำสั่งของ Claude Code ดังนั้นใน `ccam work rc --name laptop` ccam จะเห็น `rc` แต่ใน `ccam work --model opus rc` จะเห็น `opus`

## สิ่งที่ ccam ทำ

การเริ่มจะผ่านเส้นทางปกติ ([การเริ่ม Claude Code](/features/launch)) โดยมีการเปลี่ยนแปลงหนึ่งจุดในขั้นตอน MCP:

1. ccam หาบัญชี ตรวจลิงก์ที่ใช้ร่วมกันของบัญชี และหา `claude` เหมือนการเริ่มทุกครั้ง
2. ccam เห็นว่าคำสั่งคือ `rc` หรือ `remote-control` จึงไม่ใส่ `--mcp-config`
3. หากบัญชีใช้เซิร์ฟเวอร์ MCP ร่วมกัน ccam จะเปรียบเทียบ `mcp.json` ของกลุ่มกับ `sync-stamp` ใน `~/.local/state/ccam/mcp/work/` หาก `.claude.json` มีเซิร์ฟเวอร์ปัจจุบันอยู่แล้ว หรือกลุ่มไม่มีเซิร์ฟเวอร์เลย ccam จะข้ามไปขั้นตอนที่ 6 ทันที
4. มิฉะนั้น ccam จะระบุชื่อเซิร์ฟเวอร์ที่เซสชันจะไม่ได้รับ:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. เมื่อตอบ `y` ccam จะบันทึกสำเนาสำรองของ `.claude.json` ใช้ล็อกของ Claude Code เอง เพิ่มเซิร์ฟเวอร์ของกลุ่มลงใน `mcpServers` และคงเซิร์ฟเวอร์อื่นของบัญชีไว้ เมื่อตอบ `N` เมื่อไม่มีเทอร์มินัล หรือขณะที่บัญชีกำลังทำงาน ccam จะไม่เปลี่ยนแปลงอะไรและแจ้งให้ทราบ:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. ccam บันทึกเซสชัน ตั้ง `CLAUDE_CONFIG_DIR` เป็นโฟลเดอร์ของบัญชี แล้วแทนที่ตัวเองด้วย `claude rc` ผ่าน `exec`

จากนั้น Remote Control จะทำงานต่อไปในฐานะเซิร์ฟเวอร์ ทุกเซสชันที่มันเริ่มให้คุณเป็นโปรเซสลูกของโปรเซสนั้น จึงทำงานในฐานะ `work`: ใช้การเข้าสู่ระบบ การตั้งค่า และประวัติการสนทนาของ `work` และนับรวมในขีดจำกัดการใช้งานของ `work` ตราบใดที่เซิร์ฟเวอร์ยังทำงาน ccam จะถือว่าบัญชีนั้นกำลังทำงานอยู่

การทดลองรัน (dry run) แสดงให้เห็นความแตกต่าง การเริ่มปกติส่งเซิร์ฟเวอร์ของกลุ่มเป็นอาร์กิวเมนต์ แต่ Remote Control ไม่ส่ง:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## ทำไมจึงทำงานแบบนี้

**Remote Control ไม่รับ `--mcp-config`** ในการเริ่มปกติ ccam ส่งเซิร์ฟเวอร์ของกลุ่มให้ Claude Code ในรูป `--mcp-config=<path>` และไม่เคยเขียนไฟล์ของ Claude Code แต่ `claude remote-control` ปฏิเสธตัวเลือกนี้เมื่ออยู่ก่อนชื่อคำสั่ง และไม่รับเลยเมื่ออยู่หลังชื่อคำสั่ง ก่อน ccam 0.6.0 `ccam work rc` จึงล้มเหลวกับทุกบัญชีที่ใช้เซิร์ฟเวอร์ MCP ร่วมกัน ตอนนี้ ccam ไม่ใส่ตัวเลือกนี้แล้ว

**เซสชันของ Remote Control อ่านเซิร์ฟเวอร์จาก `.claude.json` เท่านั้น** เซสชันที่ Remote Control เริ่มไม่ได้สืบทอดตัวเลือกจากบรรทัดคำสั่ง ที่เดียวที่เซสชันเหล่านั้นพบเซิร์ฟเวอร์ MCP ของผู้ใช้คือ `.claude.json` ของบัญชี ไฟล์นี้ยังเก็บรายละเอียดการเข้าสู่ระบบของบัญชีและสถานะแยกตามโฟลเดอร์ด้วย ccam จึงลิงก์ไฟล์นี้ไปยังสำเนาที่ใช้ร่วมกันแบบที่ทำกับ `settings.json` ไม่ได้ การคัดลอกเซิร์ฟเวอร์ลงในไฟล์นี้จึงเป็นวิธีเดียวที่จะให้เซสชันได้รับรายการที่ใช้ร่วมกัน

**ccam ถามก่อนเขียน** ด้วยกลยุทธ์ flag ซึ่งเป็นค่าเริ่มต้น ccam สัญญาว่าจะไม่แตะ `.claude.json` Remote Control เป็นกรณีเดียวที่ ccam ต้องเขียนลงไฟล์นี้ จึงถามก่อน เก็บสำเนาสำรองไว้ และไม่แตะสิ่งใดนอกจาก `mcpServers`

**ccam ไม่เขียนขณะที่บัญชีกำลังทำงาน** Claude Code เขียน `.claude.json` ใหม่ระหว่างเซสชัน การเปลี่ยนแปลงจากภายนอกจึงอาจสูญหายหรือถูกเขียนทับ ccam จึงข้ามการคัดลอกเมื่อมีเซสชันใดของบัญชีกำลังทำงาน รวมถึงเซิร์ฟเวอร์ Remote Control อีกตัวด้วย

**ccam ถามหนึ่งครั้งต่อการเปลี่ยนแปลง** หลังจากคัดลอกแล้ว `sync-stamp` จะเก็บเวลาแก้ไขของ `mcp.json` ของกลุ่มไว้ ccam จะถามอีกครั้งเฉพาะเมื่อไฟล์นั้นเปลี่ยน เช่น หลังจาก `ccam mcp add` หากคุณตอบ `N` ก็ยังไม่มีสำเนา ccam จึงจะถามทุกครั้งที่เริ่ม Remote Control

**`ccam mcp` ทำให้สำเนาเป็นปัจจุบันอยู่เสมอ** เมื่อชื่อซ้ำกัน รายการใน `.claude.json` จะมีผลเหนือกว่า `--mcp-config` สำเนาเก่าจึงอาจบังนิยามที่ใหม่กว่าในการเริ่มปกติด้วย ด้วยเหตุนี้ `ccam mcp add`, `add-json`, `remove` และ `import` จึงอัปเดตสำเนาของทุกบัญชีที่ไม่ได้กำลังทำงาน:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**เมื่อไม่มีเทอร์มินัล คำตอบคือไม่** Remote Control มักรันจากสคริปต์ เซอร์วิส หรือ `nohup` ซึ่งไม่มีใครตอบคำถามได้ ในกรณีนั้น ccam จะเริ่มโดยไม่มีเซิร์ฟเวอร์เหล่านั้นแทนที่จะรอ

## ปัญหาและวิธีแก้

### เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกันไม่ปรากฏในเซสชันของ Remote Control

คำถามถูกตอบว่า `N` ไม่มีเทอร์มินัล หรือบัญชีกำลังทำงานอยู่ ให้หยุด Remote Control ปิดเซสชันอื่นของบัญชี แล้วรัน `ccam work rc` อีกครั้งในเทอร์มินัลและตอบ `y` หากไม่ต้องการให้ถามอีกเลย ให้ตั้ง `mcp_strategy = "sync"` ใน `~/.config/ccam/config.toml` เซิร์ฟเวอร์ที่อยู่ในเซสชันแล้วจะไม่เปลี่ยนระหว่างที่เซสชันทำงาน ดังนั้นให้เริ่มเซสชันใหม่หลังจากคัดลอก

### "Account work is running, so its .claude.json is left alone"

ยังมีเซสชันของ `work` เปิดอยู่: อาจเป็น `ccam work` ปกติ คำสั่งลัด หรือเซิร์ฟเวอร์ Remote Control อีกตัว ให้ออกจากเซสชันนั้นแล้วรัน `ccam work rc` อีกครั้ง เซสชันที่บัญชีอื่นรันอยู่ไม่นับรวม

### ccam ถามทุกครั้ง

คุณยังไม่ได้คัดลอกเซิร์ฟเวอร์ ให้ตอบ `y` หนึ่งครั้ง หรือตั้ง `mcp_strategy = "sync"` หากไม่ต้องการเซิร์ฟเวอร์เหล่านี้ใน Remote Control เลย ให้เลิกใช้ MCP ร่วมกับบัญชีนี้: `ccam share work --mcp off`

### `ccam work rc` ล้มเหลวทันทีพร้อมข้อผิดพลาดเกี่ยวกับ `--mcp-config`

ccam ของคุณเก่ากว่า 0.6.0 และส่ง `--mcp-config` ให้ Remote Control ให้รัน `ccam update` หากพบข้อผิดพลาดเดียวกันกับ ccam เวอร์ชันปัจจุบัน แปลว่ามีตัวเลือกที่มีค่าอยู่ก่อน `rc` เช่นใน `ccam work --model opus rc` ccam จึงถือว่า `opus` เป็นคำสั่ง ให้วาง `rc` ต่อจากชื่อบัญชีทันที

### การเปลี่ยนแปลงที่ทำด้วย `ccam mcp` ไม่ไปถึง Remote Control

`ccam mcp add` และ `remove` อัปเดตสำเนาเฉพาะในบัญชีที่ไม่ได้กำลังทำงาน และแจ้งให้ทราบ: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` ให้หยุด Remote Control แล้วเริ่มใหม่ ccam จะเห็นว่า `mcp.json` เปลี่ยนไปและถามอีกครั้ง

### สำเนาแทนที่เซิร์ฟเวอร์ของฉันเอง

หากบัญชีมีเซิร์ฟเวอร์ของตัวเองที่ชื่อเดียวกับเซิร์ฟเวอร์ของกลุ่ม การคัดลอกจะแทนที่เซิร์ฟเวอร์นั้น `.claude.json` ก่อนหน้าอยู่ใน `~/.local/state/ccam/mcp/work/backups/` ซึ่งเก็บไว้ 5 ชุดล่าสุด หากต้องการใช้เวอร์ชันของคุณกับทุกบัญชี ให้ย้ายเข้าไปในกลุ่มด้วย `ccam mcp import work <server> --force` ก่อนเริ่ม Remote Control ครั้งถัดไป

### ฉันต้องการลบเซิร์ฟเวอร์ที่คัดลอกมา

`ccam mcp remove <server>` ลบเซิร์ฟเวอร์ออกจากกลุ่มและจากทุกสำเนา หากต้องการกลับไปยังสถานะก่อนคัดลอก ให้กู้คืนสำเนาสำรองล่าสุดจาก `~/.local/state/ccam/mcp/work/backups/` ขณะที่ไม่มีเซสชันของบัญชีทำงานอยู่

### Remote Control ถามว่าจะเชื่อถือโฟลเดอร์หรือไม่ หรือออกไปเมื่อไม่มีเทอร์มินัล

การเชื่อถือโฟลเดอร์เก็บแยกตามบัญชีใน `.claude.json` ดังนั้นโฟลเดอร์ที่คุณเชื่อถือในฐานะ `personal` จะเป็นโฟลเดอร์ใหม่สำหรับ `work` ให้รัน `ccam work` หนึ่งครั้งในโฟลเดอร์นั้นแล้วยอมรับ จากนั้นจึงเริ่ม Remote Control ที่นั่น

### "You must be logged in with a Claude account that has a subscription"

Remote Control ต้องเข้าสู่ระบบด้วยบัญชี Claude ที่มีการสมัครสมาชิก ตรวจอีเมลของบัญชีใน `ccam list` และลงชื่อเข้าใช้อีกครั้งด้วย `ccam login work` หากจำเป็น

### เซสชันไม่ปรากฏบนโทรศัพท์

แอป Claude หรือ claude.ai/code ลงชื่อเข้าใช้ด้วยบัญชี Claude อื่น ให้เปรียบเทียบกับอีเมลที่ `ccam list` แสดงสำหรับ `work`

### ccam เตือนว่ามี Claude Code อีกตัวใช้ประวัติการสนทนาร่วมกันในโฟลเดอร์นี้

มีบัญชีอื่นในกลุ่มประวัติเดียวกันกำลังทำงานอยู่ในโฟลเดอร์เดียวกัน เซสชันที่ทำการสนทนาเดียวกันต่อจากทั้งสองบัญชีอาจทำให้ประวัติปนกัน ให้ตอบ `y` เฉพาะเมื่อทั้งสองจะทำงานกับการสนทนาคนละรายการ หรือเริ่ม Remote Control ในโฟลเดอร์อื่น

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp)
- [การใช้ร่วมกันระหว่างบัญชี](/features/share)
- [การจัดการบัญชี](/features/accounts)
