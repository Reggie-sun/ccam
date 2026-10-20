---
title: การตรวจดูส่วนขยาย
description: ccam ext แสดง skill, คำสั่ง, ปลั๊กอิน, เซิร์ฟเวอร์ MCP และ hook ที่แต่ละบัญชีโหลด แหล่งที่มาของแต่ละรายการ และสิ่งที่เสีย
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` ตอบคำถามว่า "ทำไม skill นี้ (หรือเซิร์ฟเวอร์ MCP นี้) ใช้ได้ในบัญชีหนึ่งแต่ใช้ไม่ได้ในอีกบัญชี?" โดยแสดงรายการส่วนขยายที่ Claude Code โหลดสำหรับทุกบัญชี และแหล่งที่มาของแต่ละรายการ: กลุ่มที่ใช้ร่วมกัน โฟลเดอร์ของบัญชีเอง หรือโปรเจกต์ที่คุณอยู่ ใต้ตารางจะแสดงสิ่งที่ทำให้ส่วนขยายโหลดไม่ได้อย่างที่คุณคาดไว้ คำสั่งนี้อ่านอย่างเดียว

## การใช้งาน

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

แต่ละช่องระบุแหล่งที่มาของรายการนั้นในบัญชีนั้น: กลุ่มที่ใช้ร่วมกัน (`default`) `account` เมื่อรายการอยู่ในโฟลเดอร์ของบัญชีเอง หรือ `-` เมื่อบัญชีไม่มีรายการนั้น `!` ทำเครื่องหมายรายการที่มีปัญหา โดยมีคำอธิบายอยู่ใต้ตาราง

## สิ่งที่อ่าน

| ประเภท | ตำแหน่ง | แหล่งที่แสดง |
| --- | --- | --- |
| Skill | แต่ละโฟลเดอร์ใน `skills/` ที่มี `SKILL.md` | กลุ่มที่ `skills/` ลิงก์ไป หรือ `account` |
| คำสั่ง, agent | ไฟล์ `.md` ใต้ `commands/` และ `agents/` โฟลเดอร์ย่อยแสดงเป็น `folder/name` | กลุ่ม หรือ `account` |
| ปลั๊กอิน | `plugins/installed_plugins.json` | กลุ่ม หรือ `account` |
| เซิร์ฟเวอร์ MCP | `mcp.json` ของกลุ่ม (เมื่อใช้ MCP ร่วมกัน) และ `.claude.json` ของบัญชี | กลุ่ม หรือ `account` |
| Hook | `hooks` ใน `settings.json` นับแยกตามเหตุการณ์ | กลุ่ม หรือ `account` |
| โปรเจกต์นี้ | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` และ `.mcp.json` ในโฟลเดอร์ปัจจุบัน | แสดงครั้งเดียว สำหรับทุกบัญชี |

## ปัญหาที่รายงาน

- **ลิงก์เสีย**: ลิงก์ใน `skills/`, `commands/` หรือ `agents/` ที่เป้าหมายหายไปแล้ว หรือโฟลเดอร์ที่ใช้ร่วมกัน (`skills`, `commands`, `agents`, `plugins`) ที่ลิงก์ไปยังตำแหน่งที่ไม่มีอยู่ (ซ่อมด้วย `ccam doctor --fix`)
- **โฟลเดอร์ skill ที่ไม่มี `SKILL.md`**: Claude Code จะข้ามโฟลเดอร์นั้น
- **เซิร์ฟเวอร์ MCP ที่กำหนดไว้ต่างกันในสองที่**: `mcp.json` ของกลุ่มกับ `.claude.json` ของบัญชี หรือที่ใดที่หนึ่งในสองที่นั้นกับ `.mcp.json` ของโปรเจกต์ การตั้งค่าเหมือนกันในทั้งสองที่ (ซึ่ง `mcp_strategy = "sync"` เขียนไว้) ไม่ถือเป็นปัญหา ย้ายเซิร์ฟเวอร์ของบัญชีเข้าไปในกลุ่มด้วย `ccam mcp import <account> <server>`
- **`installed_plugins.json` ที่ไม่ใช่ JSON ที่ถูกต้อง**

## ความปลอดภัยและกรณีพิเศษ

- `ccam ext` ไม่เปลี่ยนแปลงสิ่งใดและไม่ส่งคำขอผ่านเครือข่าย อ่านชื่อเซิร์ฟเวอร์และจำนวน hook เท่านั้น ไม่เคยอ่านค่าของ `env` หรือ header ในนิยาม MCP
- คำสั่งนี้แสดงสิ่งที่อยู่บนดิสก์ ปลั๊กอินที่ติดตั้งไว้แต่ปิดในการตั้งค่า หรือเซิร์ฟเวอร์ที่ถูกบล็อกโดยการตั้งค่าที่มีการจัดการ จะยังคงปรากฏอยู่
- รันในโฟลเดอร์โปรเจกต์เพื่อรวมโฟลเดอร์ `.claude` และ `.mcp.json` ของโปรเจกต์นั้นด้วย

## หัวข้อที่เกี่ยวข้อง

- [การใช้ร่วมกันระหว่างบัญชี](/features/share)
- [เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp)
- [การตรวจสอบและซ่อมแซม](/features/doctor)
