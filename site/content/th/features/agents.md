---
title: เซสชันเบื้องหลังของทุกบัญชี
description: "ccam agents แสดงรายการเอเจนต์เบื้องหลังและเซสชันที่เปิดอยู่ของ Claude Code จากทุกบัญชีในตารางเดียว และ ccam attach เปิดเซสชันหนึ่งจาก id ของเซสชันนั้น"
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code เก็บเซสชันเบื้องหลัง (`claude --bg`, agent view) แยกตามบัญชี: `claude agents` แสดงเฉพาะเซสชันของบัญชีที่ตัวมันทำงานอยู่ `ccam agents` จะถามทุกบัญชีพร้อมกันและแสดงรายการรวมกัน ส่วน `ccam attach <id>` จะเปิดเซสชันเบื้องหลังด้วยบัญชีที่เริ่มเซสชันนั้น ไม่ว่าจะเป็นบัญชีใด

## การใช้งาน

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` คือสถานะของเซสชันเบื้องหลัง (`working`, `blocked`, `done`, `failed`, `stopped`) และสิ่งที่โปรเซสของเซสชันกำลังทำ (`busy`, `waiting`, `idle`) เซสชันแบบโต้ตอบที่เปิดอยู่จะแสดง `open`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--all` | แสดงเซสชันเบื้องหลังที่เสร็จแล้ว ล้มเหลว หรือหยุดไปแล้วด้วย | ปิด |
| `--json` | พิมพ์เซสชันเป็น JSON โดยแต่ละรายการมี `account` ที่เป็นเจ้าของ | ปิด |

`ccam attach <id> [claude args]` ส่งอาร์กิวเมนต์ที่อยู่หลัง id ให้ `claude attach` เมื่อใช้ `--dry-run` จะพิมพ์คำสั่งออกมา

## วิธีการทำงาน

1. สำหรับทุกบัญชี ccam จะรัน `claude agents --json` (พร้อม `--all` เมื่อระบุ) โดยใช้ `CLAUDE_CONFIG_DIR` ของบัญชีนั้น ทุกบัญชีพร้อมกัน บัญชีละไม่เกิน 15 วินาที ขั้นตอนนี้อ่านรายการเซสชันของบัญชีเองและไม่เริ่มบริการเบื้องหลังของ Claude Code
2. แต่ละรายการจะได้ชื่อของบัญชีกำกับไว้ และรายการจะเรียงตามเวลาที่เริ่ม ใหม่สุดก่อน บัญชีที่ไม่ตอบกลับจะถูกระบุชื่อทาง standard error และไม่ถูกรวมไว้ จากนั้น ccam จะออกด้วยรหัส 1 เพื่อไม่ให้สคริปต์เข้าใจผิดว่าความล้มเหลวคือรายการว่าง
3. `ccam attach` ค้นหา `jobs/<id>` ในโฟลเดอร์ของแต่ละบัญชี ซึ่งเป็นที่ที่ Claude Code เก็บเซสชันเบื้องหลัง และรัน `claude attach <id>` ด้วยบัญชีนั้น แบบเดียวกับที่ `ccam <name>` เริ่ม Claude Code

## ความปลอดภัยและกรณีพิเศษ

- id ของเซสชันเบื้องหลังเป็นของบัญชีเดียว: `claude attach` ที่รันด้วยบัญชีอื่นจะหาไม่เจอ ซึ่ง `ccam attach` จัดการเรื่องนี้ให้
- สองบัญชีอาจมีเซสชันที่ id ขึ้นต้นเหมือนกัน เมื่อเป็นเช่นนั้น `ccam attach` จะระบุทั้งสองเซสชันและขอให้พิมพ์ id เพิ่ม
- ccam ส่งเซิร์ฟเวอร์ MCP ที่ใช้ร่วมกันให้ Claude Code ด้วย `--mcp-config=<file>` เป็นอาร์กิวเมนต์เดียว ซึ่ง `claude agents` และ `claude attach` รับได้ และเซสชันเบื้องหลังยังใช้ไฟล์นั้นต่อเมื่อเริ่มใหม่
- บริการเบื้องหลังของบัญชี API key รับคีย์จาก shell ที่เริ่มบริการนั้น ให้เริ่มเซสชันเบื้องหลังแรกของบัญชีแบบนี้ผ่าน ccam (`ccam gw --bg …`) เพื่อให้บริการมีคีย์

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [การเลือกบัญชีสำหรับสคริปต์](/features/pick)
- [แดชบอร์ด](/features/dashboard)
