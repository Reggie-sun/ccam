---
title: การคัดลอกการตั้งค่าที่มีอยู่
description: "คัดลอกการตั้งค่า, CLAUDE.md, ส่วนขยาย, เซิร์ฟเวอร์ MCP และการลงชื่อเข้าใช้ MCP จาก ~/.claude ที่มีอยู่ไปยังกลุ่มที่ใช้ร่วมกันของ ccam ด้วย ccam migrate"
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

หากคุณเคยใช้ Claude Code ก่อนใช้ ccam การตั้งค่าของคุณจะอยู่ใน `~/.claude` และ `~/.claude.json` คำสั่ง `ccam migrate` จะคัดลอกการตั้งค่าเหล่านั้นไปยังกลุ่มที่ใช้ร่วมกัน เพื่อให้ทุกบัญชีที่ใช้รายการเหล่านั้นร่วมกันผ่านกลุ่มได้เริ่มต้นด้วยการตั้งค่าที่คุณมีอยู่แล้ว โฟลเดอร์ต้นทางจะถูกอ่านอย่างเดียว และ `claude` แบบปกติยังทำงานได้เหมือนเดิม ให้รันหนึ่งครั้งหลัง `ccam add` ครั้งแรก และรันอีกครั้งทุกครั้งที่มีบัญชีเริ่มใช้เซิร์ฟเวอร์ MCP ร่วมกัน

## การใช้งาน

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

บนเทอร์มินัล หากไม่ระบุแฟล็กของรายการและไม่ใส่ `--yes` ccam จะแสดงรายการตรวจสอบ (checklist) ก่อน:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (หรือ `x`) สลับรายการที่เคอร์เซอร์อยู่ ปุ่มลูกศรหรือ `j`/`k` ใช้เลื่อน Enter แสดงแผน ส่วน Esc, `q` หรือ Ctrl+C จะพิมพ์ `Cancelled.` รายการที่ต้นทางไม่มีข้อมูลจะมีเครื่องหมาย `not found` หากคุณยกเลิกการเลือกทั้งหมด ccam จะพิมพ์ `Nothing selected.` และหยุด

จากนั้น ccam จะพิมพ์แผนและถามก่อนคัดลอก:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

กลุ่มจะไปถึงเฉพาะบัญชีที่ใช้รายการร่วมกันผ่านกลุ่มนั้น สรุปผลจึงแสดงรายชื่อบัญชีแยกตามรายการ สำหรับบัญชีที่ยังไม่ได้ใช้ร่วมกัน ccam จะพิมพ์คำสั่ง `ccam share <name> --<item> on` สำหรับเปิดการใช้ร่วมกัน

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `[folder]` | โฟลเดอร์ Claude Code ที่จะคัดลอกจาก | `~/.claude` |
| `--settings` | คัดลอก `settings.json` | ปิด |
| `--memory` | คัดลอก `CLAUDE.md` ระดับ global | ปิด |
| `--extensions` | คัดลอก `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | ปิด |
| `--mcp` | คัดลอกเซิร์ฟเวอร์ MCP ของผู้ใช้ (`mcpServers` ใน `.claude.json`) | ปิด |
| `--history` | คัดลอกประวัติการสนทนา (`projects/`, `history.jsonl`) | ปิด |
| `-g`, `--group string` | กลุ่มที่จะคัดลอกไป | `default` |
| `-y`, `--yes` | คัดลอกโดยไม่ถาม เมื่อขัดแย้งกันจะเก็บสำเนาของกลุ่มไว้ | ปิด |
| `--dry-run` (global) | พิมพ์แผนแล้วหยุดก่อนถามหรือคัดลอกใด ๆ | ปิด |

หากไม่ระบุแฟล็กของรายการ บนเทอร์มินัลรายการตรวจสอบจะเป็นตัวกำหนด นอกเทอร์มินัล หรือเมื่อใส่ `--yes` ccam จะคัดลอกทุกอย่างยกเว้นประวัติการสนทนา ประวัติไม่ถูกรวมไว้โดยค่าเริ่มต้นเพราะ `projects/` อาจมีขนาดใหญ่

`--group on` หมายถึง `default` ส่วน `--group off` และ `--group local` จะถูกปฏิเสธด้วย `--group takes a group name, not "off"` `ccam migrate` รองรับ `--dry-run` และรายการตรวจสอบยังคงแสดงในการทดลองรัน คุณจึงดูตัวอย่างสิ่งที่เลือกได้

## วิธีการทำงาน

`internal/cli/migrate.go` เลือกรายการ จากนั้น `migrate.Build` อ่านต้นทางและกลุ่มแล้วคืนแผนกลับมา ไม่มีสิ่งใดถูกเขียนจนกว่าจะถึง `Plan.Apply` หลังจากที่คุณยืนยัน

1. **หาต้นทาง** `checkSource` จะปฏิเสธพาธที่ไม่มีอยู่ ไม่ใช่โฟลเดอร์ หรือไม่มีรายการของ Claude Code ที่รู้จักเลย `migrate.ClaudeJSON` หา `.claude.json` ที่เป็นของโฟลเดอร์นั้น: อยู่ภายในโฟลเดอร์ หรือสำหรับ `~/.claude` ที่ไม่มีไฟล์นี้อยู่ภายใน คือ `~/.claude.json` ที่อยู่ข้าง ๆ ซึ่งเป็นตำแหน่งที่ Claude Code เก็บไฟล์นี้เมื่อไม่ได้ตั้ง `CLAUDE_CONFIG_DIR`
2. **จับคู่รายการกับพาธ** ทุกพาธมาจาก `share.Layout` ซึ่งเป็นรายการที่อนุญาตชุดเดียวกับที่ `ccam share` ใช้ เซิร์ฟเวอร์ MCP ถูกอ่านออกจาก `.claude.json` ด้วย `mcp.LocalServers` ซึ่งใช้ตัวสแกน `profile.TopLevelValue` แทนการพาร์สทั้งไฟล์
3. **คัดลอก ไม่ย้าย** ไฟล์จะไปที่ `~/.local/share/ccam/groups/<group>/` ด้วย `fsutil.AtomicWrite` หรือ `fsutil.CopyFile` ส่วนเซิร์ฟเวอร์ MCP จะไปอยู่ใน `mcp.json` ของกลุ่ม

แต่ละรายการจับคู่พาธต้นทางหนึ่งพาธกับพาธของกลุ่มหนึ่งพาธ:

<!-- motion:migrate-items -->

สองแถวของประวัติเป็นเส้นประ เพราะจะถูกคัดลอกเฉพาะเมื่อใส่ `--history` `.credentials.json` และ `.claude.json` อยู่ใน `share.Never` และไม่มีวันไปถึงกลุ่ม: `Build` จะคืนข้อผิดพลาดภายในหากไฟล์ใดไฟล์หนึ่งปรากฏในโครงสร้าง และการไล่ดูโฟลเดอร์จะข้ามไฟล์เหล่านี้ตามชื่อ

`buildEntry` กำหนดสถานะให้ทุกพาธ ซึ่งคือสิ่งที่แต่ละบรรทัดของแผนอธิบาย:

<!-- motion:migrate-states -->

- **Missing**: ไม่มีในต้นทาง หรือเป็นไฟล์ในตำแหน่งที่ควรเป็นโฟลเดอร์ แสดงหนึ่งครั้งต่อรายการเป็น `not found, skipped`
- **Linked**: พาธต้นทางชี้ไปยังสำเนาของกลุ่มอยู่แล้ว เช่น เมื่อต้นทางเป็นโฟลเดอร์โปรไฟล์ของ ccam
- **Same**: กลุ่มมีเนื้อหาที่เหมือนกันอยู่แล้ว
- **New**: กลุ่มยังไม่มี จึงคัดลอกทั้งหมด
- **Merge**: เป็นโฟลเดอร์ทั้งสองฝั่ง `copyMissing` เพิ่มเฉพาะไฟล์ที่กลุ่มยังไม่มี ไฟล์ที่ต่างกันจะคงไว้ตามที่กลุ่มมี และแผนจะแสดง `group copies kept` สำหรับ `plugins/` รีจิสทรี `installed_plugins.json` และ `known_marketplaces.json` จะถูกรวมเข้าด้วยกัน และพาธปลั๊กอินแบบสัมบูรณ์จะถูกเขียนใหม่ให้ชี้ไปที่โฟลเดอร์ของกลุ่ม
- **Conflict**: เป็นไฟล์ทั้งสองฝั่งแต่เนื้อหาต่างกัน คุณเลือกทีละไฟล์: `keep group`, `use mine` รวมถึง `merge JSON` สำหรับไฟล์ `.json` (การรวมแบบลึกที่ค่าของกลุ่มมีผลเหนือกว่าเมื่อค่าขัดกัน) หรือ `merge by line` สำหรับ `history.jsonl` เซิร์ฟเวอร์ MCP ที่นิยามต่างกันก็ถือเป็นความขัดแย้งเช่นกัน โดยเลือกได้ `keep group` หรือ `use mine` ตัวเลือกแรกในรายการคือค่าเริ่มต้นเมื่อคุณกด Enter ส่วน `--yes` จะเก็บสำเนาของกลุ่มไว้เสมอ

รายการตรวจสอบเองคือ `tui.MigratePicker` ซึ่งสร้างบนคอมโพเนนต์ `checklist` ที่ `ccam add` ใช้ `migrate.Found` ทำเครื่องหมายรายการที่ไม่มีด้วยการเรียก lstat และรายการเซิร์ฟเวอร์ MCP

### การลงชื่อเข้าใช้ MCP

เซิร์ฟเวอร์ MCP ที่ลงชื่อเข้าใช้ด้วย OAuth เก็บโทเค็นไว้ในไฟล์เข้าสู่ระบบ `.credentials.json` ภายใต้ `mcpOAuth` ไม่ได้อยู่ในรายการเซิร์ฟเวอร์ Claude Code ใช้คีย์ของแต่ละรายการเป็น `<name>|` ตามด้วย 16 อักขระแรกของแฮช `sha256` ของ `type`, `url` และ `headers` ของเซิร์ฟเวอร์ คีย์นี้ไม่ขึ้นกับโฟลเดอร์ รายการที่คัดลอกจึงใช้ได้ในบัญชีอื่นตราบใดที่นิยามเซิร์ฟเวอร์เหมือนกันทุกประการ ดังนั้น `buildLogins` ใน `internal/migrate/logins.go` จึงคัดลอกเฉพาะรายการของเซิร์ฟเวอร์ที่กลุ่มจะได้รับตามนิยามของต้นทาง ได้แก่ เซิร์ฟเวอร์ใหม่ เซิร์ฟเวอร์ที่เหมือนกันอยู่แล้ว หรือความขัดแย้งที่คุณเลือก `use mine` การลงชื่อเข้าใช้ของเซิร์ฟเวอร์ปลั๊กอิน (คีย์ที่ขึ้นต้นด้วย `plugin:`) จะถูกคัดลอกไปด้วยเมื่อคุณคัดลอกส่วนขยายและต้นทางมี `plugins/`

<!-- motion:migrate-signins -->

รายการเหล่านี้จะไปอยู่ใน `.credentials.json` ของแต่ละบัญชีที่ได้รับ ซึ่งเขียนด้วย `fsutil.AtomicWrite` ในโหมด `0600` และไม่มีวันไปอยู่ในกลุ่ม บัญชีจะได้รับก็ต่อเมื่อใช้ `mcp` ร่วมกัน (หรือ `extensions` สำหรับเซิร์ฟเวอร์ปลั๊กอิน) ผ่านกลุ่มเป้าหมาย และจะเพิ่มเฉพาะคีย์ที่บัญชียังไม่มี: รายการที่บัญชีมีอยู่แล้ว ซึ่งอาจใหม่กว่า จะถูกเก็บไว้ บรรทัด `sign-ins` ในแผนแสดงเซิร์ฟเวอร์ บัญชีที่ได้รับ และบัญชีที่ถูกข้าม ccam จะข้ามบัญชีและระบุชื่อพร้อมเหตุผล เมื่อบัญชีนั้น:

- ลงทะเบียนไว้ในตำแหน่งเดิมด้วย `ccam import` เพราะ `claude` แบบปกติอาจกำลังใช้ไฟล์เข้าสู่ระบบของบัญชีนั้นอยู่
- กำลังทำงาน เพราะ Claude Code เขียนไฟล์เข้าสู่ระบบใหม่ (ตรวจซ้ำอีกครั้งก่อนเขียน และรายงานเป็น `started meanwhile`)
- ยังไม่ได้เข้าสู่ระบบ
- เก็บการเข้าสู่ระบบไว้ในพวงกุญแจ (Keychain) ของ macOS ซึ่ง ccam ไม่เขียน

ให้ลงชื่อเข้าใช้เซิร์ฟเวอร์ในบัญชีที่ถูกข้ามด้วย `/mcp` บัญชีที่เริ่มใช้เซิร์ฟเวอร์ MCP ร่วมกันภายหลังจะได้รับการลงชื่อเข้าใช้เมื่อคุณรัน `ccam migrate` อีกครั้ง โดยสรุปผลจะแจ้งด้วย ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## ความปลอดภัยและกรณีพิเศษ

- ต้นทางจะถูกอ่านอย่างเดียว ไม่มีไฟล์ใดในต้นทางถูกย้าย ลิงก์ หรือเปลี่ยนแปลง และการเข้าสู่ระบบของต้นทางจะไม่ถูกเขียนเลย รวมถึงบน macOS ที่อาจอ่านจากพวงกุญแจ (Keychain)
- ก่อนที่ไฟล์ของกลุ่มจะถูกแทนที่ ไม่ว่าจะด้วย `use mine` การรวม หรือเซิร์ฟเวอร์ MCP ใหม่ใน `mcp.json` ที่มีอยู่ สำเนาของกลุ่มจะถูกบันทึกไว้ใน `<group>/.ccam-backup/<time>/` ชื่อโฟลเดอร์ไม่ซ้ำกันในแต่ละครั้งที่รัน การรันสองครั้งในวินาทีเดียวกันจึงไม่เขียนทับสำรองของกันและกัน สรุปผลจะพิมพ์พาธนั้น
- เมื่อใช้ `--yes` ทุกความขัดแย้งจะเก็บสำเนาของกลุ่มไว้ และ ccam จะบอกวิธีเลือกเอง: รันคำสั่งเดิมโดยไม่ใส่ `--yes`
- หากไม่มีสิ่งใดเปลี่ยน ccam จะพิมพ์ `Nothing to copy.` และหากทุกความขัดแย้งเก็บสำเนาของกลุ่มไว้ จะพิมพ์ `Nothing to copy: the group keeps its own copies.`
- หากการคัดลอกล้มเหลวกลางคัน ccam จะหยุดด้วย `The source folder was not changed. Group files replaced so far are saved in` ตามด้วยโฟลเดอร์สำรอง
- ต้นทางที่ไม่มีอยู่จะจบด้วย `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`

## หัวข้อที่เกี่ยวข้อง

- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [การตรวจสอบและซ่อมแซม](/features/doctor)
- [Claude Code หลายบัญชี](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
