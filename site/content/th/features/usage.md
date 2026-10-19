---
title: ขีดจำกัดการใช้งาน
description: "ccam usage แสดงขีดจำกัด 5 ชั่วโมงและรายสัปดาห์ของทุกบัญชีพร้อมกัน และ ccam statusline แสดงชื่อบัญชีและการใช้งานใน Claude Code"
updated: 2026-10-19
order: 7
translation: machine
sourceHash: 8306a8f6d950
---

`ccam usage` แสดงว่าแต่ละบัญชีใช้ขีดจำกัดแบบ 5 ชั่วโมงและรายสัปดาห์ไปแล้วเท่าไร พร้อมเวลาที่เหลือจนกว่าแต่ละช่วงจะรีเซ็ต ใช้ก่อนเริ่มงานที่ใช้เวลานาน หรือเมื่อ Claude Code แจ้งว่าใกล้ถึงขีดจำกัด เพื่อดูว่าบัญชีใดยังเหลือโควตา `ccam statusline` คือคำสั่งแถบสถานะสำหรับ Claude Code: แสดงชื่อบัญชีและการใช้งานภายใน Claude Code และบันทึกตัวเลขไว้ให้ `ccam usage`

## การใช้งาน

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

บนเทอร์มินัล ตารางจะปรากฏทันทีโดยมี `fetching…` ในแต่ละแถว และแต่ละแถวจะถูกวาดใหม่เมื่อบัญชีนั้นได้รับคำตอบ `As of` บอกอายุของตัวเลข ส่วน `(cached)` หมายความว่าตัวเลขมาจากไฟล์ที่บันทึกไว้ ไม่ใช่จาก endpoint บรรทัดสุดท้ายระบุบัญชีที่ใช้งาน 5 ชั่วโมงน้อยที่สุดและยังไม่ถึง 100% และจะแสดงเฉพาะเมื่อคุณมีมากกว่าหนึ่งบัญชี บัญชีที่ยังไม่ได้เข้าสู่ระบบจะแสดง `not logged in → ccam login <name>` และบัญชีที่ใช้ API key จะแสดง `no limits (API)`

หลังจาก `ccam statusline install` แล้ว Claude Code จะแสดงบรรทัดอย่าง `[work] · Opus · 5h 24% · 7d 41%`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--json` | พิมพ์รายงานเป็น JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | ปิด |
| `--watch[=N]` | วาดใหม่ทุก `N` วินาทีจนกว่าคุณจะกด Ctrl+C `--watch` อย่างเดียว หรือค่าที่เป็น 0 หรือน้อยกว่า หมายถึง 60 | ปิด |
| `-v`, `--verbose` | แสดงขีดจำกัดรายสัปดาห์แยกตามโมเดลและบรรทัด `source:` (`live`, `cache` หรือ `statusline`) | ปิด |
| `--no-cache` | ไม่ใช้แคช 60 วินาทีและสอบถาม endpoint ใหม่ | ปิด |
| `--source live\|cache` | `cache` จะไม่สอบถาม endpoint เลยและแสดงเฉพาะตัวเลขที่บันทึกไว้ ส่วน `live` คือพฤติกรรมเริ่มต้น | live |
| `--refresh` | ต่ออายุ access token ที่หมดอายุ เฉพาะบัญชีที่ไม่ได้ทำงานอยู่ | ปิด |
| `--dry-run` | ยอมรับ: `ccam usage` ไม่เปลี่ยนแปลงสิ่งใดนอกจากแคชของตัวเอง แต่จะถูกปฏิเสธเมื่อใช้ร่วมกับ `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | ปิด |

ค่าอื่นของ `--source` จะล้มเหลวด้วย `--source must be live or cache` `ccam usage` ไม่รับอาร์กิวเมนต์

`ccam statusline` อ่าน JSON ของแถบสถานะจาก Claude Code ทางอินพุตมาตรฐาน มีตัวเลือกเดียวคือ `--wrap <command>`: คำสั่งแถบสถานะที่มีอยู่ ซึ่งบรรทัดแรกของผลลัพธ์จะมีชื่อบัญชีนำหน้า `ccam statusline install` ไม่รับตัวเลือก ทั้งสองคำสั่งปฏิเสธ `--dry-run` ด้วย `` --dry-run is not supported by `ccam statusline`; nothing was changed ``

## วิธีการทำงาน

ตัวเลขมาจากสองแหล่ง: endpoint การใช้งานที่ไม่มีเอกสารของ Anthropic `GET https://api.anthropic.com/api/oauth/usage` ซึ่งเรียกด้วยการเข้าสู่ระบบของแต่ละบัญชีเอง และค่าที่ Claude Code ส่งให้คำสั่งแถบสถานะ ทั้งหมดอยู่ใน `internal/usage`

1. `cli.usageCmd` โหลด `config.toml` และสร้าง `usage.Account` หนึ่งรายการต่อบัญชี ได้แก่ โฟลเดอร์ `~/.local/share/ccam/profiles/<name>/` โฟลเดอร์เซสชัน (เพื่อดูว่ากำลังทำงานอยู่หรือไม่) และบัญชีนั้นถูกนำเข้าในตำแหน่งเดิมหรือไม่
2. `Fetcher.FetchAll` เริ่ม goroutine หนึ่งตัวต่อบัญชี ทุกตัวใช้ `http.Client` ตัวเดียวกัน การเชื่อมต่อไปยัง endpoint จึงถูกนำกลับมาใช้ การเชื่อมต่อใช้เวลาได้ 2 วินาที และแต่ละบัญชีมีเวลารวม 3 วินาที ไคลเอนต์ไม่ตาม redirect เพราะ redirect จะส่งโทเค็นซ้ำไปยังที่อื่น และอ่านคำตอบไม่เกิน 1 MB
3. เมื่อแต่ละบัญชีเสร็จ แถวของบัญชีนั้นจะถูกวาดใหม่ บัญชีที่ล้มเหลวหรือไม่ตอบกลับทันเวลาจะหันไปใช้ตัวเลขที่บันทึกไว้ บัญชีที่ช้าบัญชีเดียวจึงไม่ทำให้บัญชีอื่นต้องรอ

แต่ละบัญชีทำงานใน goroutine ของตัวเองกับ endpoint บัญชีที่หมดเวลาจะแสดงตัวเลขที่บันทึกไว้แทน

<!-- motion:usage-fetch -->

`personal` และ `work` ตอบกลับทันเวลาและแสดงตัวเลขสด `side` ใช้เวลานานกว่า 3 วินาที `fallback()` จึงอ่านไฟล์แคชสองไฟล์ของบัญชีนี้และใช้ไฟล์ที่เขียนล่าสุด `usage.Table()` วาดแถว และ `usage.Best()` เลือกบัญชีที่ระบุไว้ในบรรทัดสุดท้าย

สำหรับบัญชีหนึ่ง `Fetcher.Fetch` ตรวจตามลำดับต่อไปนี้:

<!-- motion:usage-cache -->

1. เมื่อใช้ `--source cache` จะไปที่ `fallback()` โดยตรง
2. หากไม่ได้ตั้ง `--no-cache` ผลลัพธ์สดใน `~/.local/state/ccam/usage/<name>.json` ที่อายุไม่ถึง 60 วินาที (`CacheTTL`) จะถูกคืนตามเดิม การรัน `ccam usage` สองครั้งติดกันจึงใช้คำขอเพียงหนึ่งครั้งต่อบัญชี
3. `readCredentials()` อ่าน `.credentials.json` (หรือพวงกุญแจ (Keychain) ของ macOS) โดยไม่เปลี่ยนแปลง และถอดรหัสเฉพาะคีย์ `claudeAiOauth` หากไม่มีการเข้าสู่ระบบแบบสมัครสมาชิก จะแสดง `not logged in` หรือ `no limits (API)` เมื่อ `.claude.json` มี API key
4. access token ที่หมดอายุจะไม่ถูกใช้ หากไม่ใส่ `--refresh` บัญชีจะแสดงตัวเลขที่บันทึกไว้พร้อมระบุว่าหมดอายุ
5. ในกรณีอื่น ccam จะเรียก endpoint สถานะ 401 หรือ 403 นับเป็นการเข้าสู่ระบบที่หมดอายุ ส่วนสถานะอื่นหรือเนื้อหาที่อ่านไม่ได้นับเป็นข้อผิดพลาด และทั้งสองกรณีจะหันไปใช้ตัวเลขที่บันทึกไว้
6. คำตอบที่ถูกต้องจะถูกเขียนลงใน `<name>.json` พร้อมเวลาที่วัด

`fallback()` อ่านทั้ง `<name>.json` และ `<name>.statusline.json` และเก็บไฟล์ที่ใหม่กว่า ก่อนแสดงตัวเลขที่บันทึกไว้ จะตรวจเวลารีเซ็ตของแต่ละช่วง: ช่วงที่รีเซ็ตไปแล้วหลังจากบันทึกตัวเลขจะแสดงเป็น 0% เพราะเปอร์เซ็นต์เดิมใช้ไม่ได้แล้ว นี่คือเหตุผลที่ `personal` ด้านบนแสดง `0%` โดยไม่มีเวลารีเซ็ต

endpoint มีรูปแบบคำตอบสองแบบ แบบใหม่มีอาร์เรย์ `limits` ที่มีรายการ `session` (5 ชั่วโมง) `weekly_all` และ `weekly_scoped` (แยกตามโมเดล) แบบเก่ามีออบเจ็กต์ `five_hour`, `seven_day`, `seven_day_opus` และ `seven_day_sonnet` `ParseLive` อ่านได้ทั้งสองแบบ: อาร์เรย์ `limits` มีผลเหนือกว่า และออบเจ็กต์แบบเก่าจะเติมส่วนที่อาร์เรย์ไม่ครอบคลุม เนื้อหาที่ไม่มีขีดจำกัดเลยจะถูกรายงานว่ารูปแบบเปลี่ยนไป ซึ่ง `ccam doctor` จะแจ้งให้คุณทราบ

### การต่ออายุโทเค็น

refresh token มีการหมุนเวียน: การต่ออายุหนึ่งครั้งจะได้ refresh token ใหม่และทำให้ตัวเดิมใช้ไม่ได้ หาก ccam ต่ออายุการเข้าสู่ระบบในขณะที่ Claude Code กำลังทำงานด้วยบัญชีเดียวกัน เซสชันนั้นจะถูกออกจากระบบ ดังนั้น `ccam usage` จึงอ่านโทเค็นแต่ไม่เคยเขียน เว้นแต่คุณจะใส่ `--refresh` จากนั้น `Fetcher.refresh` จะ:

1. ปฏิเสธบัญชีที่ลงทะเบียนไว้ในตำแหน่งเดิมด้วย `ccam import` เพราะ `claude` แบบปกติอาจใช้การเข้าสู่ระบบนั้นโดยที่ ccam ไม่รู้
2. ปฏิเสธการเข้าสู่ระบบที่เก็บไว้ในพวงกุญแจ (Keychain) ของ macOS
3. ถือไฟล์ล็อก `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (อยู่ภายใต้ `~/.local/state/ccam/locks/` เมื่อไม่ได้ตั้งตัวแปรนั้น) เพื่อไม่ให้โปรเซส ccam สองตัวต่ออายุโทเค็นเดียวกันซ้ำสองครั้ง
4. ภายใต้ล็อก จะปฏิเสธหากมีเซสชันของบัญชีกำลังทำงาน และข้ามคำขอหากโปรเซสอื่นต่ออายุโทเค็นไปแล้ว
5. ส่ง refresh token ไปยัง endpoint ของโทเค็น และเขียนโทเค็นใหม่กลับลงใน `.credentials.json` ด้วย `fsutil.AtomicWrite` โดยคงคีย์อื่นทั้งหมดในไฟล์ไว้

### ตัวเก็บข้อมูลจากแถบสถานะ

Claude Code รันคำสั่งแถบสถานะและส่งเอกสาร JSON ให้ทางไปป์ ccam อ่านชื่อโมเดลและ `rate_limits` จากเอกสารนั้น:

<!-- motion:usage-statusline -->

1. `runStatusline` หาบัญชีจาก `CCAM_PROFILE` ซึ่ง ccam ตั้งไว้เมื่อเริ่ม Claude Code หรือหากไม่มีจะจับคู่ `CLAUDE_CONFIG_DIR` กับโฟลเดอร์ของบัญชี
2. `ParseStatusLine` อ่าน `rate_limits.five_hour` และ `rate_limits.seven_day` JSON จะไม่มีขีดจำกัดอัตราสำหรับบัญชีที่ใช้ API key หรือก่อนคำตอบแรกของเซสชัน ในกรณีนั้นจะไม่บันทึกสิ่งใด
3. `SaveStatusLine` เขียน `~/.local/state/ccam/usage/<name>.statusline.json` ซึ่ง `ccam usage` จะใช้เมื่อ endpoint ไม่ตอบกลับ
4. ccam พิมพ์ `[name] · model · 5h N% · 7d N%` เมื่อใช้ `--wrap` จะรันคำสั่งของคุณผ่าน `sh -c` ด้วยอินพุตเดียวกัน แล้วพิมพ์ `[name]` ตามด้วยบรรทัดแรกของคำสั่งแทน

คำสั่งที่ถูกครอบมีเวลา 2 วินาที หลังจากนั้นจะถูกหยุด และบรรทัดจะแสดงเพียงชื่อบัญชี คำสั่งที่ช้าจึงไม่ทำให้แถบสถานะของ Claude Code ค้าง

`ccam statusline install` เขียน `statusLine` ลงใน `settings.json` ของแต่ละบัญชี เมื่อบัญชีมีคำสั่งแถบสถานะอยู่แล้ว ccam จะเก็บไว้เป็น `--wrap '<old command>'` `settings.json` ที่ใช้ร่วมกันผ่านกลุ่มจะถูกเขียนเพียงครั้งเดียวผ่านลิงก์ ทุกบัญชีในกลุ่มจึงได้รับด้วย การรันซ้ำจะพิมพ์ `(already set)` สำหรับไฟล์ที่ใช้ ccam อยู่แล้ว

## ความปลอดภัยและกรณีพิเศษ

- `ccam usage` ไม่เคยเปลี่ยนแปลงไฟล์โทเค็น เว้นแต่คุณจะใส่ `--refresh` เมื่อการเข้าสู่ระบบหมดอายุ ให้เริ่มบัญชีนั้นหนึ่งครั้งด้วย `ccam <name>` แล้ว Claude Code จะต่ออายุให้
- `--refresh` ไม่เคยแตะบัญชีที่กำลังทำงาน บัญชีที่นำเข้า หรือการเข้าสู่ระบบในพวงกุญแจ (Keychain) การต่ออายุที่ถูกปฏิเสธจะแสดงตัวเลขที่บันทึกไว้พร้อมเหตุผล
- endpoint การใช้งานไม่มีเอกสารเผยแพร่สู่สาธารณะและอาจเปลี่ยนแปลงได้ เมื่อเปลี่ยน คุณจะยังเห็นตัวเลขที่บันทึกไว้ และ `ccam doctor` จะรายงานการเปลี่ยนรูปแบบ
- เมื่อไม่มีตัวเลขที่บันทึกไว้และไม่มีคำตอบ แถวจะแสดง `lookup failed:` พร้อมเหตุผล เช่น `no saved usage` เมื่อใช้ `--source cache`
- ไฟล์แคชเก็บเพียงเปอร์เซ็นต์และเวลารีเซ็ต ไม่เคยเก็บโทเค็น

## หัวข้อที่เกี่ยวข้อง

- [ขีดจำกัดการใช้งานและการสลับบัญชี](/docs/usage-limits)
- [การสลับไปยังบัญชีถัดไป](/features/next)
- [ขั้นตอนการรับมือเมื่อถึงขีดจำกัดการใช้งาน](/guides/usage-limit-workflow)
- [คำสั่ง](/docs/commands)
