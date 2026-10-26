---
title: บัญชี API key
description: ccam add --api เพิ่มบัญชี API key สำหรับ API ของ Anthropic หรือเกตเวย์ คีย์ถูกอ่านตอนเริ่มจากตัวแปรหรือคำสั่ง และไม่เคยถูกเก็บไว้
updated: 2026-10-26
order: 13
translation: machine
sourceHash: "771d749f1f90"
---

บัญชี API key เริ่ม Claude Code ด้วย API key แทนการเข้าสู่ระบบแบบสมัครสมาชิก Claude คีย์นั้นใช้กับ API ของ Anthropic หรือกับ endpoint ใดก็ได้ที่ใช้รูปแบบ Messages เดียวกัน เช่น เกตเวย์ LiteLLM หรือ OpenRouter ccam ไม่เก็บคีย์ไว้ แต่จะอ่านคีย์ทุกครั้งที่บัญชีเริ่มทำงาน จากตัวแปรสภาพแวดล้อมหรือจากคำสั่ง เช่น โปรแกรมจัดการรหัสผ่าน และส่งให้เฉพาะโปรเซส Claude Code นั้นเท่านั้น

## การใช้งาน

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

บัญชี API key เข้าร่วมกลุ่มได้เหมือนบัญชีอื่น (ค่าเริ่มต้นคือกลุ่ม `default` ซึ่งใช้การตั้งค่า `CLAUDE.md` ส่วนขยาย และเซิร์ฟเวอร์ MCP ร่วมกัน) และใช้งานได้กับคำสั่งลัด `ccam folder` และ `ccam ext`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--api` | เพิ่มบัญชี API key แทนการเข้าสู่ระบบ | ปิด |
| `--key-from env:NAME` | อ่านคีย์จากตัวแปรสภาพแวดล้อม `NAME` | |
| `--key-from command:COMMAND` | รัน `COMMAND` (ด้วย `sh -c` หรือบน Windows ตามที่พิมพ์ไว้ใน `cmd /d /s /c "…"` โดยคงเครื่องหมายคำพูดไว้) และใช้บรรทัดแรกที่พิมพ์ออกมา | |
| `--base-url <url>` | endpoint (`ANTHROPIC_BASE_URL`) | API ของ Anthropic |
| `--auth api-key\|bearer` | ส่งคีย์เป็น `ANTHROPIC_API_KEY` (`x-api-key`) หรือเป็น `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer` ซึ่งเกตเวย์มักต้องการ) | `api-key` |
| `--model <model>` | โมเดล (`ANTHROPIC_MODEL`) | ค่าเริ่มต้นของ Claude Code |

การตั้งค่าถูกเก็บไว้ใน `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## วิธีการทำงาน

ทุกครั้งที่บัญชีเริ่มทำงาน หลังการตรวจตามปกติและก่อนที่ Claude Code จะเริ่ม:

1. ccam อ่านคีย์: `env:` จากสภาพแวดล้อม ส่วน `command:` โดยรันคำสั่งซึ่งหมดเวลาที่ 30 วินาที คำสั่งสามารถขอรหัสผ่านทางเทอร์มินัลได้ และเอาต์พุตข้อผิดพลาดของคำสั่งก็ไปที่เทอร์มินัลด้วย คีย์ที่ว่างจะหยุดการเริ่มพร้อมข้อความที่บอกว่าคาดว่าคีย์จะอยู่ที่ใด
2. ccam ลบตัวแปรที่สืบทอดมาทุกตัวที่กำหนดว่า Claude Code จะเข้าสู่ระบบอย่างไรหรือใช้ผู้ให้บริการรายใด:
   - คีย์และโทเค็น: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`
   - endpoint: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`
   - สวิตช์ผู้ให้บริการ: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`

   คีย์ที่ตั้งไว้ใน shell เพื่อใช้กับอย่างอื่นจึงไปไม่ถึงบัญชีนี้ และการเข้าสู่ระบบแบบสมัครสมาชิกก็ไม่สามารถมีลำดับความสำคัญเหนือคีย์ของบัญชีนี้ได้
3. ccam ตั้ง `ANTHROPIC_BASE_URL` และ `ANTHROPIC_MODEL` เมื่อระบุไว้ และตั้งคีย์เป็น `ANTHROPIC_API_KEY` หรือ `ANTHROPIC_AUTH_TOKEN`

บัญชีที่เข้าสู่ระบบแบบสมัครสมาชิกยังคงใช้สภาพแวดล้อมที่สืบทอดมาเหมือนเดิม ในเซสชันแบบโต้ตอบ Claude Code อาจถามหนึ่งครั้งว่าจะใช้คีย์ที่พบในสภาพแวดล้อมหรือไม่

## ความปลอดภัยและกรณีพิเศษ

- คีย์ไม่เคยถูกเขียนลงใน `config.toml` แคช หรือล็อก และ `ccam --dry-run gw` ไม่อ่านคีย์: การทดลองรันจะไม่รันแหล่ง `command:` การทดลองรันจะแสดง endpoint โมเดล และแหล่งที่มาของคีย์แทน:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  ตัวแปรที่สืบทอดมาซึ่งการเปิดลบออกจะแสดงเป็น `env -u <name>` ที่ส่วนต้น เมื่อแหล่ง `env:` มีตัวแปรที่ว่างเปล่า การทดลองรันจะล้มเหลวแบบเดียวกับการเปิดจริง
- `ccam add --api` ตรวจแหล่ง `env:` ทันทีและเตือนหากตัวแปรว่าง ส่วนแหล่ง `command:` จะยังไม่ถูกรันจนกว่าจะเริ่มบัญชี
- `ccam login` และ `ccam logout` ปฏิเสธบัญชี API key เพราะบัญชีนี้ไม่มีการเข้าสู่ระบบ `ccam doctor` แสดง endpoint และแหล่งของคีย์ และเตือนเมื่อตัวแปร `env:` ว่าง
- `ccam usage` แสดง `no limits (API)` สำหรับบัญชีนี้ `ccam next --auto`, `--fallback` และ `ccam pick` เลือกจากบัญชีแบบสมัครสมาชิกเท่านั้น เพราะ API key ไม่มีขีดจำกัด 5 ชั่วโมงหรือรายสัปดาห์ให้เปรียบเทียบ
- `settings.json` ที่ใช้ร่วมกันซึ่งมีบล็อก `env` หรือ `apiKeyHelper` จะมีผลกับทุกบัญชีในกลุ่ม รวมถึงบัญชี API key ด้วย อย่าใส่การตั้งค่าผู้ให้บริการไว้ในการตั้งค่าที่ใช้ร่วมกัน หรือนำบัญชี API key ไปไว้ในกลุ่มที่ไม่มี `settings`

## หัวข้อที่เกี่ยวข้อง

- [การจัดการบัญชี](/features/accounts)
- [การเริ่ม Claude Code](/features/launch)
- [ข้อมูลอ้างอิงคำสั่ง](/docs/commands)
