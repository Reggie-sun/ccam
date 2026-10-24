---
title: การคัดลอกการตั้งค่าที่มีอยู่
description: "คัดลอกการตั้งค่า, CLAUDE.md, ส่วนขยาย และเซิร์ฟเวอร์ MCP จาก ~/.claude ที่มีอยู่ไปยังกลุ่มของ ccam ด้วย ccam migrate"
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

หากคุณเคยใช้ Claude Code ก่อนใช้ ccam การตั้งค่าของคุณจะอยู่ใน `~/.claude` และ `~/.claude.json` คำสั่ง `ccam migrate` จะคัดลอกการตั้งค่าเหล่านั้นไปยังกลุ่มที่ใช้ร่วมกัน เพื่อให้ทุกบัญชีที่ใช้รายการเหล่านั้นร่วมกันผ่านกลุ่มได้เริ่มต้นด้วยการตั้งค่าที่คุณมีอยู่แล้ว โฟลเดอร์ต้นทางจะถูกอ่านอย่างเดียว และ `claude` แบบปกติยังทำงานได้เหมือนเดิม ให้รันหนึ่งครั้งหลัง `ccam add` ครั้งแรก เมื่อ `~/.claude.json` มีเซิร์ฟเวอร์ MCP อยู่ `ccam add` ครั้งแรกจะระบุชื่อเซิร์ฟเวอร์เหล่านั้นไว้ในคำแนะนำ

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
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

กลุ่มจะไปถึงเฉพาะบัญชีของกลุ่มเอง และเฉพาะรายการที่กลุ่มใช้ร่วมกัน สรุปผลจึงแสดงรายชื่อบัญชีของกลุ่ม เมื่อกลุ่มไม่ได้ใช้รายการใดที่คัดลอกมาร่วมกัน หรือบางบัญชีไม่ได้อยู่ในกลุ่ม ccam จะพิมพ์คำสั่ง `ccam group set <group> --add <items>` และ `ccam group join <group> <accounts>` สำหรับเปลี่ยนสิ่งนั้น

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

1. **หาต้นทาง** `checkSource` จะปฏิเสธพาธที่ไม่มีอยู่ ไม่ใช่โฟลเดอร์ หรือไม่มีรายการของ Claude Code ที่รู้จักเลย สำหรับเซิร์ฟเวอร์ MCP อย่างเดียว `~/.claude` ที่ว่างเปล่าหรือไม่มีอยู่ก็ใช้ได้ ตราบที่ `~/.claude.json` มีเซิร์ฟเวอร์อยู่ `migrate.ClaudeJSON` หา `.claude.json` ที่เป็นของโฟลเดอร์นั้น: อยู่ภายในโฟลเดอร์ หรือสำหรับ `~/.claude` ที่ไม่มีไฟล์นี้อยู่ภายใน คือ `~/.claude.json` ที่อยู่ข้าง ๆ ซึ่งเป็นตำแหน่งที่ Claude Code เก็บไฟล์นี้เมื่อไม่ได้ตั้ง `CLAUDE_CONFIG_DIR`
2. **จับคู่รายการกับพาธ** ทุกพาธมาจาก `share.Layout` ซึ่งเป็นรายการที่อนุญาตชุดเดียวกับที่กลุ่มใช้ เซิร์ฟเวอร์ MCP ถูกอ่านออกจาก `.claude.json` ด้วย `mcp.LocalServers` ซึ่งใช้ตัวสแกน `profile.TopLevelValue` แทนการพาร์สทั้งไฟล์
3. **คัดลอก ไม่ย้าย** ไฟล์จะไปที่ `~/.local/share/ccam/groups/<group>/` ด้วย `fsutil.AtomicWrite` หรือ `fsutil.CopyFile` ซึ่งเขียนไฟล์ชั่วคราวไว้ข้างไฟล์เป้าหมายแล้วเปลี่ยนชื่อ การคัดลอกที่หยุดกลางคันจึงไม่เหลือไฟล์ที่ขาดครึ่ง ส่วนเซิร์ฟเวอร์ MCP จะไปอยู่ใน `mcp.json` ของกลุ่มภายใต้ล็อกของไฟล์: เซิร์ฟเวอร์ที่ ccam อื่นเพิ่มในระหว่างนั้นจะยังคงอยู่ และสำเนาที่สร้างไว้สำหรับ [Remote Control](/guides/remote-control) จะถูกอัปเดตหลังจากนั้น

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

เซิร์ฟเวอร์ MCP ที่ลงชื่อเข้าใช้ด้วย OAuth เก็บโทเค็นไว้ในไฟล์เข้าสู่ระบบ `.credentials.json` ภายใต้ `mcpOAuth` ไม่ได้อยู่ในรายการเซิร์ฟเวอร์ ccam ไม่คัดลอกโทเค็นเหล่านั้น การลงชื่อเข้าใช้ที่ถูกคัดลอกจะถือ refresh token เดียวกับต้นฉบับ และเซิร์ฟเวอร์จะหมุนเวียนโทเค็นนั้นเมื่อต่ออายุครั้งแรก: ผู้ถือรายอื่นทั้งหมด รวมถึง `~/.claude` ต้นทาง จะหลุดจากเซิร์ฟเวอร์นั้น ดังนั้น สำหรับแต่ละเซิร์ฟเวอร์ที่คัดลอกซึ่งต้นทางลงชื่อเข้าใช้อยู่ บรรทัด `sign-ins` ของแผนจะระบุชื่อบัญชีในกลุ่มที่ต้องลงชื่อเข้าใช้เองด้วย `/mcp` และสรุปผลจะแจ้งซ้ำอีกครั้ง รายการที่ไม่มีโทเค็นไม่ถือเป็นการลงชื่อเข้าใช้และจะไม่แสดง เช่นเดียวกับบัญชีที่มีการลงชื่อเข้าใช้ของตัวเองอยู่แล้ว หรือเซิร์ฟเวอร์ที่กลุ่มเก็บนิยามของตัวเองไว้ ไม่มีไฟล์เข้าสู่ระบบใดถูกเขียน

## ความปลอดภัยและกรณีพิเศษ

- ต้นทางจะถูกอ่านอย่างเดียว ไม่มีไฟล์ใดในต้นทางถูกย้าย ลิงก์ หรือเปลี่ยนแปลง และการเข้าสู่ระบบของต้นทางจะไม่ถูกเขียนเลย รวมถึงบน macOS ที่อาจอ่านจากพวงกุญแจ (Keychain)
- ก่อนที่ไฟล์ของกลุ่มจะถูกแทนที่ ไม่ว่าจะด้วย `use mine` การรวม หรือเซิร์ฟเวอร์ MCP ใหม่ใน `mcp.json` ที่มีอยู่ สำเนาของกลุ่มจะถูกบันทึกไว้ใน `<group>/.ccam-backup/<time>/` ชื่อโฟลเดอร์ไม่ซ้ำกันในแต่ละครั้งที่รัน การรันสองครั้งในวินาทีเดียวกันจึงไม่เขียนทับสำรองของกันและกัน สรุปผลจะพิมพ์พาธนั้น
- เมื่อใช้ `--yes` ทุกความขัดแย้งจะเก็บสำเนาของกลุ่มไว้ และ ccam จะบอกวิธีเลือกเอง: รันคำสั่งเดิมโดยไม่ใส่ `--yes`
- หากไม่มีสิ่งใดเปลี่ยน ccam จะพิมพ์ `Nothing to copy.` และหากทุกความขัดแย้งเก็บสำเนาของกลุ่มไว้ จะพิมพ์ `Nothing to copy: the group keeps its own copies.`
- หากการคัดลอกล้มเหลวกลางคัน ccam จะหยุดด้วย `The source folder was not changed. Group files replaced so far are saved in` ตามด้วยโฟลเดอร์สำรอง การรันอีกครั้งจะคัดลอกไฟล์ที่มาถึงไม่ครบ
- เมื่อหลายกลุ่มใช้ MCP ร่วมกัน คำสั่งที่ `ccam doctor`, `ccam ext` และ `ccam mcp list` แนะนำจะเป็น `ccam migrate --mcp --group <a|b>`
- ต้นทางที่ไม่มีอยู่จะจบด้วย `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`

## หัวข้อที่เกี่ยวข้อง

- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [การตรวจสอบและซ่อมแซม](/features/doctor)
- [Claude Code หลายบัญชี](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
