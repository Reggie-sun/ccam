---
title: คำสั่งลัด
description: "ให้บัญชีมีคำสั่งของตัวเอง เช่น ccwork ด้วย ccam shortcut และสร้างลิงก์ใน ~/.local/bin ขึ้นใหม่ด้วย ccam shims install"
updated: 2026-10-26
order: 3
translation: machine
sourceHash: "066e7e622517"
---

คำสั่งลัดคือคำสั่งที่คุณตั้งชื่อเองเพื่อเริ่ม Claude Code ในชื่อบัญชีหนึ่ง: `ccwork --resume` ทำงานเหมือนกับ `ccam work --resume` คำสั่งลัดเป็นลิงก์ไปยัง ccam ใน `~/.local/bin` จึงใช้ได้ใน bash, zsh, สคริปต์ และทุกอย่างที่รันโปรแกรมจาก `PATH` คำสั่งลัดเป็นตัวเลือกเสริม ให้สร้างไว้สำหรับบัญชีที่คุณเริ่มบ่อย

## การใช้งาน

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccam ยืนยันการเปลี่ยนแปลงทุกครั้ง:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

บัญชีหนึ่งมีคำสั่งลัดได้ไม่เกินหนึ่งคำสั่ง ซึ่งเก็บเป็น `shortcut` ใน `~/.config/ccam/config.toml` ตัวช่วยตั้งค่าของ `ccam add` จะถาม `Create a shortcut command for <name>? [y/N]` แล้วตามด้วย `Shortcut name [cc<name>]:` นอกตัวช่วยตั้งค่า `ccam add` จะสร้างคำสั่งลัดเฉพาะเมื่อใส่ `--shortcut` ส่วน `ccam list` จะแสดงคำสั่งลัดของแต่ละบัญชี หรือแสดง `-`

คำสั่งลัดส่งทุกอาร์กิวเมนต์ให้ Claude Code แบบเดียวกับที่ `ccam <name>` ทำกับอาร์กิวเมนต์หลังชื่อบัญชี ซึ่งรวมถึง `--dry-run` และ `--help` ด้วย: หากต้องการดูตัวอย่างการเริ่ม ให้ใช้ `ccam --dry-run work`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--remove` | ลบคำสั่งลัดของบัญชีและลิงก์ของคำสั่งลัด | ปิด |
| `-h`, `--help` | ความช่วยเหลือของคำสั่ง | |
| `--dry-run` | `ccam shortcut` รองรับตัวเลือกนี้: จะตรวจชื่อและพิมพ์ `Would set the shortcut of work to ccwork.` หรือ `Would remove ccwork, the shortcut of work.` โดยไม่บันทึกหรือสร้างลิงก์ | ปิด |

`ccam shims install` ไม่รับอาร์กิวเมนต์ และปฏิเสธ `--dry-run` ด้วย `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``

กฎของชื่อ ซึ่งตรวจตามลำดับนี้:

- ยาว 1 ถึง 64 อักขระ ประกอบด้วย `a-z`, `0-9`, `.`, `-` และ `_` และขึ้นต้นด้วยตัวอักษรหรือตัวเลข
- ไม่ใช่ `ccam`, `claude` หรือคำสั่งของ ccam เช่น `run` หรือ `usage`
- ไม่ใช่คำสั่งลัดของบัญชีอื่น
- ไม่ใช่ไฟล์ที่มีอยู่แล้วใน `~/.local/bin` เว้นแต่จะเป็นลิงก์ที่ ccam สร้างไว้ (เช่น คำสั่งลัดของบัญชีนี้เอง)
- ไม่ใช่ชื่อของโปรแกรมอื่นใน `PATH` ของคุณ

`config.toml` จะถูกตรวจตามกฎเดียวกันเมื่อโหลด หากคุณแก้ไขคำสั่งลัดด้วยมือจนไม่ตรงตามรูปแบบ เป็น `claude` หรือ `ccam` หรือถูกใช้โดยสองบัญชี ทุกคำสั่งจะล้มเหลวพร้อมข้อความที่ระบุชื่อบัญชี จนกว่าคุณจะแก้ไข

การระบุคำสั่งพร้อมกับ `--remove` จะหยุดด้วย `give a command or --remove, not both` การตั้งชื่อที่บัญชีมีอยู่แล้วจะพิมพ์ `work: ccwork (unchanged)` และสร้างลิงก์ใหม่หากลิงก์หายไป

## วิธีการทำงาน

การตั้งชื่อจะรัน `checkShortcut` ซึ่งเรียก `shim.Check` พร้อมโฟลเดอร์ bin และพาธจริงของไบนารี ccam:

<!-- motion:shortcuts-check -->

ชื่อจะผ่านการตรวจลงไปทีละขั้นตามลำดับ ข้างการตรวจแต่ละขั้นคือชื่อที่ถูกปฏิเสธ: `Work` ไม่ตรงรูปแบบ `run` เป็นคำสั่งของ ccam `ccside` เป็นของบัญชีอื่น `ccfile` เป็นไฟล์ที่มีอยู่ใน `~/.local/bin` และ `git` จะบัง `/usr/bin/git` มีเพียงชื่อที่ผ่านการตรวจทุกขั้นเท่านั้นที่จะถูกบันทึกและสร้างลิงก์

1. ccam จับคู่ชื่อกับ `^[a-z0-9][a-z0-9._-]{0,63}$` จากนั้นกับ `config.Reserved` แล้วจึงกับ `cfg.FindShortcut` (`shim.Check`) จะเรียก `os.Lstat` กับ `~/.local/bin/<name>` และสำหรับทุกโฟลเดอร์อื่นใน `PATH` จะเรียก `os.Stat` และปฏิเสธไฟล์ปฏิบัติการที่ไม่ได้เป็นลิงก์ไปยัง ccam
2. ccam บันทึก `config.toml` ด้วย `fsutil.AtomicWrite` (ไฟล์ชั่วคราว `fsync` แล้ว rename)
3. ccam ลบลิงก์ของคำสั่งลัดเดิม แต่เฉพาะเมื่อเป็นลิงก์ไปยัง ccam (`shim.Remove`)
4. ccam ไล่ดูทุกคำสั่งลัดที่ตั้งค่าไว้ตามลำดับชื่อ และสร้างลิงก์ของคำสั่งลัดที่ยังไม่มี (`shim.Install`) ด้วย `os.Symlink` ไปยังไบนารี ccam ซึ่งหาพาธด้วย `os.Executable` และ `filepath.EvalSymlinks` ไฟล์ที่ขวางคำสั่งลัดหนึ่งอยู่จะไม่ทำให้คำสั่งลัดอื่นหยุดไปด้วย โดยข้อผิดพลาดจะระบุแต่ละรายการ

ccam บันทึกทุกลิงก์ที่สร้าง พร้อมไบนารีที่แต่ละลิงก์รัน ไว้ใน `~/.local/state/ccam/shortcuts.json` ลิงก์จะนับเป็นของ ccam เมื่อ `os.Readlink` คืนพาธของ ccam หรือเมื่อลิงก์และพาธของ ccam ชี้ไปที่ไฟล์เดียวกัน ลิงก์ที่บันทึกไว้ซึ่งไบนารีถูกย้ายหรือถูกแทนที่ไปแล้ว (ไบนารี ccam ที่ถูกย้ายไปโฟลเดอร์อื่น หรือการติดตั้งใหม่) ถือว่าล้าสมัย และ ccam จะชี้ลิงก์นั้นมาที่ไบนารี ccam ตัวนี้ (`shim.Install`) ลิงก์ที่บันทึกไว้ซึ่งชี้ไปยัง ccam อีกชุดที่ยังติดตั้งอยู่จะถูกปล่อยไว้ตามเดิม การลบและการตรวจก่อนเขียนทับทุกครั้งจะผ่านการทดสอบเหล่านี้ ccam จึงไม่แตะไฟล์ที่ตัวเองไม่ได้สร้าง เช่น สำเนา ccam ของคุณเองหรือลิงก์ที่คุณสร้างเอง

เมื่อคุณรันคำสั่งลัด ชื่อที่คุณพิมพ์คือสิ่งเดียวที่ ccam ต้องการ:

<!-- motion:shortcuts-link -->

เชลล์หา `ccwork` ใน `PATH` เคอร์เนลตามลิงก์และเริ่มไบนารี ccam โดย `argv[0]` ยังคงเป็น `ccwork` ใน `app.Main` ชื่อใดก็ตามที่ไม่ใช่ `ccam` จะโหลด `config.toml` หนึ่งครั้งและเรียก `shim.Lookup` ซึ่งหาบัญชีที่ `shortcut` ตรงกับ `filepath.Base(argv[0])` จากนั้น ccam จะทำขั้นตอนเดียวกับ [`ccam <name>`](/features/launch) โดยใช้การตั้งค่าที่โหลดไว้แล้ว และส่ง `argv[1:]` ให้ `claude` โดยไม่เปลี่ยนแปลง (`launcher.LaunchWith`)

ทำไมจึงใช้ลิงก์สัญลักษณ์แทน alias ของเชลล์หรือสคริปต์ครอบ: ลิงก์สัญลักษณ์ใช้ได้ในทุกเชลล์และในสคริปต์ ไม่ต้องเริ่มเชลล์เพิ่ม ซึ่งจะใช้เวลาประมาณ 1 ถึง 2 ms และการเริ่มยังคงอ่านเพียงไฟล์เดียว `scripts/perf-check.sh` วัดคำสั่งลัดด้วย `hyperfine` เทียบกับขีดจำกัด 10 ms (p95) เดียวกับ `ccam <name>`

### บน Windows

คำสั่งลัดคือ `ccwork.exe` ใน `%USERPROFILE%\.local\bin` และ ccam จะตัด `.exe` ออก (ไม่ว่าตัวพิมพ์เล็กหรือใหญ่) ก่อนค้นหา (`shim.Name`) ดังนั้น `ccwork`, `ccwork.exe` และ `CCWORK.EXE` จึงเริ่ม `work` ทั้งหมด ccam สร้างลิงก์สัญลักษณ์ก่อน ซึ่ง Windows อนุญาตใน Developer Mode หรือสำหรับผู้ดูแลระบบ หากไม่มีสิทธิ์นั้น ccam จะสร้างฮาร์ดลิงก์: ชื่อที่สองของไฟล์ `ccam.exe` ไฟล์เดียวกัน ซึ่งไม่ต้องใช้สิทธิ์พิเศษ ฮาร์ดลิงก์จะนับเป็นของ ccam เมื่อเป็นไฟล์เดียวกับ `ccam.exe`

ฮาร์ดลิงก์จะรันไฟล์ที่ใช้สร้างมันขึ้นมาเสมอ ดังนั้นเมื่อ `ccam update` วางไฟล์ใหม่แทนที่ `ccam.exe` จึงต่ออายุคำสั่งลัดที่เป็นฮาร์ดลิงก์ทุกตัวด้วย `shim.Relink` และพิมพ์ `✓ shortcut ccwork renewed` การรัน `install.ps1` อีกครั้งก็ทำแบบเดียวกันผ่าน `ccam shims install` (`OK 2 shortcuts now run the new ccam`) และ `ccam shims install` กับ `ccam doctor --fix` ก็ทำเช่นกันกับคำสั่งลัดใดก็ตามที่เป็นสำเนาของ ccam เวอร์ชันเก่า การตรวจชื่อจะมองหา `<name>` พร้อมนามสกุลแต่ละตัวใน `PATHEXT` (`.exe`, `.cmd` และอื่น ๆ) ในโฟลเดอร์อื่นของ `PATH`

ให้พิมพ์คำสั่งลัดในเชลล์ (PowerShell, cmd.exe, Git Bash) ลิงก์สัญลักษณ์ที่เปิดจาก Explorer, กล่องโต้ตอบ Run หรือ `Start-Process` จะเริ่มเป็น `ccam` ธรรมดา ซึ่งจะเปิดแดชบอร์ดแทนบัญชีนั้น

## ความปลอดภัยและกรณีพิเศษ

- **ลิงก์ที่ไม่มีบัญชีใดเป็นเจ้าของ** หากลิงก์ไปยัง ccam ใน `~/.local/bin` มีชื่อที่ไม่มีบัญชีใดใช้ เช่น หลังจากที่คุณแก้ `shortcut` ใน `config.toml` ด้วยตัวเอง ccam จะไม่หันไปใช้บัญชีเริ่มต้น แต่จะหยุดด้วย `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` ให้ตั้งชื่ออีกครั้งด้วย `ccam shortcut` หรือลบลิงก์นั้น
- **ไฟล์ที่มีอยู่แล้ว** ccam ไม่เคยแทนที่ไฟล์ที่ตัวเองไม่ได้สร้าง `ccam shims install` จะรายงาน `<path> already exists and is not a ccam link, so it was left alone` และยังคงสร้างลิงก์อื่นต่อไป ในกรณีนี้ `ccam add` และ `ccam import` จะพิมพ์ `shortcut: ccwork (not created; see above)`
- **ไบนารี ccam ที่ถูกย้าย** หากคุณย้ายไบนารี ccam ลิงก์ที่ ccam สร้างไว้จะชี้ไปยังพาธที่ไม่มีอยู่แล้ว `ccam shims install` (หรือ `ccam doctor --fix`) จะชี้ลิงก์เหล่านั้นมาที่ไบนารี ccam ที่คุณใช้รันคำสั่ง
- **การเปลี่ยนชื่อและการลบ** `ccam rename` คงชื่อคำสั่งลัดไว้ และลิงก์ยังใช้ได้ เพราะลิงก์ชี้ไปที่ ccam และ ccam ค้นหาชื่อใน `config.toml` `ccam remove` ลบลิงก์ของบัญชี ส่วน `ccam shortcut <account> <new>` จะลบลิงก์เดิมก่อนสร้างลิงก์ใหม่
- **PATH** เมื่อ ccam สร้างลิงก์และ `~/.local/bin` ไม่อยู่ใน `PATH` ของคุณ จะพิมพ์ `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` ตามด้วยบรรทัด `export PATH=...` ที่ต้องเพิ่ม บน Windows จะพิมพ์บรรทัด PowerShell ที่เพิ่มโฟลเดอร์นั้นลงใน `PATH` ของผู้ใช้แทน
- **การตรวจภายหลัง** `ccam doctor` จะรายงาน `missing shortcuts in /home/me/.local/bin: ccwork` รายงานคำสั่งลัดที่รัน ccam เวอร์ชันเก่า และเตือนเมื่อโฟลเดอร์ไม่อยู่ใน `PATH` ส่วน `ccam doctor --fix` จะสร้างลิงก์ที่หายไปและต่ออายุลิงก์ที่ล้าสมัย และจะแสดง `→ fixed` เฉพาะลิงก์ที่สร้างได้ คำสั่งลัดที่รัน ccam อีกชุดซึ่งยังติดตั้งอยู่ยังใช้งานได้ `ccam doctor` จึงแสดงเป็นเพียงข้อมูล `ccam doctor --fix` จะชี้คำสั่งลัดนั้นมาที่ไบนารี ccam ที่คุณใช้รันคำสั่ง และ `ccam add` กับ `ccam import` จะแสดงเป็น ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [การจัดการบัญชี](/features/accounts)
- [ข้อมูลอ้างอิงคำสั่ง](/docs/commands#shortcuts)
- [ติดตั้ง ccam](/docs/install)
