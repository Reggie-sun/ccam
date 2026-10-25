---
title: การเลือกบัญชีสำหรับสคริปต์
description: ccam pick พิมพ์บัญชีที่เหลือโควตามากที่สุดให้สคริปต์ ให้ worker ที่เริ่มพร้อมกันได้บัญชีต่างกัน และออกด้วย 75 เมื่อไม่มีบัญชีใดเหลือโควตา
updated: 2026-10-25
order: 12
translation: machine
sourceHash: "5c4f74db9714"
---

`ccam pick` พิมพ์ชื่อของบัญชีที่เข้าสู่ระบบแล้วซึ่งเหลือโควตามากที่สุด เพื่อให้สคริปต์ใช้เริ่ม Claude Code worker ที่เริ่มพร้อมกันจะได้บัญชีคนละบัญชีตราบใดที่มีบัญชีที่เหลือโควตามากพอ งาน `claude -p` ชุดหนึ่งจึงกระจายไปทั่วบัญชีของคุณแทนที่จะใช้บัญชีเดียวจนหมด

**ที่นั่ง (seat) ของ Team และ Enterprise** การเลือกอัตโนมัติ (`ccam next`, `--fallback`, `ccam pick`, ปุ่ม `n` ในแดชบอร์ด และ [`auto_switch`](/features/next#auto-switch) เมื่อคุณเปิดใช้) จะสลับระหว่างที่นั่งต่าง ๆ ในองค์กร Claude Team หรือ Enterprise เดียวกันด้วย ศูนย์ช่วยเหลือของ Anthropic ระบุว่า "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (แต่ละที่นั่งในแพ็กเกจ Team จะกำหนดให้กับบุคคลที่ระบุชื่อเพียงคนเดียว และมีเพียงบุคคลนั้นที่ใช้ได้) และ "Users can't route their work through another user's account to use that user's remaining usage" (ผู้ใช้ไม่สามารถส่งงานของตนผ่านบัญชีของผู้ใช้คนอื่นเพื่อใช้ปริมาณการใช้งานที่เหลืออยู่ของผู้ใช้คนนั้นได้) ([แพ็กเกจ Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ส่วน [แพ็กเกจ Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ก็เช่นเดียวกัน) ดังนั้น การใช้การเลือกอัตโนมัติระหว่างที่นั่งต่าง ๆ จึงอาจละเมิดข้อกำหนดของ Anthropic จะใช้การเลือกอัตโนมัติในกรณีนี้หรือไม่เป็นการตัดสินใจของคุณ: โปรดสอบถามองค์กรของคุณและ Anthropic ก่อนที่จะพึ่งพาการเลือกอัตโนมัติ สำหรับแพ็กเกจส่วนตัว (Pro, Max) ดูที่ [หมายเหตุเรื่องข้อกำหนดในคำถามที่พบบ่อย](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)

## การใช้งาน

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

รหัสออก: `0` พร้อมชื่อบัญชีทางเอาต์พุตมาตรฐาน `75` (`EX_TEMPFAIL`) เมื่อไม่มีบัญชีใดเหลือโควตาหรือไม่มีบัญชีใดเข้าสู่ระบบ และ `1` สำหรับข้อผิดพลาดอื่น เช่น แฟล็กไม่ถูกต้อง

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--json` | พิมพ์บัญชี การใช้งาน 5 ชั่วโมงและรายสัปดาห์ แหล่งที่มาของตัวเลข และเหตุผลที่เลือกบัญชีนั้น | ปิด |
| `--model <model>` | นับช่วงรายสัปดาห์ของโมเดลนั้นด้วย เช่น `opus` | ไม่มี |
| `--max-usage <1-100>` | เลือกบัญชีที่ใช้งานต่ำกว่าค่านี้ก่อน | `auto_max_usage` หากไม่มีคือ 100 |
| `--strategy headroom\|expiring\|balance` | การใช้งาน 5 ชั่วโมงต่ำที่สุด ช่วง 5 ชั่วโมงที่จะรีเซ็ตเร็วที่สุด หรือโควตารายสัปดาห์ที่จะถูกปล่อยทิ้งโดยไม่ได้ใช้ก่อน | `auto_strategy` หากไม่มีคือ `headroom` |
| `--cross-org` | เลือกบัญชีขององค์กรอื่นด้วย | `cross_org` หากไม่มีคือปิด |
| `--wait[=<duration>]` | เมื่อไม่มีบัญชีใดเหลือโควตา ให้รอบัญชีแรกที่รีเซ็ต นานไม่เกินค่านี้ ค่าต้องอยู่หลัง `=`: `--wait 30m` จะถูกปฏิเสธพร้อมคำแนะนำ `--wait takes its value after '=': --wait=30m` | ปิด; `5h` เมื่อระบุโดยไม่มีค่า |
| `--dry-run` | เลือกโดยไม่บันทึกการเลือก | ปิด |

## วิธีการทำงาน

1. ccam แสดงรายการบัญชีที่เข้าสู่ระบบแล้ว (บัญชี API key และการเข้าสู่ระบบ Console ไม่มีขีดจำกัดจึงไม่ถูกนับรวม) ซึ่งอยู่ในองค์กรเดียวกับบัญชีที่โฟลเดอร์นี้จะเริ่มใช้ (ดู[องค์กร](/features/next#auto)) โฟลเดอร์ที่เข้าสู่ระบบด้วยบัญชี Claude เดียวกันกับบัญชีนั้นก็จะไม่ถูกนับรวมด้วย
2. ccam ล็อก (`~/.local/state/ccam/picks.lock`) เพื่อให้การเลือกที่รันพร้อมกันเกิดขึ้นทีละรายการ
3. ccam สอบถามการใช้งานของทุกบัญชีพร้อมกันแบบเดียวกับ `ccam usage` โดยใช้แคช 60 วินาที ช่วงที่ Claude Code รายงานว่าเต็มจะนับเป็น 100% จนกว่าจะรีเซ็ต ([บันทึกขีดจำกัด](/features/usage#the-limit-record)) บัญชีที่เพิ่งถึงขีดจำกัดจึงยังถูกตัดออก แม้ตัวเลขที่บันทึกไว้จะแสดงว่ายังเหลือโควตา
4. บัญชีที่ถูกเลือกภายใน 2 นาทีที่ผ่านมา ไม่ว่าจะโดย `ccam pick` หรือโดย [`auto_switch`](/features/next#auto-switch) ที่ย้ายเซสชันไปที่นั่น (`~/.local/state/ccam/picks/<name>` โดยการเลือกที่ลงเวลาไว้หลังเวลาปัจจุบันของนาฬิกาจะไม่นับ) จะถูกตัดออกก่อน แล้ว[กฎการเลือก](/features/next#auto)ของ `ccam next --auto` จะเลือกจากบัญชีที่เหลือ: บัญชีที่ต่ำกว่า `auto_max_usage` มาก่อน ตามกลยุทธ์ และไม่เลือกบัญชีที่อยู่ที่ 100% เด็ดขาด เมื่อการตัดออกนี้ข้ามบัญชีที่กฎจะเลือกเป็นอันดับแรก `reason` ของ `--json` จะบอกไว้: `30% of its 5 hours used; work was picked a moment ago`
5. หากทุกบัญชีที่เหลือโควตาเพิ่งถูกเลือกไป กฎจะเลือกจากบัญชีทั้งหมดอีกครั้ง worker ตัวที่สี่บนสามบัญชีจึงใช้บัญชีร่วมกับตัวอื่นหนึ่งบัญชี
6. ccam บันทึกการเลือกและปลดล็อก จากนั้นจึงพิมพ์ชื่อ

ช่วง 2 นาทีมีความสำคัญ เพราะการใช้งานของงานใหม่ต้องใช้เวลาสักพักกว่าจะปรากฏในตัวเลข หากไม่มีช่วงนี้ worker ห้าตัวที่เริ่มพร้อมกันจะเห็นบัญชีเดียวกันว่าเป็นบัญชีที่เหลือโควตามากที่สุดทั้งหมด

## ความปลอดภัยและกรณีพิเศษ

- `ccam pick` อ่านแค่การใช้งานและเขียนเวลาที่เลือก ไม่เคยต่ออายุหรืออ่านโทเค็นเกินกว่าที่ `ccam usage` ทำ
- การเลือกไม่ได้จองบัญชีไว้ สคริปต์สองตัวอาจยังเริ่มบนบัญชีเดียวกันหลังผ่านไป 2 นาที ใช้ `--fallback` เพื่อกู้คืนเมื่อบัญชีนั้นถึงขีดจำกัด
- การรอด้วย `--wait` ไม่ถือล็อกไว้ การเลือกอื่นจึงดำเนินต่อไปได้ระหว่างนั้น
- `ccam pick` ไม่สนใจ [`ccam folder`](/features/folder) โดยเลือกตามการใช้งานเท่านั้น

## หัวข้อที่เกี่ยวข้อง

- [การสลับไปยังบัญชีถัดไป](/features/next)
- [รันซ้ำเมื่อถึงขีดจำกัด](/features/next#run-again-at-the-limit)
- [ขีดจำกัดการใช้งาน](/features/usage)
- [บัญชี API key](/features/api)
