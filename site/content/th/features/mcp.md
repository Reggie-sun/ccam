---
title: เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน
description: "เก็บรายการเซิร์ฟเวอร์ MCP ชุดเดียวต่อกลุ่มด้วย ccam mcp add, add-json, list, remove และ import โดยส่งผ่าน --mcp-config หรือซิงก์ลงในแต่ละบัญชี"
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code เก็บเซิร์ฟเวอร์ MCP ของผู้ใช้ไว้ใน `.claude.json` ซึ่งเป็นไฟล์เดียวกับที่เก็บการเข้าสู่ระบบของบัญชี จึงใช้ไฟล์นี้ร่วมกันระหว่างบัญชีไม่ได้ ccam เก็บเซิร์ฟเวอร์ที่ใช้ร่วมกันไว้ในไฟล์ของกลุ่ม `mcp.json` และส่งให้ทุกบัญชีในกลุ่มตอนเริ่มทำงาน ใช้ `ccam mcp` เพื่อเพิ่มเซิร์ฟเวอร์ครั้งเดียวสำหรับทุกบัญชีของคุณ และใช้ `ccam mcp import` เพื่อย้ายเซิร์ฟเวอร์ที่คุณเพิ่มไว้ภายในเซสชันหนึ่งไปยังบัญชีอื่น

## การใช้งาน

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

การเปลี่ยนแปลงจะมีผลตั้งแต่การเริ่มครั้งถัดไปของแต่ละบัญชี เซสชันที่กำลังทำงานจะใช้เซิร์ฟเวอร์ชุดเดิมที่เริ่มไว้:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

การเพิ่มชื่อที่กลุ่มมีอยู่แล้วจะแทนที่นิยามเดิม และพิมพ์ `updated` แทน `added` เซิร์ฟเวอร์จะไปถึงบัญชีก็ต่อเมื่อบัญชีนั้นใช้ MCP ร่วมกับกลุ่มนั้น (`ccam share <name> --mcp on`) หากยังไม่มีบัญชีใดใช้ MCP ร่วมกัน `add` จะแจ้งให้ทราบ

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `-g`, `--group <group>` | กลุ่มเป้าหมายของ `add`, `add-json` และ `remove` | กลุ่มเดียวที่มีบัญชีใช้ MCP ร่วมกันผ่านกลุ่มนั้น หากไม่ใช่จะเป็น `default` |
| `-- <command> [args…]` (`add`) | รันคำสั่งนี้เป็นเซิร์ฟเวอร์ stdio | |
| `--url <url>` (`add`) | URL ของเซิร์ฟเวอร์ระยะไกล (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport ของเซิร์ฟเวอร์ `--url` | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header ใส่ซ้ำได้หลายครั้ง | ไม่มี |
| `-e`, `--env <KEY=VALUE>` (`add`) | ตัวแปรสภาพแวดล้อมของเซิร์ฟเวอร์ stdio ใส่ซ้ำได้หลายครั้ง | ไม่มี |
| `--force` (`import`) | แทนที่เซิร์ฟเวอร์ของกลุ่มที่ชื่อเดียวกันแต่นิยามต่างกัน โดยสำรอง `mcp.json` เดิมไว้ | ปิด |
| `--keep-local` (`import`) | เก็บรายการไว้ใน `.claude.json` ของบัญชีด้วย | ปิด |
| `-h`, `--help` | ความช่วยเหลือของแต่ละคำสั่ง | |
| `--dry-run` (global) | รองรับเฉพาะ `mcp list` ส่วน `add`, `add-json`, `remove` และ `import` จะปฏิเสธและไม่เปลี่ยนแปลงสิ่งใด | ปิด |

`add` ต้องมีคำสั่งหลัง `--` หรือ `--url` อย่างใดอย่างหนึ่ง มิฉะนั้นจะล้มเหลวด้วย `a command to run (after --) or --url is required` `--header` ใช้กับเซิร์ฟเวอร์ `--url` และ `--env` ใช้กับเซิร์ฟเวอร์แบบคำสั่ง `add-json` รับ JSON แบบเดียวกับ `claude mcp add-json` และนิยามต้องเป็นออบเจ็กต์ที่มี `command` หรือ `url` `import` ใช้กลุ่มที่บัญชีใช้ MCP ร่วมกันเสมอ และ `list` แสดงทุกกลุ่ม ดังนั้น `-g` จึงไม่มีผลกับสองคำสั่งนี้

## วิธีการทำงาน

ไฟล์ที่เกี่ยวข้อง:

- `~/.local/share/ccam/groups/<group>/mcp.json`: เซิร์ฟเวอร์ที่ใช้ร่วมกันในรูป `{"mcpServers": {…}}` ซึ่งเป็นรูปแบบที่ `--mcp-config` อ่าน ไฟล์นี้เขียนด้วย `fsutil.AtomicWrite` และโหมด `0600` เพราะรายการเซิร์ฟเวอร์มักมี API key อยู่ด้วย
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ไฟล์ของ Claude Code เองที่มีเซิร์ฟเวอร์ของบัญชี ccam อ่านเฉพาะค่า `mcpServers` ด้วยตัวสแกน `profile.TopLevelValue` และไม่พาร์สทั้งไฟล์ในเส้นทางการเริ่ม
- `~/.local/state/ccam/mcp/<name>/`: ไฟล์ `seen`, `sync-stamp` และสำเนาสำรองของ `.claude.json`

มีสองกลยุทธ์ ซึ่งตั้งด้วย `mcp_strategy` ใน `~/.config/ccam/config.toml` ค่าเริ่มต้นคือกลยุทธ์ flag:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` หากลุ่ม MCP ของบัญชี `mcp.FlagArgs` ตรวจด้วย `stat` ครั้งเดียวว่า `mcp.json` ของกลุ่มมีอยู่ โดยไม่พาร์สไฟล์
2. `Plan.InsertArgs` วาง `--mcp-config=<path>` ไว้ต่อจาก `claude` ทันที ก่อนอาร์กิวเมนต์ของคุณ ดังนั้น `--mcp-config` หรือ `--strict-mcp-config` ของคุณเองจึงยังมีผลซ้อนทับอยู่ด้านบน
3. `mcp.NewLocal` ตรวจว่า `.claude.json` ของบัญชีเปลี่ยนไปหรือไม่ตั้งแต่การเริ่มครั้งก่อน และพาร์สเซิร์ฟเวอร์เฉพาะเมื่อเปลี่ยนไปเท่านั้น
4. `syscall.Exec` แทนที่ ccam ด้วย Claude Code

พาธถูกส่งเป็นอาร์กิวเมนต์เดียว `--mcp-config=<path>` ไม่ใช่สองอาร์กิวเมนต์ `--mcp-config` ของ Claude Code รับได้หลายค่า ดังนั้น `claude --mcp-config a.json "prompt"` จะอ่านพรอมต์เป็นไฟล์การตั้งค่าไฟล์ที่สองและล้มเหลวด้วย `MCP config file not found` รูปแบบอาร์กิวเมนต์เดียวจะหยุดอยู่ตรงนั้น พรอมต์หรือคำสั่งย่อยอย่าง `mcp list` ที่ตามมาจึงทำงานได้ Claude Code โหลดเซิร์ฟเวอร์จาก `--mcp-config` เพิ่มเติมจากเซิร์ฟเวอร์ของบัญชีเอง เมื่อชื่อซ้ำกัน รายการใน `.claude.json` จะมีผลเหนือกว่า กลยุทธ์ flag เขียน `.claude.json` เฉพาะสำหรับ Remote Control เท่านั้น

`claude remote-control` (`rc`) ไม่รับ `--mcp-config` และเซสชันของมันอ่านเซิร์ฟเวอร์จาก `.claude.json` เท่านั้น สำหรับ `ccam work rc` ccam จะไม่ใส่ flag นั้น และเมื่ออยู่บนเทอร์มินัลจะถามว่าจะคัดลอกเซิร์ฟเวอร์ของกลุ่มลงใน `.claude.json` หรือไม่ โดยสำรองไฟล์ไว้ก่อนเหมือนกลยุทธ์ `sync` ccam จะถามจนกว่าจะมีสำเนา จากนั้นจะถามอีกเฉพาะหลังจาก `mcp.json` ของกลุ่มเปลี่ยน ([Remote Control](/features/launch#remote-control)) เนื่องจากสำเนาดังกล่าวมีผลเหนือกว่า `--mcp-config` เมื่อชื่อซ้ำกัน `ccam mcp add`, `add-json`, `remove` และ `import` จึงอัปเดตสำเนาที่บัญชีอื่นในกลุ่มถืออยู่:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

เซิร์ฟเวอร์ที่เพิ่มภายในเซสชันด้วย `claude mcp add` จะอยู่ใน `.claude.json` ของบัญชีนั้นเท่านั้น ccam มองไม่เห็นเซิร์ฟเวอร์นั้นในขณะนั้น เพราะ Claude Code ได้แทนที่ ccam ไปแล้ว จึงตรวจในการเริ่มครั้งถัดไป ไฟล์ `seen` เก็บเวลาแก้ไขของ `.claude.json` และชื่อเซิร์ฟเวอร์ที่รายงานไปแล้ว เมื่อเวลาต่างกัน ccam จะเปรียบเทียบชื่อกับกลุ่มและพิมพ์คำแนะนำสำหรับเซิร์ฟเวอร์ใหม่แต่ละตัวหนึ่งครั้ง โดยไม่รออินพุต:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

เมื่อใช้ `mcp_strategy = "sync"` เซิร์ฟเวอร์จะปรากฏใน `claude mcp list` และ `/mcp` เป็นเซิร์ฟเวอร์ของผู้ใช้ตามปกติ และ ccam จะเขียนเซิร์ฟเวอร์ลงใน `.claude.json` ของแต่ละบัญชีก่อนเริ่ม:

<!-- motion:mcp-sync -->

1. `mcp.Sync` เปรียบเทียบเวลาแก้ไขของ `mcp.json` ของกลุ่มกับ `sync-stamp` หากไม่เปลี่ยน จะไม่เปิด `.claude.json` ของบัญชีเลย
2. ใช้ล็อกของ Claude Code เอง ซึ่งคือโฟลเดอร์ `.claude.json.lock` ข้างไฟล์ที่สร้างด้วย `mkdir` รอได้นานสุด 2 วินาที และถือว่าล็อกที่เก่ากว่า 10 วินาทีเป็นล็อกค้าง เหมือนกับไลบรารีล็อกของ Claude Code
3. บันทึก `.claude.json` ปัจจุบันลงใน `backups/` ในโฟลเดอร์ state และเก็บไว้ 5 ชุดล่าสุด
4. ตั้งเซิร์ฟเวอร์ทุกตัวของกลุ่มใน `mcpServers` ลบเซิร์ฟเวอร์ที่เคยเขียนไว้ซึ่งกลุ่มไม่มีแล้ว คงเซิร์ฟเวอร์ของบัญชีเองไว้ และเขียนไฟล์แบบ atomic ด้วยโหมดเดิม
5. บันทึก stamp ใหม่และชื่อที่ดูแลอยู่ จากนั้น Claude Code จะเริ่มโดยไม่มี `--mcp-config`

การซิงก์จะถูกข้ามพร้อมข้อความ `ccam: account work is already running, skipping MCP sync` เมื่อมีเซสชันของบัญชีกำลังทำงาน เพราะ Claude Code เขียน `.claude.json` ใหม่ระหว่างทำงาน การซิงก์ที่ล้มเหลวจะพิมพ์ `ccam: MCP sync failed (launching anyway): …` และเริ่มต่อไปตามปกติ คำแนะนำสำหรับเซิร์ฟเวอร์ที่เพิ่มภายในเซสชันมีเฉพาะในกลยุทธ์ flag

`ccam mcp import` ย้ายเซิร์ฟเวอร์ของบัญชีเองเข้าไปในกลุ่มของบัญชี:

<!-- motion:mcp-import -->

ccam อ่านเซิร์ฟเวอร์ของบัญชีด้วย `mcp.LocalServers` และเปรียบเทียบแต่ละตัวกับเซิร์ฟเวอร์ชื่อเดียวกันของกลุ่มด้วย `mcp.SameDef` ซึ่งไม่สนใจลำดับคีย์และช่องว่าง เซิร์ฟเวอร์ที่กลุ่มยังไม่มี หรือมีอยู่ด้วยนิยามเดียวกัน จะถูกคัดลอก ส่วนที่นิยามต่างกันจะถูกข้าม เพราะนิยามของกลุ่มคือสิ่งที่บัญชีอื่นทุกบัญชีใช้อยู่:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

เมื่อใช้ `--force` เซิร์ฟเวอร์ของกลุ่มจะถูกแทนที่ และ `mcp.json` เดิมจะถูกบันทึกไว้ที่ `.ccam-backup/<time>/mcp.json` ภายในโฟลเดอร์ของกลุ่มก่อน จากนั้น `mcp.RemoveLocal` จะลบรายการที่คัดลอกแล้วออกจาก `.claude.json` ของบัญชี ภายใต้ล็อกเดียวกันและหลังจากสำรองไว้ในโฟลเดอร์ state ขั้นตอนนี้จำเป็น เพราะเมื่อชื่อซ้ำกัน รายการของบัญชีจะมีผลเหนือรายการของกลุ่ม `--keep-local` จะเก็บรายการเหล่านั้นไว้

## ความปลอดภัยและกรณีพิเศษ

- **การลงชื่อเข้าใช้** เซิร์ฟเวอร์ที่ลงชื่อเข้าใช้ด้วย OAuth เก็บโทเค็นไว้ใน `.credentials.json` ของแต่ละบัญชี (`mcpOAuth`) ไม่ได้อยู่ในรายการเซิร์ฟเวอร์ การใช้เซิร์ฟเวอร์ร่วมกันจึงไม่ได้ใช้การลงชื่อเข้าใช้ร่วมกันไปด้วย ให้ลงชื่อเข้าใช้หนึ่งครั้งต่อบัญชีด้วย `/mcp` หรือคัดลอกการลงชื่อเข้าใช้จาก `~/.claude` ด้วย [`ccam migrate`](/features/migrate) กลุ่มจะไม่ได้รับไฟล์เข้าสู่ระบบเลย
- **บัญชีที่กำลังทำงาน** `import` คัดลอกไปยังกลุ่ม แต่จะคงรายการของบัญชีที่กำลังทำงานไว้และแจ้งให้ทราบ ให้รันอีกครั้งหลังจากเซสชันจบ การซิงก์จะข้ามบัญชีที่กำลังทำงาน
- **สำเนาสำหรับ Remote Control** สำเนาของบัญชีที่กำลังทำงานจะไม่ถูกแตะ: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` เมื่อใช้ `mcp_strategy = "sync"` ไม่ต้องอัปเดตอะไร เพราะการเริ่มทุกครั้งซิงก์อยู่แล้ว
- **ไฟล์เสียหาย** หาก `.claude.json` ไม่ใช่ JSON ที่ถูกต้อง การซิงก์จะไม่แตะไฟล์และรายงานข้อผิดพลาด หากล็อกไม่ได้ภายใน 2 วินาที ccam จะรายงาน `could not lock … (Claude Code is writing the file)`
- **บัญชีที่ไม่ได้ใช้ MCP ร่วมกัน** `import` จะปฏิเสธและบอกให้คุณรัน `ccam share side --mcp on` ก่อน
- **ชื่อกลุ่ม** `-g` รับเฉพาะชื่อกลุ่ม ชื่อนั้นจะกลายเป็นโฟลเดอร์ภายใต้ `groups/` ดังนั้นชื่ออย่าง `../x` จะถูกปฏิเสธ

## หัวข้อที่เกี่ยวข้อง

- [การใช้ร่วมกันระหว่างบัญชี](/features/share)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [คัดลอกการตั้งค่าที่มีอยู่](/features/migrate)
- [คำสั่งทั้งหมด](/docs/commands)
