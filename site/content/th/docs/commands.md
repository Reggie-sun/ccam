---
title: ข้อมูลอ้างอิงคำสั่ง
description: "ทุกคำสั่งและตัวเลือกของ ccam ตั้งแต่ ccam add และ ccs <name> ไปจนถึง ccam usage, ccam next และ ccam doctor"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` และ `ccs` เป็นโปรแกรมเดียวกัน ใช้ตัวไหนก็ได้ตามที่คุณสะดวก ใส่ `--dry-run` ไว้หน้าคำสั่งใดก็ได้เพื่อดูว่าคำสั่งนั้นจะทำอะไรโดยไม่ลงมือทำจริง เช่น `ccs --dry-run work`

<!-- commands -->

## การเริ่ม Claude Code

`ccs <name|number> [arguments]` เริ่ม Claude Code ในชื่อบัญชีนั้น ccam ไม่อ่านอาร์กิวเมนต์ที่อยู่หลังชื่อบัญชี แต่ส่งต่อให้ `claude` ตามลำดับเดิม ส่วนอินพุตมาตรฐาน รหัสออก (exit code) และสัญญาณ (signal) ทำงานเหมือนกับการใช้ `claude` โดยตรงทุกประการ จะใส่ `--` หลังชื่อบัญชีหรือไม่ก็ได้

`ccs` ที่มีอาร์กิวเมนต์แต่ไม่มีชื่อบัญชี เช่น `ccs -p "summarise"` จะใช้บัญชีเริ่มต้น ตั้งบัญชีเริ่มต้นได้ด้วย `ccam use <name>`

คำสั่งลัด `claude1`, `claude2` และต่อไปเรื่อย ๆ เป็นลิงก์ไปยัง ccam ใน `~/.local/bin` นอกจากนี้ `ccam shims install --names` ยังสร้าง `claude-<name>` ให้กับทุกบัญชีด้วย

## การเพิ่มและลบบัญชี

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` ที่ไม่ระบุโฟลเดอร์จะค้นหา `~/.claude` และ `~/.claude-*` แล้วถามคุณทีละโฟลเดอร์ โฟลเดอร์ที่นำเข้าจะถูกลงทะเบียนไว้ในตำแหน่งเดิม ccam ไม่ย้ายหรือคัดลอกการเข้าสู่ระบบใด ๆ

## การใช้ร่วมกัน

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

แต่ละรายการรับค่า `on` (กลุ่ม `default`), `off` หรือชื่อกลุ่ม ดูที่ [การใช้ร่วมกัน](/docs/sharing)

## เซิร์ฟเวอร์ MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## การใช้งานและการสลับบัญชี

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## การบำรุงรักษา

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ไฟล์

| พาธ | เนื้อหา |
| --- | --- |
| `~/.config/ccam/config.toml` | บัญชีและสิ่งที่แต่ละบัญชีใช้ร่วมกัน |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ของแต่ละบัญชี |
| `~/.local/share/ccam/groups/<group>/` | การตั้งค่า เซิร์ฟเวอร์ MCP ส่วนขยาย และประวัติที่ใช้ร่วมกัน |
| `~/.local/state/ccam/` | แคชการใช้งานและเวลาที่ใช้ล่าสุด |
| `~/.local/bin/claude1`, `ccs` | ลิงก์ไปยัง ccam |

ตัวแปร `XDG_CONFIG_HOME`, `XDG_DATA_HOME` และ `XDG_STATE_HOME` ใช้เปลี่ยนตำแหน่งของโฟลเดอร์เหล่านี้
