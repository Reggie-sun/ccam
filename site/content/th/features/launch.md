---
title: การเริ่ม Claude Code
description: เริ่ม Claude Code ในชื่อบัญชีใดก็ได้ด้วย ccam <name|number>, ccam run หรือแดชบอร์ด อาร์กิวเมนต์ส่งต่อให้ claude ตามเดิม และ --dry-run แสดงคำสั่ง
updated: 2026-10-24
order: 1
translation: machine
sourceHash: "e279341eb59b"
---

`ccam work` เริ่ม Claude Code ในชื่อบัญชี `work` และทุกอย่างหลังชื่อบัญชีจะส่งต่อให้ `claude` ตามเดิม ใช้ได้ทุกที่ที่คุณจะพิมพ์ `claude`: ในเทอร์มินัล ในสคริปต์ หรือเมื่อส่งอินพุตเข้ามาทางไปป์

## การใช้งาน

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

ระบุชื่อบัญชีหรือหมายเลขบัญชีจาก `ccam list` ccam ไม่อ่านอาร์กิวเมนต์ที่อยู่หลังชื่อบัญชี ดังนั้น `ccam work --help` จะแสดงความช่วยเหลือของ Claude Code จะใส่ `--` หลังชื่อบัญชีหรือไม่ก็ได้ อินพุตมาตรฐาน รหัสออก (exit code) และสัญญาณ (signal) ทำงานเหมือนกับ `claude`

เมื่ออาร์กิวเมนต์แรกขึ้นต้นด้วย `-` และไม่ใช่แฟล็กของ ccam เอง ccam จะเริ่มบัญชีโดยไม่ระบุชื่อพร้อมอาร์กิวเมนต์ทั้งหมด ซึ่งก็คือบัญชีที่ตั้งไว้สำหรับโฟลเดอร์ปัจจุบันด้วย [`ccam folder`](/features/folder) หากไม่มีจะใช้บัญชีเริ่มต้น: บัญชีที่ตั้งไว้ด้วย `ccam use` หรือบัญชีที่มีหมายเลขต่ำสุดหากคุณยังไม่ได้ตั้ง

### แดชบอร์ด

`ccam` ที่ไม่มีอาร์กิวเมนต์จะเปิด[แดชบอร์ด](/features/dashboard) เมื่ออินพุตและเอาต์พุตมาตรฐานเป็นเทอร์มินัล และจะพิมพ์ `ccam list` ในกรณีอื่น เลื่อนไปที่บัญชี (ด้วยปุ่มลูกศร หรือพิมพ์หมายเลขของบัญชี) แล้วกด `Enter` เพื่อเริ่มบัญชีนั้น `n` จะเริ่มบัญชีที่เหลือโควตามากที่สุด และ `o` จะเสนอ `--continue`, `--resume` หรืออาร์กิวเมนต์อื่น แดชบอร์ดจะปิดก่อนที่ Claude Code จะเริ่ม เทอร์มินัลจึงกลับสู่สภาพเดิมก่อน

**ที่นั่ง (seat) ของ Team และ Enterprise** การเลือกอัตโนมัติ (`ccam next`, `--fallback`, `ccam pick` และปุ่ม `n` ในแดชบอร์ด) จะสลับระหว่างที่นั่งต่าง ๆ ในองค์กร Claude Team หรือ Enterprise เดียวกันด้วย ศูนย์ช่วยเหลือของ Anthropic ระบุว่า "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (แต่ละที่นั่งในแพ็กเกจ Team จะกำหนดให้กับบุคคลที่ระบุชื่อเพียงคนเดียว และมีเพียงบุคคลนั้นที่ใช้ได้) และ "Users can't route their work through another user's account to use that user's remaining usage" (ผู้ใช้ไม่สามารถส่งงานของตนผ่านบัญชีของผู้ใช้คนอื่นเพื่อใช้ปริมาณการใช้งานที่เหลืออยู่ของผู้ใช้คนนั้นได้) ([แพ็กเกจ Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ส่วน [แพ็กเกจ Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ก็เช่นเดียวกัน) ดังนั้น การใช้การเลือกอัตโนมัติระหว่างที่นั่งต่าง ๆ จึงอาจละเมิดข้อกำหนดของ Anthropic โปรดสอบถามองค์กรของคุณและ Anthropic ก่อนที่จะพึ่งพาการเลือกอัตโนมัติ สำหรับแพ็กเกจส่วนตัว (Pro, Max) ดูที่ [หมายเหตุเรื่องข้อกำหนดในคำถามที่พบบ่อย](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)

เมื่อ shell ตั้งค่า `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` หรือ `ANTHROPIC_AUTH_TOKEN` ไว้ Claude Code จะใช้คีย์หรือโทเค็นนั้นแทนการเข้าสู่ระบบของบัญชีเอง ccam จึงแจ้งเรื่องนี้ทางข้อผิดพลาดมาตรฐานก่อนเริ่ม (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`) ccam เพียงอ่านตัวแปรเหล่านี้ ส่วน[บัญชี API key](/features/api)จะล้างตัวแปรเหล่านี้แทน

เมื่อตั้ง `relay = true` ที่ด้านบนสุดของ `config.toml` (เหนือ `[[group]]`, `[[profile]]` หรือ `[[folder]]` ตัวแรก) ccam จะเริ่มเซสชันแบบโต้ตอบเป็นโปรเซสลูกและอยู่ข้าง ๆ เพื่อให้ `/ccam:switch` พาการสนทนาไปต่อได้เอง พร้อมตัวเลือกของเซสชัน โหมดสิทธิ์ และระดับ effort ([รายละเอียด](/features/next#relay)) หากไม่ได้ตั้ง ccam จะแทนที่ตัวเองด้วย Claude Code และไม่มีอะไรทำงานอยู่ข้าง ๆ

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--dry-run` | พิมพ์ตัวแปรสภาพแวดล้อมและคำสั่งที่ ccam จะรัน แล้วออก | ปิด |
| `--fallback` | สำหรับการรันแบบ `-p`/`--print`: เมื่อบัญชีถึงขีดจำกัดการใช้งาน จะรันคำสั่งเดิมอีกครั้งด้วยบัญชีที่เข้าสู่ระบบอยู่และเหลือโควตามากที่สุด ([รายละเอียด](/features/next#run-again-at-the-limit)) | ปิด |
| `-v`, `--version` | พิมพ์ `ccam <version>` (เฉพาะเมื่อเป็นอาร์กิวเมนต์แรก) | |
| `-h`, `--help` | ความช่วยเหลือของ ccam ในรูป `ccam --help` หรือ `ccam run --help` | |

กฎของอาร์กิวเมนต์:

- ccam อ่านแฟล็กของตัวเองเฉพาะที่อยู่ก่อนชื่อบัญชี หลังชื่อบัญชี ทุกอาร์กิวเมนต์เป็นของ Claude Code รวมถึง `--dry-run` และ `--help`
- `--dry-run` ใช้ได้ในรูป `ccam --dry-run work`, `ccam run --dry-run work` และ `ccam --dry-run` (บัญชีที่ไม่ได้ระบุชื่อ) ตัวเลือกนี้ได้รับการรองรับ ไม่ถูกปฏิเสธ `--dry-run` และ `--fallback` จะใส่ลำดับใดก่อนก็ได้
- `ccam run <name|number> [claude args]` คือรูปแบบเต็มของ `ccam <name>` ส่วน `ccam run` ที่ไม่มีชื่อบัญชีจะใช้บัญชีของโฟลเดอร์หรือบัญชีเริ่มต้น
- ชื่อบัญชีต้องไม่ซ้ำกับชื่อคำสั่งของ ccam เช่น `add`, `list`, `run` หรือ `usage` ดังนั้น `ccam <word>` จึงไม่กำกวม

การทดลองรัน (dry run) จะพิมพ์หนึ่งบรรทัดที่คุณนำไปวางในเชลล์ได้ หรือใน PowerShell บน Windows บรรทัดนั้นเริ่มด้วยสิ่งที่การเริ่มเปลี่ยนในตัวแปรสภาพแวดล้อม: `CLAUDE_CONFIG_DIR` และ `CCAM_PROFILE` เสมอ `env -u <name>` สำหรับตัวแปรที่ถูกตัดออก (เช่น `CCAM_RELAY` ที่สืบทอดมา) และสำหรับ[บัญชี API key](/features/api) จะมี endpoint และโมเดล และคีย์ในรูปแหล่งที่มา (`ANTHROPIC_API_KEY="$WORK_KEY"`) ไม่เคยแสดงตัวคีย์เอง:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

บรรทัดของ PowerShell จะคืนค่าตัวแปรกลับหลังจากนั้น แบบเดียวกับที่ `VAR=x cmd` ทำในเชลล์ และส่งบรรทัดคำสั่งให้ Claude Code ผ่าน `Start-Process` ดังนั้น Windows PowerShell 5.1 และ PowerShell 7 จึงรันได้แบบเดียวกัน เมื่อการเริ่มตั้งตัวแปรเป็นค่าว่าง จะมีบรรทัด `#` ด้านล่างบอกว่า PowerShell ก่อนรุ่น 7.5 จะลบตัวแปรนั้นแทน สำหรับ `claude.cmd` ที่ติดตั้งด้วย npm บรรทัดนั้นจะรันคำสั่ง `cmd.exe` ที่การเริ่มจะใช้ (ดูด้านล่าง) และอาร์กิวเมนต์ที่การเริ่มจะปฏิเสธในกรณีนั้นจะทำให้การทดลองรันล้มเหลวด้วย

## วิธีการทำงาน

การเริ่มทุกครั้งผ่านขั้นตอนเดียวกัน ตั้งแต่บรรทัดคำสั่งไปจนถึง `exec`:

<!-- motion:launch -->

แต่ละกล่องจะสว่างขึ้นเมื่อขั้นตอนนั้นทำงาน และบรรทัดด้านล่างคือตัวแปรสภาพแวดล้อมและอาร์กิวเมนต์ที่ Claude Code ใช้เริ่มทำงาน

1. `cmd/ccam/main.go` เรียก `app.Main` หาก ccam ถูกเริ่มด้วยชื่ออื่น แสดงว่าเป็น[คำสั่งลัด](/features/shortcuts) และจะไปที่ตัวเริ่มโปรแกรมโดยตรง
2. `app.parse` ตรวจ `os.Args` ก่อนที่ cobra จะเข้ามาเกี่ยวข้อง มีเพียง `ccam` เดี่ยว ๆ (แดชบอร์ด) ชื่อคำสั่งของ ccam (`config.Reserved`) `-h`/`--help` และ `-v`/`--version` ที่ไปยัง `cli.Execute` นอกนั้นถือเป็นการเริ่ม Claude Code
3. `launcher.Launch` อ่าน `~/.config/ccam/config.toml` ซึ่งเป็นไฟล์เดียวที่พาร์ส
4. `launcher.Prepare` หาบัญชีจากชื่อหรือหมายเลข (หากไม่มีชื่อ `HereProfile` จะเลือกบัญชีของโฟลเดอร์ โดยเปรียบเทียบไดเรกทอรีทำงานกับรายการ `[[folder]]` ใน `config.toml` แบบสตริง โดยไม่สนตัวพิมพ์เล็กใหญ่บน macOS และ Windows) ตัด `--` ตัวแรกที่นำหน้าออกหนึ่งตัว และหาไฟล์ปฏิบัติการ `claude` ด้วย `ResolveClaude` ผลการค้นหาใน PATH จะถูกแคชไว้ใน `~/.local/state/ccam/claude-path` พร้อมเวลาแก้ไขของไฟล์และ hash ของ `PATH` ccam จึงไม่ต้องค้นหาใน `PATH` ทุกครั้งที่เริ่ม หากไฟล์เปลี่ยนหรือหายไป หรือ `PATH` เปลี่ยน จะค้นหาใหม่ พาธ `claude_bin` ใน `config.toml` (ขยาย `~/` ให้) ต้องเป็นไฟล์ที่รันได้ มิฉะนั้นการเริ่มจะหยุดด้วย `claude_bin X is not an executable file. Check claude_bin in config.toml` ก่อนที่จะเขียนสิ่งใด ccam จะปฏิเสธ `claude` ที่เป็นตัว ccam เอง (เปรียบเทียบด้วย inode) เพราะจะทำให้วนซ้ำไม่สิ้นสุด
5. `checkLinks` เปรียบเทียบโฟลเดอร์บัญชี `~/.local/share/ccam/profiles/<name>/` กับการตั้งค่าการใช้ร่วมกันด้วย `share.Check` ซึ่งเรียกเพียง `lstat` และ `readlink` หนึ่งหรือสอง system call ต่อรายการ ลิงก์ที่หายไปจะถูกสร้างทันที สิ่งใดที่อาจทำให้ข้อมูลสูญหายจะหยุดการเริ่ม
6. `applyMCP` เพิ่มเซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน (ดูด้านล่าง)
7. `session.Register` บันทึกเซสชัน ccam เขียนไฟล์เครื่องหมายที่ตั้งชื่อตาม PID ซึ่งเก็บเวลาเริ่มของโปรเซส ไว้ที่ `~/.local/state/ccam/locks/profile-<name>/<pid>` ไม่ว่า `XDG_RUNTIME_DIR` จะเป็นอะไร ทุกเชลล์จึงเห็นเครื่องหมายชุดเดียวกัน สำหรับบัญชีที่ใช้ประวัติการสนทนาร่วมกัน `historyCheck` จะค้นหาเซสชันที่กำลังทำงานใน `history-<group>-<hash>` ก่อน โดย hash คำนวณจากโฟลเดอร์ปัจจุบันหลังแปลง symbolic link เป็นพาธจริงแล้ว แล้วลงทะเบียนไว้ที่นั่นด้วย `exec` คง PID ไว้ เครื่องหมายจึงระบุถึงโปรเซสของ Claude Code เอง ที่ใช้เครื่องหมายแทนการถือล็อก เพราะโปรเซสที่ Claude Code เริ่มจะสืบทอดล็อกไปและถือไว้ต่อหลังจากเซสชันจบแล้ว เครื่องหมายที่เขียนไม่ได้จะไม่หยุดการเริ่ม: ccam จะแจ้ง `ccam: could not mark the session as running (…); starting Claude Code anyway` การตรวจอื่นของ ccam ว่าบัญชีกำลังทำงานอยู่หรือไม่ (การเปลี่ยนชื่อ การลบ การเปลี่ยนกลุ่ม `doctor --fix` การซิงก์ MCP) ยังอ่านบันทึกของ Claude Code เองในโฟลเดอร์บัญชีด้วย (`sessions/<pid>.json`, `daemon.lock`) จึงเห็น Claude Code ที่เริ่มโดยไม่ผ่าน ccam ด้วย: `claude` แบบปกติ, IDE, เซสชันเบื้องหลังและ daemon ของเซสชันเหล่านั้น
8. `Plan.Exec` ตั้งชื่อหน้าต่างเทอร์มินัลเป็น `claude · <name>` (เมื่อ standard error เป็นเทอร์มินัล) เขียนเวลาที่เริ่มลงใน `~/.local/state/ccam/last-used/<name>` โดยไม่เรียก `fsync` แล้วเรียก `syscall.Exec`

Claude Code ทำงานโดยตั้ง `CLAUDE_CONFIG_DIR` เป็นโฟลเดอร์ของบัญชี และ `CCAM_PROFILE` เป็นชื่อบัญชี ccam จะลบค่าของตัวแปรเหล่านี้ที่สืบทอดมา รวมถึง `CCAM_RELAY` ที่สืบทอดมาด้วย สำหรับโฟลเดอร์ `~/.claude` ที่ลงทะเบียนด้วย `ccam import` บน macOS ccam จะตั้ง `CLAUDE_SECURESTORAGE_CONFIG_DIR` เป็นค่าว่างด้วย เพื่อให้บัญชีนั้นและ `claude` แบบปกติใช้การเข้าสู่ระบบเดียวกันในพวงกุญแจ (Keychain) บน Windows ก็ทำเช่นเดียวกันเมื่อ Claude Code เก็บการเข้าสู่ระบบไว้ใน Credential Manager

### ทำไมจึงใช้ exec

`syscall.Exec` แทนที่โปรเซสของ ccam ด้วย `claude` ccam จึงไม่ค้างอยู่ในหน่วยความจำ การกดแป้น สัญญาณ และรหัสออกเป็นของ Claude Code เอง และอาร์กิวเมนต์ยังคงเป็นอาร์เรย์ ไม่เคยเป็นสตริงของเชลล์

Windows ไม่มี `exec` บน Windows `Plan.Exec` จะเริ่ม `claude.exe` เป็นโปรเซสลูกด้วยตัวแปรสภาพแวดล้อม อาร์กิวเมนต์ และคอนโซลเดียวกัน รอให้มันจบ แล้วออกด้วยรหัสออกของมัน `claude.cmd` ที่ npm ติดตั้งจะรันใน `cmd.exe` ด้วยพาธเต็ม และทุกอาร์กิวเมนต์จะถูกใส่เครื่องหมายคำพูดเพื่อไม่ให้ `cmd.exe` ตีความอาร์กิวเมนต์ใดเป็นไวยากรณ์ของตัวเอง (`&`, `%PATH%`) อาร์กิวเมนต์ที่มีการขึ้นบรรทัดใหม่ส่งด้วยวิธีนี้ไม่ได้และจะถูกปฏิเสธ พร้อมคำแนะนำให้ติดตั้ง Claude Code แบบเนทีฟ Ctrl+C ไปถึง Claude Code จากคอนโซลโดยตรง ccam จึงไม่สนใจสัญญาณนี้ ไฟล์เครื่องหมายของเซสชันใช้ PID ของ ccam ซึ่งมีอายุเท่ากับ Claude Code พอดี และจะถูกลบเมื่อ Claude Code จบการทำงาน ระหว่างรอ ccam ไม่ทำอย่างอื่น งานหลังจบเซสชัน เช่น การตรวจพบเซิร์ฟเวอร์ที่เพิ่มด้วย `claude mcp add` จะเกิดขึ้นในการเริ่มครั้งถัดไป: ccam เปรียบเทียบเวลาแก้ไขของ `.claude.json` กับ `~/.local/state/ccam/mcp/<name>/seen` และพาร์สไฟล์เฉพาะเมื่อไฟล์เปลี่ยนไป

### บรรทัดคำสั่งถูกส่งต่ออย่างไร

สามบรรทัดคำสั่ง สามเส้นทางผ่าน `internal/app`:

<!-- motion:launch-route -->

คำสั่งลัดจับคู่จาก `argv[0]` ส่วน `ccam <name>` จับคู่จากอาร์กิวเมนต์แรก และทั้งสองจบที่ `syscall.Exec` มีเพียงคำสั่งของ ccam เช่น `ccam usage` เท่านั้นที่ไปถึง cobra

### เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน

ด้วยกลยุทธ์ `flag` ซึ่งเป็นค่าเริ่มต้น ccam ส่งเซิร์ฟเวอร์ MCP ของกลุ่มให้ Claude Code เป็นอาร์กิวเมนต์ โดยวางไว้ก่อนอาร์กิวเมนต์ของคุณ:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` คืนอาร์กิวเมนต์เฉพาะเมื่อ `mcp.json` ของกลุ่มมีอยู่ และ `Plan.InsertArgs` วางอาร์กิวเมนต์ไว้ต่อจาก `claude` ทันที `--mcp-config` หรือ `--strict-mcp-config` ของคุณเองจึงยังมีผลซ้อนทับอยู่ด้านบน อาร์กิวเมนต์นี้เป็นอาร์กิวเมนต์เดียว `--mcp-config=<path>` เพราะ `--mcp-config` ของ Claude Code รับได้หลายค่า หากเขียนเป็นสองอาร์กิวเมนต์ ก็จะอ่านพรอมต์ที่ตามมาเป็นไฟล์การตั้งค่าไปด้วย

`claudeVerb` ตัดสินว่าแฟล็กใส่ไม่ได้ที่ใด โดยอ่านคำสั่งย่อยของ Claude Code จากอาร์กิวเมนต์แบบเดียวกับที่ Claude Code ทำ คือข้ามตัวเลือกพร้อมค่าของตัวเลือก `remote-control` หรือ `rc` จะไม่ได้รับแฟล็ก (ดูด้านล่าง) เช่นเดียวกับคำกริยาของเซสชันเบื้องหลังที่ระบุเป็นตัวแรก (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) `ccam work mcp list`, `get` และ `remove` ก็ไม่ได้รับแฟล็กเช่นกัน เพราะ Claude Code ปฏิเสธแฟล็กในคำสั่งเหล่านั้น ccam จึงแจ้งทาง standard error แทนว่าเซิร์ฟเวอร์ของกลุ่มก็โหลดในเซสชันของบัญชีนี้ด้วย:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

คำสั่งย่อยอื่นทุกคำสั่งและทุกเซสชันยังคงได้รับแฟล็ก

เมื่อตั้ง `mcp_strategy = "sync"` ใน `config.toml` ccam จะเขียนเซิร์ฟเวอร์ลงใน `.claude.json` ของบัญชีแทน และจะข้ามขั้นตอนนี้หากบัญชีกำลังทำงานอยู่แล้ว การเริ่มจะข้ามงานนี้ด้วยเมื่อทั้ง `mcp.json` ของกลุ่มและ `.claude.json` ไม่เปลี่ยนตั้งแต่การตรวจครั้งล่าสุด การซิงก์จะไม่แทนที่เซิร์ฟเวอร์ที่บัญชีมีเป็นของตัวเองภายใต้ชื่อเดียวกับเซิร์ฟเวอร์ของกลุ่ม (และจะแจ้งหนึ่งครั้ง) จะใส่สำเนาที่หายไปจาก `.claude.json` กลับคืน (`put back group default's MCP servers …`) และจะเปลี่ยนหรือลบเฉพาะสำเนาที่ตัวเองเขียนไว้และไม่มีใครเปลี่ยนตั้งแต่นั้น

### Remote Control

`ccam work rc` (หรือ `ccam work remote-control`) เริ่ม Remote Control ของ Claude Code ในฐานะ `work` ส่วน `claude remote-control` ปฏิเสธ `--mcp-config` ที่อยู่ก่อนคำกริยา (verb) ของมัน และไม่รับเลยหลังคำกริยา อีกทั้งเซสชันของมันอ่านเซิร์ฟเวอร์ MCP จาก `.claude.json` ของบัญชีเท่านั้น ดังนั้นกลยุทธ์ flag จึงทำงานต่างออกไปในกรณีนี้:

<!-- motion:launch-rc -->

1. `claudeVerb` พบ `remote-control` หรือ `rc` เป็นคำกริยาของ Claude Code แม้จะอยู่หลังตัวเลือกและค่าของตัวเลือก (`ccam work --model opus rc`) จากนั้น `mcp.FlagArgs` จะไม่เพิ่มอะไรเลย
2. `mcp.Synced` เปรียบเทียบเวลาแก้ไขของ `mcp.json` ของกลุ่มกับ `sync-stamp` ของบัญชี และตรวจว่าเซิร์ฟเวอร์ของกลุ่มทุกตัวยังอยู่ใน `.claude.json` หาก `.claude.json` มีเซิร์ฟเวอร์ปัจจุบันของกลุ่มอยู่แล้ว หรือกลุ่มไม่มีเซิร์ฟเวอร์เลย Remote Control จะเริ่มโดยไม่ถาม
3. มิฉะนั้น ccam จะระบุชื่อเซิร์ฟเวอร์ที่เซสชันจะไม่ได้รับ แล้วถามว่า:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

4. เมื่อตอบ `y` แล้ว `mcp.Sync` จะคัดลอกเซิร์ฟเวอร์เหล่านั้นแบบเดียวกับกลยุทธ์ `sync`: ภายใต้ล็อกของ Claude Code หลังจากสำรองไฟล์ และคงเซิร์ฟเวอร์ของบัญชีเองไว้ เมื่อตอบ `N` หรือเมื่อไม่มีเทอร์มินัล Remote Control จะเริ่มโดยไม่มีเซิร์ฟเวอร์เหล่านั้น และ ccam จะแนะนำ `mcp_strategy = "sync"` ซึ่งคัดลอกให้ก่อนการเริ่มทุกครั้ง

ccam จะถามทุกครั้งที่เริ่ม Remote Control จนกว่าจะคัดลอกเซิร์ฟเวอร์แล้ว และหลังจากคัดลอกแล้วจะถามอีกเฉพาะเมื่อ `mcp.json` ของกลุ่มเปลี่ยนอีกครั้ง หรือเมื่อเซิร์ฟเวอร์ที่คัดลอกไว้หายไปจาก `.claude.json` สำเนาใน `.claude.json` มีผลเหนือกว่า `--mcp-config` ในการเริ่มตามปกติครั้งต่อ ๆ ไป ดังนั้น `ccam mcp add`, `remove` และ `import` จึงอัปเดตสำเนาที่บัญชีอื่นถืออยู่ด้วย ([เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp))

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
- หากมี Claude Code อีกตัวที่ใช้ประวัติการสนทนาร่วมกับบัญชีนี้ (กลุ่มเดียวกัน) ทำงานอยู่ในโฟลเดอร์เดียวกัน ccam จะเตือนว่าประวัติอาจปนกันและถาม `Run anyway? [y/N]` โฟลเดอร์จะถูกจับคู่หลังแปลง symbolic link เป็นพาธจริงแล้ว หากไม่มีเทอร์มินัล (standard input จาก `/dev/null` ถือว่าไม่มี) หรือเมื่อรันจากงานเบื้องหลัง จะแสดงคำเตือนแล้วทำงานต่อ การรัน `-p` ใหม่จะเริ่มการสนทนาของตัวเองและไม่ได้รับคำเตือน ส่วน `--version`, `--help` และคำสั่งย่อยที่ไม่เปิดการสนทนา เช่น `mcp`, `auth` หรือ `plugin` จะไม่ถูกตรวจเลย
- `--dry-run` ข้ามการซิงก์ MCP การแจ้งเซิร์ฟเวอร์ใหม่ หมายเหตุของ `mcp list` การตรวจประวัติ เครื่องหมายเซสชัน ชื่อหน้าต่าง และเวลาที่ใช้ล่าสุด และไม่เรียก exec แต่ยังคงตรวจลิงก์ที่ใช้ร่วมกัน และสร้างลิงก์ที่หายไปเหมือนการเริ่มปกติ
- `ccam work rc` ขณะที่ `work` กำลังทำงานอยู่แล้ว จะไม่แตะ `.claude.json` แจ้งให้ทราบ และเริ่ม Remote Control โดยไม่มีเซิร์ฟเวอร์ของกลุ่ม ให้รันอีกครั้งหลังจากเซสชันจบเพื่อคัดลอก เมื่อใช้ `--dry-run` จะพิมพ์ข้อความแต่ไม่ถามและไม่คัดลอก
- เมื่อใช้ `mcp_strategy = "sync"` การซิงก์ที่ล้มเหลวจะพิมพ์ `ccam: MCP sync failed (launching anyway): …` และ Claude Code จะเริ่มด้วยเซิร์ฟเวอร์ที่มีอยู่เดิม

## หัวข้อที่เกี่ยวข้อง

- [คำสั่งลัด](/features/shortcuts)
- [การใช้ Remote Control กับ ccam](/guides/remote-control)
- [การจัดการบัญชี](/features/accounts)
- [ข้อมูลอ้างอิงคำสั่ง](/docs/commands)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [CLAUDE_CONFIG_DIR ทำงานอย่างไร](/guides/claude-config-dir)
