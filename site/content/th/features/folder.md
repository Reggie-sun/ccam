---
title: บัญชีตามโฟลเดอร์
description: "ccam folder ทำให้โฟลเดอร์หนึ่งและทุกโฟลเดอร์ย่อยใช้บัญชีเดียวเมื่อคุณเริ่ม Claude Code โดยไม่ระบุชื่อบัญชี เช่น ใช้ work กับโค้ดของบริษัท"
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` ทำให้ทุกการเปิดที่ไม่ระบุชื่อบัญชีภายใน `~/code/company` และในทุกโฟลเดอร์ย่อย เริ่ม Claude Code ในชื่อบัญชี `work` นอกโฟลเดอร์นั้น บัญชีเริ่มต้นจะเริ่มเหมือนเดิม การระบุชื่อบัญชีด้วย `ccam personal` หรือคำสั่งลัดอย่าง `ccpersonal` จะเริ่มบัญชีนั้นเสมอ

## การใช้งาน

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

การเปิดต่อไปนี้ใช้บัญชีของโฟลเดอร์:

- `ccam -p "…"` และการเปิดอื่นใดที่อาร์กิวเมนต์แรกเป็นตัวเลือกของ Claude Code
- `ccam run` ที่ไม่มีชื่อบัญชี และ `ccam --dry-run`
- `ccam --fallback -p "…"` สำหรับการลองครั้งแรก
- ตัวเลือกบัญชี (`ccam` อย่างเดียว) ซึ่งจะเริ่มโดยวางเคอร์เซอร์ไว้ที่บัญชีของโฟลเดอร์

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `<account>` | ชื่อหรือหมายเลขของบัญชีที่จะใช้ | |
| `[folder]` | โฟลเดอร์ แบบพาธสัมพัทธ์หรือพาธสัมบูรณ์ก็ได้ | โฟลเดอร์ปัจจุบัน |
| `--remove` | ลบรายการของโฟลเดอร์นั้นเท่านั้น | ปิด |
| `--dry-run` | รองรับ: พิมพ์การเปลี่ยนแปลงโดยไม่บันทึก | ปิด |

## วิธีการทำงาน

แต่ละโฟลเดอร์คือรายการ `[[folder]]` หนึ่งรายการใน `~/.config/ccam/config.toml` ซึ่งเป็นไฟล์ที่ทุกการเปิดอ่านอยู่แล้ว:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` แปลงโฟลเดอร์เป็นพาธสัมบูรณ์โดยแก้ลิงก์สัญลักษณ์ทั้งหมดแล้วจึงบันทึก การตั้งโฟลเดอร์เดิมอีกครั้งจะแทนที่บัญชีของโฟลเดอร์นั้น
2. การเปิดที่ไม่มีชื่อบัญชีจะเรียก `launcher.HereProfile` หากไม่มีรายการใดเลย จะใช้บัญชีเริ่มต้นทันที
3. มิฉะนั้นจะเปรียบเทียบไดเรกทอรีทำงานกับ `path` แต่ละรายการแบบสตริง: ต้องเป็นโฟลเดอร์นั้นเองหรือโฟลเดอร์ที่อยู่ข้างใต้ ไม่ใช่โฟลเดอร์ที่แค่ขึ้นต้นด้วยตัวอักษรเดียวกัน (`/code/company` ไม่ครอบคลุม `/code/companyx`) รายการที่ลึกที่สุดจะถูกเลือก ดังนั้น `~/code` จึงเป็น `personal` ได้ ขณะที่ `~/code/company` เป็น `work`
4. เฉพาะเมื่อไดเรกทอรีทำงานตามที่เป็นอยู่ไม่ตรงกับรายการใดเลย ccam จึงจะแก้ลิงก์สัญลักษณ์ของไดเรกทอรีนั้นแล้วเปรียบเทียบอีกครั้ง
5. บน Windows การเปรียบเทียบจะไม่สนใจตัวพิมพ์เล็กและตัวพิมพ์ใหญ่ เช่นเดียวกับที่ Windows ทำ

ไม่มีการอ่านไฟล์อื่นนอกจาก `config.toml` การเปิดภายในโฟลเดอร์จึงไม่เสียเวลาที่วัดได้ (4.0 ms ทั้งเมื่อมีและไม่มีรายการ, `hyperfine`, 300 รอบ)

## ความปลอดภัยและกรณีพิเศษ

- `ccam rename` ย้ายรายการโฟลเดอร์ไปยังชื่อใหม่ และ `ccam remove` ลบรายการของบัญชีที่ถูกลบ
- รายการที่บัญชีไม่มีอยู่แล้วจะหยุดการเปิดที่ไม่มีชื่อบัญชีในโฟลเดอร์นั้น: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder` ส่วน `ccam doctor` จะรายงานรายการเหล่านี้ และรายการที่โฟลเดอร์หายไปแล้ว
- `ccam folder --remove` ภายในโฟลเดอร์ที่ได้บัญชีมาจากโฟลเดอร์ด้านบน จะบอกชื่อโฟลเดอร์ด้านบนนั้นแทนการลบสิ่งใด
- `ccam next` เลือกบัญชีตามการหมุนเวียนหรือการใช้งาน และไม่ดูโฟลเดอร์

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [คำสั่งลัด](/features/shortcuts)
- [การจัดการบัญชี](/features/accounts)
- [ข้อมูลอ้างอิงคำสั่ง](/docs/commands)
