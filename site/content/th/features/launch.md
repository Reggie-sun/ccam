---
title: การเริ่ม Claude Code
description: "เริ่ม Claude Code ในชื่อบัญชีใดก็ได้ด้วย ccam <name|number>, ccam run หรือตัวเลือกบัญชีของ ccam อาร์กิวเมนต์จะส่งต่อให้ claude โดยไม่เปลี่ยนแปลง และ --dry-run จะแสดงคำสั่ง"
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

`ccam work` เริ่ม Claude Code ในชื่อบัญชี `work` และทุกอย่างหลังชื่อบัญชีจะส่งต่อให้ `claude` ตามเดิม ใช้ได้ทุกที่ที่คุณจะพิมพ์ `claude`: ในเทอร์มินัล ในสคริปต์ หรือเมื่อส่งอินพุตเข้ามาทางไปป์

## การใช้งาน

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

ระบุชื่อบัญชีหรือหมายเลขบัญชีจาก `ccam list` ccam ไม่อ่านอาร์กิวเมนต์ที่อยู่หลังชื่อบัญชี ดังนั้น `ccam work --help` จะแสดงความช่วยเหลือของ Claude Code จะใส่ `--` หลังชื่อบัญชีหรือไม่ก็ได้ อินพุตมาตรฐาน รหัสออก (exit code) และสัญญาณ (signal) ทำงานเหมือนกับ `claude`

เมื่ออาร์กิวเมนต์แรกขึ้นต้นด้วย `-` และไม่ใช่แฟล็กของ ccam เอง ccam จะเริ่มบัญชีเริ่มต้นพร้อมอาร์กิวเมนต์ทั้งหมด บัญชีเริ่มต้นคือบัญชีที่ตั้งไว้ด้วย `ccam use` หรือบัญชีที่มีหมายเลขต่ำสุดหากคุณยังไม่ได้ตั้ง

### ตัวเลือกบัญชี

`ccam` ที่ไม่มีอาร์กิวเมนต์จะเปิดตัวเลือกบัญชี (picker) เมื่ออินพุตและเอาต์พุตมาตรฐานเป็นเทอร์มินัล และจะพิมพ์ `ccam list` ในกรณีอื่น ตัวเลือกบัญชีแสดงหมายเลข อีเมล การใช้งาน 5 ชั่วโมง เวลาที่ใช้ล่าสุด และรายการที่ใช้ร่วมกันของแต่ละบัญชี ตัวเลขการใช้งานจะปรากฏขึ้นเมื่อแต่ละบัญชีได้รับคำตอบ

| ปุ่ม | สิ่งที่ทำ |
| --- | --- |
| `Enter` | เริ่ม Claude Code ในชื่อบัญชีที่เลือก |
| `1` ถึง `9` | เริ่มบัญชีที่มีหมายเลขนั้น |
| `Up`/`Down`, `k`/`j` | เลื่อนตัวเลือก |
| `a` | ปิดตัวเลือกบัญชีแล้วเปิดตัวช่วยตั้งค่าของ `ccam add` |
| `s` | แก้ไขรายการที่ใช้ร่วมกันของบัญชีที่เลือก (`Space` หรือ `x` สลับ, `Enter` นำไปใช้, `Esc` ย้อนกลับ) |
| `d` | ตั้งบัญชีที่เลือกเป็นบัญชีเริ่มต้น |
| `r` | ดึงการใช้งานอีกครั้งโดยไม่ใช้แคชระยะสั้น |
| `q`, `Esc`, `Ctrl+C` | ออก |

ตัวเลือกบัญชีจะปิดตัวเองก่อนทำงาน เทอร์มินัลจึงกลับสู่สภาพเดิมก่อนที่ Claude Code จะเริ่ม หลังจากกด `d` หรือเปลี่ยนการใช้ร่วมกัน ตัวเลือกบัญชีจะเปิดขึ้นอีกครั้ง การเปลี่ยนการใช้ร่วมกันจะแสดงแผนและถามก่อน เหมือนกับที่ `ccam share` ทำ

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--dry-run` | พิมพ์ตัวแปรสภาพแวดล้อมและคำสั่งที่ ccam จะรัน แล้วออก | ปิด |
| `-v`, `--version` | พิมพ์ `ccam <version>` (เฉพาะเมื่อเป็นอาร์กิวเมนต์แรก) | |
| `-h`, `--help` | ความช่วยเหลือของ ccam ในรูป `ccam --help` หรือ `ccam run --help` | |

กฎของอาร์กิวเมนต์:

- ccam อ่านแฟล็กของตัวเองเฉพาะที่อยู่ก่อนชื่อบัญชี หลังชื่อบัญชี ทุกอาร์กิวเมนต์เป็นของ Claude Code รวมถึง `--dry-run` และ `--help`
- `--dry-run` ใช้ได้ในรูป `ccam --dry-run work`, `ccam run --dry-run work` และ `ccam --dry-run` (บัญชีเริ่มต้น) ตัวเลือกนี้ได้รับการรองรับ ไม่ถูกปฏิเสธ
- `ccam run <name|number> [claude args]` คือรูปแบบเต็มของ `ccam <name>` ส่วน `ccam run` ที่ไม่มีชื่อบัญชีจะใช้บัญชีเริ่มต้น
- ชื่อบัญชีต้องไม่ซ้ำกับชื่อคำสั่งของ ccam เช่น `add`, `list`, `run` หรือ `usage` ดังนั้น `ccam <word>` จึงไม่กำกวม

การทดลองรัน (dry run) จะพิมพ์หนึ่งบรรทัดที่คุณนำไปวางในเชลล์ได้:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## วิธีการทำงาน

การเริ่มทุกครั้งผ่านขั้นตอนเดียวกัน ตั้งแต่บรรทัดคำสั่งไปจนถึง `exec`:

<!-- motion:launch -->

แต่ละกล่องจะสว่างขึ้นเมื่อขั้นตอนนั้นทำงาน และบรรทัดด้านล่างคือตัวแปรสภาพแวดล้อมและอาร์กิวเมนต์ที่ Claude Code ใช้เริ่มทำงาน

1. `cmd/ccam/main.go` เรียก `app.Main` หาก ccam ถูกเริ่มด้วยชื่ออื่น แสดงว่าเป็น[คำสั่งลัด](/features/shortcuts) และจะไปที่ตัวเริ่มโปรแกรมโดยตรง
2. `app.parse` ตรวจ `os.Args` ก่อนที่ cobra จะเข้ามาเกี่ยวข้อง มีเพียง `ccam` เดี่ยว ๆ (ตัวเลือกบัญชี) ชื่อคำสั่งของ ccam (`config.Reserved`) `-h`/`--help` และ `-v`/`--version` ที่ไปยัง `cli.Execute` นอกนั้นถือเป็นการเริ่ม Claude Code
3. `launcher.Launch` อ่าน `~/.config/ccam/config.toml` ซึ่งเป็นไฟล์เดียวที่พาร์ส
4. `launcher.Prepare` หาบัญชีจากชื่อหรือหมายเลข ตัด `--` ตัวแรกที่นำหน้าออกหนึ่งตัว และหาไฟล์ปฏิบัติการ `claude` ด้วย `ResolveClaude` ผลการค้นหาใน PATH จะถูกแคชไว้ใน `~/.local/state/ccam/claude-path` พร้อมเวลาแก้ไขของไฟล์ ccam จึงไม่ต้องค้นหาใน `PATH` ทุกครั้งที่เริ่ม หากไฟล์เปลี่ยนหรือหายไปจะค้นหาใหม่ พาธ `claude_bin` ใน `config.toml` จะถูกใช้ตามที่ระบุ ccam จะปฏิเสธ `claude` ที่เป็นตัว ccam เอง (เปรียบเทียบด้วย inode) เพราะจะทำให้วนซ้ำไม่สิ้นสุด
5. `checkLinks` เปรียบเทียบโฟลเดอร์บัญชี `~/.local/share/ccam/profiles/<name>/` กับการตั้งค่าการใช้ร่วมกันด้วย `share.Check` ซึ่งเรียกเพียง `lstat` และ `readlink` หนึ่งหรือสอง system call ต่อรายการ ลิงก์ที่หายไปจะถูกสร้างทันที สิ่งใดที่อาจทำให้ข้อมูลสูญหายจะหยุดการเริ่ม
6. `applyMCP` เพิ่มเซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน (ดูด้านล่าง)
7. `session.Register` บันทึกเซสชัน ccam เขียนไฟล์เครื่องหมายที่ตั้งชื่อตาม PID ซึ่งเก็บเวลาเริ่มของโปรเซส ไว้ที่ `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (หรือ `~/.local/state/ccam/locks/` หากไม่มี `XDG_RUNTIME_DIR`) สำหรับบัญชีที่ใช้ประวัติการสนทนาร่วมกัน `historyCheck` จะค้นหาเซสชันที่กำลังทำงานใน `history-<group>-<hash>` ก่อน โดย hash คำนวณจากโฟลเดอร์ปัจจุบัน แล้วลงทะเบียนไว้ที่นั่นด้วย `exec` คง PID ไว้ เครื่องหมายจึงระบุถึงโปรเซสของ Claude Code เอง ที่ใช้เครื่องหมายแทนการถือล็อก เพราะโปรเซสที่ Claude Code เริ่มจะสืบทอดล็อกไปและถือไว้ต่อหลังจากเซสชันจบแล้ว
8. `Plan.Exec` ตั้งชื่อหน้าต่างเทอร์มินัลเป็น `claude · <name>` (เมื่อ standard error เป็นเทอร์มินัล) เขียนเวลาที่เริ่มลงใน `~/.local/state/ccam/last-used/<name>` โดยไม่เรียก `fsync` แล้วเรียก `syscall.Exec`

Claude Code ทำงานโดยตั้ง `CLAUDE_CONFIG_DIR` เป็นโฟลเดอร์ของบัญชี และ `CCAM_PROFILE` เป็นชื่อบัญชี ccam จะลบค่าของตัวแปรเหล่านี้ที่สืบทอดมา สำหรับโฟลเดอร์ `~/.claude` ที่ลงทะเบียนด้วย `ccam import` บน macOS ccam จะตั้ง `CLAUDE_SECURESTORAGE_CONFIG_DIR` เป็นค่าว่างด้วย เพื่อให้บัญชีนั้นและ `claude` แบบปกติใช้การเข้าสู่ระบบเดียวกันในพวงกุญแจ (Keychain)

### ทำไมจึงใช้ exec

`syscall.Exec` แทนที่โปรเซสของ ccam ด้วย `claude` ccam จึงไม่ค้างอยู่ในหน่วยความจำ การกดแป้น สัญญาณ และรหัสออกเป็นของ Claude Code เอง และอาร์กิวเมนต์ยังคงเป็นอาร์เรย์ ไม่เคยเป็นสตริงของเชลล์ งานหลังจบเซสชัน เช่น การตรวจพบเซิร์ฟเวอร์ที่เพิ่มด้วย `claude mcp add` จะเกิดขึ้นในการเริ่มครั้งถัดไป: ccam เปรียบเทียบเวลาแก้ไขของ `.claude.json` กับ `~/.local/state/ccam/mcp/<name>/seen` และพาร์สไฟล์เฉพาะเมื่อไฟล์เปลี่ยนไป

### บรรทัดคำสั่งถูกส่งต่ออย่างไร

สามบรรทัดคำสั่ง สามเส้นทางผ่าน `internal/app`:

<!-- motion:launch-route -->

คำสั่งลัดจับคู่จาก `argv[0]` ส่วน `ccam <name>` จับคู่จากอาร์กิวเมนต์แรก และทั้งสองจบที่ `syscall.Exec` มีเพียงคำสั่งของ ccam เช่น `ccam usage` เท่านั้นที่ไปถึง cobra

### เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน

ด้วยกลยุทธ์ `flag` ซึ่งเป็นค่าเริ่มต้น ccam ส่งเซิร์ฟเวอร์ MCP ของกลุ่มให้ Claude Code เป็นอาร์กิวเมนต์ โดยวางไว้ก่อนอาร์กิวเมนต์ของคุณ:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` คืนอาร์กิวเมนต์เฉพาะเมื่อ `mcp.json` ของกลุ่มมีอยู่ และ `Plan.InsertArgs` วางอาร์กิวเมนต์ไว้ต่อจาก `claude` ทันที `--mcp-config` หรือ `--strict-mcp-config` ของคุณเองจึงยังมีผลซ้อนทับอยู่ด้านบน อาร์กิวเมนต์นี้เป็นอาร์กิวเมนต์เดียว `--mcp-config=<path>` เพราะ `--mcp-config` ของ Claude Code รับได้หลายค่า หากเขียนเป็นสองอาร์กิวเมนต์ ก็จะอ่านพรอมต์ที่ตามมาเป็นไฟล์การตั้งค่าไปด้วย

เมื่อตั้ง `mcp_strategy = "sync"` ใน `config.toml` ccam จะเขียนเซิร์ฟเวอร์ลงใน `.claude.json` ของบัญชีแทน และจะข้ามขั้นตอนนี้หากบัญชีกำลังทำงานอยู่แล้ว

### Remote Control

`ccam work rc` (หรือ `ccam work remote-control`) เริ่ม Remote Control ของ Claude Code ในฐานะ `work` ส่วน `claude remote-control` ปฏิเสธ `--mcp-config` ที่อยู่ก่อนคำกริยา (verb) ของมัน และไม่รับเลยหลังคำกริยา อีกทั้งเซสชันของมันอ่านเซิร์ฟเวอร์ MCP จาก `.claude.json` ของบัญชีเท่านั้น ดังนั้นกลยุทธ์ flag จึงทำงานต่างออกไปในกรณีนี้:

<!-- motion:launch-rc -->

1. `remoteControl` ถือว่าอาร์กิวเมนต์แรกหลังชื่อบัญชีที่ไม่ขึ้นต้นด้วย `-` คือคำกริยาของ Claude Code เมื่อเป็น `remote-control` หรือ `rc` แล้ว `mcp.FlagArgs` จะไม่เพิ่มอะไรเลย
2. `mcp.Synced` เปรียบเทียบเวลาแก้ไขของ `mcp.json` ของกลุ่มกับ `sync-stamp` ของบัญชี หาก `.claude.json` มีเซิร์ฟเวอร์ปัจจุบันของกลุ่มอยู่แล้ว หรือกลุ่มไม่มีเซิร์ฟเวอร์เลย Remote Control จะเริ่มโดยไม่ถาม
3. มิฉะนั้น ccam จะระบุชื่อเซิร์ฟเวอร์ที่เซสชันจะไม่ได้รับ แล้วถามว่า:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. เมื่อตอบ `y` แล้ว `mcp.Sync` จะคัดลอกเซิร์ฟเวอร์เหล่านั้นแบบเดียวกับกลยุทธ์ `sync`: ภายใต้ล็อกของ Claude Code หลังจากสำรองไฟล์ และคงเซิร์ฟเวอร์ของบัญชีเองไว้ เมื่อตอบ `N` หรือเมื่อไม่มีเทอร์มินัล Remote Control จะเริ่มโดยไม่มีเซิร์ฟเวอร์เหล่านั้น และ ccam จะแนะนำ `mcp_strategy = "sync"` ซึ่งคัดลอกให้ก่อนการเริ่มทุกครั้ง

ccam จะถามทุกครั้งที่เริ่ม Remote Control จนกว่าจะคัดลอกเซิร์ฟเวอร์แล้ว และหลังจากคัดลอกแล้วจะถามอีกเฉพาะเมื่อ `mcp.json` ของกลุ่มเปลี่ยนอีกครั้ง สำเนาใน `.claude.json` มีผลเหนือกว่า `--mcp-config` ในการเริ่มตามปกติครั้งต่อ ๆ ไป ดังนั้น `ccam mcp add`, `remove` และ `import` จึงอัปเดตสำเนาที่บัญชีอื่นถืออยู่ด้วย ([เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp))

### งบเวลา 10 ms

ccam ใช้เวลาเพิ่มไม่เกิน 10 ms (p95) ก่อนที่ Claude Code จะเริ่ม เส้นทางการเริ่มไม่เคยสร้างโครงสร้างคำสั่งของ cobra ไม่เริ่ม Bubble Tea และไม่เปิดการเชื่อมต่อเครือข่าย และ `TestPackageInitIsCheap` จะล้มเหลวหากฟังก์ชัน `init()` ของแพ็กเกจใช้เวลารวมเกิน 5 ms `scripts/perf-check.sh` วัด `ccam --dry-run work`, `ccam work` และคำสั่งลัดด้วย `hyperfine` และล้มเหลวเมื่อเกิน 10 ms ส่วน `BenchmarkLaunchPath` ใน `internal/launcher` รันขั้นตอนต่าง ๆ ซ้ำในลูป หากต้องการดูเวลาของแต่ละขั้นตอน ให้ตั้ง `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## ความปลอดภัยและกรณีพิเศษ

- ชื่อหรือหมายเลขที่ไม่รู้จักจะหยุดด้วย `no account "<name>"` และแนะนำ `ccam list` หากยังไม่มีบัญชีเลย ccam จะขอให้คุณเพิ่มบัญชีด้วย `ccam add <name>`
- หากพาธที่ใช้ร่วมกันกลายเป็นไฟล์ธรรมดาหรือชี้ไปผิดที่ ccam จะหยุดด้วย `stopped because of share link problems:` แสดงรายการปัญหา และขอให้คุณรัน `ccam doctor --fix` ccam ไม่ย้ายข้อมูลในเส้นทางการเริ่ม
- หากมี Claude Code อีกตัวที่ใช้ประวัติการสนทนาร่วมกับบัญชีนี้ (กลุ่มเดียวกัน) ทำงานอยู่ในโฟลเดอร์เดียวกัน ccam จะเตือนว่าประวัติอาจปนกันและถาม `Run anyway? [y/N]` หากไม่มีเทอร์มินัลจะแสดงคำเตือนแล้วทำงานต่อ
- `--dry-run` ข้ามการซิงก์ MCP การแจ้งเซิร์ฟเวอร์ใหม่ การตรวจประวัติ เครื่องหมายเซสชัน ชื่อหน้าต่าง และเวลาที่ใช้ล่าสุด และไม่เรียก exec แต่ยังคงตรวจลิงก์ที่ใช้ร่วมกัน และสร้างลิงก์ที่หายไปเหมือนการเริ่มปกติ
- `ccam work rc` ขณะที่ `work` กำลังทำงานอยู่แล้ว จะไม่แตะ `.claude.json` แจ้งให้ทราบ และเริ่ม Remote Control โดยไม่มีเซิร์ฟเวอร์ของกลุ่ม ให้รันอีกครั้งหลังจากเซสชันจบเพื่อคัดลอก เมื่อใช้ `--dry-run` จะพิมพ์ข้อความแต่ไม่ถามและไม่คัดลอก
- เมื่อใช้ `mcp_strategy = "sync"` การซิงก์ที่ล้มเหลวจะพิมพ์ `ccam: MCP sync failed (launching anyway): …` และ Claude Code จะเริ่มด้วยเซิร์ฟเวอร์ที่มีอยู่เดิม

## หัวข้อที่เกี่ยวข้อง

- [คำสั่งลัด](/features/shortcuts)
- [การจัดการบัญชี](/features/accounts)
- [ข้อมูลอ้างอิงคำสั่ง](/docs/commands)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [CLAUDE_CONFIG_DIR ทำงานอย่างไร](/guides/claude-config-dir)
